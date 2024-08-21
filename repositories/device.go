package repositories

import (
	"github.com/Kelado/DeviceService/models"
)

type DeviceRepo interface {
	Add(*models.DeviceModel) error
	GetById(string) (*models.DeviceModel, error)
	ListAll() ([]models.DeviceModel, error)
	Update(*models.DeviceModel) error
	Delete(string) error
	SearchByBrand(string) ([]models.DeviceModel, error)
}
