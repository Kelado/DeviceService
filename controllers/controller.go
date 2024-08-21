package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Kelado/DeviceService/middleware"
	"github.com/Kelado/DeviceService/models"
	repositories "github.com/Kelado/DeviceService/repositories"
	"github.com/go-chi/chi/v5"
)

type DeviceController struct {
	repo repositories.DeviceRepo
}

func NewDeviceController(deviceRepo repositories.DeviceRepo) *DeviceController {
	return &DeviceController{
		repo: deviceRepo,
	}
}

func (c *DeviceController) InitRouter(r *chi.Mux) {
	r.Route("/api/v1/devices", func(r chi.Router) {
		r.Use(middleware.Filters())

		r.Post("/", c.AddDevice)
		r.Get("/{device-id}", c.GetById)
		r.Get("/", c.ListAll)
		r.Put("/{device-id}", c.Update)
		r.Delete("/{device-id}", c.Delete)
	})
}

func (c *DeviceController) AddDevice(w http.ResponseWriter, r *http.Request) {
	var device models.DeviceModel
	json.NewDecoder(r.Body).Decode(&device)

	device.ID = models.GenerateUUID()
	device.CreatedAt = models.GetCurrentFormatedTime()

	if err := validateDevice(&device); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.repo.Add(&device)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, device)
}

func (c *DeviceController) GetById(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "device-id")
	device, err := c.repo.GetById(deviceId)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, device)
}

func (c *DeviceController) ListAll(w http.ResponseWriter, r *http.Request) {
	var devices []models.DeviceModel

	filters := middleware.GetFilterFromCtx(r)
	if len(filters) == 0 {
		allDevices, err := c.repo.ListAll()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		devices = allDevices
	} else {
		for _, f := range filters {
			switch filter := f.(type) {
			case models.BrandFilter:
				filteredDevices, err := c.repo.SearchByBrand(filter.GetValue())
				if err != nil {
					respondError(w, http.StatusInternalServerError, err.Error())
					return
				}
				devices = filteredDevices
			default:
				log.Println("Uknown filter, ignoring it!")
			}
		}
	}
	respondJSON(w, http.StatusOK, devices)
}

func (c *DeviceController) Update(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Update a device")
	deviceId := chi.URLParam(r, "device-id")

	var device models.DeviceModel
	json.NewDecoder(r.Body).Decode(&device)
	device.ID = deviceId

	if err := validateDevice(&device); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.repo.Update(&device)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, device)
}

func (c *DeviceController) Delete(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Delete device")
	deviceId := chi.URLParam(r, "device-id")
	err := c.repo.Delete(deviceId)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondSuccess(w)
}

func validateDevice(d *models.DeviceModel) error {
	if err := validateDeviceName(d.Name); err != nil {
		return err
	}
	if err := validateDeviceBrand(d.Brand); err != nil {
		return err
	}
	return nil
}

func validateDeviceName(name string) error {
	empty := name == ""
	if empty {
		return errors.New("name must not be empty")
	}
	return nil
}

func validateDeviceBrand(brand string) error {
	isNokia := brand == "nokia"
	if isNokia {
		return errors.New("you can't have nokia! This is considered a weapon")
	}
	return nil
}

func respondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func respondError(w http.ResponseWriter, statusCode int, msg string) {
	data := map[string]interface{}{
		"error": msg,
	}
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
