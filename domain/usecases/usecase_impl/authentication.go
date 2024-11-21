package usecase_impl

import (
	"context"
	"github.com/golang-jwt/jwt"
	entities "kraken/domain/entities"
	"kraken/domain/entities/rules"
	"kraken/domain/usecases"
	"kraken/infrastructure/modules/impl/http_error"
	"kraken/infrastructure/repositories"
	"kraken/settings_loader"
	"kraken/utils"
	"log"
	"strings"
	"time"
)

func NewAuthenticationUseCase(
	repo repositories.AuthenticationRepository,
	settings settings_loader.SettingsLoader,
) usecases.AuthUseCase {
	return &authUseCase{
		repo:     repo,
		settings: settings,
	}
}

type authUseCase struct {
	repo     repositories.AuthenticationRepository
	settings settings_loader.SettingsLoader
}

func (u authUseCase) Login(
	ctx context.Context,
	credential entities.LoginCredentials,
) (*entities.User, string, error) {
	if credential.Login == "" {
		return nil, "", http_error.NewBadRequestError(http_error.LoginCannotBeEmpty)
	}

	if credential.Password == "" {
		return nil, "", http_error.NewBadRequestError(http_error.EmptyPasswordField)
	}

	if !strings.Contains(credential.Login, "@") {
		credential.Login = strings.Replace(credential.Login, ".", "", -1)
		credential.Login = strings.Replace(credential.Login, "-", "", -1)
	}

	exists, err := u.repo.UserExists(ctx, credential)
	if err != nil {
		log.Println("[Login] Error UserExists", err)
		return nil, "", err
	}

	if !exists {
		log.Println("[Login] Error Not Found exists", err)
		return nil, "", http_error.NewForbiddenError(http_error.Forbidden)
	}

	user, err := u.repo.GetUserByLogin(ctx, credential.Login)
	if err != nil {
		log.Println("[Login] Error Not Found", err)
		return nil, "", http_error.NewForbiddenError(http_error.Forbidden)
	}

	//credential.Password = strings.ToLower(credential.Password)
	credential.Password, err = utils.RemoveAccents(credential.Password)
	if err != nil {
		return nil, "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	passwordCheck, err := u.repo.ComparePasswordHash(ctx, credential.Login, credential.Password)
	if err != nil {
		return nil, "", err
	}
	if !passwordCheck {
		return nil, "", http_error.NewForbiddenError(http_error.Forbidden)
	}

	atClaims := jwt.MapClaims{}
	atClaims["id"] = user.ID
	atClaims["exp"] = time.Now().Add(time.Hour * 720).Unix()

	securityConfig := u.settings.GetSecurityConfig()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenKraken, err := at.SignedString([]byte(securityConfig.JWTSecret))
	if err != nil {
		return nil, "", err
	}

	return user, tokenKraken, nil
}

func (u authUseCase) RegisterUser(ctx context.Context, user entities.User) error {
	err := rules.ValidateUserRegister(&user)
	if err != nil {
		log.Println("[RegisterUser] Error ValidateUserRegister")
		return err
	}

	existEmail, err := u.repo.EmailExists(ctx, user)
	if err != nil {
		return err
	}

	if existEmail {
		return http_error.NewBadRequestError(http_error.EmailExistError)
	}

	if user.IsForeigner == true {
		if *user.Document != "" {
			documentExist, err := u.repo.DocumentExists(ctx, user)
			if err != nil {
				return err
			}

			if documentExist {
				return http_error.NewUnexpectedError(http_error.DocumentExistError)
			}
		}
	} else {
		documentExist, err := u.repo.DocumentExists(ctx, user)
		if err != nil {
			return err
		}

		if documentExist {
			return http_error.NewUnexpectedError(http_error.DocumentExistError)
		}
	}

	user.Name = utils.CapitalizeWords(user.Name)

	return u.repo.RegisterUser(ctx, user)
}
