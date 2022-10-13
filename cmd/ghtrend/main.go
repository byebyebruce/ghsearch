package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/byebyebruce/ghsearch"
	"github.com/byebyebruce/ghsearch/util"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
)

var (
	flagSpokenLang = flag.String("spoken", "", "spoken language[zh/en/de/fr]. empty means any")
	flagLang       = flag.String("lang", "go,rust,c,c++,java,c#,js", "program languages")
)

func main() {
	flag.Parse()
	var err error
	for {
		// select program language. skip selection if only one type
		langs := strings.Split(*flagLang, ",")
		lang := langs[0]
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
			time.Sleep(time.Second * 5)
			continue
		}

		// render data
		bf := &bytes.Buffer{}
		table := tablewriter.NewWriter(bf)
		table.SetBorder(false)
		table.SetAutoWrapText(false)
		for i, v := range ret {
			desc := v.Desc
			const maxDesc = 16
			if len([]rune(v.Desc)) > maxDesc {
				desc = (string)([]rune(v.Desc)[:maxDesc]) + "..."
			}
			table.Append([]string{fmt.Sprintf("%2d ⭐%-6d %s", i+1, v.Stars, v.Link), desc})
		}
		table.Render()

		// read line
		reader := bufio.NewReader(bf)
		it := []string{}
		for {
			l, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			it = append(it, string(l))
		}

		// show result
		prompt = promptui.Select{
			Label: fmt.Sprintf("%s %s trending", lang, date),
			Size:  20,
			Items: it,
		}

		// select to open
		var (
			selectIdx int
			scrollPos int
		)
		for {
			i, _, err := prompt.RunCursorAt(selectIdx, scrollPos)
			if err != nil {
				break
			}
			link := ret[i].Link
			fmt.Println("open", link)
			util.OpenWebBrowser(link)
			selectIdx = prompt.ScrollPosition()
			selectIdx = i
		}
	}
}
