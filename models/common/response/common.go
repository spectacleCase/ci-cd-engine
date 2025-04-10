package response

import "math"

type PageInfo struct {
	Page      int `json:"page"`       // 当前页码
	PageSize  int `json:"page_size"`  // 每页数量
	Total     int `json:"total"`      // 总记录数
	TotalPage int `json:"total_page"` // 总页数
}

func (p *PageInfo) Verify() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10 // 默认每页10条
	}
}

// SetTotal 设置总数并计算总页数
func (p *PageInfo) SetTotal(total int) {
	p.Verify()
	p.Total = total
	p.TotalPage = int(math.Ceil(float64(total) / float64(p.PageSize)))
}

// Offset 计算偏移量
func (p *PageInfo) Offset() int {
	p.Verify()
	return (p.Page - 1) * p.PageSize
}

// Limit 获取每页数量
func (p *PageInfo) Limit() int {
	p.Verify()
	return p.PageSize
}
