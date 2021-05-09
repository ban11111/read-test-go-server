package exporter

import "read-test-server/model"

type ExportUser struct {
	model.User
	Score string `json:"score"`
}