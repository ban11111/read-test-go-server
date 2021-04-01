package mysql

import (
	"gorm.io/gorm"
)

const (
	defaultPerPage     = 10
	defaultIdFieldName = "id"
)

func NewListPage(lp ListPage) ListPageEmbed {
	return ListPageEmbed{PageSetting: lp}
}

func NewPageResult(result ListPageResult) PageResultEmbed {
	return PageResultEmbed{PageResult: result}
}

type ListPageEmbed struct {
	PageSetting ListPage `json:"page_setting"`
}

type PageResultEmbed struct {
	PageResult ListPageResult `json:"page_result"`
}

type ListPage struct {
	PageNo   uint   `json:"page_no" name:"页码"`
	PageSize int    `json:"page_size" name:"每页数量, ps:如果传-1则获取全部数据"`
	LastId   uint   `json:"last_id" name:"上一轮查询的最后一个id(要么传 page_no 要么传last_id, 两个都传的话取 page_no)"`
	idField  string // 指定 id field name
	noCount  bool   // 是否返回 总页数, 总条数
	desc     bool   // 是否倒叙
}

func (p *ListPage) Desc() *ListPage {
	p.desc = true
	return p
}

func (p *ListPage) Asc() *ListPage {
	p.desc = false
	return p
}

func (p *ListPage) NoCount() *ListPage {
	p.noCount = true
	return p
}

func (p *ListPage) Count() *ListPage {
	p.noCount = false
	return p
}

func (p *ListPage) Id(idFieldName string) *ListPage {
	p.idField = idFieldName
	return p
}

type ListPageResult struct {
	TotalCount uint  `json:"total_count" name:"总条数"`
	TotalPage  uint  `json:"total_page" name:"总页数"`
	CountErr   error `json:"-"`
	Err        error `json:"-"`
}

func (p *ListPage) Find(sql *gorm.DB, result interface{}) *ListPageResult {
	if p.PageNo > 0 {
		return p.FindByPage(sql, result)
	}
	return p.FindByLastId(sql, result)
}

func (p *ListPage) Scan(sql *gorm.DB, result interface{}, countSql ...*gorm.DB) *ListPageResult {
	if p.PageNo > 0 {
		return p.ScanByPage(sql, result, countSql...)
	}
	return p.ScanByLastId(sql, result, countSql...)
}

func (p *ListPage) FindByPage(sql *gorm.DB, result interface{}) *ListPageResult {
	r := new(ListPageResult)
	r.TotalPage, r.TotalCount, r.CountErr, r.Err = FindDataByPageAndPerPage(sql, int(p.PageNo), p.PageSize, result, p.noCount)
	return r
}

func (p *ListPage) FindByLastId(sql *gorm.DB, result interface{}) *ListPageResult {
	r := new(ListPageResult)
	r.TotalPage, r.TotalCount, r.CountErr, r.Err = FindDataByLastIdAndPerPage(sql, p.idField, p.LastId, p.PageSize, result, p.noCount, p.desc)
	return r
}

func (p *ListPage) ScanByPage(sql *gorm.DB, result interface{}, countSql ...*gorm.DB) *ListPageResult {
	r := new(ListPageResult)
	r.TotalPage, r.TotalCount, r.CountErr, r.Err = ScanDataByPageAndPerPage(sql, int(p.PageNo), p.PageSize, result, p.noCount, countSql...)
	return r
}

func (p *ListPage) ScanByLastId(sql *gorm.DB, result interface{}, countSql ...*gorm.DB) *ListPageResult {
	r := new(ListPageResult)
	r.TotalPage, r.TotalCount, r.CountErr, r.Err = ScanDataByLastIdAndPerPage(sql, p.idField, p.LastId, p.PageSize, result, p.noCount, p.desc, countSql...)
	return r
}

// 获取offset
func SetOffsetAndLimit(page, perPage int) func(*gorm.DB) *gorm.DB {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(perPage)
	}
}

// 获取分页页数总数及数据列表。m是查询的表的model，result是列表结果传指针进来！
func FindDataByPageAndPerPage(db *gorm.DB, page, perPage int, result interface{}, noCount bool) (totalPages, totalCount uint, countErr, err error) {
	if perPage == 0 {
		perPage = defaultPerPage
	}
	if perPage < 0 {
		sql := db.Find(result)
		return 1, uint(sql.RowsAffected), nil, sql.Error
	}
	if err = db.Scopes(SetOffsetAndLimit(page, perPage)).Find(result).Error; err != nil {
		return
	}
	if noCount {
		return
	}

	sql := db.Offset(-1).Limit(-1)
	var tmpTotalCount int64
	if countErr = sql.Count(&tmpTotalCount).Error; countErr != nil {
		return
	}
	if perPage <= 0 {
		return
	}
	totalCount = uint(tmpTotalCount)
	totalPages = (totalCount + uint(perPage) - 1) / uint(perPage)
	return
}

func FindDataByLastIdAndPerPage(db *gorm.DB, idField string, lastId uint, perPage int, result interface{}, noCount, desc bool) (totalPages, totalCount uint, countErr, err error) {
	if perPage == 0 {
		perPage = defaultPerPage
	}
	if perPage < 0 {
		sql := db.Find(result)
		return 1, uint(sql.RowsAffected), nil, sql.Error
	}
	if idField == "" {
		idField = defaultIdFieldName
	}
	sql := db
	if lastId > 0 {
		if desc {
			sql = sql.Where(idField+"<?", lastId)
		} else {
			sql = sql.Where(idField+">?", lastId)
		}
	}
	if err = sql.Limit(perPage).Find(result).Error; err != nil {
		return
	}
	if noCount {
		return
	}

	sql = sql.Offset(-1).Limit(-1)
	var tmpTotalCount int64
	if countErr = sql.Count(&tmpTotalCount).Error; countErr != nil {
		return
	}
	if perPage <= 0 {
		return
	}
	totalCount = uint(tmpTotalCount)
	totalPages = (totalCount + uint(perPage) - 1) / uint(perPage)
	return
}

func ScanDataByPageAndPerPage(db *gorm.DB, page, perPage int, result interface{}, noCount bool, countSql ...*gorm.DB) (totalPages, totalCount uint, countErr, err error) {
	if perPage == 0 {
		perPage = defaultPerPage
	}
	if perPage < 0 {
		sql := db.Scan(result)
		return 1, uint(sql.RowsAffected), nil, sql.Error
	}

	if err = db.Scopes(SetOffsetAndLimit(page, perPage)).Scan(result).Error; err != nil {
		return
	}

	if noCount {
		return
	}

	var tmpTotalCount int64
	if len(countSql) > 0 && countSql[0] != nil {
		if countErr = countSql[0].Count(&tmpTotalCount).Error; countErr != nil {
			return
		}
	} else {
		sql := db.Offset(-1).Limit(-1)
		if countErr = sql.Count(&tmpTotalCount).Error; countErr != nil {
			return
		}
	}
	totalCount = uint(tmpTotalCount)

	if perPage <= 0 {
		return
	}
	totalPages = (totalCount + uint(perPage) - 1) / uint(perPage)
	return
}

func ScanDataByLastIdAndPerPage(db *gorm.DB, idField string, lastId uint, perPage int, result interface{}, noCount, desc bool, countSql ...*gorm.DB) (totalPages, totalCount uint, countErr, err error) {
	if perPage == 0 {
		perPage = defaultPerPage
	}
	if perPage < 0 {
		sql := db.Scan(result)
		return 1, uint(sql.RowsAffected), nil, sql.Error
	}
	if idField == "" {
		idField = defaultIdFieldName
	}

	sql := db
	if lastId > 0 {
		if desc {
			sql = sql.Where(idField+"<?", lastId)
		} else {
			sql = sql.Where(idField+">?", lastId)
		}
	}

	if err = sql.Limit(perPage).Scan(result).Error; err != nil {
		return
	}

	if noCount {
		return
	}

	var tmpTotalCount int64
	if len(countSql) > 0 && countSql[0] != nil {
		if countErr = countSql[0].Count(&tmpTotalCount).Error; countErr != nil {
			return
		}
	} else {
		sql = sql.Offset(-1).Limit(-1)
		if countErr = db.Count(&tmpTotalCount).Error; countErr != nil {
			return
		}
	}
	totalCount = uint(tmpTotalCount)

	if perPage <= 0 {
		return
	}
	totalPages = (totalCount + uint(perPage) - 1) / uint(perPage)
	return
}
