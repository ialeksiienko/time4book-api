package companyresourcetype

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName    = errors.New("name is required")
	ErrInvalidIconKey = errors.New("icon key is not allowed")
)

// AllowedIconKeys matches the picker on the frontend (80 distinct lucide-derived slugs).
var AllowedIconKeys = map[string]struct{}{
	"room": {}, "building": {}, "armchair": {},
	"car": {}, "truck": {}, "bus": {}, "ambulance": {}, "tractor": {}, "bike": {}, "cable_car": {}, "plane": {}, "train_front": {}, "ship": {},
	"equipment": {}, "monitor": {}, "laptop": {}, "tablet": {}, "smartphone": {}, "watch": {},
	"keyboard": {}, "mouse": {},
	"headphones": {}, "mic": {}, "speaker": {},
	"cpu": {}, "hard_drive": {}, "usb": {}, "bluetooth": {}, "wifi": {}, "router": {}, "plug": {}, "battery_full": {},
	"database": {}, "server_cog": {}, "cloud_gear": {},
	"printer": {}, "projector": {}, "camera": {},
	"multimedia": {}, "video": {}, "clapperboard": {}, "videotape": {},
	"gamepad": {}, "scan_line": {}, "qr_code": {},
	"disc": {}, "notepad": {},
	"toolbox": {}, "wrench": {}, "hammer": {}, "ruler": {}, "scissors": {},
	"microscope": {}, "glasses": {},
	"paintbrush": {}, "palette": {},
	"backpack": {}, "briefcase": {}, "package": {},
	"clipboard_list": {}, "mail": {}, "inbox": {}, "archive": {},
	"compass": {}, "map": {}, "globe": {},
	"flower": {}, "trees": {}, "mountain": {}, "cloud": {}, "sun": {}, "moon": {},
	"umbrella": {}, "sparkles": {}, "droplets": {}, "flame": {},
	"zap": {}, "lightbulb": {}, "fan": {},
	"coffee": {},
}

type CompanyResourceType struct {
	id        uuid.UUID
	companyID uuid.UUID
	name      string
	iconKey   string
	createdAt time.Time
	updatedAt time.Time
}

func NewCompanyResourceType(companyID uuid.UUID, name, iconKey string) (*CompanyResourceType, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidName
	}
	iconKey = strings.TrimSpace(strings.ToLower(iconKey))
	if _, ok := AllowedIconKeys[iconKey]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidIconKey, iconKey)
	}

	now := time.Now().UTC()
	return &CompanyResourceType{
		id:        uuid.New(),
		companyID: companyID,
		name:      name,
		iconKey:   iconKey,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Reconstitute(id, companyID uuid.UUID, name, iconKey string, createdAt, updatedAt time.Time) *CompanyResourceType {
	return &CompanyResourceType{
		id:        id,
		companyID: companyID,
		name:      name,
		iconKey:   iconKey,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (t *CompanyResourceType) Update(name, iconKey string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidName
	}
	iconKey = strings.TrimSpace(strings.ToLower(iconKey))
	if _, ok := AllowedIconKeys[iconKey]; !ok {
		return fmt.Errorf("%w: %s", ErrInvalidIconKey, iconKey)
	}

	t.name = name
	t.iconKey = iconKey
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *CompanyResourceType) ID() uuid.UUID        { return t.id }
func (t *CompanyResourceType) CompanyID() uuid.UUID { return t.companyID }
func (t *CompanyResourceType) Name() string         { return t.name }
func (t *CompanyResourceType) IconKey() string      { return t.iconKey }
func (t *CompanyResourceType) CreatedAt() time.Time { return t.createdAt }
func (t *CompanyResourceType) UpdatedAt() time.Time { return t.updatedAt }
