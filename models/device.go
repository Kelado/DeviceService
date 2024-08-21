package models

import (
	"net/http"
)

type DeviceModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	CreatedAt string
}

func (dm DeviceModel) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
