package ipc

import (
	"fmt"
	"os"

	"github.com/LightningDev1/LB-Selfbot/internal/services/eventbroker"
)

func (i *ipcService) handleMessage(msg *WSMessage) error {
	switch msg.Op {
	case OpCodeHello:
		return nil

	case OpCodeExit:
		os.Exit(0)
		return nil

	case OpCodeConfigUpdate:
		i.eventBroker.Publish(eventbroker.EventSendConfigUpdateSync, msg.Data.Config)
		i.configSvc.SetConfigFromUI(msg.Data.Config)
		return nil

	default:
		return fmt.Errorf("unexpected message op: %d", msg.Op)
	}
}
