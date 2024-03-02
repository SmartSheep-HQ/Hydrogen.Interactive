package server

import (
	"net/http"
	"strings"
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var A *fiber.App

func NewServer() {
	A = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Hydrogen.Interactive",
		AppName:               "Hydrogen.Interactive",
		ProxyHeader:           fiber.HeaderXForwardedFor,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		BodyLimit:             50 * 1024 * 1024,
		EnablePrintRoutes:     viper.GetBool("debug"),
	})

	A.Use(idempotency.New())
	A.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodOptions,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	A.Use(logger.New(logger.Config{
		Format: "${status} | ${latency} | ${method} ${path}\n",
		Output: log.Logger,
	}))

	A.Get("/.well-known", getMetadata)

	api := A.Group("/api").Name("API")
	{
		api.Get("/users/me", authMiddleware, getUserinfo)
		api.Get("/users/:accountId", getOthersInfo)
		api.Get("/users/:accountId/follow", authMiddleware, getAccountFollowed)
		api.Post("/users/:accountId/follow", authMiddleware, doFollowAccount)

		api.Get("/attachments/o/:fileId", cache.New(cache.Config{
			Expiration:   365 * 24 * time.Hour,
			CacheControl: true,
		}), openAttachment)
		api.Post("/attachments", authMiddleware, uploadAttachment)

		// TODO Feed (aka. Union source)

		moments := api.Group("/moments").Name("Moments API")
		{
			moments.Get("/", listMoment)
			moments.Get("/:momentId", getMoment)
			moments.Post("/", authMiddleware, createMoment)
			moments.Post("/:momentId/react/:reactType", authMiddleware, reactMoment)
			moments.Put("/:momentId", authMiddleware, editMoment)
			moments.Delete("/:momentId", authMiddleware, deleteMoment)
		}

		api.Get("/categories", listCategories)
		api.Post("/categories", authMiddleware, newCategory)
		api.Put("/categories/:categoryId", authMiddleware, editCategory)
		api.Delete("/categories/:categoryId", authMiddleware, deleteCategory)

		api.Get("/realms", listRealm)
		api.Get("/realms/me", authMiddleware, listOwnedRealm)
		api.Get("/realms/me/available", authMiddleware, listAvailableRealm)
		api.Get("/realms/:realmId", getRealm)
		api.Post("/realms", authMiddleware, createRealm)
		api.Post("/realms/:realmId/invite", authMiddleware, inviteRealm)
		api.Post("/realms/:realmId/kick", authMiddleware, kickRealm)
		api.Put("/realms/:realmId", authMiddleware, editRealm)
		api.Delete("/realms/:realmId", authMiddleware, deleteRealm)
	}

	A.Use("/", cache.New(cache.Config{
		Expiration:   24 * time.Hour,
		CacheControl: true,
	}), filesystem.New(filesystem.Config{
		Root:         http.FS(views.FS),
		PathPrefix:   "dist",
		Index:        "index.html",
		NotFoundFile: "dist/index.html",
	}))
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
