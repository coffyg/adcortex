package adcortex_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coffyg/adcortex"
)

func TestNewAdCortexUserInfo(t *testing.T) {
	userInfo, err := adcortex.NewAdCortexUserInfo(
		"user123",
		30,
		adcortex.AdCortexGenderMale,
		"US",
		"en",
		[]adcortex.AdCortexInterest{adcortex.AdCortexInterestGaming, adcortex.AdCortexInterestSports},
	)
	assert.NoError(t, err)
	assert.NotNil(t, userInfo)

	_, err = adcortex.NewAdCortexUserInfo("user123", 30, adcortex.AdCortexGenderMale, "XX", "en", nil)
	assert.Error(t, err, "invalid location should fail")
}

func TestNewAdCortexSessionInfo(t *testing.T) {
	userInfo, _ := adcortex.NewAdCortexUserInfo("user123", 30, adcortex.AdCortexGenderMale, "US", "en", nil)
	sessionInfo, err := adcortex.NewAdCortexSessionInfo(
		"session123",
		"MyCharacter",
		map[string]interface{}{"description": "A fun character"},
		userInfo,
		adcortex.AdCortexPlatform{Name: "iOS", Version: "16.0"},
	)
	assert.NoError(t, err)
	assert.NotNil(t, sessionInfo)
}
