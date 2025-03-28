package api

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/http"
)

type TranslateResult struct {
	Original            string
	Result              string
	SourceLanguage      string
	DestinationLanguage string
	Confidence          float64
	Error               error
}

func Translate(text, language string) *TranslateResult {
	language = strings.ToLower(language)
	language = strings.Split(language, "_")[0]

	if _, ok := constants.LANGUAGES[language]; !ok {
		if _, ok := constants.LANGUAGE_SPECIAL_CASES[language]; ok {
			language = constants.LANGUAGE_SPECIAL_CASES[language]
		} else if _, ok := constants.LANGUAGE_CODES[language]; ok {
			language = constants.LANGUAGE_CODES[language]
		} else {
			return &TranslateResult{
				Error: fmt.Errorf("invalid destination language"),
			}
		}
	}

	translateUrl := fmt.Sprintf(
		"https://translate.google.com/translate_a/single?client=gtx&sl=auto&tl=%s&dt=t&q=%s",
		url.QueryEscape(language),
		url.QueryEscape(text),
	)

	resp, err := http.DefaultClient().Get(translateUrl)
	if err != nil {
		return &TranslateResult{Error: err}
	}

	result := resp.Json().Get(0, 0)
	return &TranslateResult{
		Original:            result.Get(1).MustString(),
		Result:              result.Get(0).MustString(),
		SourceLanguage:      resp.Json().Get(2).MustString(),
		DestinationLanguage: language,
		Confidence:          resp.Json().Get(6).MustFloat64(),
	}
}

func TTS(text, language string) (*bytes.Reader, error) {
	url := fmt.Sprintf(
		"https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s",
		url.QueryEscape(text),
		url.QueryEscape(language),
	)

	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	return bytes.NewReader(resp.MustGetBody()), nil
}
