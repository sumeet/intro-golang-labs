package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var DB *sql.DB = setupGlobalDatabase()

func setupGlobalDatabase() *sql.DB {
	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}
	return db
}

type User struct {
	ID       int
	Name     string
	Language string
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		err2 := fmt.Errorf("Couldn't set up database: %w", err)
		return nil, err2
	}

	sqlStmt := `
	create table users (id integer not null primary key AUTOINCREMENT, name text, language text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		err2 := fmt.Errorf("Couldn't create table: %w", err)
		return nil, err2
	}

	_, err = db.Exec("insert into users(name, language) values('jason', 'ruby'), ('sue', 'python')")
	if err != nil {
		panic(err)
	}

	return db, nil
}

func main() {
	http.HandleFunc("/", Handler)
	fmt.Println("Listening on port http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func Handler(w http.ResponseWriter, r *http.Request) {
	urls := ScrapeURLs("https://slashdot.org/")

	urlsByHostname := map[string][]string{}
	for _, u := range urls {
		parsed, err := url.Parse(u)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't parse URL: %s (%w), skipping...\n", u, err)
		}
		if parsed.Path == "" {
			continue
		}
		urlsByHostname[parsed.Hostname()] = append(urlsByHostname[parsed.Hostname()], parsed.Path)
	}

	for k, v := range urlsByHostname {
		urlsByHostname[k] = removeDuplicate(v)
	}

	fmt.Fprintf(w, `{"urls": `)
	if err := json.NewEncoder(w).Encode(urlsByHostname); err != nil {
		panic(err)
	}
	fmt.Fprintf(w, `}`)
}

func ScrapeURLs(entrypoint string) []string {
	var results []string

	// Request the HTML page.
	res, err := http.Get(entrypoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get URL %s: %w", entrypoint, err)
		return results
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "status code error: %d %s\n", res.StatusCode, res.Status)
		return results
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse doc: %s", err)
		return results
	}

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		url := hrefToURL(entrypoint, href)
		if url != "" {
			results = append(results, url)
		}
	})
	return results
}

func hrefToURL(origin, href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	} else if strings.HasPrefix(href, "//") {
		return "https:" + href
	} else if strings.HasPrefix(href, "/") {
		originURL, err := url.Parse(origin)
		if err != nil {
			panic(fmt.Errorf("origin isn't a valid URL: %w", err))
		}
		originURL.Path = href
		return originURL.String()
	} else if strings.HasPrefix(href, "#") {
		return ""
	} else {
		fmt.Fprintf(os.Stderr, "don't know how to deal with URL %#v, ignoring...\n", href)
		return ""
	}
}
