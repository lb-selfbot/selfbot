//go:build lbmobile

package mobile

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/LightningDev1/LB-Selfbot/mobile/gui"
)

func Main() {
	g := gui.New("LightningBot")

	fmt.Println("Starting LightningBot...")

	go SBMain()

	g.Main()
}

func SBMain() {

	fmt.Println("Continuing LightningBot...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
