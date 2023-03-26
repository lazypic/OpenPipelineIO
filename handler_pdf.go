package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func handlerAPIPdfToJson(w http.ResponseWriter, r *http.Request) {
	// 파일 읽기
	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 파일 데이터 읽기
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// unidoc model에서는 io.Reader 타입이 아니라. io.ReadSeeker 타입이 필요하다 ReadSeeker를 만든다.
	readSeekerData := make([]byte, len(data))
	copy(readSeekerData, data)
	dataSeeker := bytes.NewReader(readSeekerData)

	pdfReader, err := model.NewPdfReader(dataSeeker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ex, err := extractor.New(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		text, err := ex.ExtractText()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("------------------------------")
		fmt.Printf("Page %d:\n", pageNum)
		fmt.Printf("\"%s\"\n", text)
		fmt.Println("------------------------------")
	}
}
