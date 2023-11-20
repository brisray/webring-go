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

func readConfig(config_str []byte) Config {
	var config = Config{}
	unmarshallerr := toml.Unmarshal(config_str, &config)
	if unmarshallerr != nil {
		panic(unmarshallerr)
	}
	return config
}

func isURLDomainTheSame(a, b string) bool {
	a_url, a_err := url.Parse(a)
	b_url, b_err := url.Parse(b)
	if a_err != nil || b_err != nil {
		return false
	}
	return a_url.Host == b_url.Host
}

func findWebsiteIndexInList(websites []Website, url string) (int, error) {
	// find current page in config
	for i, website := range websites {
		if isURLDomainTheSame(website.Url, url) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("website not found in list")
}

func nextOrPrevEndpoint(w http.ResponseWriter, r *http.Request, next bool) {
	// get referrer
	referrer := r.Header.Get("Referer")
	if referrer == "" {
		fmt.Fprintf(w, "no referrer")
		return
	}
	// strip trailing slash
	if referrer[len(referrer)-1] == '/' {
		referrer = referrer[:len(referrer)-1]
	}
	// get config
	configstr, filereaderr := os.ReadFile("webring.toml")
	if filereaderr != nil {
		panic(filereaderr)
	}
	config := readConfig(configstr)

	// if referrer is "{config.Root}/home", redirect to first page
	if referrer == config.Root+"/home" {
		// redirect to first/last page
		if next {
			http.Redirect(w, r, config.Websites[0].Url, http.StatusFound)
		} else {
			http.Redirect(w, r, config.Websites[len(config.Websites)-1].Url, http.StatusFound)
		}
		return
	}

	// find current page in config
	index, finderr := findWebsiteIndexInList(config.Websites, referrer)
	if finderr != nil {
		fmt.Fprintf(w, "the site you came from is not in the webring!")
		return
	}
	// find next page in config
	var nextindex int
	if next {
		nextindex = (index + 1) % len(config.Websites)
	} else {
		nextindex = (index - 1) % len(config.Websites)
	}
	for nextindex < 0 {
		nextindex += len(config.Websites)
	}
	nextpage := config.Websites[nextindex]
	// redirect to next page
	http.Redirect(w, r, nextpage.Url, http.StatusFound)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		configstr, filereaderr := os.ReadFile("webring.toml")
		if filereaderr != nil {
			panic(filereaderr)
		}
		config := readConfig(configstr)

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
		config := readConfig(configstr)

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
		nextOrPrevEndpoint(w, r, true)
	})
	http.HandleFunc("/previous", func(w http.ResponseWriter, r *http.Request) {
		nextOrPrevEndpoint(w, r, false)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
