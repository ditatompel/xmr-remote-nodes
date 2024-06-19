package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
	"xmr-remote-nodes/internal/config"
	"xmr-remote-nodes/internal/monero"

	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
)

const RPCUserAgent = "ditatombot/0.0.1 (Monero RPC Monitoring; https://github.com/ditatompel/xmr-remote-nodes)"

const (
	errNoEndpoint         = errProber("no SERVER_ENDPOINT was provided")
	errNoTorSocks         = errProber("no TOR_SOCKS was provided")
	errNoAPIKey           = errProber("no API_KEY was provided")
	errInvalidCredentials = errProber("invalid API_KEY credentials")
)

type errProber string

func (err errProber) Error() string {
	return string(err)
}

type proberClient struct {
	endpoint  string // server endpoint
	apiKey    string // prober api key
	acceptTor bool   // accept tor
	torSOCKS  string // IP:Port of tor socks
	message   string // message to include when reporting back to server
}

func newProber() *proberClient {
	cfg := config.AppCfg()
	return &proberClient{
		endpoint:  cfg.ServerEndpoint,
		apiKey:    cfg.APIKey,
		acceptTor: cfg.AcceptTor,
		torSOCKS:  cfg.TorSOCKS,
	}
}

var ProbeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Probe remote nodes",
	Run: func(cmd *cobra.Command, args []string) {
		prober := newProber()
		if e, _ := cmd.Flags().GetString("endpoint"); e != "" {
			prober.SetEndpoint(e)
		}
		if t, _ := cmd.Flags().GetBool("no-tor"); t {
			prober.SetAcceptTor(false)
		}

		if err := prober.Run(); err != nil {
			switch err.(type) {
			case errProber:
				slog.Error(fmt.Sprintf("[PROBE] %s", err.Error()))
				os.Exit(1)
			default:
				slog.Warn(fmt.Sprintf("[PROBE] %s", err.Error()))
			}
		}
	},
}

func (p *proberClient) SetEndpoint(endpoint string) {
	p.endpoint = endpoint
}

func (p *proberClient) SetAcceptTor(acceptTor bool) {
	p.acceptTor = acceptTor
}

// Fetch a new job from the server, fetches node info, and sends it to the server
func (p *proberClient) Run() error {
	if err := p.validateConfig(); err != nil {
		return err
	}

	node, err := p.fetchJob()
	if err != nil {
		return err
	}

	fetchNode, err := p.fetchNode(node)
	if err != nil {
		return err
	}
	slog.Debug(fmt.Sprintf("[PROBE] fetchNode: %s", prettyPrint(fetchNode)))
	return nil
}

// checks if all required environment variables are set
func (p *proberClient) validateConfig() error {
	if p.endpoint == "" {
		return errNoEndpoint
	}
	if p.apiKey == "" {
		return errNoAPIKey
	}
	if p.acceptTor && p.torSOCKS == "" {
		return errNoTorSocks
	}

	return nil
}

// Get monero node info to fetch from the server
func (p *proberClient) fetchJob() (monero.Node, error) {
	queryParams := ""
	if p.acceptTor {
		queryParams = "?accept_tor=1"
	}

	var node monero.Node

	uri := fmt.Sprintf("%s/api/v1/job%s", p.endpoint, queryParams)
	slog.Info(fmt.Sprintf("[PROBE] Getting node from %s", uri))

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return node, err
	}
	req.Header.Add(monero.ProberAPIKey, p.apiKey)
	req.Header.Set("User-Agent", RPCUserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return node, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		return node, errInvalidCredentials
	default:
		return node, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	response := struct {
		Data monero.Node `json:"data"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return node, err
	}

	node = response.Data
	slog.Info(fmt.Sprintf("[PROBE] Got node: %s://%s:%d", node.Protocol, node.Hostname, node.Port))

	return node, nil
}

func (p *proberClient) fetchNode(node monero.Node) (monero.Node, error) {
	startTime := time.Now()
	endpoint := fmt.Sprintf("%s://%s:%d/json_rpc", node.Protocol, node.Hostname, node.Port)
	rpcParam := []byte(`{"jsonrpc": "2.0","id": "0","method": "get_info"}`)
	slog.Info(fmt.Sprintf("[PROBE] Fetching node info from %s", endpoint))
	slog.Debug(fmt.Sprintf("[PROBE] RPC param: %s", string(rpcParam)))

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(rpcParam))
	if err != nil {
		return node, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", RPCUserAgent)
	req.Header.Set("Origin", "https://xmr.ditatompel.com")

	var client http.Client
	if p.acceptTor && node.IsTor {
		dialer, err := proxy.SOCKS5("tcp", p.torSOCKS, nil, proxy.Direct)
		if err != nil {
			return node, err
		}
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		transport := &http.Transport{
			DialContext:       dialContext,
			DisableKeepAlives: true,
		}
		client.Transport = transport
		client.Timeout = 60 * time.Second
	}

	// reset the default node struct
	node.IsAvailable = false

	resp, err := client.Do(req)
	if err != nil {
		p.message = err.Error()
		if err := p.reportResult(node, time.Since(startTime).Seconds()); err != nil {
			return node, err
		}
		return node, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		p.message = fmt.Sprintf("status code: %d", resp.StatusCode)
		if err := p.reportResult(node, time.Since(startTime).Seconds()); err != nil {
			return node, err
		}
		return node, errors.New(p.message)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.message = err.Error()
		if err := p.reportResult(node, time.Since(startTime).Seconds()); err != nil {
			return node, err
		}
		return node, err
	}

	reportNode := struct {
		monero.Node `json:"result"`
	}{}

	if err := json.Unmarshal(body, &reportNode); err != nil {
		p.message = err.Error()
		if err := p.reportResult(node, time.Since(startTime).Seconds()); err != nil {
			return node, err
		}
	}
	if reportNode.Status == "OK" {
		node.IsAvailable = true
	}
	node.Nettype = reportNode.Nettype
	node.AdjustedTime = reportNode.AdjustedTime
	node.DatabaseSize = reportNode.DatabaseSize
	node.Difficulty = reportNode.Difficulty
	node.Height = reportNode.Height
	node.Version = reportNode.Version

	if resp.Header.Get("Access-Control-Allow-Origin") == "*" || resp.Header.Get("Access-Control-Allow-Origin") == "https://xmr.ditatompel.com" {
		node.CORSCapable = true
	}

	if !node.IsTor {
		hostIp, err := net.LookupIP(node.Hostname)
		if err != nil {
			fmt.Println("Warning: Could not resolve hostname: " + node.Hostname)
		} else {
			node.IP = hostIp[0].String()
		}
	}

	// Sleeping 1 second to avoid too many request on host behind CloudFlare
	// time.Sleep(1 * time.Second)

	// check fee
	fee, err := p.fetchFee(client, endpoint)
	if err != nil {
		return node, err
	}
	node.EstimateFee = fee

	tookTime := time.Since(startTime).Seconds()

	slog.Info(fmt.Sprintf("[PROBE] Took %f seconds", tookTime))
	if err := p.reportResult(node, tookTime); err != nil {
		return node, err
	}
	return node, nil
}

// get estimate fee from remote node
func (p *proberClient) fetchFee(client http.Client, endpoint string) (uint, error) {
	rpcParam := []byte(`{"jsonrpc": "2.0","id": "0","method": "get_fee_estimate"}`)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(rpcParam))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", RPCUserAgent)

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, fmt.Errorf("status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	f := struct {
		Result struct {
			Fee uint `json:"fee"`
		} `json:"result"`
	}{}

	if err := json.Unmarshal(body, &f); err != nil {
		return 0, err
	}

	return f.Result.Fee, nil
}

func (p *proberClient) reportResult(node monero.Node, tookTime float64) error {
	jsonData, err := json.Marshal(monero.ProbeReport{
		TookTime: tookTime,
		Message:  p.message,
		Node:     node,
	})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/api/v1/job", p.endpoint)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add(monero.ProberAPIKey, p.apiKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", RPCUserAgent)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return nil
}

// for debug purposes
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
