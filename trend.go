package ghsearch

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// thanks https://github.com/darjun/ghtrending

// Repository represent a repository in the trending list.
type Repository struct {
	Author  string
	Name    string
	Link    string
	Desc    string
	Lang    string
	Stars   int
	Forks   int
	Add     int
	BuiltBy []string
}

// Developer represent a developer in the developer trending list.
type Developer struct {
	Name        string
	Username    string
	PopularRepo string
	Desc        string
}

const GitHubURL = "https://github.com"

// TrendingRepos fetch all repositories from  GitHub trending.
// lang go/c/c++/java/c#/rust...
// spokenLang [zh/en/de/fr...] empty means any
// dataRange daily/weekly/monthly
func TrendingRepos(lang string, dateRange string, spokenLang string) ([]*Repository, error) {
	url := fmt.Sprintf("%s/trending/%s?spoken_language_code=%s&since=%s", GitHubURL, lang, spokenLang, dateRange)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	repos := make([]*Repository, 0, 10)
	doc.Find(".Box .Box-row").Each(func(i int, s *goquery.Selection) {
		repo := &Repository{}

		// author name link
		titleSel := s.Find("h1 a")
		repo.Author = strings.Trim(titleSel.Find("span").Text(), "/\n ")
		repo.Name = strings.TrimSpace(titleSel.Contents().Last().Text())
		relativeLink, _ := titleSel.Attr("href")
		if len(relativeLink) > 0 {
			repo.Link = GitHubURL + relativeLink
		}

		// desc
		repo.Desc = strings.TrimSpace(s.Find("p").Text())

		var langIdx, addIdx, builtByIdx int
		spanSel := s.Find("div>span")
		if spanSel.Size() == 2 {
			// language not exist
			langIdx = -1
			addIdx = 1
		} else {
			builtByIdx = 1
			addIdx = 2
		}

		// language
		if langIdx >= 0 {
			repo.Lang = strings.TrimSpace(spanSel.Eq(langIdx).Text())
		} else {
			repo.Lang = "unknown"
		}

		// add
		addParts := strings.SplitN(strings.TrimSpace(spanSel.Eq(addIdx).Text()), " ", 2)
		repo.Add, _ = strconv.Atoi(addParts[0])

		// builtby
		spanSel.Eq(builtByIdx).Find("a>img").Each(func(i int, img *goquery.Selection) {
			src, _ := img.Attr("src")
			repo.BuiltBy = append(repo.BuiltBy, src)
		})

		// stars forks
		aSel := s.Find("div>a")
		starStr := strings.TrimSpace(aSel.Eq(-2).Text())
		star, _ := strconv.Atoi(strings.Replace(starStr, ",", "", -1))
		repo.Stars = star
		forkStr := strings.TrimSpace(aSel.Eq(-1).Text())
		fork, _ := strconv.Atoi(strings.Replace(forkStr, ",", "", -1))
		repo.Forks = fork

		repos = append(repos, repo)
	})

	return repos, nil
}
