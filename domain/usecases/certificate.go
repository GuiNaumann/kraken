package usecases

import (
	"context"
	"kraken/domain/entities"
)

type CertificateUseCase interface {
	//CreateCertificateUseCase - Create Certificate and return id of Certificate
	CreateCertificateUseCase(ctx context.Context, user entities.User, certificate entities.Certificate) (int64, error)

	//ListCertificatesUseCase Return a list of all Certificate with status code 0
	ListCertificateUseCase(
		ctx context.Context,
		user entities.User,
		filter entities.GeneralFilter,
	) (*entities.PaginatedListUpdated[entities.Certificate], error)

	//GetCertificateByIdUseCase Get a Certificate by id and return the Certificate
	GetCertificateByIdUseCase(ctx context.Context, user entities.User, certificateID int64) (*entities.Certificate, error)

	//Editertificate Edit information about Certificate
	EditCertificateUseCase(ctx context.Context, user entities.User, certificate entities.Certificate) error

	//DeleteCertificate - delete an Clertificate
	DeleteCertificateUseCase(ctx context.Context, user entities.User, certificateID int64) error
}
