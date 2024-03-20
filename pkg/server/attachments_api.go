package server

import (
	"path/filepath"

	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func readAttachment(c *fiber.Ctx) error {
	id := c.Params("fileId")
	basepath := viper.GetString("content")

	return c.SendFile(filepath.Join(basepath, id))
}

func uploadAttachment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	file, err := c.FormFile("attachment")
	if err != nil {
		return err
	}

	attachment, err := services.NewAttachment(user, file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.SaveFile(file, attachment.GetStoragePath()); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"info": attachment,
		"url":  attachment.GetAccessPath(),
	})
}
