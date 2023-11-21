package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/aymerick/raymond"
	"github.com/pelletier/go-toml/v2"
)

type Website struct {
	Name        string
	Url         string
	Description string
}

type Config struct {
	Name        string
	Description string
	Root        string
	Websites    []Website
}

func readConfig(config_str []byte) (Config, error) {
	var config = Config{}
	err := toml.Unmarshal(config_str, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func isURLDomainTheSame(a, b string) bool {
	a_url, a_err := url.Parse(a)
	b_url, b_err := url.Parse(b)
	if a_err != nil || b_err != nil {
		return false
	}
	return a_url.Host == b_url.Host
}

func findWebsiteIndexInList(websites []Website, url string) int {
	// find current page in config
	for i, website := range websites {
		if isURLDomainTheSame(website.Url, url) {
			return i
		}
	}
	return -1
}

func nextOrPrev(w http.ResponseWriter, r *http.Request, nextOrPrev string) {
	// load config
	configstr, err := os.ReadFile("webring.toml")
	if err != nil {
		panic(err)
	}
	config, err := readConfig(configstr)
	if err != nil {
		panic(err)
	}

	// get referer
	referrer := r.Header.Get("Referer")
	if referrer == "" {
		fmt.Fprintf(w, "no referrer")
		return
	}

	// if referer is "{config.Root}", redirect to first or last page
	if isURLDomainTheSame(referrer, config.Root) {
		if nextOrPrev == "next" {
			http.Redirect(w, r, config.Websites[0].Url, http.StatusFound)
		} else {
			http.Redirect(w, r, config.Websites[len(config.Websites)-1].Url, http.StatusFound)
		}
		return
	}

	// find current page in config
	index := findWebsiteIndexInList(config.Websites, referrer)

	// if referer not in config, return error
	if index == -1 {
		fmt.Fprintf(w, "the site you came from is not in the webring!")
		return
	}

	// find next page in config
	var diff int
	if nextOrPrev == "next" {
		diff = 1
	} else {
		diff = -1
	}

	nextindex := (index + len(config.Websites) + diff) % len(config.Websites)
	nextpage := config.Websites[nextindex]

	// redirect to next page
	http.Redirect(w, r, nextpage.Url, http.StatusFound)
}

func main() {
	// handle /static as fileserver
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// load config
		configstr, err := os.ReadFile("webring.toml")
		if err != nil {
			panic(err)
		}
		config, err := readConfig(configstr)
		if err != nil {
			panic(err)
		}

		// render html
		htmltemplate, err := os.ReadFile("templates/homepage.html.template")
		if err != nil {
			panic(err)
		}
		htmlfile, err := raymond.Render(string(htmltemplate), config)
		if err != nil {
			panic(err)
		}

		// send response
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprint(w, htmlfile)
	})

	http.HandleFunc("/webring.js", func(w http.ResponseWriter, r *http.Request) {
		configstr, filereaderr := os.ReadFile("webring.toml")
		if filereaderr != nil {
			panic(filereaderr)
		}
		config, err := readConfig(configstr)
		if err != nil {
			panic(err)
		}

		// render html
		htmltemplate, err := os.ReadFile("templates/webring.html.template")
		if err != nil {
			panic(err)
		}
		htmlfile, err := raymond.Render(string(htmltemplate), config)
		if err != nil {
			panic(err)
		}

		// render js
		jstemplate, err := os.ReadFile("templates/webring.js.template")
		if err != nil {
			panic(err)
		}
		jsfile, err := raymond.Render(string(jstemplate), map[string]string{
			"webring_html": htmlfile,
		})
		if err != nil {
			panic(err)
		}

		// send response
		w.Header().Add("Content-Type", "application/javascript")
		fmt.Fprint(w, jsfile)
	})

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		nextOrPrev(w, r, "next")
	})
	http.HandleFunc("/previous", func(w http.ResponseWriter, r *http.Request) {
		nextOrPrev(w, r, "previous")
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
