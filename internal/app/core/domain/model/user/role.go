package user

import (
	"errors"
	"fmt"
	"time"
)

type RoleKey string

const (
	RoleOwnerKey     RoleKey = "owner"
	RoleAdminKey     RoleKey = "admin"
	RoleEmployeeKey  RoleKey = "employee"
	RoleDeveloperKey RoleKey = "developer"
)

func (r RoleKey) FriendlyName() string {
	switch r {
	case RoleOwnerKey:
		return "Owner"
	case RoleAdminKey:
		return "Admin"
	case RoleEmployeeKey:
		return "Employee"
	case RoleDeveloperKey:
		return "Developer"
	default:
		return ""
	}
}

type Role struct {
	id           RoleKey
	friendlyName string
	createdAt    time.Time
	updatedAt    *time.Time
}

func NewRole(
	key RoleKey,
	friendlyName string,
) (*Role, error) {
	if err := validateRoleKey(key); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &Role{
		id:           key,
		friendlyName: friendlyName,
	}, nil
}

func (r *Role) IsPrivileged() bool {
	return r.IsOwner() || r.IsAdmin() || r.IsDeveloper()
}

func (r *Role) Change(key RoleKey, friendlyName string) error {
	if err := validateRoleKey(key); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	now := time.Now().UTC()

	r.friendlyName = friendlyName
	r.updatedAt = &now

	return nil
}

func validateRoleKey(key RoleKey) error {
	if key != RoleOwnerKey &&
		key != RoleAdminKey &&
		key != RoleEmployeeKey &&
		key != RoleDeveloperKey {
		return errors.New("role key")
	}

	return nil
}

func (r *Role) IsAdmin() bool     { return r.id == RoleAdminKey }
func (r *Role) IsOwner() bool     { return r.id == RoleOwnerKey }
func (r *Role) IsEmployee() bool  { return r.id == RoleEmployeeKey }
func (r *Role) IsDeveloper() bool { return r.id == RoleDeveloperKey }

func (r *Role) Key() RoleKey         { return r.id }
func (r *Role) FriendlyName() string { return r.friendlyName }

func (r *Role) String() string { return string(r.Key()) }

func RoleKeyFromString(s string) RoleKey { return RoleKey(s) }
