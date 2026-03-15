package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   interface{}     `json:"error,omitempty"`
}

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type CallToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func sanitizeName(name string) string {
	if name == "inform_about_unavailable_entities_and_sent_restart_action" {
		return "inform_about_unavail_entities_restart"
	}

	reg := regexp.MustCompile(`[^a-zA-Z0-9_.:]`)
	sanitized := reg.ReplaceAllString(name, "_")
	if len(sanitized) > 64 {
		sanitized = sanitized[:64]
	}
	return sanitized
}

func main() {
	hassServer := os.Getenv("HASS_SERVER")
	if hassServer == "" {
		hassServer = "http://homeassistant.local:8123"
	}
	url := hassServer + "/api/mcp"
	token := os.Getenv("HASS_TOKEN")

	if token == "" {
		log.Fatal("HASS_TOKEN environment variable is required")
	}

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 10*1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	client := &http.Client{}

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var rpcReq JSONRPCRequest
		if err := json.Unmarshal(line, &rpcReq); err == nil {
			if rpcReq.Method == "tools/call" {
				var params CallToolParams
				if err := json.Unmarshal(rpcReq.Params, &params); err == nil {
					if params.Name == "inform_about_unavail_entities_restart" {
						params.Name = "inform_about_unavailable_entities_and_sent_restart_action"
						newParams, _ := json.Marshal(params)
						rpcReq.Params = newParams
						line, _ = json.Marshal(rpcReq)
					}
				}
			}
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(line))
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to send request: %v", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var rpcResp JSONRPCResponse
			if err := json.Unmarshal(respBody, &rpcResp); err == nil {
				if rpcReq.Method == "tools/list" && rpcResp.Result != nil {
					var listResult ListToolsResult
					if err := json.Unmarshal(rpcResp.Result, &listResult); err == nil {
						for i := range listResult.Tools {
							listResult.Tools[i].Name = sanitizeName(listResult.Tools[i].Name)
						}
						newResult, _ := json.Marshal(listResult)
						rpcResp.Result = newResult
						respBody, _ = json.Marshal(rpcResp)
					}
				}
			}
			fmt.Printf("%s\n", string(respBody))
		} else {
			log.Printf("HA API error: %d %s", resp.StatusCode, string(respBody))
		}
	}
}
