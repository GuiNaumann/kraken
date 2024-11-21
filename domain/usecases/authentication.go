package usecases

import (
	"context"
	entities "kraken/domain/entities"
)

// AuthUseCase standard contract defining the valid use cases for authentication
type AuthUseCase interface {
	//Login Check user CPF and Password, if user haves CPF and Password in local database,
	//a token is created to  the user log in. > REQUEST ORIGIN (http_authentication / login)
	Login(ctx context.Context, credential entities.LoginCredentials) (*entities.User, string, error)

	//RegisterUser Create a new user.
	RegisterUser(ctx context.Context, user entities.User) error
}
