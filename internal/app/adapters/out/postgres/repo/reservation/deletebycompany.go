package reservationrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"

	"github.com/google/uuid"
)

func (r *ReservationRepo) DeleteByCompanyID(ctx context.Context, companyID uuid.UUID) error {
	q := `DELETE FROM reservations WHERE company_id = $1`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, companyID)
	if err != nil {
		return fmt.Errorf("exec delete reservations by company id: %w", err)
	}
	return nil
}
