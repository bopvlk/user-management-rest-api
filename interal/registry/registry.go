package registry

import (
	"git.foxminded.com.ua/3_REST_API/interal/config"
	"git.foxminded.com.ua/3_REST_API/interal/interface/controller"
	"gorm.io/gorm"
)

type registry struct {
	db     *gorm.DB
	config *config.Config
}

type Registry interface {
	NewAppController() *controller.AppController
}

func NewRegistry(db *gorm.DB, config *config.Config) Registry {
	return &registry{db, config}
}

func (r *registry) NewAppController() *controller.AppController {
	return &controller.AppController{
		UserController: r.NewUserController(),
	}
}
