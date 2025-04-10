package request

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PageInfo struct {
	Page     int `json:"page"`     // 当前页码
	PageSize int `json:"pageSize"` // 每页数量

}

func (p *PageInfo) Verify() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10 // 默认每页10条
	}
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

func NewPageInfo(c *gin.Context) PageInfo {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	return PageInfo{
		Page:     page,
		PageSize: pageSize,
	}
}
