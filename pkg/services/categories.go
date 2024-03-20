package services

import (
	"errors"

	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"gorm.io/gorm"
)

func ListCategory() ([]models.Category, error) {
	var categories []models.Category

	err := database.C.Find(&categories).Error

	return categories, err
}

func GetCategory(alias string) (models.Category, error) {
	var category models.Category
	if err := database.C.Where(models.Category{Alias: alias}).First(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func GetCategoryWithID(id uint) (models.Category, error) {
	var category models.Category
	if err := database.C.Where(models.Category{
		BaseModel: models.BaseModel{ID: id},
	}).First(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func NewCategory(alias, name, description string) (models.Category, error) {
	category := models.Category{
		Alias:       alias,
		Name:        name,
		Description: description,
	}

	err := database.C.Save(&category).Error

	return category, err
}

func EditCategory(category models.Category, alias, name, description string) (models.Category, error) {
	category.Alias = alias
	category.Name = name
	category.Description = description

	err := database.C.Save(&category).Error

	return category, err
}

func DeleteCategory(category models.Category) error {
	return database.C.Delete(category).Error
}

func GetTagOrCreate(alias, name string) (models.Tag, error) {
	var tag models.Tag
	if err := database.C.Where(models.Category{Alias: alias}).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = models.Tag{
				Alias: alias,
				Name:  name,
			}
			err := database.C.Save(&tag).Error
			return tag, err
		}
		return tag, err
	}
	return tag, nil
}
