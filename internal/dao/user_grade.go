package dao

import (
	"jh_admin_service/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// userGradeDao is the data access object for the table user_grade.
type userGradeDao struct {
	table   string
	group   string
	columns UserGradeColumns
}

// UserGradeColumns 用户等级表字段
type UserGradeColumns struct {
	Id                   string
	SiteId               string
	Name                 string
	PointsUpgrade        string
	BonusUpgrade         string
	BonusBirthday        string
	RebatePercentSports  string
	RebatePercentLottery string
	RebatePercentLive    string
	RebatePercentEgame   string
	RebatePercentPoker   string
	FieldsDisable        string
	AutoProviding        string
	Status               string
	CreatedAt            string
	UpdatedAt            string
}

var (
	// UserGrade is a globally accessible object for table user_grade operations.
	UserGrade = userGradeDao{
		table: "user_grade",
		group: "default",
		columns: UserGradeColumns{
			Id:                   "id",
			SiteId:               "site_id",
			Name:                 "name",
			PointsUpgrade:        "points_upgrade",
			BonusUpgrade:         "bonus_upgrade",
			BonusBirthday:        "bonus_birthday",
			RebatePercentSports:  "rebate_percent_sports",
			RebatePercentLottery: "rebate_percent_lottery",
			RebatePercentLive:    "rebate_percent_live",
			RebatePercentEgame:   "rebate_percent_egame",
			RebatePercentPoker:   "rebate_percent_poker",
			FieldsDisable:        "fields_disable",
			AutoProviding:        "auto_providing",
			Status:               "status",
			CreatedAt:            "created_at",
			UpdatedAt:            "updated_at",
		},
	}
)

// Table 返回表名
func (dao *userGradeDao) Table() string {
	return dao.table
}

// Group 返回数据库组
func (dao *userGradeDao) Group() string {
	return dao.group
}

// Columns 返回字段信息
func (dao *userGradeDao) Columns() UserGradeColumns {
	return dao.columns
}

// Entity 返回实体类型
func (dao *userGradeDao) Entity() *entity.UserGrade {
	return &entity.UserGrade{}
}

// DB 返回数据库连接
func (dao *userGradeDao) DB(dbName ...string) gdb.DB {
	if len(dbName) > 0 {
		return g.DB(dbName[0])
	}
	return g.DB(dao.group)
}
