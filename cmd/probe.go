package cmd

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
	"xmr-remote-nodes/internal/repo"

	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
)

const RPCUserAgent = "ditatombot/0.0.1 (Monero RPC Monitoring; https://github.com/ditatompel/xmr-remote-nodes)"

type proberClient struct {
	config  *config.App
	message string
}

func newProber(cfg *config.App) *proberClient {
	return &proberClient{config: cfg}
}

var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Run Monero node prober",
	Run: func(_ *cobra.Command, _ []string) {
		runProbe()
	},
}

func init() {
	rootCmd.AddCommand(probeCmd)
}

func runProbe() {
	cfg := config.AppCfg()
	if cfg.ServerEndpoint == "" {
		fmt.Println("Please set SERVER_ENDPOINT in .env")
		os.Exit(1)
	}

	if cfg.AcceptTor && cfg.TorSocks == "" {
		fmt.Println("Please set TOR_SOCKS in .env")
		os.Exit(1)
	}

	probe := newProber(cfg)

	node, err := probe.getJob()
	if err != nil {
		slog.Error(fmt.Sprintf("[PROBE] getJob: %s", err.Error()))
		os.Exit(1)
	}

	fetchNode, err := probe.fetchNode(node)
	if err != nil {
		slog.Error(fmt.Sprintf("[PROBE] fetchNode: %s", err.Error()))
		os.Exit(1)
	}
	slog.Debug(fmt.Sprintf("[PROBE] fetchNode: %s", prettyPrint(fetchNode)))
}

func (p *proberClient) getJob() (repo.MoneroNode, error) {
	queryParams := ""
	if p.config.AcceptTor {
		queryParams = "?accept_tor=1"
	}

	node := repo.MoneroNode{}

	endpoint := fmt.Sprintf("%s/api/v1/job%s", p.config.ServerEndpoint, queryParams)
	slog.Info(fmt.Sprintf("[PROBE] Getting node from %s", endpoint))

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return node, err
	}
	req.Header.Add("X-Prober-Api-Key", p.config.ApiKey)
	req.Header.Set("User-Agent", RPCUserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return node, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return node, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	response := struct {
		Data repo.MoneroNode `json:"data"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return node, err
	}

	node = response.Data
	slog.Info(fmt.Sprintf("[PROBE] Got node: %s://%s:%d", node.Protocol, node.Hostname, node.Port))

	return node, nil
}

func (p *proberClient) fetchNode(node repo.MoneroNode) (repo.MoneroNode, error) {
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
	if p.config.AcceptTor && node.IsTor {
		dialer, err := proxy.SOCKS5("tcp", p.config.TorSocks, nil, proxy.Direct)
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
		p.reportResult(node, time.Since(startTime).Seconds())
		return node, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		p.message = fmt.Sprintf("status code: %d", resp.StatusCode)
		p.reportResult(node, time.Since(startTime).Seconds())
		return node, errors.New(p.message)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.message = err.Error()
		p.reportResult(node, time.Since(startTime).Seconds())
		return node, err
	}

	reportNode := struct {
		repo.MoneroNode `json:"result"`
	}{}

	if err := json.Unmarshal(body, &reportNode); err != nil {
		p.message = err.Error()
		p.reportResult(node, time.Since(startTime).Seconds())
		return node, err
	}
	if reportNode.Status == "OK" {
		node.IsAvailable = true
	}
	node.NetType = reportNode.NetType
	node.AdjustedTime = reportNode.AdjustedTime
	node.DatabaseSize = reportNode.DatabaseSize
	node.Difficulty = reportNode.Difficulty
	node.Height = reportNode.Height
	node.Version = reportNode.Version

	if resp.Header.Get("Access-Control-Allow-Origin") == "*" || resp.Header.Get("Access-Control-Allow-Origin") == "https://xmr.ditatompel.com" {
		node.CorsCapable = true
	}

	if !node.IsTor {
		hostIp, err := net.LookupIP(node.Hostname)
		if err != nil {
			fmt.Println("Warning: Could not resolve hostname: " + node.Hostname)
		} else {
			node.Ip = hostIp[0].String()
		}
	}

	// Sleeping 1 second to avoid too many request on host behind CloudFlare
	// time.Sleep(1 * time.Second)

	// check fee
	rpcCheckFeeParam := []byte(`{"jsonrpc": "2.0","id": "0","method": "get_fee_estimate"}`)
	reqCheckFee, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(rpcCheckFeeParam))
	if err != nil {
		return node, err
	}
	reqCheckFee.Header.Set("Content-Type", "application/json; charset=UTF-8")
	reqCheckFee.Header.Set("User-Agent", RPCUserAgent)

	checkFee, err := client.Do(reqCheckFee)
	if err != nil {
		return node, err
	}
	defer checkFee.Body.Close()

	if checkFee.StatusCode != 200 {
		return node, fmt.Errorf("status code: %d", checkFee.StatusCode)
	}

	bodyCheckFee, err := io.ReadAll(checkFee.Body)
	if err != nil {
		return node, err
	}

	feeEstimate := struct {
		Result struct {
			Fee uint `json:"fee"`
		} `json:"result"`
	}{}

	if err := json.Unmarshal(bodyCheckFee, &feeEstimate); err != nil {
		return node, err
	}

	tookTime := time.Since(startTime).Seconds()
	node.EstimateFee = feeEstimate.Result.Fee

	slog.Info(fmt.Sprintf("[PROBE] Took %f seconds", tookTime))
	if err := p.reportResult(node, tookTime); err != nil {
		return node, err
	}
	return node, nil
}

func (p *proberClient) reportResult(node repo.MoneroNode, tookTime float64) error {
	jsonData, err := json.Marshal(repo.ProbeReport{
		TookTime: tookTime,
		Message:  p.message,
		NodeInfo: node,
	})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/api/v1/job", p.config.ServerEndpoint)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add("X-Prober-Api-Key", p.config.ApiKey)
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
