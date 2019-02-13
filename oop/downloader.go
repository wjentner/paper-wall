package oop

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type downloader struct {
	url          string
	year         int
	outputFolder string
}

func New(url string, year int, outputFolder string) downloader {
	d := downloader{url, year, outputFolder}
	return d
}

// This will get called for each HTML element found
func (d downloader) processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	if exists && strings.HasSuffix(href, ".pdf") {
		fmt.Println(href)

		splitres := strings.Split(href, "/")
		file := fmt.Sprintf("%03d", index) + "-" + splitres[len(splitres)-1]

		defer d.downloadFile(file, href)
	}
}

func (d downloader) DownloadPapers() {
	// remove

	// remove folder if it exists
	if _, err := os.Stat(d.outputFolder); os.IsExist(err) {
		os.Remove(d.outputFolder)
	}

	// create the folder again
	os.Mkdir(d.outputFolder, 0777)

	// Make HTTP request
	response, err := http.Get(d.url + strconv.Itoa(d.year))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find all links and process them with the function
	// defined earlier
	document.Find("a").Each(d.processElement)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func (d downloader) downloadFile(filename string, url string) error {

	// Create the file
	out, err := os.Create(d.outputFolder + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
