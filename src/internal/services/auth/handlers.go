package auth

import (
	"fmt"
	"os"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
)

func (a *AuthServiceImpl) onMessage(msg *WSMessage) error {
	switch msg.Op {
	case OpCodeHelloAck:
		return a.onHelloAck(msg)
	case OpCodeFreeTrialStart:
		return a.onFreeTrialStart(msg)
	case OpCodeExit:
		return a.onExit(msg)
	case OpCodeRestart:
		return a.onRestart(msg)
	case OpCodeWebsiteConnect:
		return a.onWebsiteConnect(msg)
	case OpCodeConfigUpdate:
		return a.onConfigUpdate(msg)
	case OpCodePong:
		return nil
	case OpCodeAnotherSBConnected:
		// Handled by error handler
		return nil
	default:
		return fmt.Errorf("unexpected message op: %s", msg.Op)
	}
}

func (a *AuthServiceImpl) onHelloAck(msg *WSMessage) error {
	a.receivedNonce = msg.Data.Nonce
	a.CheckNonce()

	close(a.readyChan)

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
		Op:   OpCodeConsoleSet,
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
