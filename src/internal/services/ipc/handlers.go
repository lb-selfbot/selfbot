package ipc

import (
	"fmt"
	"os"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
)

func (i *IPCServiceImpl) onMessage(msg *WSMessage) error {
	switch msg.Op {
	case OpCodeHello:
		return nil

	case OpCodeExit:
		os.Exit(0)
		return nil

	case OpCodeConfigUpdate:
		i.authSvc.SendConfigUpdate()
		config.SaveConfigFromUI(msg.Data.Config)
		return nil

	default:
		return fmt.Errorf("unexpected message op: %d", msg.Op)
	}
}
