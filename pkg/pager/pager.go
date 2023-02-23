package pager

type Pager struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Total    int `json:"total"`
}

func NewPager(page, pageSize int) *Pager {
	return &Pager{
		Page:     page,
		PageSize: pageSize,
		Total:    0,
	}
}

func (p *Pager) Offset() int {
	if p.Page <= 0 || p.Page >= 1000 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Limit()
}

func (p *Pager) Limit() int {
	if p.PageSize <= 0 || p.PageSize >= 1000 {
		p.PageSize = 20
	}
	return p.PageSize
}
