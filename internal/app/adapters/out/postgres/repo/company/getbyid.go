package companyrepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/company"

	"github.com/google/uuid"
)

func (r *CompanyRepo) ByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	q := `SELECT id, owner_id, name, nip, address, industry, status, created_at, updated_at FROM companies WHERE id = $1`

	var row struct {
		ID        uuid.UUID
		OwnerID   uuid.UUID
		Name      string
		NIP       *string
		Address   *string
		Industry  *string
		Status    string
		CreatedAt time.Time
		UpdatedAt *time.Time
	}

	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, q, id).Scan(
		&row.ID,
		&row.OwnerID,
		&row.Name,
		&row.NIP,
		&row.Address,
		&row.Industry,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan company: %w", err)
	}

	return company.Reconstitute(&company.Props{
		ID:        row.ID,
		OwnerID:   row.OwnerID,
		Name:      row.Name,
		NIP:       row.NIP,
		Address:   row.Address,
		Industry:  row.Industry,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}), nil
}
