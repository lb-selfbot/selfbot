package startup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/assets"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

// Initialize sets up basic requirements for both desktop and hosting modes
func Initialize(verbose bool) error {
	if err := utils.IO.CreateFolders(); err != nil {
		return fmt.Errorf("create folders: %w", err)
	}

	// Initialize icon
	iconFile := filepath.Join(utils.IO.GetMainFolder(), "Logo.png")
	if !utils.IO.Exists(iconFile) && !constants.HOSTING {
		icon, err := assets.Assets.ReadFile("Logo.png")
		if err != nil {
			logging.UI.Error("Error reading icon file: " + err.Error())
		}

		_ = utils.IO.WriteBytes(iconFile, icon)
	}

	// Initialize logging
	logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")
	if err := utils.IO.Write(logPath, ""); err != nil {
		return fmt.Errorf("clear log file: %w", err)
	}

	discordgo.Logger = func(msgL, caller int, format string, a ...any) {
		pc, file, line, _ := runtime.Caller(caller + 1)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		if strings.Contains(msg, "unknown event") && !verbose {
			return
		}

		msg = fmt.Sprintf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)

		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return
		}
		defer f.Close()

		_, _ = f.WriteString(msg)

		if msgL == discordgo.LogError {
			logging.UI.Error("DGC error: " + msg)
		}
	}

	return nil
}
