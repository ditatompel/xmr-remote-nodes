package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ditatompel/xmr-nodes/internal/config"
	"github.com/ditatompel/xmr-nodes/internal/repo"

	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
)

const RPCUserAgent = "ditatombot/0.0.1 (Monero RPC Monitoring; Contact: ditatombot@ditatompel.com)"

type proberClient struct {
	config *config.App
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
	fmt.Printf("Accept Tor: %t\n", cfg.AcceptTor)

	if cfg.AcceptTor && cfg.TorSocks == "" {
		fmt.Println("Please set TOR_SOCKS in .env")
		os.Exit(1)
	}

	probe := newProber(cfg)

	node, err := probe.getJob()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fetchNode, err := probe.fetchNode(node)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(prettyPrint(fetchNode))
}

func (p *proberClient) getJob() (repo.MoneroNode, error) {
	queryParams := ""
	if p.config.ApiKey != "" {
		queryParams = "?api_key=" + p.config.ApiKey
	}

	node := repo.MoneroNode{}

	endpoint := fmt.Sprintf("%s/api/v1/job%s", p.config.ServerEndpoint, queryParams)

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

	return node, nil
}

func (p *proberClient) fetchNode(node repo.MoneroNode) (repo.MoneroNode, error) {
	startTime := time.Now()
	endpoint := fmt.Sprintf("%s://%s:%d/json_rpc", node.Protocol, node.Hostname, node.Port)
	rpcParam := []byte(`{"jsonrpc": "2.0","id": "0","method": "get_info"}`)

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
		// TODO: Post report to server
		return node, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO: Post report to server
		return node, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// TODO: Post report to server
		return node, err
	}

	reportNode := struct {
		repo.MoneroNode `json:"result"`
	}{}

	if err := json.Unmarshal(body, &reportNode); err != nil {
		// TODO: Post report to server
		return node, err
	}
	node.IsAvailable = true
	node.NetType = reportNode.NetType
	node.AdjustedTime = reportNode.AdjustedTime
	node.DatabaseSize = reportNode.DatabaseSize
	node.Difficulty = reportNode.Difficulty
	node.NodeVersion = reportNode.NodeVersion

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

	fmt.Printf("Took %f seconds\n", tookTime)
	return node, nil
}

// for debug purposes
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}