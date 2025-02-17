package auth

import (
	"fmt"
	"os"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
)

func (a *AuthServiceImpl) onMessage(msg *WSMessage) error {
	switch msg.Op {
	case OP_HELLO_ACK:
		return a.onHelloAck(msg)
	case OP_FREE_TRIAL_START:
		return a.onFreeTrialStart(msg)
	case OP_EXIT:
		return a.onExit(msg)
	case OP_RESTART:
		return a.onRestart(msg)
	case OP_WEBSITE_CONNECT:
		return a.onWebsiteConnect(msg)
	case OP_CONFIG_UPDATE:
		return a.onConfigUpdate(msg)
	case OP_PONG:
		return nil
	case OP_ANOTHER_SELFBOT_CONNECTED:
		// Handled by error handler
		return nil
	default:
		return fmt.Errorf("unexpected message op: %s", msg.Op)
	}
}

func (a *AuthServiceImpl) onHelloAck(msg *WSMessage) error {
	a.receivedNonce = msg.Data.Nonce
	a.CheckNonce()

	a.readyChan <- struct{}{}

	return a.SendConfigUpdate()
}

func (a *AuthServiceImpl) onFreeTrialStart(_ *WSMessage) error {
	return a.uiSvc.ShowNotification("Your free trial has started! You will be able to use LightningBot for free for the next 7 days.")
}

func (a *AuthServiceImpl) onExit(_ *WSMessage) error {
	os.Exit(0)
	return nil
}

func (a *AuthServiceImpl) onRestart(_ *WSMessage) error {
	return a.ipcSvc.SendRestart()
}

func (a *AuthServiceImpl) onWebsiteConnect(_ *WSMessage) error {
	err := a.write(&WSMessage{
		Op:   OP_CONSOLE_SET,
		Data: WSData{Text: a.ipcSvc.GetHTMLConsole()},
	})
	if err != nil {
		return fmt.Errorf("failed to send console state: %w", err)
	}

	err = a.SendConfigUpdate()
	if err != nil {
		return fmt.Errorf("failed to send config update: %w", err)
	}

	return nil
}

func (a *AuthServiceImpl) onConfigUpdate(msg *WSMessage) error {
	if msg.Data.Config == nil {
		return nil
	}

	// Restore sensitive data
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	msg.Data.Config.Account = cfg.Account

	if err = msg.Data.Config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return a.ipcSvc.SendConfigUpdate(msg.Data.Config)
}
