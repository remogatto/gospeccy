package main

import (
	"fmt"
	"bufio"
	"path"
	"regexp"
	"http"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"spectrum"
	"ftp"
)

var (
	patterns = []string{"*.tap.zip", "*.sna.zip", "*.z80.zip"}
	queryBaseURL = "http://www.worldofspectrum.org/infoseek.cgi?regexp="
)

func knownType(filename string) bool {
	for _, t := range patterns {
		matched, err := path.Match(t, filename)
		if err != nil {
			panic(err)
		}
		if matched {
			return true
		}
	}
	return false
}

func query(app *spectrum.Application, q string) []string {
	re, _ := regexp.Compile("ftp://([a-zA-Z0-9\\-\\/\\.\\?_]+)")
	client := new(http.Client)
	query :=  queryBaseURL + q
	if app.Verbose {
		app.PrintfMsg("Query string:", query)
		app.PrintfMsg("Fetching...")
	}
	response, err := client.Get(query)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	if app.Verbose {
		app.PrintfMsg(response.Status)
	}
	uris := re.FindAllString(string(body), -1)
	matches := make([]string, 0)
	for _, uri := range uris {
		if knownType(path.Base(uri)) {
			matches = append(matches, uri)
			app.PrintfMsg("[%d] - %s", len(matches) - 1, uri)
		}
	}
	app.PrintfMsg("Found %d matching resources", len(matches))
	return matches
}

func choice(app *spectrum.Application, matches []string) (url string) {
	fmt.Print("Please select a number (press ENTER to exit): ")
	in := bufio.NewReader(os.Stdin)
        if input, err := in.ReadString('\n'); err != nil {
                panic(err)
	} else if input == "\n" {
		os.Exit(0)
	} else {
		if id, err := strconv.Atoi(strings.TrimRight(input, "\n")); err != nil {
			panic(err)
		} else {
			if app.Verbose {
				app.PrintfMsg("You've selected %s", matches[id])
			}
			url = matches[id]
		}
	}
	return url
}

func get(app *spectrum.Application, url string) string {
	filename := path.Base(url)
	userDirPath := path.Join(spectrum.DefaultUserDir, "zip", filename)
	ftpURL := url[6:len(url)]
	if app.Verbose {
		ftp.Log = true
	}
	if f, err := os.Create(userDirPath); err != nil {
		panic(err)
	} else {
		if err = ftp.Get(ftpURL, f); err != nil {
			panic(err)
		}
	}
	return userDirPath
}

