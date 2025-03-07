# adcortex
golang client for AdCortex (by 3102labs)
---

## Features

- **AdCortexClient**: Low-level client for fetching ads based on user session and conversation messages.
- **AdCortexChatClient**: Chat-style client that automatically fetches ads after a configurable number of messages.
- **Typed API**: Strongly-typed structs for user info, session info, messages, ads, etc.
- **Easy Integration**: Bring your own `http.Client` or use the default.  
- **Customizable**: Provide your own base URL, context templates, thresholds for ad fetching, etc.

---

## Installation

```bash
go get github.com/coffyg/adcortex
```

--- 

## Quick Start
Below is a minimal example of how to use the AdCortexClient and AdCortexChatClient.

### 1. Create a User and Session
```golang
package main

import (
    "fmt"
    "log"

    "github.com/coffyg/adcortex"
)

func main() {
    // 1) Create user info
    userInfo, err := adcortex.NewAdCortexUserInfo(
        "user123",
        25,
        adcortex.AdCortexGenderMale,
        "US",
        "en",
        []adcortex.AdCortexInterest{
            adcortex.AdCortexInterestGaming,
            adcortex.AdCortexInterestSports,
        },
    )
    if err != nil {
        log.Fatal("Error creating user info:", err)
    }

    // 2) Create session info
    sessionInfo, err := adcortex.NewAdCortexSessionInfo(
        "session-abc",
        "MyAssistant",
        map[string]interface{}{"description": "An AI companion"},
        userInfo,
        adcortex.AdCortexPlatform{Name: "linux", Version: "1.0"},
    )
    if err != nil {
        log.Fatal("Error creating session info:", err)
    }

    // 3) Create a client
    client, err := adcortex.NewAdCortexClient(
        sessionInfo,
        "Special Offer: {ad_title} - {ad_description}, link: {link}",
        "MY_SECRET_KEY", // Your ADCortex API key
        "",             // Leave empty to use default base URL
        nil,            // Optional custom http.Client
    )
    if err != nil {
        log.Fatal("Error creating ADCortexClient:", err)
    }

    // 4) Fetch an ad with a set of messages
    ad, err := client.AdCortexFetchAd([]adcortex.AdCortexMessage{
        {Role: adcortex.AdCortexRoleUser, Content: "Hello AI!"},
        {Role: adcortex.AdCortexRoleAI,   Content: "Hello user, how can I help you today?"},
    })
    if err != nil {
        log.Fatal("Error fetching ad:", err)
    }

    // 5) Generate a context string (e.g., to add to conversation)
    if ad != nil {
        contextStr := client.AdCortexGenerateContext(ad)
        fmt.Println("Ad Context:", contextStr)
    } else {
        fmt.Println("No ad returned.")
    }
}
```
### 1. Create a User and Session
```golang
package main

import (
    "fmt"
    "log"

    "github.com/coffyg/adcortex"
)

func main() {
    // 1) Create user info
    userInfo, err := adcortex.NewAdCortexUserInfo(
        "user123",
        25,
        adcortex.AdCortexGenderMale,
        "US",
        "en",
        []adcortex.AdCortexInterest{
            adcortex.AdCortexInterestGaming,
            adcortex.AdCortexInterestSports,
        },
    )
    if err != nil {
        log.Fatal("Error creating user info:", err)
    }

    // 2) Create session info
    sessionInfo, err := adcortex.NewAdCortexSessionInfo(
        "session-abc",
        "MyAssistant",
        map[string]interface{}{"description": "An AI companion"},
        userInfo,
        adcortex.AdCortexPlatform{Name: "linux", Version: "1.0"},
    )
    if err != nil {
        log.Fatal("Error creating session info:", err)
    }

    // 3) Create a client
    client, err := adcortex.NewAdCortexClient(
        sessionInfo,
        "Special Offer: {ad_title} - {ad_description}, link: {link}",
        "MY_SECRET_KEY", // Your ADCortex API key
        "",             // Leave empty to use default base URL
        nil,            // Optional custom http.Client
    )
    if err != nil {
        log.Fatal("Error creating ADCortexClient:", err)
    }

    // 4) Fetch an ad with a set of messages
    ad, err := client.AdCortexFetchAd([]adcortex.AdCortexMessage{
        {Role: adcortex.AdCortexRoleUser, Content: "Hello AI!"},
        {Role: adcortex.AdCortexRoleAI,   Content: "Hello user, how can I help you today?"},
    })
    if err != nil {
        log.Fatal("Error fetching ad:", err)
    }

    // 5) Generate a context string (e.g., to add to conversation)
    if ad != nil {
        contextStr := client.AdCortexGenerateContext(ad)
        fmt.Println("Ad Context:", contextStr)
    } else {
        fmt.Println("No ad returned.")
    }
}
```
### 2. Chat-Style Usage
If you want to accumulate messages and only fetch ads automatically after a certain number of messages, use AdCortexChatClient:

```golang
package main

import (
    "fmt"
    "log"

    "github.com/coffyg/adcortex"
)

func main() {
    // 1) Create user & session info (same as above) ...
    userInfo, _ := adcortex.NewAdCortexUserInfo("user123", 25, adcortex.AdCortexGenderMale, "US", "en", nil)
    sessionInfo, _ := adcortex.NewAdCortexSessionInfo("session-abc", "MyAssistant", nil, userInfo, adcortex.AdCortexPlatform{Name: "linux", Version: "1.0"})

    // 2) Create a Chat Client
    chatClient, err := adcortex.NewAdCortexChatClient(
        sessionInfo,
        "Amazing Product: {ad_title} - {link}",
        "MY_SECRET_KEY", // Your ADCortex API key
        "",             // Default base URL
        nil,            // Optional custom http.Client
        3,              // numMessagesBeforeAd (first ad)
        2,              // numMessagesBetweenAds (subsequent ads)
    )
    if err != nil {
        log.Fatal("Error creating AdCortexChatClient:", err)
    }

    // 3) Add messages
    if _, err := chatClient.AdCortexAddMessage(adcortex.AdCortexRoleUser, "Hello AI!"); err != nil {
        log.Fatal(err)
    }
    if _, err := chatClient.AdCortexAddMessage(adcortex.AdCortexRoleAI, "Hi, how can I help?"); err != nil {
        log.Fatal(err)
    }

    // No ad yet, as we only have 2 messages and need at least 3.

    // 4) Add a third message - triggers first ad
    ad, err := chatClient.AdCortexAddMessage(adcortex.AdCortexRoleUser, "Tell me something interesting.")
    if err != nil {
        log.Fatal(err)
    }

    // Check if we got an ad
    if ad != nil {
        fmt.Println("Got Ad:", ad.AdTitle)
        fmt.Println("Ad Context:", chatClient.AdCortexCreateContext())
    } else {
        fmt.Println("No ad returned.")
    }
}
```

## Tests
`go test -v`

