package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kelado/DeviceService/middleware"
	"github.com/Kelado/DeviceService/models"
	"github.com/Kelado/DeviceService/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

const (
	Addr               = ":8000"
	defaultDeviceName  = "p40"
	defaultDeviceBrand = "huawei"
)

func TestE2E(t *testing.T) {
	router := initDeviceService()

	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()

	device := createDevice()

	requestBody, _ := json.Marshal(device)

	var createdDevice models.DeviceModel

	t.Run("TestAddDevice", func(t *testing.T) {
		resp, err := client.Post(server.URL+"/api/v1/devices", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		json.NewDecoder(resp.Body).Decode(&createdDevice)

		assert.Equal(t, device.Name, createdDevice.Name)
		assert.Equal(t, device.Brand, createdDevice.Brand)
		assert.NotEmpty(t, createdDevice.ID)
		assert.NotZero(t, createdDevice.CreatedAt)

	})

	t.Run("TestGetById", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/devices/" + createdDevice.ID)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var retrievedDevice models.DeviceModel
		json.NewDecoder(resp.Body).Decode(&retrievedDevice)

		assert.Equal(t, createdDevice.ID, retrievedDevice.ID)
		assert.Equal(t, createdDevice.Name, retrievedDevice.Name)
		assert.Equal(t, createdDevice.Brand, retrievedDevice.Brand)
	})

	t.Run("TestListAll", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/devices")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var devices []models.DeviceModel
		json.NewDecoder(resp.Body).Decode(&devices)

		assert.GreaterOrEqual(t, len(devices), 1)
	})

	t.Run("TestSearchByBrand", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/api/v1/devices?s=brand:" + defaultDeviceBrand)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var devices []models.DeviceModel
		json.NewDecoder(resp.Body).Decode(&devices)

		assert.GreaterOrEqual(t, len(devices), 1)
	})

	t.Run("TestUpdateDevice", func(t *testing.T) {
		updatedDevice := models.DeviceModel{
			Name:  "p40 Pro",
			Brand: "huawei",
		}
		updateBody, _ := json.Marshal(updatedDevice)

		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/devices/"+createdDevice.ID, bytes.NewBuffer(updateBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedDeviceResponse models.DeviceModel
		json.NewDecoder(resp.Body).Decode(&updatedDeviceResponse)

		assert.Equal(t, createdDevice.ID, updatedDeviceResponse.ID)
		assert.Equal(t, updatedDevice.Name, updatedDeviceResponse.Name)
		assert.Equal(t, updatedDevice.Brand, updatedDeviceResponse.Brand)
	})

	t.Run("TestDeleteDevice", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/devices/"+createdDevice.ID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		getResp, err := client.Get(server.URL + "/api/v1/devices/" + createdDevice.ID)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer getResp.Body.Close()

		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

}

func initDeviceService() *chi.Mux {
	deviceRepo := repositories.NewSQLiteDeviceRepo(&repositories.SQLiteRepoConfig{DSN: ":memory:"})
	c := NewDeviceController(deviceRepo)
	router := chi.NewRouter()

	router.Route("/api/v1/devices", func(r chi.Router) {
		r.Use(middleware.Filters())

		r.Post("/", c.AddDevice)
		r.Get("/{device-id}", c.GetById)
		r.Get("/", c.ListAll)
		r.Put("/{device-id}", c.Update)
		r.Delete("/{device-id}", c.Delete)
	})

	return router
}

func createDevice() *models.DeviceModel {
	id := models.GenerateUUID()
	return &models.DeviceModel{
		ID:        id,
		Name:      defaultDeviceName,
		Brand:     defaultDeviceBrand,
		CreatedAt: models.GetCurrentFormatedTime(),
	}
}
