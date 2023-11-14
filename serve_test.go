package main

import (
	"reflect"
	"testing"
)

const workingConfig = `name = "alifee's webring"
description = "For the glory of the Indieweb!"

[[websites]]
name = "localhost"
url = "https://localhost:8080"
description = "The root page"

[[websites]]
name = "a deeper page"
url = "https://localhost:8080/deeper"
description = "One page deep"
`

// test readConfig function with a string
func TestReadConfig(t *testing.T) {
	config_str := []byte(workingConfig)
	config := readConfig(config_str)
	config_obj := Config{
		Name:        "alifee's webring",
		Description: "For the glory of the Indieweb!",
		Websites: []Website{
			{
				Name:        "localhost",
				Url:         "https://localhost:8080",
				Description: "The root page",
			},
			{
				Name:        "a deeper page",
				Url:         "https://localhost:8080/deeper",
				Description: "One page deep",
			},
		},
	}
	if !reflect.DeepEqual(config, config_obj) {
		t.Errorf("readConfig() = %v, want %v", config, config_obj)
	}
}

// test findWebsiteIndexInList function with a string
func TestFindWebsiteInWebring(t *testing.T) {
	config_str := []byte(workingConfig)
	config := readConfig(config_str)
	index, _ := findWebsiteIndexInList(config.Websites, "https://localhost:8080")
	if index != 0 {
		t.Errorf("findWebsiteInWebring() = %v, want %v", index, 0)
	}
}

// test findWebsiteIndexInList with a trailing /
func TestFindWebsiteInWebringTrailingSlash(t *testing.T) {
	config_str := []byte(workingConfig)
	config := readConfig(config_str)
	index, _ := findWebsiteIndexInList(config.Websites, "https://localhost:8080/")
	if index != 0 {
		t.Errorf("findWebsiteInWebring() = %v, want %v", index, 0)
	}
}

// test findWebsiteIndexInList returns an error when the website is not found
func TestFindWebsiteInWebringNotFound(t *testing.T) {
	config_str := []byte(workingConfig)
	config := readConfig(config_str)
	_, err := findWebsiteIndexInList(config.Websites, "https://localhost:8080/notfound")
	if err == nil {
		t.Errorf("findWebsiteInWebring() = %v, want %v", err, "error")
	}
}
