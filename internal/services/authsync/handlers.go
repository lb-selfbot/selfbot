package authsync

import (
	"fmt"
	"os"

	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
)

func (a *authSyncService) handleMessage(msg *WSMessage) error {
	switch msg.Op {
	case OpCodeHelloAck:
		return a.handleHelloAck(msg)
	case OpCodeFreeTrialStart:
		return a.handleFreeTrialStart(msg)
	case OpCodeExit:
		return a.handleExit(msg)
	case OpCodeRestart:
		return a.handleRestart(msg)
	case OpCodeWebsiteConnect:
		return a.handleWebsiteConnect(msg)
	case OpCodeConfigUpdate:
		return a.handleConfigUpdate(msg)
	case OpCodePong:
		return nil
	case OpCodeAnotherSBConnected:
		// Handled by error handler
		return nil
	default:
		return fmt.Errorf("unexpected message op: %s", msg.Op)
	}
}

func (a *authSyncService) handleHelloAck(msg *WSMessage) error {
	a.receivedNonce = msg.Data.Nonce
	a.CheckNonce()

	close(a.readyChan)

	cfg := a.configSvc.GetConfig()
	return a.SendSync(SyncTypeConfigUpdate, cfg)
}

func (a *authSyncService) handleFreeTrialStart(_ *WSMessage) error {
	a.eventBroker.Publish(eventbroker.EventShowNotification, "Your free trial has started! You will be able to use LightningBot for free for the next 7 days.")
	return nil
}

func (a *authSyncService) handleExit(_ *WSMessage) error {
	os.Exit(0)
	return nil
}

func (a *authSyncService) handleRestart(_ *WSMessage) error {
	a.eventBroker.Publish(eventbroker.EventRestart, nil)
	return nil
}

func (a *authSyncService) handleWebsiteConnect(_ *WSMessage) error {
	err := a.SendSync(SyncTypeConsoleSet, a.consoleSvc.GetHTMLConsole())
	if err != nil {
		return fmt.Errorf("failed to send console state: %w", err)
	}

	cfg := a.configSvc.GetConfig()
	err = a.SendSync(SyncTypeConfigUpdate, cfg)
	if err != nil {
		return fmt.Errorf("failed to send config update: %w", err)
	}

	return nil
}

func (a *authSyncService) handleConfigUpdate(msg *WSMessage) error {
	if msg.Data.Config == nil {
		return nil
	}

	// Restore sensitive data
	cfg := a.configSvc.GetConfig()
	msg.Data.Config.Account = cfg.Account

	if err := a.configSvc.SaveConfig(msg.Data.Config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return a.SendSync(SyncTypeConfigUpdate, msg.Data.Config)
}
