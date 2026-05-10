package user

import (
	"errors"
	"time"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
	ErrUnauthorized = errors.New("security violation: unauthorized action")
)

type User struct {
	id        uuid.UUID
	companyID uuid.UUID
	firstname string
	lastname  string
	email     string
	role      *Role
	status    UserStatus
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(
	firstname string,
	lastname string,
	email string,
	role *Role,
	companyID uuid.UUID,
) (*User, error) {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	return &User{
		id:        uuid.New(),
		companyID: companyID,
		firstname: firstname,
		lastname:  lastname,
		email:     email,
		role:      role,
		status:    StatusActive,
		createdAt: time.Now().UTC(),
	}, nil
}

func (u *User) Deactivate() {
	u.status = StatusInactive
	u.updatedAt = time.Now().UTC()
}

func (u *User) IsActive() bool       { return u.status == StatusActive }
func (u *User) ID() uuid.UUID        { return u.id }
func (u *User) CompanyID() uuid.UUID { return u.companyID }
func (u *User) Firstname() string    { return u.firstname }
func (u *User) Lastname() string     { return u.lastname }
func (u *User) Email() string        { return u.email }
func (u *User) Role() *Role          { return u.role }
func (u *User) Status() UserStatus   { return u.status }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

func (u *User) UpdateProfile(firstname, lastname *string) {
	if firstname != nil {
		u.firstname = *firstname
	}
	if lastname != nil {
		u.lastname = *lastname
	}
	u.updatedAt = time.Now().UTC()
}

func (u *User) ChangeStatus(status UserStatus) {
	u.status = status
	u.updatedAt = time.Now().UTC()
}

func (u *User) ChangeRole(newRole *Role) {
	u.role = newRole
	u.updatedAt = time.Now().UTC()
}

func (u *User) SetCompanyID(companyID uuid.UUID) {
	if u.companyID == uuid.Nil {
		u.companyID = companyID
	}
}

func (u *User) CanCreateUserWithRole(target *Role) bool {
	if u.role.IsOwner() {
		return true
	}
	if u.role.IsAdmin() {
		return target.IsEmployee()
	}
	if u.role.IsDeveloper() {
		return target.IsEmployee() || target.IsDeveloper()
	}
	return false
}

type Props struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Firstname string
	Lastname  string
	Email     string
	RoleKey   string
	Status    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func Reconstitute(p *Props) *User {
	u := &User{
		id:        p.ID,
		companyID: p.CompanyID,
		firstname: p.Firstname,
		lastname:  p.Lastname,
		email:     p.Email,
		role:      &Role{id: RoleKeyFromString(p.RoleKey)},
		status:    UserStatus(p.Status),
		createdAt: p.CreatedAt,
	}
	if p.UpdatedAt != nil {
		u.updatedAt = *p.UpdatedAt
	}
	return u
}
