package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/byebyebruce/ghsearch"
	"github.com/byebyebruce/ghsearch/util"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/manifoldco/promptui"
)

var (
	flagSpokenLang = flag.String("spoken", "", "spoken language[zh/en/de/fr]. empty means any")
	flagLang       = flag.String("lang", "go,rust,c,c++,java,c#,js", "program languages")
)

func main() {
	flag.Parse()

	if err := ui.Init(); err != nil {
		fmt.Println("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	for {
		ui.Clear()

		var (
			err   error
			langs = strings.Split(*flagLang, ",")
			lang  = langs[0]
		)
		if len(langs) > 1 {
			prompt := promptui.Select{
				Label: "select language",
				Items: langs,
			}
			_, lang, err = prompt.Run()
			if err != nil {
				return
			}
		}

		// select date range
		prompt := promptui.Select{
			Label: "select date range",
			Items: []string{"daily", "weekly", "monthly"},
		}
		_, date, err := prompt.Run()
		if err != nil {
			return
		}

		// get data
		ret, err := util.AsyncTaskAndShowLoadingBar("loading", func() ([]*ghsearch.Repository, error) {
			return ghsearch.TrendingRepos(lang, date, *flagSpokenLang)
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		// list
		l := widgets.NewList()
		l.Title = "Trend"
		for i, v := range ret {
			l.Rows = append(l.Rows, fmt.Sprintf("%2d ⭐%-6d %s/%s", i+1, v.Stars, v.Author, v.Name))
		}
		l.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorCyan)
		l.TextStyle = ui.NewStyle(ui.ColorWhite)
		l.WrapText = false

		// desc
		p := widgets.NewParagraph()
		p.Title = "Desc"
		p.TextStyle.Fg = ui.ColorGreen
		p.BorderStyle.Fg = ui.ColorCyan
		showDesc := func(idx int) {
			current := ret[idx]
			p.Text = fmt.Sprintf(`[Project: %s](fg:white,mod:bold)
[Author: %s](fg:red)
[Link: %s](fg:blue)
Desc: 
    %s
`, current.Name, current.Author, current.Link, current.Desc)
		}
		showDesc(l.SelectedRow)

		grid := ui.NewGrid()
		termWidth, termHeight := ui.TerminalDimensions()
		grid.SetRect(0, 0, termWidth, termHeight)
		grid.Set(
			ui.NewCol(0.4, l),
			ui.NewCol(0.6, p),
		)

		ui.Render(grid)

		previousKey := ""
		uiEvents := ui.PollEvents()
	LOOP:
		for {
			e := <-uiEvents
			switch e.ID {
			case "q":
				break LOOP
			case "<C-c>":
				return
			case "j", "<Down>":
				l.ScrollDown()
				showDesc(l.SelectedRow)
			case "k", "<Up>":
				l.ScrollUp()
				showDesc(l.SelectedRow)
			case "<C-d>":
				l.ScrollHalfPageDown()
			case "<C-u>":
				l.ScrollHalfPageUp()
			case "<C-f>":
				l.ScrollPageDown()
			case "<C-b>":
				l.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					l.ScrollTop()
				}
			case "<Home>":
				l.ScrollTop()
			case "G", "<End>":
				l.ScrollBottom()
			case "<Enter>":
				link := ret[l.SelectedRow].Link
				util.OpenWebBrowser(link)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

			ui.Render(grid)
		}
	}
}
