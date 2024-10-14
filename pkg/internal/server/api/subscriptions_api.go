package api

import (
	"fmt"
	"strconv"

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
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get subscription: %v", err))
	} else if subscription == nil {
		return fiber.NewError(fiber.StatusNotFound, "subscription does not exist")
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
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get subscription: %v", err))
	} else if subscription == nil {
		return fiber.NewError(fiber.StatusNotFound, "subscription does not exist")
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
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get subscription: %v", err))
	} else if subscription == nil {
		return fiber.NewError(fiber.StatusNotFound, "subscription does not exist")
	}

	return c.JSON(subscription)
}

func getSubscriptionOnRealm(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	realmId, err := c.ParamsInt("realmId", 0)
	realm, err := services.GetRealmWithID(uint(realmId))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get realm: %v", err))
	}

	subscription, err := services.GetSubscriptionOnRealm(user, realm)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get subscription: %v", err))
	} else if subscription == nil {
		return fiber.NewError(fiber.StatusNotFound, "subscription does not exist")
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

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.subscribe.users",
		strconv.Itoa(int(otherUser.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

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

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.subscribe.tags",
		strconv.Itoa(int(tag.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

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

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.subscribe.categories",
		strconv.Itoa(int(category.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

	return c.JSON(subscription)
}

func subscribeToRealm(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	realmId, err := c.ParamsInt("realmId", 0)
	realm, err := services.GetRealmWithID(uint(realmId))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get realm: %v", err))
	}

	subscription, err := services.SubscribeToRealm(user, realm)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to subscribe to realm: %v", err))
	}

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.subscribe.realms",
		strconv.Itoa(int(realm.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

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

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.unsubscribe.users",
		strconv.Itoa(int(otherUser.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

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

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.unsubscribe.tags",
		strconv.Itoa(int(tag.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

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

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.unsubscribe.categories",
		strconv.Itoa(int(category.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

	return c.SendStatus(fiber.StatusOK)
}

func unsubscribeFromRealm(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	realmId, err := c.ParamsInt("realmId", 0)
	realm, err := services.GetRealmWithID(uint(realmId))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get realm: %v", err))
	}

	err = services.UnsubscribeFromRealm(user, realm)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to unsubscribe from realm: %v", err))
	}

	_ = gap.H.RecordAuditLog(
		user.ID,
		"posts.unsubscribe.realms",
		strconv.Itoa(int(realm.ID)),
		c.IP(),
		c.Get(fiber.HeaderUserAgent),
	)

	return c.SendStatus(fiber.StatusOK)
}
