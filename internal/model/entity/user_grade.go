package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserGrade 用户等级实体
type UserGrade struct {
	Id                   int         `json:"id"                     orm:"id,primary"                 description:"等级ID"`
	SiteId               int         `json:"site_id"                orm:"site_id"                    description:"站点ID"`
	Name                 string      `json:"name"                   orm:"name"                       description:"等级名称"`
	PointsUpgrade        int         `json:"points_upgrade"         orm:"points_upgrade"             description:"升级所需积分"`
	BonusUpgrade         float64     `json:"bonus_upgrade"          orm:"bonus_upgrade"              description:"升级赠送彩金"`
	BonusBirthday        float64     `json:"bonus_birthday"         orm:"bonus_birthday"             description:"生日彩金"`
	RebatePercentSports  float64     `json:"rebate_percent_sports"  orm:"rebate_percent_sports"      description:"体育返水比例"`
	RebatePercentLottery float64     `json:"rebate_percent_lottery" orm:"rebate_percent_lottery"     description:"彩票返水比例"`
	RebatePercentLive    float64     `json:"rebate_percent_live"    orm:"rebate_percent_live"        description:"真人视讯返水比例"`
	RebatePercentEgame   float64     `json:"rebate_percent_egame"   orm:"rebate_percent_egame"       description:"电子游戏返水比例"`
	RebatePercentPoker   float64     `json:"rebate_percent_poker"   orm:"rebate_percent_poker"       description:"扑克返水比例"`
	FieldsDisable        string      `json:"fields_disable"         orm:"fields_disable"             description:"禁用字段配置"`
	AutoProviding        string      `json:"auto_providing"         orm:"auto_providing"             description:"自动发放配置"`
	Status               int         `json:"status"                 orm:"status"                     description:"状态 1=正常 0=禁用"`
	CreatedAt            *gtime.Time `json:"created_at"             orm:"created_at"                 description:"创建时间"`
	UpdatedAt            *gtime.Time `json:"updated_at"             orm:"updated_at"                 description:"更新时间"`
}
