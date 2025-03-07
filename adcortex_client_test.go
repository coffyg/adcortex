package adcortex_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coffyg/adcortex"
)

func TestAdCortexClientFetchAd(t *testing.T) {
	// Mock server returning one ad
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ads":[{"idx":1,"ad_title":"Test Ad","ad_description":"Description","placement_template":"Template","link":"http://example.com"}]}`))
	}))
	defer server.Close()

	userInfo, _ := adcortex.NewAdCortexUserInfo("user123", 30, adcortex.AdCortexGenderMale, "US", "en", nil)
	sessionInfo, _ := adcortex.NewAdCortexSessionInfo("session123", "ChatBot", nil, userInfo, adcortex.AdCortexPlatform{Name: "ios", Version: "16.0"})

	client, err := adcortex.NewAdCortexClient(
		sessionInfo,
		"Buy now: {ad_title} - {link}",
		"TEST_API_KEY",
		server.URL,
		nil,
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	ad, err := client.AdCortexFetchAd([]adcortex.AdCortexMessage{
		{Role: adcortex.AdCortexRoleUser, Content: "Hello!"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, ad)
	assert.Equal(t, 1, ad.Idx)
	assert.Equal(t, "Test Ad", ad.AdTitle)

	contextString := client.AdCortexGenerateContext(ad)
	assert.Contains(t, contextString, "Test Ad")
	assert.Contains(t, contextString, "http://example.com")
}
