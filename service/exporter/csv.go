package exporter

import "github.com/gocarina/gocsv"

type CsvExporter struct {

}

func (e *CsvExporter)Export() {
	gocsv.MarshalCSV()
}