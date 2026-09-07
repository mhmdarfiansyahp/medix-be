package user

import (
	"medix-be/internal/user/handler"
	"medix-be/internal/user/repository"
	"medix-be/internal/user/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB         *gorm.DB
	Logger     *logrus.Logger
	Router     *gin.RouterGroup
	AuthRouter *gin.RouterGroup
}

func StartApp(cfg *UserHandler) {
	cfg.Logger.Info("User module starting...")

	userRepo := repository.NewUserRepository(cfg.DB)
	jwtSecret := os.Getenv("JWT_SECRET")

	userSvc := service.NewUserService(
		userRepo,
		jwtSecret,
	)

	handlerContract := &handler.HandlerContract{
		Logger: cfg.Logger,
		Router: cfg.Router,
	}
	authContract := &handler.HandlerContract{
		Logger: cfg.Logger,
		Router: cfg.AuthRouter,
	}

	handler.StartUserHandler(handlerContract, authContract, &handler.UserHandlerProps{
		UserService: userSvc,
	})
}
