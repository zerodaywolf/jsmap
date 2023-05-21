package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var httpClient = &http.Client{}

type MapData struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
}

func ValidateMapURL(siteURL string) ([]string, bool) {
	// Check if the URL is likely a JavaScript file
	if !(strings.HasSuffix(siteURL, ".js") || strings.Contains(siteURL, ".js?")) {
		return nil, false
	}
    
    // If URL ends with ".js?", strip everything after "?" and add ".map" to it
    if strings.Contains(siteURL, ".js?") {
        siteURL = strings.Split(siteURL, "?")[0] + ".map"
    } else {
        siteURL += ".map"
    }

	mapURL := siteURL
	validatedURLs := make([]string, 0)

	request, err := http.NewRequest("GET", mapURL, nil)
	if err != nil {
		return validatedURLs, false
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109) Gecko/20100101 Firefox/112.0")

	response, err := httpClient.Do(request)
	if err != nil {
		return validatedURLs, false
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return validatedURLs, false
		}

		if strings.Contains(string(body), "\"mappings\"") {
			validatedURLs = append(validatedURLs, mapURL)
		}
	}

	return validatedURLs, true
}

func ExtractMap(mapURL string) (data MapData, err error) {
	request, err := http.NewRequest("GET", mapURL, nil)
	if err != nil {
		return
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109) Gecko/20100101 Firefox/112.0")

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &data)
	return
}

func WriteToFile(path string, content string) error {
	path = filepath.Clean(path)

	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), 0700)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(path, []byte(content), 0600)
}

func main() {
	inputFile := flag.String("f", "", "File with URLs")
	outputDir := flag.String("o", "", "Directory to output files")
	flag.Parse()

	var reader io.Reader
	if *inputFile == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s", file)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		siteURL := scanner.Text()
		validatedURLs, isValid := ValidateMapURL(siteURL)
		if !isValid {
			continue
		}

		for _, mapURL := range validatedURLs {
			mapData, err := ExtractMap(mapURL)
			if err != nil {
				fmt.Println("Error in extracting map: ", err)
				continue
			}

			for i, sourceURL := range mapData.Sources {
				sourceContent := mapData.SourcesContent[i]

				path := sourceURL
				if strings.HasPrefix(path, "webpack:///") {
					path = strings.TrimPrefix(path, "webpack:///")
				}
				path = strings.ReplaceAll(path, "/", string(os.PathSeparator))

				siteDirName := strings.ReplaceAll(siteURL, "/", "_")
				if *outputDir != "" {
					path = filepath.Join(*outputDir, siteDirName, path)
				} else {
					path = filepath.Join(siteDirName, path)
				}

				err = WriteToFile(path, sourceContent)
				if err != nil {
					fmt.Println("Error in writing file: ", err)
				}
			}
		}
	}
}

