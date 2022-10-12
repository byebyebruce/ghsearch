package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/byebyebruce/ghsearch"
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
				ret, err := ghsearch.SearchRepo(token, page, lang, args...)
				if err != nil {
					return err
				}
				table := tablewriter.NewWriter(cmd.OutOrStdout())
				table.SetHeader([]string{"#", "repo", "star", "last push", "Description"})
				table.SetBorder(false) // Set Border to false
				table.SetAutoWrapText(false)

				for i, v := range ret {
					updateAt := v.PushedAt.Format("2006-01-02")
					desc := v.Description
					const maxDesc = 32
					if len([]rune(v.Description)) > maxDesc {
						desc = (string)([]rune(v.Description)[:maxDesc]) + "..."
					}
					table.Append([]string{strconv.Itoa(i + 1) /*v.Name,*/, v.HTMLURL, strconv.Itoa(v.StargazersCount), updateAt, desc})
				}
				table.Render()
			} else {
				ret, err := ghsearch.SearchCode(token, page, lang, args...)
				if err != nil {
					return err
				}
				table := tablewriter.NewWriter(cmd.OutOrStdout())
				table.SetHeader([]string{"#", "file", "url"})
				table.SetBorder(false) // Set Border to false
				table.SetAutoWrapText(false)

				for i, v := range ret {
					table.Append([]string{strconv.Itoa(i + 1), v.Name, v.HTMLURL})
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
