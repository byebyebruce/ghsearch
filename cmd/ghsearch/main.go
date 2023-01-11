package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/byebyebruce/ghsearch"
	"github.com/byebyebruce/ghsearch/util"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
)

const (
	COUNT_PER_PAGE = 30
)

// GITHUB_TOKEN api token go build -ldflags "-X main.Version=$(TOKEN)"
var GITHUB_TOKEN string

func main() {
	var (
		lang  string
		token string
		code  bool // false:search repo, true:search code
	)
	rootCmd := &cobra.Command{
		Use:          "ghsearch",
		Short:        "github repo search",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, osArgs []string) error {
			if len(token) == 0 {
				token = GITHUB_TOKEN
				if env := os.Getenv("GITHUB_TOKEN"); len(env) > 0 {
					token = env
				}
			}

			if token == "" {
				return fmt.Errorf("token is emtpy. please set env GITHUB_TOKEN or use --token=xx")
			}

			var (
				err error
			)

			args := osArgs
			for {
				if len(args) == 0 {
					fmt.Println("Please input key words, split by space. Press q to exit.")
					var words string
					reader := bufio.NewReader(os.Stdin)
					bytes, _, _ := reader.ReadLine()
					words = strings.TrimSpace(string(bytes))
					if len(words) == 0 {
						continue
					}
					if words == "q" {
						return nil
					}
					args = strings.Split(words, " ")
				}
				if code {
					err = searchCode(token, lang, args...)
				} else {
					err = searchRepo(token, lang, args...)
				}
				args = args[:0]
				if err != nil {
					return err
				}
			}
		},
	}
	rootCmd.Flags().StringVar(&lang, "lang", "go", "language")
	rootCmd.Flags().StringVar(&token, "token", "", "github api token")
	rootCmd.Flags().BoolVar(&code, "code", false, "false:search repo, true:search code")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func searchRepo(token, lang string, args ...string) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	defer ui.Clear()

LOOP:
	for {
		var (
			page = 1
		)
		ui.Clear()

		ret, err := util.AsyncTaskAndShowLoadingBar("loading", func() ([]ghsearch.SearchRepoResultItems, error) {
			return ghsearch.SearchRepo(token, page, lang, args...)
		})
		if err != nil {
			return err
		}
		if len(ret) == 0 {
			fmt.Println("no result")
			time.Sleep(time.Second * 2)
			return nil
		}

		// title
		title := widgets.NewParagraph()
		title.Title = "Search"
		title.TextStyle.Fg = ui.ColorWhite
		title.Text = strings.Join(args, " ")

		// input
		input := widgets.NewParagraph()
		input.Title = fmt.Sprintf("j/k: up/down, enter: open, q: quit, ctrl+n/p: next/privious page")
		input.TitleStyle = ui.NewStyle(ui.ColorCyan)
		input.Border = false

		// list
		l := widgets.NewList()
		l.Title = "Repo"
		for i, v := range ret {
			l.Rows = append(l.Rows, fmt.Sprintf("%2d ⭐%-6d %s/%s", (page-1)*COUNT_PER_PAGE+i+1, v.StargazersCount, v.Owner.Login, v.Name))
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
			l.Title = fmt.Sprintf("Page:%d", page)
			current := ret[idx]
			updateAt := ret[idx].PushedAt.Format(time.RFC3339)
			p.Text = fmt.Sprintf(`[Project: %s](fg:white,mod:bold)
[Author: %s](fg:red)
[Link: %s](fg:blue)
[Last Update: %s](fg:blue)
Desc: 
    %s
`, current.Name, current.Owner.Login, current.HTMLURL, updateAt, current.Description)
		}

		grid := ui.NewGrid()
		termWidth, termHeight := ui.TerminalDimensions()
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

		for {
			e := <-uiEvents
			switch e.ID {
			case "q":
				return nil
			case "<C-c>":
				ui.Close()
				os.Exit(0)
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
			case "<C-n>":
				page += 1
				continue LOOP
			case "<C-p>":
				if page > 1 {
					page -= 1
					continue LOOP
				}
			case "g":
				if previousKey == "g" {
					l.ScrollTop()
				}
			case "<Home>":
				l.ScrollTop()
			case "G", "<End>":
				l.ScrollBottom()
			case "<Enter>":
				link := ret[l.SelectedRow].HTMLURL
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
	return nil
}
func searchCode(token, lang string, args ...string) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	defer ui.Clear()

LOOP:
	for {
		var (
			page = 1
		)
		ui.Clear()

		ret, err := util.AsyncTaskAndShowLoadingBar("loading", func() ([]ghsearch.SearchCodeResultItems, error) {
			return ghsearch.SearchCode(token, 0, lang, args...)
		})
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		if len(ret) == 0 {
			fmt.Println("no result")
			time.Sleep(time.Second * 2)
			return nil
		}

		// title
		title := widgets.NewParagraph()
		title.Title = "Search"
		title.TextStyle.Fg = ui.ColorWhite
		title.Text = strings.Join(args, " ")

		// input
		input := widgets.NewParagraph()
		input.Title = fmt.Sprintf("j/k: up/down, enter: open, q: quit, ctrl+n/p: next/privious page")
		input.TitleStyle = ui.NewStyle(ui.ColorCyan)
		input.Border = false

		// list
		l := widgets.NewList()
		l.Title = "Repo"
		for i, v := range ret {
			l.Rows = append(l.Rows, fmt.Sprintf("%2d %s [%s](fg:yellow)", (page-1)*COUNT_PER_PAGE+i+1, v.Repository.Name, v.Name))
		}
		l.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorCyan)
		l.TextStyle = ui.NewStyle(ui.ColorWhite)
		l.WrapText = false

		// desc
		p := widgets.NewParagraph()
		p.Title = "Desc"
		p.TextStyle.Fg = ui.ColorGreen
		p.BorderStyle.Fg = ui.ColorCyan
		p.WrapText = true
		showDesc := func(idx int) {
			l.Title = fmt.Sprintf("Page:%d", page)
			current := ret[idx]
			p.Text = fmt.Sprintf(`%s
[Project: %s](fg:white,mod:bold)
[Author: %s](fg:red)
[File: %s](fg:yellow)
[Score: %02f](fg:blue)
Desc:
    %s
`,
				current.HTMLURL,
				current.Repository.Name,
				current.Repository.Owner.Login,
				current.Path,
				current.Score,
				current.Repository.Description)
		}

		grid := ui.NewGrid()
		termWidth, termHeight := ui.TerminalDimensions()
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

		for {
			e := <-uiEvents
			switch e.ID {
			case "q":
				return nil
			case "<C-c>":
				ui.Close()
				os.Exit(0)
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
			case "<C-n>":
				page += 1
				continue LOOP
			case "<C-p>":
				if page > 1 {
					page -= 1
					continue LOOP
				}
			case "g":
				if previousKey == "g" {
					l.ScrollTop()
				}
			case "<Home>":
				l.ScrollTop()
			case "G", "<End>":
				l.ScrollBottom()
			case "<Enter>":
				link := ret[l.SelectedRow].HTMLURL
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
	return nil
}
