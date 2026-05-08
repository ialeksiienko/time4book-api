package company

type CompanyStatus string

const (
	StatusActive  CompanyStatus = "active"
	StatusBlocked CompanyStatus = "blocked"
)

func (s CompanyStatus) String() string { return string(s) }
