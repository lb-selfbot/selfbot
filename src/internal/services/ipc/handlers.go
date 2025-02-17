package ipc

import (
	"fmt"
	"os"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
)

func (i *IPCServiceImpl) onMessage(msg *WSMessage) error {
	switch msg.Op {
	case OP_HELLO:
		return nil

	case OP_EXIT:
		os.Exit(0)
		return nil

	case OP_CONFIG_UPDATE:
		config.SaveConfigFromUI(msg.Data.Config)
		// TODO: config callbacks
		return nil

	default:
		return fmt.Errorf("unexpected message op: %d", msg.Op)
	}
}
