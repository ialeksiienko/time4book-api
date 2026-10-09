package resourcerepo

import (
	"database/sql"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Reader struct {
	db *pgxpool.Pool
}

func NewReader(db *postgres.Datastore) *Reader {
	return &Reader{db: db.Pool()}
}

type scanner interface {
	Scan(...interface{}) error
}

func scanResource(row scanner) (*queries.Resource, error) {
	var result queries.Resource
	var companyTypeID uuid.NullUUID
	var companyTypeName sql.NullString
	var companyTypeIcon sql.NullString
	err := row.Scan(
		&result.ID,
		&result.CompanyID,
		&result.Name,
		&result.Type,
		&result.Description,
		&result.Location,
		&result.MaxReservationMinutes,
		&result.AvailableFrom,
		&result.AvailableTo,
		&result.Status,
		&result.UnavailableFrom,
		&result.UnavailableTo,
		&result.UnavailableReason,
		&companyTypeID,
		&companyTypeName,
		&companyTypeIcon,
	)
	if err != nil {
		return nil, err
	}
	if companyTypeID.Valid {
		result.CompanyResourceType = &queries.CompanyResourceType{
			ID:      companyTypeID.UUID,
			Name:    companyTypeName.String,
			IconKey: companyTypeIcon.String,
		}
	}
	return &result, nil
}

const selectResources = `SELECT r.id, r.company_id, r.name, r.type, r.description, r.location, r.max_reservation_minutes, r.available_from, r.available_to, r.status, r.unavailable_from, r.unavailable_to, r.unavailable_reason, r.company_resource_type_id, crt.name, crt.icon_key
FROM resources r
LEFT JOIN company_resource_types crt ON crt.id = r.company_resource_type_id AND crt.company_id = r.company_id`
