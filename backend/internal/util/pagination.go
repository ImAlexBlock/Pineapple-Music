package util

// PaginationParams holds pagination query parameters.
type PaginationParams struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// Normalize sets defaults and caps for pagination.
func (p *PaginationParams) Normalize() {
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 50
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
