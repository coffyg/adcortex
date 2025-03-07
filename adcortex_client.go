package adcortex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DefaultAdCortexBaseURL is the default endpoint for ads.
const DefaultAdCortexBaseURL = "https://adcortex.3102labs.com/ads/match"

// AdCortexClient is a client for fetching ads from the ADCortex API.
type AdCortexClient struct {
	sessionInfo     *AdCortexSessionInfo
	contextTemplate string
	apiKey          string
	baseURL         string
	httpClient      *http.Client
}

// NewAdCortexClient creates a new client for the ADCortex API. If baseURL is empty,
// it defaults to "https://adcortex.3102labs.com/ads/match".
func NewAdCortexClient(
	sessionInfo *AdCortexSessionInfo,
	contextTemplate string,
	apiKey string,
	baseURL string,
	httpClient *http.Client,
) (*AdCortexClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey must be provided")
	}
	if sessionInfo == nil {
		return nil, fmt.Errorf("sessionInfo must not be nil")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if baseURL == "" {
		baseURL = DefaultAdCortexBaseURL
	}

	client := &AdCortexClient{
		sessionInfo:     sessionInfo,
		contextTemplate: contextTemplate,
		apiKey:          apiKey,
		baseURL:         baseURL,
		httpClient:      httpClient,
	}
	return client, nil
}

// adCortexRequestPayload is the request payload for fetching ads.
type adCortexRequestPayload struct {
	RGUID       string               `json:"RGUID"`
	SessionInfo *AdCortexSessionInfo `json:"session_info"`
	UserData    *AdCortexUserInfo    `json:"user_data"`
	Messages    []AdCortexMessage    `json:"messages"`
	Platform    AdCortexPlatform     `json:"platform"`
}

// adCortexResponsePayload is the response payload from the ads API.
type adCortexResponsePayload struct {
	Ads []AdCortexAd `json:"ads"`
}

// AdCortexFetchAd sends a request to fetch a single advertisement. It returns
// an AdCortexAd and a possible error.
func (c *AdCortexClient) AdCortexFetchAd(messages []AdCortexMessage) (*AdCortexAd, error) {
	payload := adCortexRequestPayload{
		RGUID:       c.sessionInfo.SessionID,
		SessionInfo: c.sessionInfo,
		UserData:    c.sessionInfo.UserInfo,
		Messages:    messages,
		Platform:    c.sessionInfo.Platform,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseData adCortexResponsePayload
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		stringBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("failed to decode response: %w: '%s'", err, stringBody)
	}

	if len(responseData.Ads) == 0 {
		// No ads returned, return nil with no error or a custom error if desired
		return nil, nil
	}

	return &responseData.Ads[0], nil // Return the first ad
}

// AdCortexGenerateContext uses the contextTemplate to produce a personalized
// string containing ad details.
func (c *AdCortexClient) AdCortexGenerateContext(ad *AdCortexAd) string {
	if ad == nil {
		return ""
	}

	// Simple naive replacement. For a more robust solution, consider text/template or strings.Builder.
	result := c.contextTemplate
	result = adCortexStrReplace(result, "{ad_title}", ad.AdTitle)
	result = adCortexStrReplace(result, "{ad_description}", ad.AdDescription)
	result = adCortexStrReplace(result, "{placement_template}", ad.PlacementTemplate)
	result = adCortexStrReplace(result, "{link}", ad.Link)
	return result
}

// adCortexStrReplace replaces all occurrences of old with new in s.
func adCortexStrReplace(s, old, new string) string {
	// As of Go 1.20+, strings.ReplaceAll is standard, or we can implement it manually
	// return strings.ReplaceAll(s, old, new)

	// Implementation using standard library:
	return adCortexReplaceAll(s, old, new)
}

// adCortexReplaceAll is a private helper function. You could also just use strings.ReplaceAll.
func adCortexReplaceAll(s, old, new string) string {
	for {
		index := -1
		index = findIndex(s, old)
		if index == -1 {
			break
		}
		s = s[:index] + new + s[index+len(old):]
	}
	return s
}

// findIndex is a simplistic substring search for demonstration.
func findIndex(s, sub string) int {
	subLen := len(sub)
	for i := 0; i+subLen <= len(s); i++ {
		if s[i:i+subLen] == sub {
			return i
		}
	}
	return -1
}
