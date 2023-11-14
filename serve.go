package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// from the toml
// name = "alifee's webring"
// description = "For the glory of the Indieweb!"

// [[websites]]
// name = "localhost"
// url = "https://localhost:8080"
// description = "The root page"

// [[websites]]
// name = "a deeper page"
// url = "https://localhost:8080/deeper"
// description = "One page deep"

type Website struct {
	Name        string
	Url         string
	Description string
}

type Config struct {
	Name        string
	Description string
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

func findWebsiteIndexInList(websites []Website, url string) (int, error) {
	// strip trailing slash
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	// find current page in config
	for i, website := range websites {
		if website.Url == url {
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
	// config := readConfig()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/webring.html", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "webring.html should be here")
	})
	http.HandleFunc("/webring.js", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "webring.js should be here")
	})
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		nextOrPrevEndpoint(w, r, true)
	})
	http.HandleFunc("/previous", func(w http.ResponseWriter, r *http.Request) {
		nextOrPrevEndpoint(w, r, false)
	})

	// debug
	http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
		// print request headers back to client
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Fprintf(w, "%v: %v\n", name, value)
			}
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
