// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Admin is the golang structure for table admin.
type Admin struct {
	Id                 uint        `json:"id"                 orm:"id"                   description:""`
	SiteId             int         `json:"siteId"             orm:"site_id"              description:""`
	Username           string      `json:"username"           orm:"username"             description:""`
	Nickname           string      `json:"nickname"           orm:"nickname"             description:""`
	Password           string      `json:"password"           orm:"password"             description:""`
	AdminRoleId        int         `json:"adminRoleId"        orm:"admin_role_id"        description:""`
	Status             int         `json:"status"             orm:"status"               description:""`
	SwitchGoogle2Fa    int         `json:"switchGoogle2Fa"    orm:"switch_google2fa"     description:"二次验证开关。1=打开；0=关闭"`
	Google2FaSecret    string      `json:"google2FaSecret"    orm:"google2fa_secret"     description:"二次验证密钥"`
	TransferAuditSound int         `json:"transferAuditSound" orm:"transfer_audit_sound" description:"转账.审核提示声音控制。0=关闭 1=播放一次；2=循环播放"`
	SoundLoopTime      string      `json:"soundLoopTime"      orm:"sound_loop_time"      description:"声音循环时间 单位秒"`
	PaymentSound       int         `json:"paymentSound"       orm:"payment_sound"        description:"第三方支付提示声音控制。0=关闭 1=播放一次"`
	LastLoginIp        string      `json:"lastLoginIp"        orm:"last_login_ip"        description:""`
	LastLoginTime      *gtime.Time `json:"lastLoginTime"      orm:"last_login_time"      description:""`
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:""`
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:""`
	DeleteAt           *gtime.Time `json:"deleteAt"           orm:"delete_at"            description:""`
}
