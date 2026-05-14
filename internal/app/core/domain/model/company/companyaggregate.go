package company

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	id        uuid.UUID
	ownerID   uuid.UUID
	name      string
	nip       *string
	address   *string
	industry  *string
	status    CompanyStatus
	createdAt time.Time
	updatedAt *time.Time
}

func NewCompany(
	ownerID uuid.UUID,
	name string,
	nip *string,
	address *string,
	industry *string,
) (*Company, error) {
	if name == "" {
		return nil, fmt.Errorf("company name is required")
	}

	return &Company{
		id:        uuid.New(),
		ownerID:   ownerID,
		name:      name,
		nip:       nip,
		address:   address,
		industry:  industry,
		status:    StatusActive,
		createdAt: time.Now().UTC(),
	}, nil
}

func (c *Company) Block() {
	c.status = StatusBlocked
	now := time.Now().UTC()
	c.updatedAt = &now
}

func (c *Company) Unblock() {
	c.status = StatusActive
	now := time.Now().UTC()
	c.updatedAt = &now
}

func (c *Company) IsBlocked() bool       { return c.status == StatusBlocked }
func (c *Company) ID() uuid.UUID         { return c.id }
func (c *Company) OwnerID() uuid.UUID    { return c.ownerID }
func (c *Company) Name() string          { return c.name }
func (c *Company) NIP() *string          { return c.nip }
func (c *Company) Address() *string      { return c.address }
func (c *Company) Industry() *string     { return c.industry }
func (c *Company) Status() CompanyStatus { return c.status }
func (c *Company) CreatedAt() time.Time  { return c.createdAt }
func (c *Company) UpdatedAt() *time.Time { return c.updatedAt }

type Props struct {
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

func Reconstitute(p *Props) *Company {
	return &Company{
		id:        p.ID,
		ownerID:   p.OwnerID,
		name:      p.Name,
		nip:       p.NIP,
		address:   p.Address,
		industry:  p.Industry,
		status:    CompanyStatus(p.Status),
		createdAt: p.CreatedAt,
		updatedAt: p.UpdatedAt,
	}
}
