package handler

import (
	"strings"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/middleware"
	"github.com/bayupaths/bypur-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	profilePath           = "/profile"
	experiencesPath       = "/experiences"
	experienceIdPath      = "/experiences/:id"
	offeringsPath         = "/offerings"
	offeringIdPath        = "/offerings/:id"
	skillsPath            = "/skills"
	skillIdPath           = "/skills/:id"
	projectsPath          = "/projects"
	projectIdPath         = "/projects/:id"
	contactMessagesIdPath = "/contact/messages/:id"
	storageFilesPath      = "/storage/files"
)

type Router struct {
	App         *fiber.App
	Cfg         *config.Config
	AuthH       *AuthHandler
	PublicPortH *PublicPortfolioHandler
	AdminH      *AdminHandler
	StorH       *StorageHandler
}

func (r *Router) Setup() {
	app := r.App
	cfg := r.Cfg

	app.Use(middleware.RecoveryMiddleware())

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	allowOrigins := strings.Join(cfg.Server.ParsedCorsOrigins, ", ")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, x-api-key",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return !cfg.IsProduction()
		},
		Max:        100,
		Expiration: 15 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return response.SendError(c, "Too many requests, please try again later.", "Rate limit exceeded", fiber.StatusTooManyRequests)
		},
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return response.SendSuccess(c, fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		}, "Health OK")
	})

	api := app.Group("/api")

	api.Get("/version", func(c *fiber.Ctx) error {
		return response.SendSuccess(c, fiber.Map{
			"app":         cfg.App.Name,
			"version":     cfg.App.Version,
			"environment": cfg.App.Env,
		}, "Version retrieved successfully")
	})

	apiKeyAuth := middleware.AuthenticateApiKey(cfg)
	jwtAuth := middleware.AuthenticateJWT(cfg)

	// Public Routes (Protected by API Key)
	r.registerPublicRoutes(api.Group("/public", apiKeyAuth))

	// Auth Routes (Public & Protected)
	auth := api.Group("/auth")
	r.registerAuthRoutes(auth, auth.Group("", jwtAuth))

	// Admin Routes (Protected by JWT)
	r.registerAdminRoutes(api.Group("/admin", jwtAuth))
}

func (r *Router) registerPublicRoutes(public fiber.Router) {
	h := r.PublicPortH

	public.Get(profilePath, h.GetProfile)

	public.Get(experiencesPath, h.GetExperiences)
	public.Get(experienceIdPath, h.GetExperienceByID)

	public.Get(offeringsPath, h.GetOfferings)
	public.Get("/offerings/:slug", h.GetOfferingBySlug)

	public.Get(skillsPath, h.GetSkills)
	public.Get("/skills/grouped/category", h.GetSkillsByCategory)
	public.Get(skillIdPath, h.GetSkillByID)

	public.Get("/projects/featured", h.GetFeaturedProjects)
	public.Get(projectsPath, h.GetProjects)
	public.Get("/projects/:slug", h.GetProjectBySlug)

	public.Post("/contact", h.SubmitContact)
}

func (r *Router) registerAuthRoutes(auth fiber.Router, authProtected fiber.Router) {
	h := r.AuthH

	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)
	auth.Post("/logout", h.Logout)

	authProtected.Get("/me", h.Me)
	authProtected.Put(profilePath, h.UpdateProfile)
	authProtected.Post("/change-password", h.ChangePassword)
}

func (r *Router) registerAdminRoutes(admin fiber.Router) {
	h := r.AdminH
	storH := r.StorH

	admin.Put(profilePath, h.UpdateProfile)

	admin.Post(experiencesPath, h.CreateExperience)
	admin.Put(experienceIdPath, h.UpdateExperience)
	admin.Delete(experienceIdPath, h.DeleteExperience)

	admin.Get("/offerings/all", h.GetAllOfferings)
	admin.Get(offeringIdPath, h.GetOfferingByID)
	admin.Post(offeringsPath, h.CreateOffering)
	admin.Put(offeringIdPath, h.UpdateOffering)
	admin.Delete(offeringIdPath, h.DeleteOffering)
	admin.Patch("/offerings/:id/toggle-status", h.ToggleOfferingStatus)
	admin.Post("/offerings/reorder", h.ReorderOfferings)

	admin.Post(skillsPath, h.CreateSkill)
	admin.Put(skillIdPath, h.UpdateSkill)
	admin.Delete(skillIdPath, h.DeleteSkill)

	admin.Post(projectsPath, h.CreateProject)
	admin.Put(projectIdPath, h.UpdateProject)
	admin.Delete(projectIdPath, h.DeleteProject)

	admin.Get("/contact/messages", h.GetMessages)
	admin.Get(contactMessagesIdPath, h.GetMessageByID)
	admin.Put("/contact/messages/:id/status", h.UpdateMessageStatus)
	admin.Delete(contactMessagesIdPath, h.DeleteMessage)
	admin.Get("/contact/unread", h.GetUnreadMessages)
	admin.Get("/contact/stats", h.GetMessageStats)
	admin.Post("/contact/messages/:id/mark-as-read", h.MarkAsRead)

	admin.Get("/storage/check", storH.CheckConnection)
	admin.Post("/storage/upload", storH.UploadFile)
	admin.Get(storageFilesPath, storH.ListFiles)
	admin.Get("/storage/info", storH.GetStorageInfo)
	admin.Delete(storageFilesPath, storH.DeleteFile)
	admin.Delete(storageFilesPath+"/*", storH.DeleteFile)
}
