package sql

type Pages struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type PageResult struct {
	Page   *Pages      `json:"page"`
	Result interface{} `json:"result"`
}

func (p *Pages) Offset() (offset int) {
	if p.Page > 0 {
		offset = (p.Page - 1) * p.PerPage
	}
	return
}

func (p *Pages) TotalPage() (totalPage int) {
	if p.Total == 0 || p.PerPage == 0 {
		totalPage = 0
	}
	totalPage = p.Total / p.PerPage
	if p.Total%p.PerPage > 0 {
		totalPage = totalPage + 1
	}
	return
}
