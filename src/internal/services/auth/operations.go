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
		Op:   OpCodeConfigUpdate,
		Data: WSData{Config: &cfg},
	})
}

func (a *AuthServiceImpl) AppendConsole(html string) error {
	return a.write(&WSMessage{
		Op:   OpCodeConsoleApend,
		Data: WSData{Text: html},
	})
}

func (a *AuthServiceImpl) ClearConsole() error {
	return a.write(&WSMessage{Op: OpCodeConsoleClear})
}

func (a *AuthServiceImpl) SetConsole(html string) error {
	return a.write(&WSMessage{
		Op:   OpCodeConsoleSet,
		Data: WSData{Text: html},
	})
}
