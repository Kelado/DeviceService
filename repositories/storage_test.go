package repositories

import (
	"log"
	"reflect"
	"testing"

	"github.com/Kelado/DeviceService/models"
)

const (
	defaultDeviceName  = "p40"
	defaultDeviceBrand = "huawei"
)

func TestSQLiteAddDevice(t *testing.T) {
	deviceRepo := NewSQLiteDeviceRepo(&SQLiteRepoConfig{DSN: ":memory:"})

	expectedDevice := createDevice()

	err := deviceRepo.Add(expectedDevice)
	if err != nil {
		t.Errorf("Add device error: %v", err)
	}

	insertedDevice, err := deviceRepo.GetById(expectedDevice.ID)
	if err != nil {
		t.Errorf("Get device error: %v", err)
	}

	if !reflect.DeepEqual(insertedDevice, expectedDevice) {
		t.Errorf("Expected device \n\t(exp):\t%+v, \n\t(got):\t%+v", expectedDevice, insertedDevice)
	}
	log.Println(insertedDevice)
}

func TestSQLiteListDevices(t *testing.T) {
	deviceRepo := NewSQLiteDeviceRepo(&SQLiteRepoConfig{DSN: ":memory:"})

	numOfDevices := 4
	for i := 0; i < numOfDevices; i++ {
		d := createDevice()
		deviceRepo.Add(d)
	}

	devices, err := deviceRepo.ListAll()
	if err != nil {
		t.Errorf("List devices error: %v", err)
	}

	if len(devices) != numOfDevices {
		t.Errorf("Wrong numb er of devices in db\n\t(exp):\t%+v, \n\t(got):\t%+v", numOfDevices, len(devices))
	}
}

func TestSQLiteUpdateDevice(t *testing.T) {
	deviceRepo := NewSQLiteDeviceRepo(&SQLiteRepoConfig{DSN: ":memory:"})

	expectedDevice := createDevice()
	err := deviceRepo.Add(expectedDevice)
	if err != nil {
		t.Errorf("Add device error: %v", err)
	}

	expectedDevice.Name = "poco"
	expectedDevice.Brand = "Xiaomi"

	err = deviceRepo.Update(expectedDevice)
	if err != nil {
		t.Errorf("Update device error: %v", err)
	}

	updatedDevice, err := deviceRepo.GetById(expectedDevice.ID)
	if err != nil {
		t.Errorf("Get device error: %v", err)
	}

	if !reflect.DeepEqual(updatedDevice, expectedDevice) {
		t.Errorf("Expected device \n\t(exp):\t%+v, \n\t(got):\t%+v", expectedDevice, updatedDevice)
	}
}

func TestSQLiteDeleteDevice(t *testing.T) {
	deviceRepo := NewSQLiteDeviceRepo(&SQLiteRepoConfig{DSN: ":memory:"})

	expectedDevice := createDevice()

	err := deviceRepo.Add(expectedDevice)
	if err != nil {
		t.Errorf("Add device error: %v", err)
	}

	err = deviceRepo.Delete(expectedDevice.ID)
	if err != nil {
		t.Errorf("Delete device error: %v", err)
	}

	_, err = deviceRepo.GetById(expectedDevice.ID)
	if err == nil {
		t.Errorf("Device was not deleted from db")
	}
}

func TestSQLiteDSearchDevicesByBrand(t *testing.T) {
	deviceRepo := NewSQLiteDeviceRepo(&SQLiteRepoConfig{DSN: ":memory:"})

	numOfDefaultDevices := 4
	for i := 0; i < numOfDefaultDevices; i++ {
		d := createDevice()
		deviceRepo.Add(d)
	}

	numOfIphoneDevices := 2
	for i := 0; i < numOfIphoneDevices; i++ {
		id := models.GenerateUUID()
		d := &models.DeviceModel{
			ID:        id,
			Name:      "14",
			Brand:     "iphone",
			CreatedAt: models.GetCurrentFormatedTime(),
		}
		deviceRepo.Add(d)
	}

	devices, err := deviceRepo.SearchByBrand(defaultDeviceBrand)
	if err != nil {
		t.Error("Search error: ", err)
	}

	if len(devices) != numOfDefaultDevices {
		t.Errorf("Expected number of devices searched \n\t(exp):\t%+v, \n\t(got):\t%+v", numOfDefaultDevices, len(devices))
	}
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
