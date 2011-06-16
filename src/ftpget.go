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

var query_regexp = regexp.MustCompile("ftp://([a-zA-Z0-9\\-\\/\\.\\?_]+)")

func query(app *spectrum.Application, q string) []string {
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	if app.Verbose {
		app.PrintfMsg(response.Status)
	}
	uris := query_regexp.FindAllString(string(body), -1)
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

func choice(app *spectrum.Application, matches []string) string {
	app.PrintfMsg("")
	fmt.Printf("Select a number from the above list (press ENTER to exit GoSpeccy): ")
	in := bufio.NewReader(os.Stdin)

        input, err := in.ReadString('\n')
	if err != nil {
                panic(err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		os.Exit(0)
	}

	id, err := strconv.Atoi(strings.TrimRight(input, "\n"))
	if err != nil {
		panic(err)
	}
	if (id < 0) || (id >= len(matches)) {
		panic(os.NewError("Invalid selection"))
	}

	url := matches[id]
	if app.Verbose {
		app.PrintfMsg("You've selected %s", url)
	}
	return url
}

func get(app *spectrum.Application, url string) string {
	filename := path.Base(url)
	dir := path.Join(spectrum.DefaultUserDir, "zip")
	filePath := path.Join(dir, filename)
	ftpURL := url[6:len(url)]

	if app.Verbose {
		ftp.Log = true
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		panic(err)
	}

	// Already downloaded in the past ?
	if _, err := os.Stat(filePath); err == nil {
		app.PrintfMsg("Not downloading, file " + filePath + " already exists");
		return filePath
	}

	// The actual download
	app.PrintfMsg("Downloading into " + filePath);
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if err = ftp.Get(ftpURL, f); err != nil {
		os.Remove(filePath)
		panic(err)
	}

	return filePath
}

