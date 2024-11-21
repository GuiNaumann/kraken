package repositories

import (
	"context"
	"kraken/domain/entities"
)

type CertificateRepository interface {
	//CreateCertificateRepository - Create certificate and return id of certificate
	CreateCertificateRepository(ctx context.Context, certificate entities.Certificate, user entities.User) (int64, error)

	SetCertificateStatusCode(ctx context.Context, certificateID int64, statusCode entities.StatusCode) error

	//ListCertificateRepository Return a list of all Certificate with status code 0
	ListCertificateRepository(
		ctx context.Context,
		filter entities.GeneralFilter,
		user entities.User,
	) (*entities.PaginatedListUpdated[entities.Certificate], error)

	//GetCertificateByIdRepository Get a certificate by id
	GetCertificateByIdRepository(ctx context.Context, certificateID int64, user entities.User) (*entities.Certificate, error)

	//EditCertificateRepository - Edit the instructor
	EditCertificateRepository(ctx context.Context, certificate entities.Certificate, user entities.User) error

	//DeleteCertificate - Set status code of Certificate to StatusDeleted
	DeleteCertificate(ctx context.Context, certificateID int64) error
}
