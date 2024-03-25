package llama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LlamaAPIHandler struct {
	URL string
}

func (lh *LlamaAPIHandler) newDefaultLlamaRequest(username string, botname string) *LlamaCppRequest {
	stopStrings := []string{
		"<s>",
		fmt.Sprintf("%s:", username),
		fmt.Sprintf("%s:", botname),
	}
	return &LlamaCppRequest{
		Stream:           false,
		RepeatPenalty:    1.18,
		TopK:             40,
		TopP:             0.95,
		MinP:             0.05,
		TfsZ:             1,
		TypicalP:         1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		Mirostat:         0,
		MirostatTau:      5,
		MirostatEta:      0.1,
		Grammar:          "",
		NProbs:           0,
		MinKeep:          0,
		CachePrompt:      false,
		SlotID:           0,
		Temperature:      0.7,
		Stop:             stopStrings,
	}
}

func NewLlamaAPIHandler(port string) *LlamaAPIHandler {
	return &LlamaAPIHandler{
		URL: fmt.Sprintf("http://localhost:%s/completion", port),
	}
}

func (lh *LlamaAPIHandler) MakeRequestAndDecode(ctx context.Context, prompt string, sysPrompt string, username string, botname string) (*LlamaApiResponse, error) {
	// Create a new request
	requestPayload := lh.newDefaultLlamaRequest(username, botname)
	requestPayload.Prompt = fmt.Sprintf("%s <s> %s: %s %s:", sysPrompt, username, prompt, botname)
	requestPayload.NPredict = len(prompt) + 10 // trying something...
	requestPayload.Stop = append(requestPayload.Stop, []string{username, botname}...)
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	// Make the POST request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, lh.URL, body) // bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 360, // TODO - configure a better timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Read and decode the response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse LlamaApiResponse
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
