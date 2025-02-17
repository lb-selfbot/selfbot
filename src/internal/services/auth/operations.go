package auth

import "github.com/LightningDev1/LB-Selfbot/selfbot/config"

func (a *AuthServiceImpl) SendConfigUpdate() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Remove sensitive data
	cfg.Account.LightningBot.Password = ""
	cfg.Account.Discord.Token = ""

	return a.write(&WSMessage{
		Op:   OP_CONFIG_UPDATE,
		Data: WSData{Config: &cfg},
	})
}

func (a *AuthServiceImpl) AppendConsole(html string) error {
	return a.write(&WSMessage{
		Op:   OP_CONSOLE_APPEND,
		Data: WSData{Text: html},
	})
}

func (a *AuthServiceImpl) ClearConsole() error {
	return a.write(&WSMessage{Op: OP_CONSOLE_CLEAR})
}

func (a *AuthServiceImpl) SetConsole(html string) error {
	return a.write(&WSMessage{
		Op:   OP_CONSOLE_SET,
		Data: WSData{Text: html},
	})
}
