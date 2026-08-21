package app

import (
	"context"
	"fmt"
	"teaching_assistant/internal/config"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/internal/infrastructure/database"
	"teaching_assistant/internal/repository/mongodb"
	"teaching_assistant/internal/usecase"
	"teaching_assistant/pkg/jwt"

	httpRouter "teaching_assistant/internal/delivery/http"
	httpHandler "teaching_assistant/internal/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	fiberApp     *fiber.App
	cfg          *config.Config
	client       *mongo.Client
	db           *mongo.Database
	jwtManager   *jwt.Manager
	repositories Repositories
	services     Services
	handlers     Handlers
}

type Repositories struct {
	UserRepository user.UserRepository
}

type Services struct {
	UserService user.UserService
}

type Handlers struct {
	UserHandler *httpHandler.UserHandler
}

func NewApplication(ctx context.Context, cfg *config.Config) (*Application, error) {
	a := &Application{
		cfg: cfg,
	}

	a.initDatabase(ctx)
	a.initJwtManager()
	a.initRepositories()
	a.initServices()
	a.initHandlers()
	a.initRouter()

	return a, nil
}

func (a *Application) initDatabase(ctx context.Context) error {
	db, client, err := database.Connect(ctx, a.cfg.MongoDB.URI, a.cfg.MongoDB.DBName)
	if err != nil {
		return err
	}
	a.db = db
	a.client = client
	return nil
}

func (a *Application) initRepositories() {
	a.repositories.UserRepository = mongodb.NewUserRepository(a.db)
}

func (a *Application) initServices() {
	a.services.UserService = usecase.NewUserService(a.repositories.UserRepository, a.jwtManager)
}

func (a *Application) initHandlers() {
	a.handlers.UserHandler = httpHandler.NewUserHandler(a.services.UserService)
}

func (a *Application) initJwtManager() {
	a.jwtManager = jwt.NewManager(a.cfg.JWT.Secret, a.cfg.JWT.ExpireHours)
}

func (a *Application) initRouter() {
	a.fiberApp = fiber.New()
	httpRouter.NewRouter(a.fiberApp, a.handlers.UserHandler)
}

func (a *Application) Run() error {
	return a.fiberApp.Listen(fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port))
}

func (a *Application) Shutdown(ctx context.Context) error {
	if err := a.fiberApp.ShutdownWithContext(ctx); err != nil {
		return err
	}
	return a.client.Disconnect(ctx)
}
