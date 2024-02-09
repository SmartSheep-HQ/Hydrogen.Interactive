package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/view"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
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
		api.Get("/auth", doLogin)
		api.Get("/auth/callback", postLogin)
		api.Post("/auth/refresh", doRefreshToken)

		api.Get("/users/me", auth, getUserinfo)
		api.Get("/users/:accountId", getOthersInfo)
		api.Get("/users/:accountId/follow", auth, getAccountFollowed)
		api.Post("/users/:accountId/follow", auth, doFollowAccount)

		api.Get("/attachments/o/:fileId", openAttachment)
		api.Post("/attachments", auth, uploadAttachment)

		api.Get("/posts", listPost)
		api.Get("/posts/:postId", getPost)
		api.Post("/posts", auth, createPost)
		api.Post("/posts/:postId/react/:reactType", auth, reactPost)
		api.Put("/posts/:postId", auth, editPost)
		api.Delete("/posts/:postId", auth, deletePost)

		api.Get("/realms", listRealm)
		api.Get("/realms/me", auth, listOwnedRealm)
		api.Get("/realms/:realmId", getRealm)
		api.Post("/realms", auth, createRealm)
		api.Post("/realms/:realmId/invite", auth, inviteRealm)
		api.Put("/realms/:realmId", auth, editRealm)
		api.Delete("/realms/:realmId", auth, deleteRealm)
	}

	A.Use("/", cache.New(cache.Config{
		Expiration:   24 * time.Hour,
		CacheControl: true,
	}), filesystem.New(filesystem.Config{
		Root:         http.FS(view.FS),
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
