package impl

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"kraken/domain/entities"
	"kraken/infrastructure/modules/impl/http_error"
	"kraken/infrastructure/repositories"
	"kraken/settings_loader"
	"kraken/utils"
	"log"
	"math"
	"strconv"
	"strings"
)

func NewCertificateRepository(
	db *sql.DB,
	settings settings_loader.SettingsLoader,
) repositories.CertificateRepository {
	return &certificateRepository{
		conn:     db,
		settings: settings,
	}
}

type certificateRepository struct {
	settings settings_loader.SettingsLoader
	conn     *sql.DB
}

func (c certificateRepository) CreateCertificateRepository(ctx context.Context, certificate entities.Certificate, user entities.User) (int64, error) {
	//language=sql
	query := `
	INSERT INTO certificate (
	    id_user,
	    name,
	    image_url,
	    is_active,
	    street,
	    address_number,
	    district,
		zip_code,
		city,
	    state,
		cpf,
		cnpj,
		phone,
		email,
		status_code,
		created_at,
		modified_at,
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	RETURNING id
	`

	tx, err := c.conn.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	var certificateID int64
	err = stmt.QueryRowContext(
		ctx,
		query,
		user.ID,
		certificate.Name,
		certificate.ImageURL,
		certificate.IsActive,
		certificate.Address.Street,
		certificate.Address.AddressNumber,
		certificate.Address.District,
		certificate.Address.ZipCode,
		certificate.Address.City,
		certificate.Address.StateID,
		certificate.CPF,
		certificate.CNPJ,
		certificate.PHONE,
		certificate.Email,
		certificate.StatusCode,
		certificate.CreatedAt,
		certificate.ModifiedAt,
	).Scan(&certificateID)

	if err != nil {
		log.Println("[CreateCertificateRepository] Error QueryRowContext", err)
		return 0, err
	}

	return certificateID, nil
}

func (c certificateRepository) SetCertificateStatusCode(ctx context.Context, certificateID int64, statusCode entities.StatusCode) error {
	//language=sql
	query := `
	UPDATE certificate 
	SET status_code = $1
	WHERE id = $2
	`

	_, err := c.conn.ExecContext(ctx, query, statusCode, certificateID)
	if err != nil {
		log.Println("[SetCertificateStatusCode] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c certificateRepository) ListCertificateRepository(
	ctx context.Context,
	filter entities.GeneralFilter,
	user entities.User,
) (*entities.PaginatedListUpdated[entities.Certificate], error) {
	////language=sql
	queryCount := `
	    SELECT COUNT(*)
	  FROM certificate
	 WHERE status_code != 2
	 AND id_user = $1`

	//language=sql
	query := `
	SELECT DISTINCT c.id,
	                c.name,
					c.image_url,
					c.is_active,
					c.street,
					c.address_number,
					c.district,
					c.zip_code,
					c.city,
					c.state,
					c.cpf,
					c.cnpj,
					c.phone,
					c.email,
					c.last_visit_date,
					c.status_code,
					c.created_at,
					c.modified_at
	FROM certificate c
	WHERE c.status_code != 2
	AND c.id_user = $1 `

	trimSearch := strings.TrimSpace(filter.Search)
	var formattedSearch string
	var searchStartingWith string
	var searchContaining string
	var searchEndingWith string

	if trimSearch != "" {
		cleanedString := utils.CleanMySQLRegexp(trimSearch)
		formattedSearch = cleanedString

		searchStartingWith = cleanedString + "%"
		searchContaining = "%" + cleanedString + "%"
		searchEndingWith = "%" + cleanedString

		query += ` AND LOWER(c.name) LIKE LOWER($2)
		ORDER BY
		CASE
        WHEN c.name LIKE $3 THEN 1 
        WHEN c.name LIKE $4 THEN 2 
        WHEN c.name LIKE $5 THEN 3 
        ELSE 4
        END,
		`
		queryCount += "\n AND name REGEXP $2 "
	}
	var ordinationAsc string
	if filter.OrdinationAsc {
		ordinationAsc = "ASC"
	} else {
		ordinationAsc = "DESC"
	}

	if trimSearch == "" {
		query += ` ORDER BY `
	}
	switch filter.Column {
	case "name":
		query += ` c.name ` + ordinationAsc
		break
	default:
		column := " c.name %s"
		if trimSearch == "" {
			column = " c.modified_at %s"
		}
		query += fmt.Sprintf(column, ordinationAsc)
	}

	if filter.Limit != 0 {
		if filter.Page > 0 {
			filter.Page--
		}
		firstItem := filter.Page * filter.Limit
		query += ` LIMIT ` + strconv.Itoa(int(firstItem)) + `, ` + strconv.Itoa(int(filter.Limit))
	} else {
		filter.Limit = math.MaxInt
	}

	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[ListCertificateRepository] Error PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmt.Close()

	var rows *sql.Rows
	if formattedSearch != "" {
		rows, err = stmt.QueryContext(ctx, user.ID, searchContaining, searchStartingWith, searchContaining, searchEndingWith)
	} else {
		rows, err = stmt.QueryContext(ctx, user.ID)
	}
	if err != nil {
		log.Println("[ListCertificateRepository] Error QueryContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer rows.Close()

	var certificates = make([]entities.Certificate, 0)
	for rows.Next() {
		var certificate entities.Certificate
		err = rows.Scan(
			&certificate.Id,
			&certificate.Name,
			&certificate.ImageURL,
			&certificate.IsActive,
			&certificate.Address.Street,
			&certificate.Address.AddressNumber,
			&certificate.Address.District,
			&certificate.Address.ZipCode,
			&certificate.Address.City,
			&certificate.Address.StateID,
			&certificate.CPF,
			&certificate.CNPJ,
			&certificate.PHONE,
			&certificate.Email,
			&certificate.LastVisitDate,
			&certificate.StatusCode,
			&certificate.CreatedAt,
			&certificate.ModifiedAt,
		)
		if err != nil {
			log.Println("[ListCertificateRepository] Error Scan", err)
			return nil, http_error.NewUnexpectedError(http_error.Unexpected)
		}
		certificates = append(certificates, certificate)
	}

	stmtCount, err := c.conn.PrepareContext(ctx, queryCount)
	if err != nil {
		log.Println("[ListCertificateRepository] Error stmtCount PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmtCount.Close()

	var totalCount int64
	if formattedSearch != "" {
		err = stmtCount.QueryRowContext(ctx, user.ID, formattedSearch).Scan(&totalCount)
	} else {
		err = stmtCount.QueryRowContext(ctx, user.ID).Scan(&totalCount)
	}
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		log.Println("[ListCertificateRepository] Error stmtCount Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	mathPage := float64(totalCount) / float64(filter.Limit)
	page := int64(math.Ceil(mathPage))

	return &entities.PaginatedListUpdated[entities.Certificate]{
		Items:          certificates,
		RequestedCount: int64(len(certificates)),
		TotalCount:     totalCount,
		Page:           page,
	}, nil
}

func (c certificateRepository) GetCertificateByIdRepository(ctx context.Context, certificateID int64, user entities.User) (*entities.Certificate, error) {
	//language=sql
	query := `
	SELECT c.id,
		   c.name,
		   c.image_url,
		   c.is_active,
		   c.street,
		   c.address_number,
		   c.district,
		   c.zip_code,
		   c.city,
		   c.state,
		   c.cpf,
		   c.cnpj,
		   c.phone,
		   c.email,
		   c.last_visit_dat
		   c.status_code,
		   c.created_at,
		   c.modified_at,
	FROM certificate c
	WHERE c.id = $1
	  AND c.status_code != $2
	  AND c.id_user = $3`

	var certificate entities.Certificate
	err := c.conn.QueryRowContext(ctx, query, certificateID, entities.StatusDeleted, user.ID).Scan(
		&certificate.Id,
		&certificate.Name,
		&certificate.ImageURL,
		&certificate.IsActive,
		&certificate.Address.Street,
		&certificate.Address.AddressNumber,
		&certificate.Address.District,
		&certificate.Address.ZipCode,
		&certificate.Address.City,
		&certificate.Address.StateID,
		&certificate.CPF,
		&certificate.CNPJ,
		&certificate.PHONE,
		&certificate.Email,
		&certificate.LastVisitDate,
		&certificate.StatusCode,
		&certificate.ModifiedAt,
		&certificate.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[GetCertificateByIdRepository] Error sql.ErrNoRows", err)
			return nil, http_error.NewBadRequestError(http_error.CertificateNotFound)
		}
		log.Println("[GetCertificateByIdRepository] Error Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return &certificate, nil
}

func (c certificateRepository) EditCertificateRepository(ctx context.Context, certificate entities.Certificate, user entities.User) error {
	//language=sql
	command := `
	UPDATE certificate
	SET name = $1,
		image_url = $2,
		is_active = $3,
		street = $4,
		address_number = $5,
		district = $6,
		zip_code = $7,
		city = $8,
		state = $9,
		cpf = $10,
		cnpj = $11,
		phone = $12,
		email = $13,
		id_user = $14
	WHERE id = $15`

	_, err := c.conn.ExecContext(ctx,
		command,
		certificate.Name,
		certificate.ImageURL,
		certificate.IsActive,
		certificate.Address.Street,
		certificate.Address.AddressNumber,
		certificate.Address.District,
		certificate.Address.ZipCode,
		certificate.Address.City,
		certificate.Address.StateID,
		certificate.CPF,
		certificate.CNPJ,
		certificate.PHONE,
		certificate.Email,
		user.ID,
		certificate.Id,
	)
	if err != nil {
		log.Println("[EditCertificateRepository] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c certificateRepository) DeleteCertificate(ctx context.Context, certificateID int64) error {
	//language=sql
	command := `
	UPDATE certificate
	SET status_code = 2
	WHERE id = $1`

	_, err := c.conn.ExecContext(
		ctx,
		command,
		certificateID,
	)
	if err != nil {
		log.Println("[DeleteCertificate] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}
