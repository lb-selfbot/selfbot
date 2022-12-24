package setup

import (
	"os"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

func GetLogin() (string, string) {
	username := utils.Input.GetString("Please enter your LightningBot username")
	password := utils.Input.GetString("Please enter your LightningBot password")

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

func GetToken(session *discordgo.Session) string {
	token := utils.Input.GetString("Please enter your Discord token")

	session.Token = token

	_, err := api.GetTokenData(session, token)
	if err != nil {
		utils.Logging.Error("That token is invalid! Error:", err)
		return GetToken(session)
	}

	return token
}
