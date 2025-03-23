package slashcommand

import (
	"errors"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

var errInvalidSubCommand = errors.New("invalid subcommand")

func optionError(name string, err error) error {
	return errors.New("Error in option " + name + ": " + err.Error())
}

func ArgumentsToCommandOptions(options []*discordgo.ApplicationCommandOption, args *router.Arguments, offset int) ([]*discordgo.ApplicationCommandOption, error) {
	if len(options) == 0 {
		// No need to do anything because there are no options
		return options, nil
	}

	// Subcommand handling
	isSubCommand := options[0].Type == discordgo.ApplicationCommandOptionSubCommand ||
		options[0].Type == discordgo.ApplicationCommandOptionSubCommandGroup

	if isSubCommand {
		subCommands := utils.Filter(options, func(option *discordgo.ApplicationCommandOption) bool {
			return option.Name == args.Get(0+offset).Raw()
		})

		if len(subCommands) == 0 {
			return nil, errInvalidSubCommand
		}

		subCommand := subCommands[0]

		subOptions, err := ArgumentsToCommandOptions(subCommand.Options, args, offset+1)
		if err != nil {
			return nil, err
		}

		return []*discordgo.ApplicationCommandOption{
			{
				Type:    subCommand.Type,
				Name:    subCommand.Name,
				Options: subOptions,
			},
		}, nil
	}

	// Normal command option handling
	for i, option := range options {
		arg := args.Get(i + offset)

		switch option.Type {
		case discordgo.ApplicationCommandOptionString:
			option.Value = arg.Raw()
		case discordgo.ApplicationCommandOptionInteger, discordgo.ApplicationCommandOptionNumber:
			Int, err := arg.Int()
			if err != nil {
				return nil, optionError(option.Name, err)
			}

			option.Value = Int
		case discordgo.ApplicationCommandOptionBoolean:
			asBool, err := arg.Bool()
			if err != nil {
				return nil, optionError(option.Name, err)
			}

			option.Value = asBool
		case discordgo.ApplicationCommandOptionUser, discordgo.ApplicationCommandOptionChannel, discordgo.ApplicationCommandOptionRole, discordgo.ApplicationCommandOptionMentionable:
			option.Value = arg.UserMentionID()
		}
	}

	return options, nil
}
