package pdfs

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	//unicommon "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

type pdfs struct {
	inputFolder string
	outputFile  string
}

func New(inputFolder string, outputFile string) pdfs {
	if !strings.HasSuffix(inputFolder, "/") {
		inputFolder = inputFolder + "/"
	}

	pdfs := pdfs{inputFolder, outputFile}
	return pdfs
}

func (pdfs pdfs) SplitAndMerge() error {
	files, err := ioutil.ReadDir(pdfs.inputFolder)
	if err != nil {
		log.Fatal(err)
	}

	pdfWriter := pdf.NewPdfWriter()

	for _, file := range files {

		page, err := getFirstPage(pdfs.inputFolder + file.Name())
		if err != nil {
			return err
		}

		err2 := pdfWriter.AddPage(page)
		if err2 != nil {
			return err2
		}
	}

	err3 := pdfs.writeMergedPdf(pdfWriter)
	if err3 != nil {
		return err3
	}

	return nil
}

func getFirstPage(file string) (*pdf.PdfPage, error) {
	pageIdx := 1

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return nil, err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	if numPages < pageIdx {
		return nil, err
	}

	page, err := pdfReader.GetPage(pageIdx)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (pdfs pdfs) writeMergedPdf(pdfWriter pdf.PdfWriter) error {
	fWrite, err := os.Create(pdfs.outputFile)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
