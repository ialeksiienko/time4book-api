package reservationrepo

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestConcurrentReservationsCannotOverlap(t *testing.T) {
	databaseURL := os.Getenv("TIME4BOOK_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TIME4BOOK_TEST_DATABASE_URL for PostgreSQL integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	companyID, userID, resourceID := uuid.New(), uuid.New(), uuid.New()
	if _, err := pool.Exec(ctx, "INSERT INTO companies (id, owner_id, name) VALUES ($1, $2, $3)", companyID, userID, companyID.String()); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, company_id, firstname, lastname, email, role_id) VALUES ($1, $2, 'Test', 'User', $3, 'owner')", userID, companyID, userID.String()+"@example.com"); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO resources (id, company_id, name, type) VALUES ($1, $2, 'Room', 'room')", resourceID, companyID); err != nil {
		t.Fatal(err)
	}

	repo := New(postgres.New(pool))
	start := time.Date(2026, time.October, 9, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	first, err := reservation.NewReservation(userID, companyID, resourceID, start, end, nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := reservation.NewReservation(userID, companyID, resourceID, start.Add(30*time.Minute), end.Add(30*time.Minute), nil)
	if err != nil {
		t.Fatal(err)
	}

	ready := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for _, booking := range []*reservation.Reservation{first, second} {
		workers.Add(1)
		go func(booking *reservation.Reservation) {
			defer workers.Done()
			<-ready
			results <- repo.Create(ctx, booking)
		}(booking)
	}
	close(ready)
	workers.Wait()
	close(results)

	successes, conflicts := 0, 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, reservationout.ErrTimeSlotBooked):
			conflicts++
		default:
			t.Fatalf("unexpected create error: %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("got %d successful bookings and %d conflicts, want one of each", successes, conflicts)
	}

	adjacent, err := reservation.NewReservation(userID, companyID, resourceID, end.Add(30*time.Minute), end.Add(90*time.Minute), nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, adjacent); err != nil {
		t.Fatalf("adjacent booking failed: %v", err)
	}

	if _, err := pool.Exec(ctx, "UPDATE reservations SET status = 'cancelled' WHERE resource_id = $1", resourceID); err != nil {
		t.Fatal(err)
	}
	replacement, err := reservation.NewReservation(userID, companyID, resourceID, start, end, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, replacement); err != nil {
		t.Fatalf("booking cancelled time failed: %v", err)
	}
}
