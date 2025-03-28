package embed

import (
	"fmt"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fn"
)

const (
	EmbedFormatCodeblock  = "Codeblock"
	EmbedFormatText       = "Text"
	EmbedFormatIndent     = "Indent"
	EmbedFormatURLEmbed   = "URLEmbed"
)

type EmbedField struct {
	Name  string
	Value string
}

type EmbedSendCallback func(*Embed) []*discordgo.Message

type Embed struct {
	Title       string
	Description string
	Subtext     string
	Fields      []EmbedField
	ImageURL    string

	sendCallback EmbedSendCallback
}

func NewEmbed(sendCallback EmbedSendCallback) *Embed {
	return &Embed{sendCallback: sendCallback}
}

func (e *Embed) SetTitle(format string, a ...any) *Embed {
	e.Title = fmt.Sprintf(format, a...)
	return e
}

func (e *Embed) SetDescription(format string, a ...any) *Embed {
	e.Description = fmt.Sprintf(format, a...)
	return e
}

func (e *Embed) SetSubtext(format string, a ...any) *Embed {
	e.Subtext = fmt.Sprintf(format, a...)
	return e
}

func (e *Embed) AddField(name, value string, cond ...bool) *Embed {
	if len(cond) > 0 && !cond[0] {
		return e
	}

	e.Fields = append(e.Fields, EmbedField{Name: name, Value: value})
	return e
}

func (e *Embed) SetImage(url string) *Embed {
	e.ImageURL = url
	return e
}

func (e *Embed) Send() []*discordgo.Message {
	return e.sendCallback(e)
}

func (e *Embed) ToEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       e.Title,
		Description: e.Description,
		Fields: fn.Map(e.Fields, func(field EmbedField) *discordgo.MessageEmbedField {
			return &discordgo.MessageEmbedField{Name: field.Name, Value: field.Value}
		}),
		Footer: &discordgo.MessageEmbedFooter{
			Text: e.Subtext,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: e.ImageURL,
		},
		Color: 0x3990c1,
	}
}

func (e *Embed) ToString(format string, maxLen int) []string {
	switch format {
	case EmbedFormatCodeblock:
		return e.toCodeblock(maxLen)
	case EmbedFormatText:
		return e.toText(maxLen)
	case EmbedFormatIndent:
		return e.toIndent(maxLen)
	case EmbedFormatURLEmbed:
		return e.toURLEmbed(maxLen)
	default:
		return e.toCodeblock(maxLen)
	}
}

func (e *Embed) toCodeblock(maxLen int) []string {
	contentBuilder := new(strings.Builder)

	if e.Title != "" {
		contentBuilder.WriteString("```ini\n[" + e.Title + "]\n```")
	}

	if e.Description != "" {
		contentBuilder.WriteString("```md\n" + e.Description + "\n```")
	}

	for _, field := range e.Fields {
		contentBuilder.WriteString("```md\n# " + field.Name + "\n" + field.Value + "\n```")
	}

	if e.Subtext != "" {
		contentBuilder.WriteString("```ansi\n" + e.Subtext + "\n```")
	}

	content := contentBuilder.String()

	if len(content) < maxLen {
		return []string{content}
	}

	msgs := make([]string, len(content)/maxLen+1)

	currentIndex := 0
	messageIndex := 0

	for currentIndex < len(content) {
		// If we're at a codeblock
		if currentIndex+3 <= len(content) && content[currentIndex:currentIndex+3] == "```" {
			// Find the end of this codeblock
			endIndex := currentIndex + 3
			for endIndex < len(content) {
				if endIndex+3 <= len(content) && content[endIndex:endIndex+3] == "```" {
					endIndex += 3
					break
				}
				endIndex++
			}

			codeblock := content[currentIndex:endIndex]

			// Check if the codeblock fits in the current message
			if len(msgs[messageIndex])+len(codeblock) <= maxLen {
				msgs[messageIndex] += codeblock
			} else {
				// Move to next message
				messageIndex++
				msgs[messageIndex] = codeblock
			}

			currentIndex = endIndex
		} else {
			// Regular text, add char by char until next codeblock or maxLen
			if len(msgs[messageIndex])+1 > maxLen {
				messageIndex++
				msgs[messageIndex] = ""
			}

			msgs[messageIndex] += string(content[currentIndex])
			currentIndex++
		}
	}

	if e.ImageURL != "" {
		msgs = append(msgs, e.ImageURL)
	}

	return fn.Filter(msgs, func(line string) bool {
		return line != ""
	})
}

func (e *Embed) toText(maxLen int) []string {
	contentBuilder := new(strings.Builder)

	if e.Title != "" {
		contentBuilder.WriteString("```ini\n[" + e.Title + "]```\n\n")
	}

	if e.Description != "" {
		contentBuilder.WriteString(e.Description + "\n\n")
	}

	for _, field := range e.Fields {
		contentBuilder.WriteString("> **" + field.Name + "**\n" + field.Value + "\n\n")
	}

	if e.Subtext != "" {
		contentBuilder.WriteString("```ansi\n" + e.Subtext + "\n```")
	}

	content := contentBuilder.String()
	var msgs []string
	newStr := ""

	for _, line := range strings.Split(content, "\n") {
		if len(newStr+line)+1 > maxLen {
			msgs = append(msgs, newStr)
			newStr = ""
		}
		newStr += line + "\n"
	}

	msgs = append(msgs, newStr)

	if e.ImageURL != "" {
		msgs = append(msgs, e.ImageURL)
	}

	return fn.Filter(msgs, func(line string) bool {
		return line != ""
	})
}

func (e *Embed) toIndent(maxLen int) []string {
	contentBuilder := new(strings.Builder)

	if e.Title != "" {
		contentBuilder.WriteString("```ini\n[" + e.Title + "]```\n\n")
	}
	if e.Description != "" {
		contentBuilder.WriteString(e.Description + "\n\n")
	}

	for _, v := range e.Fields {
		contentBuilder.WriteString("**" + v.Name + "**\n" + v.Value + "\n\n")
	}

	if e.Subtext != "" {
		contentBuilder.WriteString("```ansi\n" + e.Subtext + "\n```")
	}

	content := contentBuilder.String()
	splitContent := strings.Split(content, "\n")

	var msgs []string

	msgBuilder := new(strings.Builder)

	for _, line := range splitContent {
		if msgBuilder.Len()+len(line)+3 > maxLen {
			msgs = append(msgs, msgBuilder.String())
			msgBuilder.Reset()
		}

		msgBuilder.WriteString("> " + line + "\n")
	}

	msgs = append(msgs, msgBuilder.String())

	if e.ImageURL != "" {
		msgs = append(msgs, e.ImageURL)
	}

	return fn.Filter(msgs, func(line string) bool {
		return line != ""
	})
}

func (e *Embed) toURLEmbed(maxLen int) []string {
	// TODO: urlembed
	return []string{}

	// contentBuilder := new(strings.Builder)

	// if e.Description != "" {
	// 	contentBuilder.WriteString(e.Description)
	// }

	// for _, v := range e.Fields {
	// 	contentBuilder.WriteString("\n\n• " + v.Name + "\n" + v.Value)
	// }

	// if e.Subtext != "" {
	// 	contentBuilder.WriteString("\n\n" + constants.ANSI_REGEX.ReplaceAllString(e.Subtext, ""))
	// }

	// content := contentBuilder.String()
	// var descriptions []string
	// var current string
	// for _, line := range strings.Split(content, "\n") {
	// 	if len(current+line)+1 > 350 {
	// 		descriptions = append(descriptions, current)
	// 		current = ""
	// 	}
	// 	current += line + "\n"
	// }

	// descriptions = append(descriptions, current)

	// return fn.Filter(descriptions, func(line string) bool {
	// 	return line != ""
	// })
}
