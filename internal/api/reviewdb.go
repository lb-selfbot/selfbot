package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/http"
)

var reviewDBToken string

func GetReviewDBToken(session *discordgo.Session) (string, error) {
	if reviewDBToken != "" {
		return reviewDBToken, nil
	}

	authorizeUrl := "https://discord.com/api/v9/oauth2/authorize?client_id=915703782174752809&response_type=code&redirect_uri=https%3A%2F%2Fmanti.vendicated.dev%2Fapi%2Freviewdb%2Fauth&scope=identify"

	rawData, err := session.Request("POST", authorizeUrl, map[string]any{
		"authorize":    true,
		"permissions":  "0",
		"user_install": false,
	})
	if err != nil {
		return "", err
	}

	var data map[string]any
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return "", err
	}

	location, ok := data["location"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response: no location: %s", string(rawData))
	}

	resp, err := http.DefaultClient().Get(location + "&clientMod=vencord")
	if err != nil {
		return "", err
	}

	token := resp.Json().Get("token").MustString()
	if token == "" {
		return "", fmt.Errorf("invalid response: no token")
	}

	return token, nil
}

type Review struct {
	ID        int64
	Username  string
	Comment   string
	Timestamp time.Time
}

func GetReviewDBReviews(userID string) ([]Review, error) {
	resp, err := http.DefaultClient().Get("https://	api/reviewdb/users/" + userID + "/reviews?flags=0&offset=0")
	if err != nil {
		return nil, err
	}

	data := resp.Json()

	if !data.Get("success").MustBool() {
		return nil, fmt.Errorf("invalid response: %s", data.Get("message").MustString())
	}

	reviewsData := data.Get("reviews")
	reviews := make([]Review, 0, reviewsData.Len())

	for review := range reviewsData.Iterator() {
		reviews = append(reviews, Review{
			ID:        review.Get("id").MustInt64(),
			Username:  review.Get("sender", "username").MustString(),
			Comment:   review.Get("comment").MustString(),
			Timestamp: time.Unix(review.Get("timestamp").MustInt64(), 0),
		})
	}

	return reviews, nil
}

func AddReviewDBReview(session *discordgo.Session, userID, review string) error {
	token, err := GetReviewDBToken(session)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient().Put("https://manti.vendicated.dev/api/reviewdb/users/"+userID+"/reviews", map[string]any{
		"token":   strings.TrimSpace(token),
		"comment": review,
		"userid":  userID,
	})

	return err
}
