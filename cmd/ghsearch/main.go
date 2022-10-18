package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/byebyebruce/ghsearch"
	"github.com/byebyebruce/ghsearch/util"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// GITHUB_TOKEN api token go build -ldflags "-X main.Version=$(TOKEN)"
var GITHUB_TOKEN string

func main() {
	var (
		lang  string
		token string
		page  int
		code  bool // false:search repo, true:search code
	)
	rootCmd := &cobra.Command{
		Use:          "ghsearch keyword...",
		Short:        "github repo search",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(token) == 0 {
				token = GITHUB_TOKEN
				if env := os.Getenv("GITHUB_TOKEN"); len(env) > 0 {
					token = env
				}
			}

			if token == "" {
				return fmt.Errorf("token is emtpy. please set env GITHUB_TOKEN or use --token=xx")
			}

			if !code {
				ret, err := util.AsyncTaskAndShowLoadingBar("loading", func() ([]ghsearch.SearchRepoResultItems, error) {
					return ghsearch.SearchRepo(token, page, lang, args...)
				})
				if err != nil {
					return err
				}
				var (
					link []string
					bf   = &bytes.Buffer{}
				)
				table := tablewriter.NewWriter(bf)
				//table.SetHeader([]string{"#", "repo", "star", "last push", "Description"})
				table.SetBorder(false) // Set Border to false
				table.SetAutoWrapText(false)

				for i, v := range ret {
					updateAt := v.PushedAt.Format("2006-01-02")
					desc := v.Description
					const maxDesc = 16
					if len([]rune(v.Description)) > maxDesc {
						desc = (string)([]rune(v.Description)[:maxDesc]) + "..."
					}
					co1 := fmt.Sprintf("%2d ⭐%-6d %s", i+1, v.StargazersCount, v.HTMLURL)
					table.Append([]string{co1, updateAt, desc})
					link = append(link, v.HTMLURL)
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
				prompt := promptui.Select{
					Label: append([]string{"lang:" + lang}, args...), //fmt.Sprintf("search %s %s", args...),
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
					link := link[i]
					fmt.Println("open", link)
					util.OpenWebBrowser(link)
					selectIdx = prompt.ScrollPosition()
					selectIdx = i
				}

			} else {
				ret, err := util.AsyncTaskAndShowLoadingBar("loading", func() ([]ghsearch.SearchCodeResultItems, error) {
					return ghsearch.SearchCode(token, page, lang, args...)
				})
				if err != nil {
					return err
				}
				table := tablewriter.NewWriter(cmd.OutOrStdout())
				//table.SetHeader([]string{"#", "file", "url"})
				table.SetBorder(false) // Set Border to false
				table.SetAutoWrapText(false)

				const replaceToken = "blob/"
				for i, v := range ret {
					_ = i
					url := v.HTMLURL
					/*
						//table.Append([]string{strconv.Itoa(i + 1), v.Name, v.HTMLURL})
						idx := strings.Index(v.HTMLURL, replaceToken)
						if idx != -1 {
							last := strings.Index(v.HTMLURL[idx+len(replaceToken):], "/")
							if last != -1 {
								url = v.HTMLURL[:idx] + "..." + v.HTMLURL[idx+len(replaceToken)+last:]
							}
						}
					*/
					table.Append([]string{strconv.Itoa(i + 1), url})
					//link = append(link, v.HTMLURL)
				}
				table.Render()
			}

			return nil
		},
	}
	rootCmd.Flags().StringVar(&lang, "lang", "go", "language")
	rootCmd.Flags().StringVar(&token, "token", "", "github api token")
	rootCmd.Flags().IntVar(&page, "page", 1, "page index")
	rootCmd.Flags().BoolVar(&code, "code", false, "false:search repo, true:search code")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
