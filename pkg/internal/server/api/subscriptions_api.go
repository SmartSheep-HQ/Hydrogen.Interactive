package api

import (
	"fmt"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getSubscriptionOnUser(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	otherUserId, err := c.ParamsInt("userId", 0)
	otherUser, err := services.GetAccountWithID(uint(otherUserId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get user: %v", err))
	}

	subscription, err := services.GetSubscriptionOnUser(user, otherUser)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get subscription: %v", err))
	}

	return c.JSON(subscription)
}

func getSubscriptionOnTag(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	tagId, err := c.ParamsInt("tagId", 0)
	tag, err := services.GetTagWithID(uint(tagId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get tag: %v", err))
	}

	subscription, err := services.GetSubscriptionOnTag(user, tag)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get subscription: %v", err))
	}

	return c.JSON(subscription)
}

func getSubscriptionOnCategory(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	categoryId, err := c.ParamsInt("categoryId", 0)
	category, err := services.GetCategoryWithID(uint(categoryId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get category: %v", err))
	}

	subscription, err := services.GetSubscriptionOnCategory(user, category)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get subscription: %v", err))
	}

	return c.JSON(subscription)
}

func subscribeToUser(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	otherUserId, err := c.ParamsInt("userId", 0)
	otherUser, err := services.GetAccountWithID(uint(otherUserId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get user: %v", err))
	}

	subscription, err := services.SubscribeToUser(user, otherUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to subscribe to user: %v", err))
	}

	return c.JSON(subscription)
}

func subscribeToTag(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	tagId, err := c.ParamsInt("tagId", 0)
	tag, err := services.GetTagWithID(uint(tagId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get tag: %v", err))
	}

	subscription, err := services.SubscribeToTag(user, tag)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to subscribe to tag: %v", err))
	}

	return c.JSON(subscription)
}

func subscribeToCategory(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	categoryId, err := c.ParamsInt("categoryId", 0)
	category, err := services.GetCategoryWithID(uint(categoryId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get category: %v", err))
	}

	subscription, err := services.SubscribeToCategory(user, category)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to subscribe to category: %v", err))
	}

	return c.JSON(subscription)
}

func unsubscribeFromUser(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	otherUserId, err := c.ParamsInt("userId", 0)
	otherUser, err := services.GetAccountWithID(uint(otherUserId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get user: %v", err))
	}

	err = services.UnsubscribeFromUser(user, otherUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to unsubscribe from user: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}

func unsubscribeFromTag(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	tagId, err := c.ParamsInt("tagId", 0)
	tag, err := services.GetTagWithID(uint(tagId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get tag: %v", err))
	}

	err = services.UnsubscribeFromTag(user, tag)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to unsubscribe from tag: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}

func unsubscribeFromCategory(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	categoryId, err := c.ParamsInt("categoryId", 0)
	category, err := services.GetCategoryWithID(uint(categoryId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unable to get category: %v", err))
	}

	err = services.UnsubscribeFromCategory(user, category)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to unsubscribe from category: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
