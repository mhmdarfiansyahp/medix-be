package user

import (
	"medix-be/internal/user/handler"
	"medix-be/internal/user/repository"
	"medix-be/internal/user/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB     *gorm.DB
	Logger *logrus.Logger
	Router *gin.RouterGroup
}

func StartApp(cfg *UserHandler) {
	cfg.Logger.Info("User module starting...")

	userRepo := repository.NewUserRepository(cfg.DB)
	userSvc := service.NewUserService(userRepo)
	handlerContract := &handler.HandlerContract{
		Logger: cfg.Logger,
		Router: cfg.Router,
	}

	handler.StartUserHandler(handlerContract, &handler.UserHandlerProps{
		UserService: userSvc,
	})
}
