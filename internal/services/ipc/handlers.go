package ipc

import (
	"fmt"
	"os"

	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
)

func (i *ipcService) handleMessage(msg *IPCMessage) error {
	switch msg.Op {
	case OpCodePing:
		i.readyChan <- struct{}{}
		return nil

	case OpCodeExit:
		os.Exit(0)
		return nil

	case OpCodeConfigUpdate:
		if msg.Data == nil || msg.Data.Config == nil {
			return fmt.Errorf("missing config data")
		}

		i.eventBroker.Publish(eventbroker.EventSendConfigUpdateSync, msg.Data.Config)
		i.configSvc.SetConfigFromUI(msg.Data.Config)
		return nil

	default:
		return fmt.Errorf("unexpected message op: %d", msg.Op)
	}
}
