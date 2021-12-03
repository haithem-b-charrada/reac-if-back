package repositories

import (
	"inventory/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	Storage *gorm.DB
}

func (r *PostRepository) Save(drone *models.Post) {
	r.Storage.Save(drone)
}

func (r *PostRepository) Find(drone *models.Post) {
	r.Storage.First(drone)
}

func (r *PostRepository) FindById(drone *models.Post) {
	r.Storage.Where("id = ?", drone.ID).First(drone)
}

func (r *PostRepository) All(drone *[]models.Post) {
	r.Storage.Find(drone)
}

func (r *PostRepository) Remove(drone *models.Post) {
	r.Storage.Delete(&drone)
}
