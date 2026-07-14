package bootstrap

import (
	"bayupur-portofolio-be/internal/config"
	"bayupur-portofolio-be/internal/handler"
	"bayupur-portofolio-be/internal/repository"
	"bayupur-portofolio-be/internal/service"

	"gorm.io/gorm"
)

type container struct {
	repos    *repositories
	services *services
	handlers *handlers
}

type repositories struct {
	userRepo       repository.UserRepository
	tokenRepo      repository.RefreshTokenRepository
	profileRepo    repository.ProfileRepository
	offeringRepo   repository.OfferingRepository
	skillRepo      repository.SkillRepository
	experienceRepo repository.ExperienceRepository
	projectRepo    repository.ProjectRepository
	messageRepo    repository.ContactMessageRepository
}

type services struct {
	storageService   *service.StorageService
	mailService      *service.MailService
	authService      *service.AuthService
	portfolioService *service.PortfolioService
}

type handlers struct {
	authH       *handler.AuthHandler
	publicPortH *handler.PublicPortfolioHandler
	adminH      *handler.AdminHandler
	storageH    *handler.StorageHandler
}

func initContainer(db *gorm.DB, cfg *config.Config) *container {
	// 1. Inisialisasi Repositories
	repos := &repositories{
		userRepo:       repository.NewUserRepository(db),
		tokenRepo:      repository.NewRefreshTokenRepository(db),
		profileRepo:    repository.NewProfileRepository(db),
		offeringRepo:   repository.NewOfferingRepository(db),
		skillRepo:      repository.NewSkillRepository(db),
		experienceRepo: repository.NewExperienceRepository(db),
		projectRepo:    repository.NewProjectRepository(db),
		messageRepo:    repository.NewContactMessageRepository(db),
	}

	// 2. Inisialisasi Services
	services := &services{
		storageService: service.NewStorageService(cfg),
		mailService:    service.NewMailService(cfg),
	}
	services.authService = service.NewAuthService(repos.userRepo, repos.tokenRepo, cfg)
	services.portfolioService = service.NewPortfolioService(
		repos.profileRepo,
		repos.offeringRepo,
		repos.skillRepo,
		repos.experienceRepo,
		repos.projectRepo,
		repos.messageRepo,
	)

	// 3. Inisialisasi Handlers
	handlers := &handlers{
		authH:       handler.NewAuthHandler(services.authService, cfg),
		publicPortH: handler.NewPublicPortfolioHandler(services.portfolioService, services.mailService, cfg),
		adminH:      handler.NewAdminHandler(services.portfolioService),
		storageH:    handler.NewStorageHandler(services.storageService, cfg),
	}

	return &container{
		repos:    repos,
		services: services,
		handlers: handlers,
	}
}
