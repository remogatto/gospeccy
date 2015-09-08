package spectrum

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	fileSuffixes = []string{".tap.zip", ".sna.zip", ".z80.zip"}
	baseURL      = "http://www.worldofspectrum.org/"    
	queryBaseURL = baseURL + "infoseek.cgi?"
)

func knownType(filename string) bool {
	for _, suffix := range fileSuffixes {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}

// Information about a game/demo/utility located at [www.worldofspectrum.org]
type WosRecord struct {
	Title       string
	MachineType string   // "ZX Spectrum 48K", "ZX Spectrum 48K/128K", ...
	Publication string   // "Freeware", "Commercial", "unknown", ...
	FtpFiles    []string // formerly ftp, now http://www.worldofspectrum.org/...
	Score       string   // "7.59 (24 votes)", "No votes yet"
}

// Query [www.worldofspectrum.org] for matching files
func WosQuery(app *Application, query string) ([]WosRecord, error) {
	client := new(http.Client)
	query = queryBaseURL + query
	if app.Verbose {
		app.PrintfMsg("WOS query: %s", query)
		app.PrintfMsg("Fetching...")
	}
	response, err := client.Get(query)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	if app.Verbose {
		app.PrintfMsg("Response status: %s", response.Status)
	}

	var matches []WosRecord
	{
		regexp_title := regexp.MustCompile(">Full title<")
		regexp1 := regexp.MustCompile("<A[^>]*>" + "([^<]+)" + "</A>")
		regexp2 := regexp.MustCompile("<FONT[^>]*>(<I>)?" + "([^<]+)")
		regexp_score := regexp.MustCompile(("<FONT[^>]*>" + "([^<]+)" + "<" + "[^<]+") + ("<I[^>]*>" + "([^<]+)" + "</I>"))
		regexp_file := regexp.MustCompile("/pub/sinclair/[a-zA-Z0-9\\-\\/\\.\\?_]+")

		var a [][]int = regexp_title.FindAllIndex(body, -1)
		for i := 0; i < len(a); i++ {
			var startIndex, endIndex int
			startIndex = a[i][0]
			if i+1 < len(a) {
				endIndex = a[i+1][0]
			} else {
				endIndex = len(body)
			}

			body2 := body[startIndex:endIndex]

			// Full title
			var title string
			{
				title = ""

				var b2 []int = regexp1.FindSubmatchIndex(body2)

				if 3 < len(b2) {
					startIndex3 := b2[2]
					endIndex3 := b2[3]

					title = string(body2[startIndex3:endIndex3])
					title = strings.TrimSpace(title)
				}
			}

			// Machine type
			var machineType string
			{
				machineType = "unknown"

				startIndex2 := strings.Index(string(body2), ">Machine type<")
				if startIndex2 != -1 {
					var b2 []int = regexp2.FindSubmatchIndex(body2[startIndex2:])

					if 5 < len(b2) {
						startIndex3 := b2[4]
						endIndex3 := b2[5]

						machineType = string(body2[startIndex2:][startIndex3:endIndex3])
						machineType = strings.TrimSpace(machineType)
					}
				}
			}

			// Publication (freeware, commercial, unknown, ...)
			var publication string
			{
				publication = "unknown"

				startIndex2 := strings.Index(string(body2), ">Original publication<")
				if startIndex2 != -1 {
					var b2 []int = regexp2.FindSubmatchIndex(body2[startIndex2:])

					if 5 < len(b2) {
						startIndex3 := b2[4]
						endIndex3 := b2[5]

						publication = string(body2[startIndex2:][startIndex3:endIndex3])
						publication = strings.TrimSpace(publication)
					}
				}
			}

			// Score
			var score string
			{
				score = ""

				startIndex2 := strings.Index(string(body2), ">Score<")
				if startIndex2 != -1 {
					var b2 []int = regexp_score.FindSubmatchIndex(body2[startIndex2:])

					if 5 < len(b2) {
						startIndex3 := b2[2]
						endIndex3 := b2[3]
						startIndex4 := b2[4]
						endIndex4 := b2[5]

						// 'score1' is the score, 'score2' is the rate of confidence
						score1 := string(body2[startIndex2:][startIndex3:endIndex3])
						score2 := string(body2[startIndex2:][startIndex4:endIndex4])
						score1 = strings.TrimSpace(score1)
						score2 = strings.TrimSpace(score2)
						if len(score1) > 0 {
							score = score1 + " " + score2
						} else {
							score = score2
						}
					}
				}
			}

			var fileUrls []string
			{
				urls := regexp_file.FindAllString(string(body2), -1)
				for _, url := range urls {
					if knownType(path.Base(url)) {
						fileUrls = append(fileUrls, url)
					}
				}
			}

			matches = append(matches, WosRecord{title, machineType, publication, fileUrls, score})
		}
	}

	return matches, nil
}


// Fetch the specified URL and store to the given open file
func httpGet(url string, outfile io.Writer) error {
   	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	
	_, err = io.Copy(outfile, resp.Body)
	return err
}


// Download from [ftp.worldofspectrum.org].
// An URL can be obtained by calling function WosQuery.
func WosGet(app *Application, stdout io.Writer, url string) (string, error) {
	filename := path.Base(url)
	dir := DownloadPath()
	filePath := path.Join(dir, filename)
	httpURL := baseURL + url

	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}

	// Already downloaded in the past ?
	if _, err := os.Stat(filePath); err == nil {
		fmt.Fprintf(stdout, "Not downloading, file %s already exists\n", filePath)
		return filePath, nil
	}

	// The actual download:

	fmt.Fprintf(stdout, "Downloading into %s\n", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	if err = httpGet(httpURL, f); err != nil {
		os.Remove(filePath)
		return "", err
	}

	err = f.Close()
	if err != nil {
		os.Remove(filePath)
		return "", err
	}

	return filePath, nil
}
