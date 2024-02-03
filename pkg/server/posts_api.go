package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"strings"
)

func listPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	if err := database.C.
		Where(&models.Post{RealmID: nil}).
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	posts, err := services.ListPost(take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  posts,
	})

}

func createPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias      string            `json:"alias"`
		Title      string            `json:"title"`
		Content    string            `json:"content" validate:"required"`
		Tags       []models.Tag      `json:"tags"`
		Categories []models.Category `json:"categories"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	post, err := services.NewPost(user, data.Alias, data.Title, data.Content, data.Categories, data.Tags)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(post)
}

func reactPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("postId", 0)

	var post models.Post
	if err := database.C.Where(&models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
	}).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch strings.ToLower(c.Params("reactType")) {
	case "like":
		if positive, err := services.LikePost(user, post); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	case "dislike":
		if positive, err := services.DislikePost(user, post); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported reaction")
	}
}
