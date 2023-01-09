package main

import (
	"flag"
	"fmt"

	"github.com/byebyebruce/ghsearch"
	"github.com/byebyebruce/ghsearch/util"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/manifoldco/promptui"
)

var (
	spokenLang = ""
	lang       = ""
)

func main() {
	flag.StringVar(&spokenLang, "spoken", "", "spoken language[zh/en/de/fr]. empty means any")
	flag.StringVar(&lang, "lang", "go", "program languages:go,rust,c,c++,java,c#,js")
	flag.Parse()

	if err := ui.Init(); err != nil {
		fmt.Println("failed to initialize termui", err)
		return
	}

	defer ui.Close()

	for {
		ui.Clear()

		var (
			err error
		)

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
			return ghsearch.TrendingRepos(lang, date, spokenLang)
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		// title
		title := widgets.NewParagraph()
		title.Title = "Trend"
		title.TextStyle.Fg = ui.ColorWhite
		title.Text = fmt.Sprintf("%s:%s", lang, date)

		// input
		input := widgets.NewParagraph()
		input.Title = fmt.Sprintf("j/k: up/down, enter: open, q: quit, ctrl+n/p: next/privious page")
		input.TitleStyle = ui.NewStyle(ui.ColorCyan)
		input.Border = false

		// list
		l := widgets.NewList()
		l.Title = "List"
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

		onResize := func(w, h int) {
			const titleOffset = 3
			const inputOffset = 1
			title.SetRect(0, 0, w, titleOffset)
			grid.SetRect(0, titleOffset, w, h-inputOffset)
			input.SetRect(0, h-inputOffset, w, h)
		}

		showDesc(l.SelectedRow)
		onResize(termWidth, termHeight)
		ui.Render(title, input, grid)

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
				onResize(payload.Width, payload.Height)
				ui.Clear()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

			ui.Render(title, input, grid)
		}
	}
}
