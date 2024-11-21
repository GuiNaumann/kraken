package perm_impl

import (
	"context"
	"kraken/domain/entities"
	"kraken/domain/usecases"
	"kraken/infrastructure/modules/impl/http_error"
)

type certificatePermUseCase struct {
	perm usecases.CertificateUseCase
}

func NewPermCertificateUseCase(certificateUseCase usecases.CertificateUseCase) usecases.CertificateUseCase {
	return &certificatePermUseCase{
		perm: certificateUseCase,
	}
}

func (c certificatePermUseCase) CreateCertificateUseCase(ctx context.Context, user entities.User, certificate entities.Certificate) (int64, error) {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return 0, http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.CreateCertificateUseCase(ctx, user, certificate)
}

func (c certificatePermUseCase) ListCertificateUseCase(
	ctx context.Context,
	user entities.User,
	filter entities.GeneralFilter,
) (*entities.PaginatedListUpdated[entities.Certificate], error) {
	if user.IsMaster() || user.IsFlat3() || user.IsFlat2() || user.IsFlat1() {
		return c.perm.ListCertificateUseCase(ctx, user, filter)
	}

	return nil, http_error.NewUnauthorizedError(http_error.Unauthorized)
}

func (c certificatePermUseCase) GetCertificateByIdUseCase(ctx context.Context, user entities.User, certificateID int64) (*entities.Certificate, error) {
	if user.IsMaster() || user.IsFlat3() || user.IsFlat2() || user.IsFlat1() {
		return c.perm.GetCertificateByIdUseCase(ctx, user, certificateID)
	}

	return nil, http_error.NewUnauthorizedError(http_error.Unauthorized)
}

func (c certificatePermUseCase) EditCertificateUseCase(ctx context.Context, user entities.User, certificate entities.Certificate) error {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.EditCertificateUseCase(ctx, user, certificate)
}

func (c certificatePermUseCase) DeleteCertificateUseCase(ctx context.Context, user entities.User, certificateID int64) error {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.DeleteCertificateUseCase(ctx, user, certificateID)
}
