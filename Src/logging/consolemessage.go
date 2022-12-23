package logging

import (
	"container/list"
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/utils"
)

type ConsoleMessage struct {
	Title string
	Color string
	Info  *list.List
}

func (c *ConsoleMessage) SetTitle(title string) *ConsoleMessage {
	c.Title = title
	return c
}

func (c *ConsoleMessage) SetColor(color string) *ConsoleMessage {
	c.Color = color
	return c
}

func (c *ConsoleMessage) AddInfo(key, value string) *ConsoleMessage {
	c.Info.PushBack([]string{key, value})
	return c
}

func (c *ConsoleMessage) AddInfoConditional(key, value string, condition bool) *ConsoleMessage {
	if condition {
		c.AddInfo(key, value)
	}
	return c
}

func (c *ConsoleMessage) Show() {
	if constants.CONSOLE {
		r, g, b := utils.Misc.HexToRGB(c.Color)
		timeFormat := Console.Colored(time.Now().Local().Format("15:04:05"), 193, 239, 255, false)
		message := Console.Colored("[", 238, 238, 238, false) + timeFormat + Console.Colored("] ", 238, 238, 238, false) + Console.Colored(c.Title, r, g, b, false) + "\n"

		for e := c.Info.Front(); e != nil; e = e.Next() {
			info := e.Value.([]string)
			if info[0] != "" {
				message += Console.Colored("  "+info[0]+": ", 238, 238, 238, false) + Console.Colored(info[1], 193, 239, 255, false) + "\n"
			} else {
				message += Console.Colored("  "+info[1], 193, 239, 255, false) + "\n"
			}
		}

		fmt.Println(message)
	}

	timeFormat := UI.Colored(time.Now().Local().Format("15:04:05"), "#C1EFFF")
	message := UI.Colored("[", "#eeeeee") + timeFormat + UI.Colored("] ", "#eeeeee") + UI.Colored(c.Title, c.Color) + "\n"

	for e := c.Info.Front(); e != nil; e = e.Next() {
		info := e.Value.([]string)
		if info[0] != "" {
			message += UI.Colored("  "+info[0]+": ", "#eeeeee") + UI.Colored(info[1], "#C1EFFF") + "\n"
		} else {
			message += UI.Colored("  "+info[1], "#C1EFFF") + "\n"
		}
	}

	UI.Append(message)
}

func NewConsoleMessage() *ConsoleMessage {
	return &ConsoleMessage{
		Info: list.New(),
	}
}
