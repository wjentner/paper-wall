package main

import (
	"go-paper-teaser/oop"
	"go-paper-teaser/pdfs"
	"strconv"
)

func main() {
	year := 2019

	prefix := "./papers-"

	outputFolder := prefix + strconv.Itoa(year) + "/"

	pdfFile := prefix + strconv.Itoa(year) + ".pdf"

	d := oop.New("https://bib.dbvis.de/publications/index/home/filterByYear=", year, outputFolder)
	d.DownloadPapers()
	//
	p := pdfs.New(outputFolder, pdfFile)
	p.SplitAndMerge()
}
