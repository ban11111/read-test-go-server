package exporter

import "io"

type ReaderWrapper struct {
	io.Reader
	Len int
	ContentType string
}

type Exporter interface {
	Ext() string
	Export(data interface{}) (*ReaderWrapper, error)
}

var ImplementedExporters map[string]Exporter
var ImplementedExportDataGetters map[string]ExportDataGetter

func init() {
	ImplementedExporters = make(map[string]Exporter)
	ImplementedExportDataGetters = make(map[string]ExportDataGetter)
}

func RegisterExporter(exporters ...Exporter) {
	for _, x := range exporters {
		ImplementedExporters[x.Ext()] = x
	}
}

func RegisterExportDataGetter(getters ...ExportDataGetter) {
	for _, x := range getters {
		ImplementedExportDataGetters[x.Table()] = x
	}
}

type GetterCtx interface {
	GetIds() []uint
	GetPaperId() uint
}

type ExportDataGetter interface {
	Table() string
	Getter(ctx GetterCtx) (data interface{}, err error)
}