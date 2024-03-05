package server

import (
	"strings"
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func contextMoment() *services.PostTypeContext {
	return &services.PostTypeContext{
		Tx:         database.C,
		TableName:  "moments",
		ColumnName: "moment",
		CanReply:   false,
		CanRepost:  true,
	}
}

func createMoment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias       string              `json:"alias"`
		Content     string              `json:"content" validate:"required"`
		Hashtags    []models.Tag        `json:"hashtags"`
		Categories  []models.Category   `json:"categories"`
		Attachments []models.Attachment `json:"attachments"`
		PublishedAt *time.Time          `json:"published_at"`
		RealmID     *uint               `json:"realm_id"`
		RepostTo    uint                `json:"repost_to"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	item := &models.Moment{
		PostBase: models.PostBase{
			Alias:       data.Alias,
			Attachments: data.Attachments,
			PublishedAt: data.PublishedAt,
			AuthorID:    user.ID,
		},
		Hashtags:   data.Hashtags,
		Categories: data.Categories,
		Content:    data.Content,
		RealmID:    data.RealmID,
	}

	var relatedCount int64
	if data.RepostTo > 0 {
		if err := database.C.Where("id = ?", data.RepostTo).
			Model(&models.Moment{}).Count(&relatedCount).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if relatedCount <= 0 {
			return fiber.NewError(fiber.StatusNotFound, "related post was not found")
		} else {
			item.RepostID = &data.RepostTo
		}
	}

	var realm *models.Realm
	if data.RealmID != nil {
		if err := database.C.Where(&models.Realm{
			BaseModel: models.BaseModel{ID: *data.RealmID},
		}).First(&realm).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	item, err := services.NewPost(item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func editMoment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("momentId", 0)

	var data struct {
		Alias       string              `json:"alias" validate:"required"`
		Content     string              `json:"content" validate:"required"`
		PublishedAt *time.Time          `json:"published_at"`
		Hashtags    []models.Tag        `json:"hashtags"`
		Categories  []models.Category   `json:"categories"`
		Attachments []models.Attachment `json:"attachments"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var item *models.Moment
	if err := database.C.Where(models.Comment{
		PostBase: models.PostBase{
			BaseModel: models.BaseModel{ID: uint(id)},
			AuthorID:  user.ID,
		},
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.Alias = data.Alias
	item.Content = data.Content
	item.PublishedAt = data.PublishedAt
	item.Hashtags = data.Hashtags
	item.Categories = data.Categories
	item.Attachments = data.Attachments

	if item, err := services.EditPost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(item)
	}
}

func deleteMoment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("momentId", 0)

	var item *models.Moment
	if err := database.C.Where(models.Comment{
		PostBase: models.PostBase{
			BaseModel: models.BaseModel{ID: uint(id)},
			AuthorID:  user.ID,
		},
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeletePost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
