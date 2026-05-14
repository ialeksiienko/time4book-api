package resource

import (
	"time"

	domain "time4book/internal/app/core/domain/model/resource"

	"github.com/google/uuid"
)

type CompanyResourceTypeRef struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	IconKey string    `json:"iconKey"`
}

type ResourceBody struct {
	ID                    uuid.UUID               `json:"id"`
	CompanyID             uuid.UUID               `json:"companyId"`
	Name                  string                  `json:"name"`
	Type                  string                  `json:"type"`
	Description           string                  `json:"description"`
	Location              string                  `json:"location"`
	MaxReservationMinutes *int                    `json:"maxReservationMinutes,omitempty"`
	AvailableFrom         *string                 `json:"availableFrom,omitempty"`
	AvailableTo           *string                 `json:"availableTo,omitempty"`
	ResourceStatus        string                  `json:"resourceStatus"`
	UnavailableFrom       *time.Time              `json:"unavailableFrom,omitempty"`
	UnavailableTo         *time.Time              `json:"unavailableTo,omitempty"`
	UnavailableReason     *string                 `json:"unavailableReason,omitempty"`
	CompanyResourceType   *CompanyResourceTypeRef `json:"companyResourceType,omitempty"`
}

func toResourceBody(r *domain.Resource) ResourceBody {
	var ref *CompanyResourceTypeRef
	if id := r.CompanyResourceTypeID(); id != nil {
		ref = &CompanyResourceTypeRef{
			ID: *id,
		}
		if n := r.CustomTypeName(); n != nil {
			ref.Name = *n
		}
		if k := r.CustomTypeIconKey(); k != nil {
			ref.IconKey = *k
		}
	}
	return ResourceBody{
		ID:                    r.ID(),
		CompanyID:             r.CompanyID(),
		Name:                  r.Name(),
		Type:                  r.Type().String(),
		Description:           r.Description(),
		Location:              r.Location(),
		MaxReservationMinutes: r.MaxReservationMinutes(),
		AvailableFrom:         r.AvailableFrom(),
		AvailableTo:           r.AvailableTo(),
		ResourceStatus:        r.Status().String(),
		UnavailableFrom:       r.UnavailableFrom(),
		UnavailableTo:         r.UnavailableTo(),
		UnavailableReason:     r.UnavailableReason(),
		CompanyResourceType:   ref,
	}
}
