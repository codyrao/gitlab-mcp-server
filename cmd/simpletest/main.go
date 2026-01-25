package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Message struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	log.Println("GitLab MCP Server Test Client")
	log.Println("================================")

	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	log.Println("\n=== Test 1: Initialize ===")
	initMsg := Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendMessage(initMsg, stdout); err != nil {
		log.Fatalf("Failed to send initialize: %v", err)
	}

	initResp, err := readMessage(stdin)
	if err != nil {
		log.Fatalf("Failed to read initialize response: %v", err)
	}

	log.Printf("Initialize response: %+v", initResp)
	if initResp.Error != nil {
		log.Fatalf("Initialize failed: %s", initResp.Error.Message)
	}

	log.Println("\n=== Test 2: List Tools ===")
	toolsMsg := Message{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	if err := sendMessage(toolsMsg, stdout); err != nil {
		log.Fatalf("Failed to send tools/list: %v", err)
	}

	toolsResp, err := readMessage(stdin)
	if err != nil {
		log.Fatalf("Failed to read tools/list response: %v", err)
	}

	log.Printf("Tools list response: %+v", toolsResp)
	if toolsResp.Error != nil {
		log.Fatalf("Tools list failed: %s", toolsResp.Error.Message)
	}

	if result, ok := toolsResp.Result.(map[string]interface{}); ok {
		if tools, ok := result["tools"].([]interface{}); ok {
			log.Printf("Found %d tools:", len(tools))
			for _, tool := range tools {
				if toolMap, ok := tool.(map[string]interface{}); ok {
					name := ""
					if n, ok := toolMap["name"].(string); ok {
						name = n
					}
					desc := ""
					if d, ok := toolMap["description"].(string); ok {
						desc = d
					}
					log.Printf("  - %s: %s", name, desc)
				}
			}
		}
	}

	log.Println("\n=== All Tests Completed ===")
}

func sendMessage(msg Message, stdout *bufio.Writer) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Fprintln(stdout, string(data))
	return nil
}

func readMessage(stdin *bufio.Reader) (*Message, error) {
	line, err := stdin.ReadString('\n')
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
