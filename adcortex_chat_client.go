package adcortex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AdCortexChatClient is a stateful "chat" client that accumulates messages and
// fetches ads when certain thresholds are reached.
type AdCortexChatClient struct {
	sessionInfo           *AdCortexSessionInfo
	contextTemplate       string
	apiKey                string
	httpClient            *http.Client
	numMessagesBeforeAd   int
	numMessagesBetweenAds int
	baseURL               string

	messages []AdCortexMessage
	latestAd *AdCortexAd
	shownAds []adCortexShownAd // Records ads and at which message index they were shown
}

// adCortexShownAd tracks when an ad was shown in the conversation.
type adCortexShownAd struct {
	ad           AdCortexAd
	messageIndex int
}

// NewAdCortexChatClient constructs an AdCortexChatClient.
func NewAdCortexChatClient(
	sessionInfo *AdCortexSessionInfo,
	contextTemplate string,
	apiKey string,
	baseURL string,
	httpClient *http.Client,
	numMessagesBeforeAd int,
	numMessagesBetweenAds int,
) (*AdCortexChatClient, error) {
	if sessionInfo == nil {
		return nil, fmt.Errorf("sessionInfo cannot be nil")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey must be provided")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if baseURL == "" {
		baseURL = DefaultAdCortexBaseURL
	}

	return &AdCortexChatClient{
		sessionInfo:           sessionInfo,
		contextTemplate:       contextTemplate,
		apiKey:                apiKey,
		httpClient:            httpClient,
		numMessagesBeforeAd:   numMessagesBeforeAd,
		numMessagesBetweenAds: numMessagesBetweenAds,
		baseURL:               baseURL,
		messages:              make([]AdCortexMessage, 0),
		shownAds:              make([]adCortexShownAd, 0),
	}, nil
}

// AdCortexAddMessage adds a user or AI message. If conditions are met,
// it automatically fetches an ad. Returns the ad if fetched, or nil otherwise.
func (c *AdCortexChatClient) AdCortexAddMessage(role AdCortexRole, content string) (*AdCortexAd, error) {
	c.messages = append(c.messages, AdCortexMessage{
		Role:    role,
		Content: content,
	})

	if c.adCortexShouldShowAd() {
		return c.adCortexFetchAd()
	}

	return nil, nil
}

// AdCortexCreateContext returns a user-facing string containing the details of
// the most recently fetched ad (if any).
func (c *AdCortexChatClient) AdCortexCreateContext() string {
	if c.latestAd == nil {
		return ""
	}
	result := c.contextTemplate
	result = adCortexStrReplace(result, "{ad_title}", c.latestAd.AdTitle)
	result = adCortexStrReplace(result, "{ad_description}", c.latestAd.AdDescription)
	result = adCortexStrReplace(result, "{placement_template}", c.latestAd.PlacementTemplate)
	result = adCortexStrReplace(result, "{link}", c.latestAd.Link)
	return result
}

// adCortexShouldShowAd determines whether an ad should be shown based on the configured thresholds.
func (c *AdCortexChatClient) adCortexShouldShowAd() bool {
	totalMsgs := len(c.messages)
	// If no ad has been shown yet
	if len(c.shownAds) == 0 {
		// Only show if message count meets or exceeds numMessagesBeforeAd
		return totalMsgs >= c.numMessagesBeforeAd
	}

	// If we have shown ads before, we see how many messages have been added since the last ad
	lastShown := c.shownAds[len(c.shownAds)-1]
	msgsSinceLastAd := totalMsgs - lastShown.messageIndex
	return msgsSinceLastAd >= c.numMessagesBetweenAds
}

// adCortexFetchAd is a private function that fetches an ad from the server
// using only the latest numMessagesBeforeAd messages.
func (c *AdCortexChatClient) adCortexFetchAd() (*AdCortexAd, error) {
	// If not enough messages, skip
	if len(c.messages) < c.numMessagesBeforeAd {
		return nil, nil
	}

	startIndex := len(c.messages) - c.numMessagesBeforeAd
	if startIndex < 0 {
		startIndex = 0
	}

	payload := adCortexRequestPayload{
		RGUID:       c.sessionInfo.SessionID,
		SessionInfo: c.sessionInfo,
		UserData:    c.sessionInfo.UserInfo,
		Messages:    c.messages[startIndex:],
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
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK status %d: %s", resp.StatusCode, string(b))
	}

	var res adCortexResponsePayload
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		stringBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("failed to decode response: %w: '%s'", err, stringBody)
	}

	if len(res.Ads) == 0 {
		return nil, nil
	}

	c.latestAd = &res.Ads[0]
	c.shownAds = append(c.shownAds, adCortexShownAd{
		ad:           *c.latestAd,
		messageIndex: len(c.messages),
	})

	return c.latestAd, nil
}
