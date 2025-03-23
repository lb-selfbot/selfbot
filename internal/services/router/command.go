package router

type CommandFlag int

const (
	// CmdFlagDangerous is a flag for dangerous commands.
	CmdFlagDangerous CommandFlag = iota

	// CmdFlagHostingDangerous is a flag for commands that
	// are dangerous when hosting the bot.
	CmdFlagHostingDangerous

	// CmdFlagGuildOnly is a flag for commands that are only available in guilds.
	CmdFlagGuildOnly
)

// Category represents a command category.
type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Command represents a command.
type Command struct {
	Name        string           `json:"name"`
	Aliases     []string         `json:"aliases,omitempty"`
	Description string           `json:"description,omitempty"`
	Arguments   string           `json:"arguments,omitempty"`
	Flags       []CommandFlag    `json:"flags,omitempty"`
	SubCommands []*Command       `json:"subcommands,omitempty"`
	Category    *Category        `json:"-"`
	Handler     ExecutionHandler `json:"-"`
}

// GetIdentifiers returns the identifiers of the command.
func (c *Command) GetIdentifiers() []string {
	return append(c.Aliases, c.Name)
}

// GetSubCmd returns the subcommand with the given name.
func (r *Command) GetSubCmd(name string, ignoreCase bool) *Command {
	for _, subCmd := range r.SubCommands {
		if stringArrayContains(subCmd.GetIdentifiers(), name, ignoreCase) {
			return subCmd
		}
	}

	return nil
}
