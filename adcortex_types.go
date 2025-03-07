package adcortex

import (
	"errors"
	"fmt"
)

// AdCortexGender is a string-based enum representing a user's gender.
type AdCortexGender string

const (
	AdCortexGenderMale   AdCortexGender = "male"
	AdCortexGenderFemale AdCortexGender = "female"
	AdCortexGenderOther  AdCortexGender = "other"
)

// AdCortexRole is a string-based enum representing the role of the message sender.
type AdCortexRole string

const (
	AdCortexRoleUser AdCortexRole = "user"
	AdCortexRoleAI   AdCortexRole = "ai"
)

// AdCortexInterest is a string-based enum representing user interests.
type AdCortexInterest string

const (
	AdCortexInterestFlirting   AdCortexInterest = "flirting"
	AdCortexInterestGaming     AdCortexInterest = "gaming"
	AdCortexInterestSports     AdCortexInterest = "sports"
	AdCortexInterestMusic      AdCortexInterest = "music"
	AdCortexInterestTravel     AdCortexInterest = "travel"
	AdCortexInterestTechnology AdCortexInterest = "technology"
	AdCortexInterestArt        AdCortexInterest = "art"
	AdCortexInterestCooking    AdCortexInterest = "cooking"
	AdCortexInterestAll        AdCortexInterest = "all"
)

// AdCortexUserInfo stores user information for the ADCortex API.
type AdCortexUserInfo struct {
	UserID    string             `json:"user_id"`
	Age       int                `json:"age"`
	Gender    AdCortexGender     `json:"gender"`
	Location  string             `json:"location"` // ISO 3166-1 alpha-2 code
	Language  string             `json:"language"` // Must be "en" (in this example)
	Interests []AdCortexInterest `json:"interests"`
}

// NewAdCortexUserInfo is a constructor that validates user info before returning.
func NewAdCortexUserInfo(
	userID string,
	age int,
	gender AdCortexGender,
	location string,
	language string,
	interests []AdCortexInterest,
) (*AdCortexUserInfo, error) {
	if err := adCortexValidateLocation(location); err != nil {
		return nil, err
	}
	if err := adCortexValidateLanguage(language); err != nil {
		return nil, err
	}
	return &AdCortexUserInfo{
		UserID:    userID,
		Age:       age,
		Gender:    gender,
		Location:  location,
		Language:  language,
		Interests: interests,
	}, nil
}

// adCortexValidateLocation checks if the location is a valid ISO code (minimal set here).
func adCortexValidateLocation(location string) error {
	validCodes := map[string]bool{
		"US": true,
		"GB": true,
		"CA": true,
		"FR": true,
		"DE": true,
		"IN": true,
		// Expand this list as needed...
	}
	if !validCodes[location] {
		return fmt.Errorf("%s is not a valid or supported country code in this demo", location)
	}
	return nil
}

// adCortexValidateLanguage ensures the language is "en" only in this example.
func adCortexValidateLanguage(lang string) error {
	if lang != "en" {
		return errors.New("language must be 'en'")
	}
	return nil
}

// AdCortexPlatform contains platform-related metadata.
type AdCortexPlatform struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// AdCortexSessionInfo stores session details including user and platform information.
type AdCortexSessionInfo struct {
	SessionID         string                 `json:"session_id"`
	CharacterName     string                 `json:"character_name"`
	CharacterMetadata map[string]interface{} `json:"character_metadata"`
	UserInfo          *AdCortexUserInfo      `json:"user_info"`
	Platform          AdCortexPlatform       `json:"platform"`
}

// NewAdCortexSessionInfo constructs a valid session info object.
func NewAdCortexSessionInfo(
	sessionID string,
	characterName string,
	characterMetadata map[string]interface{},
	userInfo *AdCortexUserInfo,
	platform AdCortexPlatform,
) (*AdCortexSessionInfo, error) {
	if sessionID == "" {
		return nil, errors.New("sessionID cannot be empty")
	}
	if userInfo == nil {
		return nil, errors.New("userInfo cannot be nil")
	}
	return &AdCortexSessionInfo{
		SessionID:         sessionID,
		CharacterName:     characterName,
		CharacterMetadata: characterMetadata,
		UserInfo:          userInfo,
		Platform:          platform,
	}, nil
}

// AdCortexMessage represents a single message in a conversation.
type AdCortexMessage struct {
	Role    AdCortexRole `json:"role"`
	Content string       `json:"content"`
}

// AdCortexAd represents an advertisement fetched via the ADCortex API.
type AdCortexAd struct {
	Idx               int    `json:"idx"`
	AdTitle           string `json:"ad_title"`
	AdDescription     string `json:"ad_description"`
	PlacementTemplate string `json:"placement_template"`
	Link              string `json:"link"`
}
