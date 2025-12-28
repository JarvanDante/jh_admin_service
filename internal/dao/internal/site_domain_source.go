// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteDomainSourceDao is the data access object for the table site_domain_source.
type SiteDomainSourceDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  SiteDomainSourceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// SiteDomainSourceColumns defines and stores column names for the table site_domain_source.
type SiteDomainSourceColumns struct {
	Id        string // ID
	SiteId    string // 商户id
	Domain    string // 域名
	Status    string // 状态1开0关
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// siteDomainSourceColumns holds the columns for the table site_domain_source.
var siteDomainSourceColumns = SiteDomainSourceColumns{
	Id:        "id",
	SiteId:    "site_id",
	Domain:    "domain",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewSiteDomainSourceDao creates and returns a new DAO object for table data access.
func NewSiteDomainSourceDao(handlers ...gdb.ModelHandler) *SiteDomainSourceDao {
	return &SiteDomainSourceDao{
		group:    "default",
		table:    "site_domain_source",
		columns:  siteDomainSourceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteDomainSourceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteDomainSourceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteDomainSourceDao) Columns() SiteDomainSourceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteDomainSourceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteDomainSourceDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SiteDomainSourceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
