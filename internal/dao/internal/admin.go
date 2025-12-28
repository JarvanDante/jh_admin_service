// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminDao is the data access object for the table admin.
type AdminDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  AdminColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// AdminColumns defines and stores column names for the table admin.
type AdminColumns struct {
	Id                 string //
	SiteId             string //
	Username           string //
	Nickname           string //
	Password           string //
	AdminRoleId        string //
	Status             string //
	SwitchGoogle2Fa    string // 二次验证开关。1=打开；0=关闭
	Google2FaSecret    string // 二次验证密钥
	TransferAuditSound string // 转账.审核提示声音控制。0=关闭 1=播放一次；2=循环播放
	SoundLoopTime      string // 声音循环时间 单位秒
	PaymentSound       string // 第三方支付提示声音控制。0=关闭 1=播放一次
	LastLoginIp        string //
	LastLoginTime      string //
	CreatedAt          string //
	UpdatedAt          string //
	DeleteAt           string //
}

// adminColumns holds the columns for the table admin.
var adminColumns = AdminColumns{
	Id:                 "id",
	SiteId:             "site_id",
	Username:           "username",
	Nickname:           "nickname",
	Password:           "password",
	AdminRoleId:        "admin_role_id",
	Status:             "status",
	SwitchGoogle2Fa:    "switch_google2fa",
	Google2FaSecret:    "google2fa_secret",
	TransferAuditSound: "transfer_audit_sound",
	SoundLoopTime:      "sound_loop_time",
	PaymentSound:       "payment_sound",
	LastLoginIp:        "last_login_ip",
	LastLoginTime:      "last_login_time",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeleteAt:           "delete_at",
}

// NewAdminDao creates and returns a new DAO object for table data access.
func NewAdminDao(handlers ...gdb.ModelHandler) *AdminDao {
	return &AdminDao{
		group:    "default",
		table:    "admin",
		columns:  adminColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AdminDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AdminDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AdminDao) Columns() AdminColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AdminDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AdminDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AdminDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
