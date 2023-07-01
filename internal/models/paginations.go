package models

import "gorm.io/gorm"

type Pagination struct {
	Limit   int         `form:"limit,omitempty;query:limit"`
	Page    int         `form:"page,omitempty;query:page"`
	Total   int64       `form:"total"`
	Results interface{} `form:"results"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.Getlimit()
}

func (p *Pagination) Getlimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(*gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	pagination.Total = totalRows
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.Getlimit())
	}
}
