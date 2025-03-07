package adcortex_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coffyg/adcortex"
)

func TestAdCortexChatClientBasic(t *testing.T) {
	// Mock server returning one ad
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ads":[{"idx":123,"ad_title":"Mock Chat Ad","ad_description":"Description Chat","placement_template":"Template Chat","link":"http://chat.example.com"}]}`))
	}))
	defer server.Close()

	userInfo, _ := adcortex.NewAdCortexUserInfo("user123", 20, adcortex.AdCortexGenderFemale, "US", "en", nil)
	sessionInfo, _ := adcortex.NewAdCortexSessionInfo(
		"chat-session-456",
		"CompanionBot",
		map[string]interface{}{"description": "Testing chat bot"},
		userInfo,
		adcortex.AdCortexPlatform{Name: "android", Version: "13"},
	)

	chatClient, err := adcortex.NewAdCortexChatClient(
		sessionInfo,
		"Check out this product: {ad_title} - {link}",
		"TEST_API_KEY",
		server.URL,
		nil,
		3, // number of messages before first ad
		2, // number of messages between ads
	)
	assert.NoError(t, err)

	// Add some messages
	ad, err := chatClient.AdCortexAddMessage(adcortex.AdCortexRoleUser, "Hello")
	assert.NoError(t, err)
	assert.Nil(t, ad, "No ad yet because not enough messages")

	ad, err = chatClient.AdCortexAddMessage(adcortex.AdCortexRoleAI, "Hi there!")
	assert.NoError(t, err)
	assert.Nil(t, ad, "Still no ad; only 2 messages so far")

	ad, err = chatClient.AdCortexAddMessage(adcortex.AdCortexRoleUser, "How are you?")
	assert.NoError(t, err)
	// Now we have 3 total messages; that meets numMessagesBeforeAd
	assert.NotNil(t, ad)
	assert.Equal(t, 123, ad.Idx)
	assert.Equal(t, "Mock Chat Ad", ad.AdTitle)
	assert.Equal(t, 1, callCount, "Should have called server once")

	// Check the context
	contextStr := chatClient.AdCortexCreateContext()
	assert.Contains(t, contextStr, "Mock Chat Ad")

	// Next messages...
	ad, err = chatClient.AdCortexAddMessage(adcortex.AdCortexRoleAI, "I'm just a bot, but good!")
	assert.NoError(t, err)
	assert.Nil(t, ad, "We only have 1 new message since last ad")

	ad, err = chatClient.AdCortexAddMessage(adcortex.AdCortexRoleUser, "Great!")
	assert.NoError(t, err)
	// Now 2 new messages since last ad => Should show ad again
	assert.NotNil(t, ad)
	assert.Equal(t, 123, ad.Idx) // same mocked ad from server
	assert.Equal(t, 2, callCount, "Should have called server twice by now")
}
