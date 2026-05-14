package companyresourcetypecommands

import "errors"

var ErrTypeInUse = errors.New("company resource type is assigned to existing resources")
