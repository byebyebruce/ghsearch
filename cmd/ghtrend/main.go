package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/byebyebruce/ghsearch"
	"github.com/byebyebruce/ghsearch/util"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/sync/errgroup"
)

var (
	dates      = []string{"daily", "weekly", "monthly"}
	spokenLang = ""
	lang       = ""
)

type Repos []*ghsearch.Repository

func main() {
	flag.StringVar(&spokenLang, "spoken", "", "spoken language[zh/en/de/fr]. empty means any")
	flag.StringVar(&lang, "lang", "go", "program languages:go,rust,c,c++,java,c#,js")
	flag.Parse()

	if err := ui.Init(); err != nil {
		fmt.Println("failed to initialize termui", err)
		return
	}

	defer ui.Close()

	// get data
	ret, err := util.AsyncTaskAndShowLoadingBar("loading", func() ([]Repos, error) {
		reps := make([]Repos, len(dates))
		eg, _ := errgroup.WithContext(context.Background())
		for i, v := range dates {
			idx := i
			date := v
			eg.Go(func() error {
				time.Sleep(time.Millisecond * time.Duration(idx*50)) // 避免github api限制
				ret, err := ghsearch.TrendingRepos(lang, date, spokenLang)
				if err != nil {
					return err
				}
				reps[idx] = ret
				return nil
			})
		}
		return reps, eg.Wait()
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// tab
	tabpane := widgets.NewTabPane(dates...)
	tabpane.Title = "Date"
	tabpane.Border = true
	tabpane.BorderStyle.Fg = ui.ColorYellow

	// list
	l := widgets.NewList()
	l.Title = "List"
	l.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorCyan)
	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.WrapText = false

	// desc
	p := widgets.NewParagraph()
	p.Title = "Desc"
	p.TextStyle.Fg = ui.ColorGreen
	p.BorderStyle.Fg = ui.ColorCyan

	// help
	help := widgets.NewParagraph()
	help.Title = fmt.Sprintf("j↑/k↓: down/up, h←/l→: tab left/regiht, enter: open, q: quit")
	help.TitleStyle = ui.NewStyle(ui.ColorCyan)
	help.Border = false

	// grid
	grid := ui.NewGrid()
	grid.Set(
		ui.NewCol(0.4, l),
		ui.NewCol(0.6, p),
	)

	showList := func(repos Repos) {
		l.Rows = l.Rows[:0]
		for i, v := range repos {
			l.Rows = append(l.Rows, fmt.Sprintf("%2d ⭐%-6d %s/%s", i+1, v.Stars, v.Author, v.Name))
		}
	}

	showDesc := func(current *ghsearch.Repository) {
		p.Text = fmt.Sprintf(`[Project: %s](fg:white,mod:bold)
[Author: %s](fg:red)
[Link: %s](fg:blue)
Desc: 
    %s
`, current.Name, current.Author, current.Link, current.Desc)
	}

	currentList := ret[0]
	onTab := func(index int) {
		currentList = ret[index]
		showList(currentList)
		showDesc(currentList[l.SelectedRow])
	}
	onResize := func(w, h int) {
		ui.Clear()
		const tabOffset = 3
		const helpOffset = 1
		tabpane.SetRect(0, 0, w, tabOffset)
		grid.SetRect(0, tabOffset, w, h-helpOffset)
		help.SetRect(0, h-helpOffset, w, h)

	}
	render := func() {
		ui.Render(tabpane, help, grid)
	}

	termWidth, termHeight := ui.TerminalDimensions()
	onTab(tabpane.ActiveTabIndex)
	onResize(termWidth, termHeight)
	render()

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
			showDesc(currentList[l.SelectedRow])
		case "k", "<Up>":
			l.ScrollUp()
			showDesc(currentList[l.SelectedRow])
		case "h", "<Left>":
			tabpane.FocusLeft()
			onTab(tabpane.ActiveTabIndex)
		case "l", "<Right>":
			tabpane.FocusRight()
			onTab(tabpane.ActiveTabIndex)
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
			link := currentList[l.SelectedRow].Link
			util.OpenWebBrowser(link)
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			onResize(payload.Width, payload.Height)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		render()
	}
}
