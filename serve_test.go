package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// ** this section **
// from https://medium.com/@benbjohnson/structuring-tests-in-go-46ddee7a25c
// https://github.com/benbjohnson/testing

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func assertFalse(tb testing.TB, condition bool, msg string, v ...interface{}) {
	assert(tb, !condition, msg, v...)
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// ** this section **

const workingConfig = `name = "alifee's webring"
description = "For the glory of the Indieweb!"

[[websites]]
name = "Random URL"
url = "https://randomurl.com"
description = "Random URL"

[[websites]]
name = "a page about dwarves"
url = "https://digging.allaboutdwarves.co.uk/home"
description = "Dig! Dig! Dig!"
`

func TestReadConfig(t *testing.T) {
	expected_config := Config{
		Name:        "alifee's webring",
		Description: "For the glory of the Indieweb!",
		Websites: []Website{
			{
				Name:        "Random URL",
				Url:         "https://randomurl.com",
				Description: "Random URL",
			},
			{
				Name:        "a page about dwarves",
				Url:         "https://digging.allaboutdwarves.co.uk/home",
				Description: "Dig! Dig! Dig!",
			},
		},
	}

	config_str := []byte(workingConfig)
	config, _ := readConfig(config_str)
	assert(t, reflect.DeepEqual(config, expected_config), "readConfig() = %v, want %v", config, expected_config)
}

func TestReadRealConfig(t *testing.T) {
	config_str, err := os.ReadFile("webring.toml")
	ok(t, err)
	_, err = readConfig(config_str)
	ok(t, err)
}

func Test_findWebsiteInWebring(t *testing.T) {
	config_str := []byte(workingConfig)
	config, _ := readConfig(config_str)
	// TEST 1: url is in list
	index := findWebsiteIndexInList(config.Websites, "https://randomurl.com")
	assert(t, index == 0, "findWebsiteInWebring() = %v, want %v", index, 0)
	// TEST 2: url has trailing slash
	index = findWebsiteIndexInList(config.Websites, "https://randomurl.com/")
	assert(t, index == 0, "findWebsiteInWebring() = %v, want %v", index, 0)
	// TEST 3: url not in list
	index = findWebsiteIndexInList(config.Websites, "https://notfound.com")
	assert(t, index == -1, "findWebsiteInWebring() = %v, want %v", index, -1)
}

func Test_isURLDomainTheSame(t *testing.T) {
	shouldBeSame := isURLDomainTheSame("https://example.com", "https://example.com")
	assert(t, shouldBeSame, "compareURLDomain() = %v, want %v", false, true)

	shouldBeSame = isURLDomainTheSame("https://example.com/sub-page/sub-sub-page", "https://example.com")
	assert(t, shouldBeSame, "compareURLDomain() = %v, want %v", false, true)

	shouldBeSame = isURLDomainTheSame("https://example.com/sub-page/sub-sub-page", "https://example.com/thing#header?query=string")
	assert(t, shouldBeSame, "compareURLDomain() = %v, want %v", false, true)

	shouldBeSame = isURLDomainTheSame("https://example.com/page1/", "https://example.com/secondpage")
	assert(t, shouldBeSame, "compareURLDomain() = %v, want %v", false, true)

	shouldBeDifferent := isURLDomainTheSame("https://example.com", "https://example.org")
	assertFalse(t, shouldBeDifferent, "compareURLDomain() = %v, want %v", true, false)
}
