package main

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/spf13/cast"
	"os"
	"strings"
	"time"

	"github.com/caseymrm/menuet"
)

const (
	targetPlaceholder  = "2022-02-18 16:35"
	elapsedPlaceholder = "1h"
)

var (
	targetTime    = time.Time{}
	allowNegative = false
)

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

func clockLoop() {
	choose()

	ticker := time.NewTicker(time.Second)

	for {

		select {
		case <-ticker.C:
			duration := targetTime.Sub(time.Now())
			duration -= duration % time.Second

			if duration == 0 {
				menuet.App().Notification(menuet.Notification{
					Title:      "Time's up!",
					Message:    fmt.Sprintf("Target is %s", targetTime.String()),
					Identifier: "countdown",
				})
			}

			if duration < 0 && !allowNegative {
				menuet.App().SetMenuState(&menuet.MenuState{
					Image: "icon.icns",
				})

				targetTime = time.Time{}

				continue
			}

			menuet.App().SetMenuState(&menuet.MenuState{
				Title: shortDur(duration),
			})

			menuet.App().MenuChanged()
		}
	}
}

func setTimeElapsed() {
	alert := menuet.Alert{
		MessageText:     "Input time elapsed.",
		InformativeText: "Default is 1 hour later",
		Buttons:         []string{"OK", "Cancel"},
		Inputs:          []string{elapsedPlaceholder},
	}

	alertClicked := menuet.App().Alert(alert)

	if alertClicked.Button == 1 {
		if targetTime.IsZero() {
			os.Exit(0)
		}

		return
	}

	dateStr := alertClicked.Inputs[0]

	if alertClicked.Button == 0 && dateStr == "" {
		dateStr = targetPlaceholder
	}

	duration, err := cast.ToDurationE(dateStr)
	if err != nil {
		panic(err)
	}

	targetTime = time.Now().Add(duration)
}

func setTargetTime() {
	alert := menuet.Alert{
		MessageText:     "Input target time.",
		InformativeText: "Default is 2022-02-18 16:35\n\n2022/2/18 16:35\n2/18/2022 16:35\n18/2/2022 16:35\nAre all fine.",
		Buttons:         []string{"OK", "Cancel"},
		Inputs:          []string{targetPlaceholder},
	}

	alertClicked := menuet.App().Alert(alert)

	if alertClicked.Button == 1 {
		if targetTime.IsZero() {
			os.Exit(0)
		}

		return
	}

	dateStr := alertClicked.Inputs[0]

	if alertClicked.Button == 0 && dateStr == "" {
		dateStr = targetPlaceholder
	}

	tz, _ := time.LoadLocation("Asia/Shanghai")
	t, err := dateparse.ParseIn(dateStr, tz)
	if err != nil {
		panic(err)
	}

	targetTime = t
}

func choose() {
	alert := menuet.Alert{
		MessageText: "Choose mode.",
		Buttons:     []string{"Target mode", "Countdown mode", "Cancel"},
	}

	alertClicked := menuet.App().Alert(alert)

	switch alertClicked.Button {
	case 0:
		setTargetTime()
	case 1:
		setTimeElapsed()
	default:
		os.Exit(0)
	}
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{
		{
			Text: "Made with love",
			Clicked: func() {
				menuet.App().Alert(menuet.Alert{
					MessageText: "mua~",
				})
			},
		},
	}

	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})

	if !targetTime.IsZero() {
		items = append(items, menuet.MenuItem{
			Text: fmt.Sprintf("Target: %s", targetTime.String()),
		})
	}

	items = append(items, menuet.MenuItem{
		Text:    "Set target time",
		Clicked: setTargetTime,
	})

	items = append(items, menuet.MenuItem{
		Text:    "Set elapsed duration",
		Clicked: setTimeElapsed,
	})

	items = append(items, menuet.MenuItem{
		Text:  "Allow negative time",
		State: allowNegative,
		Clicked: func() {
			allowNegative = !allowNegative
		},
	})

	return items
}
func main() {
	go menuet.App().Notification(menuet.Notification{
		Title: "Test notification!",
	})

	go clockLoop()
	menuet.App().Children = menuItems
	menuet.App().Label = "me.shanicky.countdown"
	menuet.App().RunApplication()
}
