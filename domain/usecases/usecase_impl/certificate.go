package usecase_impl

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"kraken/domain/entities"
	"kraken/domain/entities/rules"
	"kraken/domain/usecases"
	"kraken/domain/usecases/perm_impl"
	"kraken/infrastructure/modules/impl/http_error"
	"kraken/infrastructure/repositories"
	"kraken/infrastructure/storage"
	"kraken/settings_loader"
	"log"
	"path/filepath"
	"strings"
)

func NewCertificateUseCase(
	repo repositories.CertificateRepository,
	settings settings_loader.SettingsLoader,
	fileStorage storage.FileStorageRepositoryNew,
) usecases.CertificateUseCase {
	return perm_impl.NewPermCertificateUseCase(
		&certificateUseCase{
			repo:        repo,
			settings:    settings,
			fileStorage: fileStorage,
		})
}

type certificateUseCase struct {
	repo        repositories.CertificateRepository
	settings    settings_loader.SettingsLoader
	fileStorage storage.FileStorageRepositoryNew
}

func (c certificateUseCase) CreateCertificateUseCase(ctx context.Context, user entities.User, certificate entities.Certificate) (int64, error) {
	err := rules.CertificateRules(&certificate)
	if err != nil {
		log.Println("[CreateCertificateUseCase] Error certificateRules", err)
		return 0, err
	}

	if certificate.ImageBase64 != "" && !storage.IsURL(certificate.ImageBase64) {
		generated, err := uuid.NewUUID()
		if err != nil {
			log.Println("[CreateCertificateUseCase] Error NewUUID", err)
			return 0, http_error.NewUnexpectedError(http_error.Unexpected)
		}

		filePath := fmt.Sprintf("/certificates/%s", generated)

		certificate.ImageURL, err = c.fileStorage.SaveBase64(certificate.ImageBase64, filePath)
		if err != nil {
			log.Println("[CreateCertificateUseCase] Error SaveBase64", err)
			return 0, http_error.NewUnexpectedError(http_error.Unexpected)
		}
	}

	id, err := c.repo.CreateCertificateRepository(ctx, certificate, user)
	if err != nil {
		log.Println("[CreateCertificateUseCase] Error CreateCertificateRepository", err)
		return 0, err
	}

	err = c.repo.SetCertificateStatusCode(ctx, certificate.Id, entities.StatusExist)
	if err != nil {
		log.Println("[CreateCertificateUseCase] Error SetCertificateStatusCode", err)
		return 0, err
	}

	return id, nil
}

func (c certificateUseCase) ListCertificateUseCase(
	ctx context.Context,
	user entities.User,
	filter entities.GeneralFilter,
) (*entities.PaginatedListUpdated[entities.Certificate], error) {
	return c.repo.ListCertificateRepository(ctx, filter, user)
}

func (c certificateUseCase) GetCertificateByIdUseCase(ctx context.Context, user entities.User, certificateID int64) (*entities.Certificate, error) {
	return c.repo.GetCertificateByIdRepository(ctx, certificateID, user)
}

func (c certificateUseCase) EditCertificateUseCase(ctx context.Context, user entities.User, certificate entities.Certificate) error {
	err := rules.CertificateRules(&certificate)
	if err != nil {
		log.Println("[EditCertificateUseCase] Error CertificateRules", err)
		return err
	}

	oldCertificate, err := c.repo.GetCertificateByIdRepository(ctx, certificate.Id, user)
	if err != nil {
		log.Println("[EditCertificateUseCase] Error GetCertificateByIdRepository", err)
		return err
	}

	if certificate.ImageBase64 != "" && !storage.IsURL(certificate.ImageBase64) {
		if oldCertificate.ImageURL != "" {
			_, fileName := filepath.Split(strings.Split(oldCertificate.ImageURL, "?")[0])

			err = c.fileStorage.DeletePath(filepath.Join("images", "certificates", fileName))
			if err != nil {
				log.Println("[EditCertificateUseCase] Error DeletePath", err)
				return http_error.NewUnexpectedError(http_error.Unexpected)
			}
		}

		generated, err := uuid.NewUUID()
		if err != nil {
			log.Println("[EditCertificateUseCase] Error NewUUID", err)
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}

		filePath := fmt.Sprintf("/certificates/%s", generated)

		certificate.ImageURL, err = c.fileStorage.SaveBase64(certificate.ImageBase64, filePath)
		if err != nil {
			log.Println("[EditCertificateUseCase] Error SaveBase64", err)
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}
	}

	err = c.repo.SetCertificateStatusCode(ctx, certificate.Id, entities.StatusExist)
	if err != nil {
		log.Println("[EditCertificateUseCase] Error SetCertificateStatusCode", err)
		return err
	}

	return c.repo.EditCertificateRepository(ctx, certificate, user)
}

func (c certificateUseCase) DeleteCertificateUseCase(ctx context.Context, user entities.User, certificateID int64) error {
	oldCertificate, err := c.repo.GetCertificateByIdRepository(ctx, certificateID, user)
	if err != nil {
		log.Println("[DeleteCertificateUseCase] Error GetCertificateByIdRepository")
		return err
	}

	if oldCertificate.ImageBase64 != "" {
		_, fileName := filepath.Split(strings.Split(oldCertificate.ImageURL, "?")[0])

		err = c.fileStorage.DeletePath(filepath.Join("images", "certificateImages", fileName))
		if err != nil {
			log.Println("[DeleteCertificateUseCase] Error DeletePath")
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}
	}

	return c.repo.DeleteCertificate(ctx, certificateID)
}
