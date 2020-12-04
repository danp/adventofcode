package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	now := time.Now()
	year := flag.Int("year", now.Year(), "year to work in")
	day := flag.Int("day", now.Day(), "day to work on")
	flag.Parse()

	base := os.Getenv("ADVENT_OF_CODE_BASE")
	if base == "" {
		base = filepath.Join(build.Default.GOPATH, "src", "github.com", "danp", "adventofcode")
	}

	yearBase := filepath.Join(base, strconv.Itoa(*year))
	dayBase := filepath.Join(yearBase, "cmd", fmt.Sprintf("day%02d", *day))

	if _, err := os.Stat(dayBase); os.IsNotExist(err) {
		log.Println(dayBase, "creating")
		if err := os.MkdirAll(dayBase, 0700); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Println(dayBase, "exists")
	}

	dayMain := filepath.Join(dayBase, "main.go")
	if _, err := os.Stat(dayMain); os.IsNotExist(err) {
		log.Println(dayMain, "creating")
		mainTmpl, err := ioutil.ReadFile(filepath.Join(base, "cmd", "aoc-init", "main.go.tmpl"))
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(filepath.Join(dayMain), mainTmpl, 0600); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Println(dayMain, "exists")
	}

	dayInput := filepath.Join(dayBase, "input")
	session, err := loadSession()
	if err != nil {
		log.Println(dayInput, "unable to load session for fetching:", err)
	} else {
		if _, err := os.Stat(dayInput); os.IsNotExist(err) {
			log.Println(dayInput, "fetching")

			input, err := fetchInput(session, *year, *day)
			if err != nil {
				log.Println(dayInput, "unable to fetch:", err)
			}

			if err := ioutil.WriteFile(filepath.Join(dayInput), input, 0600); err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Fatal(err)
		} else {
			log.Println(dayInput, "exists")
		}
	}

}

func loadSession() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	sb, err := ioutil.ReadFile(filepath.Join(config, "adventofcode", "session"))
	return string(sb), err
}

func fetchInput(session string, year, day int) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if st := resp.StatusCode; st != 200 {
		return nil, fmt.Errorf("bad status %d: %s", st, b)
	}

	return b, nil
}
