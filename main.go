package main

import (
	"fmt"
	"github.com/araddon/dateparse"
	"os"
	"strings"
	"time"

	"github.com/caseymrm/menuet"
)

const (
	placeholder = "2022-02-18 16:35"
)

var (
	targetTime time.Time
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

func helloClock() {
	setTargetTime()

	ticker := time.NewTicker(time.Second)

	for {

		select {
		case <-ticker.C:
			duration := targetTime.Sub(time.Now())
			duration -= duration % time.Second

			menuet.App().SetMenuState(&menuet.MenuState{
				Title: shortDur(duration),
			})
			menuet.App().MenuChanged()
		}
	}
}

func setTargetTime() {

	alert := menuet.Alert{
		MessageText:     "Input target time.",
		InformativeText: "Default is 2022-02-18 16:35\n\n2022/2/18 16:35\n2/18/2022 16:35\n18/2/2022 16:35\nAre all fine.",
		Buttons:         []string{"OK", "Quit"},
		Inputs:          []string{placeholder},
	}

	alertClicked := menuet.App().Alert(alert)

	if alertClicked.Button == 1 {
		os.Exit(0)
	}

	dateStr := alertClicked.Inputs[0]

	if alertClicked.Button == 0 && dateStr == "" {
		dateStr = placeholder
	}

	tz, _ := time.LoadLocation("Asia/Shanghai")
	t, err := dateparse.ParseIn(dateStr, tz)
	if err != nil {
		panic(err)
	}

	targetTime = t
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}

	if !targetTime.IsZero() {
		items = append(items, menuet.MenuItem{
			Text: fmt.Sprintf("Target: %s", targetTime.String()),
		})
	}

	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})

	items = append(items, menuet.MenuItem{
		Text:    "Set target",
		Clicked: setTargetTime,
	})

	return items
}

func main() {
	go helloClock()
	menuet.App().Children = menuItems
	menuet.App().Label = "me.shanicky.countdown"
	menuet.App().RunApplication()
}