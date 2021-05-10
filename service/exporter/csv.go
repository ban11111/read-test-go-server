package exporter

import (
	"bytes"
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"mime"
)

type CsvExporter struct {
}

func (e *CsvExporter) Ext() string {
	return "csv"
}

func (e *CsvExporter) Export(data interface{}) (*ReaderWrapper, error) {
	file := bytes.NewBuffer(nil)
	if err := gocsv.MarshalCSV(data, gocsv.NewSafeCSVWriter(csv.NewWriter(file))); err != nil {
		return nil, err
	}
	contentType := mime.TypeByExtension(".csv")
	if contentType == "" {
		contentType = "application/vnd.ms-excel"
	}
	contentType += ";charset=utf-8"
	return &ReaderWrapper{file, file.Len(), contentType}, nil
}
