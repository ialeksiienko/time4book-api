package companyrepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/company"

	"github.com/google/uuid"
)

func (r *CompanyRepo) List(ctx context.Context, page, limit int) ([]*company.Company, int64, error) {
	var total int64
	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, "SELECT count(id) FROM companies").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count companies: %w", err)
	}

	q := `SELECT id, owner_id, name, nip, address, industry, status, created_at, updated_at FROM companies ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := postgres.ExtractQuerier(ctx, r.db).Query(ctx, q, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, fmt.Errorf("query companies: %w", err)
	}
	defer rows.Close()

	var res []*company.Company
	for rows.Next() {
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

		if err := rows.Scan(
			&row.ID,
			&row.OwnerID,
			&row.Name,
			&row.NIP,
			&row.Address,
			&row.Industry,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan company: %w", err)
		}

		res = append(res, company.Reconstitute(&company.Props{
			ID:        row.ID,
			OwnerID:   row.OwnerID,
			Name:      row.Name,
			NIP:       row.NIP,
			Address:   row.Address,
			Industry:  row.Industry,
			Status:    row.Status,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}))
	}

	return res, total, nil
}
