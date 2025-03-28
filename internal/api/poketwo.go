package api

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/lb-selfbot/selfbot/internal/http"
)

type Pokemon struct {
	ID            string `json:"id"`
	Region        string `json:"region"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	NameJapanese  string `json:"name.ja"`
	NameJapaneseT string `json:"name.ja_t"`
	NameEnglish   string `json:"name.en"`
	NameGerman    string `json:"name.de"`
	NameFrench    string `json:"name.fr"`
	Type          string `json:"type.0"`
	Event         string `json:"event"`
	Height        string `json:"height"`
	Weight        string `json:"weight"`
	HP            string `json:"base.hp"`
	Attack        string `json:"base.atk"`
	Defense       string `json:"base.def"`
}

func PokeTwo(image string) (*Pokemon, float64, error) {
	resp, err := http.DefaultClient().Post(
		"https://lightningdev1-poketwo.hf.space/run/predict",
		map[string]any{"data": []any{image}},
	)
	if err != nil {
		return nil, 0, err
	}

	var respData struct {
		Data []struct {
			Pokemon    *Pokemon `json:"pokemon"`
			Confidence float64  `json:"confidence"`
		}
	}

	if err := json.Unmarshal(resp.MustGetBody(), &respData); err != nil {
		return nil, 0, err
	}

	if len(respData.Data) == 0 {
		return nil, 0, fmt.Errorf("no data")
	}

	return respData.Data[0].Pokemon, respData.Data[0].Confidence, nil
}
