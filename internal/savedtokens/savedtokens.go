package savedtokens

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"golang.org/x/exp/slices"
)

var plainTokenRegex = regexp.MustCompile(`[\w-]{24}\.[\w-]{6}\.[\w-]{25,110}|mfa\.[\w-]{84}`)

type SavedTokens struct {
	RawTokens []string
	Tokens    []SavedToken
}

type SavedToken struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	ID       string `json:"id"`
	Avatar   string `json:"avatar"`
}

func Gets() []string {
	switch runtime.GOOS {
	case "windows":
		local := os.Getenv("LOCALAPPDATA")
		roaming := os.Getenv("APPDATA")

		return []string{
			roaming + "/Lightcord",
			roaming + "/discord",
			roaming + "/discordcanary",
			roaming + "/discordptb",
			local + "/Google/Chrome/User Data/Default",
			roaming + "/Opera Software/Opera Stable",
			roaming + "/Opera Software/Opera GX Stable",
			local + "/BraveSoftware/Brave-Browser/User Data/Default",
			local + "/Yandex/YandexBrowser/User Data/Default",
		}
	case "darwin":
		homedir, _ := os.UserHomeDir()
		dataPath := homedir + "/Library/Application Support"

		return []string{
			dataPath + "/Lightcord",
			dataPath + "/discord",
			dataPath + "/discordcanary",
			dataPath + "/discordptb",
			dataPath + "/Google/Chrome/Default",
			dataPath + "/Opera Software/Opera Stable",
			dataPath + "/Opera Software/Opera GX Stable",
			dataPath + "/BraveSoftware/Brave-Browser/Default",
			dataPath + "/Yandex/YandexBrowser/Default",
		}
	default: // Unix
		homedir, _ := os.UserHomeDir()

		return []string{
			homedir + "/.config/discord",
			homedir + "/.config/discordcanary",
			homedir + "/.config/discordptb",
			homedir + "/.config/google-chrome/Default",
			homedir + "/.config/opera",
			homedir + "/.config/BraveSoftware",
			homedir + "/.config/yandex-browser-beta/Default",
			homedir + "/.config/yandex-browser/Default",
		}
	}
}

func searchPlainTokens(content []byte, savedTokens *SavedTokens) {
	for _, match := range plainTokenRegex.FindAll(content, -1) {
		token := string(match)

		// Add the token if it doesn't exist yet
		if slices.Contains(savedTokens.RawTokens, token) {
			continue
		}

		savedTokens.RawTokens = append(savedTokens.RawTokens, token)
	}
}

func handleFile(path string, file os.DirEntry, savedTokens *SavedTokens) {
	if file.IsDir() {
		return
	}

	if !strings.HasSuffix(file.Name(), ".ldb") && !strings.HasSuffix(file.Name(), ".log") {
		return
	}

	content, err := os.ReadFile(filepath.Join(path, file.Name()))
	if err != nil {
		return
	}

	if runtime.GOOS == "windows" {
		searchEncryptedTokens(content, savedTokens)
	}

	searchPlainTokens(content, savedTokens)
}

func GetSavedTokens() []SavedToken {
	paths := Gets()

	savedTokens := &SavedTokens{}

	for _, path := range paths {
		// Check if the path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		path += "/Local Storage/leveldb"

		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, file := range files {
			handleFile(path, file, savedTokens)
		}
	}

	for _, rawToken := range savedTokens.RawTokens {
		session, err := discordgo.New(rawToken)
		if err != nil {
			continue
		}

		user, err := session.User("@me")
		if err != nil {
			continue
		}

		alreadyAdded := fn.Any(savedTokens.Tokens, func(token SavedToken) bool { return token.Token == rawToken })
		if alreadyAdded {
			continue
		}

		savedTokens.Tokens = append(savedTokens.Tokens, SavedToken{
			Token:    rawToken,
			Username: user.Username,
			ID:       user.ID,
			Avatar:   user.AvatarURL(""),
		})
	}

	return savedTokens.Tokens
}
