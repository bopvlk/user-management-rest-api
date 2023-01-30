package registry

import (
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/interface/controller"
	ir "git.foxminded.com.ua/3_REST_API/interal/interface/repository"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/repository"
)

func (r *registry) NewUserController() controller.UserController {
	return controller.NewUserController(r.NewUserInteractor())
}

func (r *registry) NewUserInteractor() interactor.UserInteractor {
	return interactor.NewUserInteractor(
		r.NewUserRepository(),
		r.config.HashSalt,
		[]byte(r.config.SigningKey),
		time.Duration(r.config.TokenTtl))
}

func (r *registry) NewUserRepository() repository.UserRepository {
	return ir.NewUserRepository(r.db)
}
