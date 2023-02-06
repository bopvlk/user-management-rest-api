package registry

import (
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/interface/controller"
	ir "git.foxminded.com.ua/3_REST_API/interal/interface/repository"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
)

func (r *registry) NewUserController() controller.UserController {
	return controller.NewUserController(
		interactor.NewUserInteractor(
			ir.NewUserRepository(r.db),
			r.config.HashSalt,
			[]byte(r.config.SigningKey),
			time.Duration(r.config.TokenTtl)))
}
