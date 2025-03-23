package setup

import (
	"os"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/tui"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

func GetLogin() (string, string) {
	if username, ok := os.LookupEnv("LB_USERNAME"); ok {
		if password, ok := os.LookupEnv("LB_PASSWORD"); ok {
			return username, password
		}
	}

	var username, password string

	if tui.IsSupported() {
		username, password = tui.GetLogin("Please enter your LightningBot credentials:")
	} else {
		username = utils.Input.GetString("Please enter your LightningBot username")
		password = utils.Input.GetString("Please enter your LightningBot password")
	}

	if username == "" || password == "" {
		utils.Logging.Error("You must enter a username and password!")
		time.Sleep(5 * time.Second)

		return GetLogin()
	}

	httpResponse := api.Auth(username, password)
	if httpResponse.Error != nil {
		utils.Logging.Error("Error connecting to API:", httpResponse.Error)
		os.Exit(1)
	}

	json := httpResponse.Json()

	errorString := json.Get("error").MustString()

	if errorString == "INVALID_CREDENTIALS" {
		utils.Logging.Error("Invalid credentials! Try again...")
		return GetLogin()
	} else if errorString != "" {
		utils.Logging.Error("Error logging in:", errorString, "status", json.Get("status").MustString())
		os.Exit(1)
	}

	return username, password
}

func GetToken() string {
	if token, ok := os.LookupEnv("LB_TOKEN"); ok {
		return token
	}

	var token string

	if tui.IsSupported() {
		token = tui.GetToken()
	} else {
		token = utils.Input.GetString("Please enter your Discord token")
	}

	// This fixes a bug where the bot tries to DoS Discord when there is no stdin
	if token == "" {
		utils.Logging.Error("You must enter a token!")
		time.Sleep(5 * time.Second)

		return GetToken()
	}

	session, err := discordgo.New(token)
	if err != nil {
		utils.Logging.Error("Error initializing session:", err)
		return GetToken()
	}
	session.SetIdentify(discordgo.IdentifyDiscordClient)

	_, err = session.User("@me")
	if err != nil {
		utils.Logging.Error("That token is invalid! Error:", err)
		return GetToken()
	}

	return token
}
