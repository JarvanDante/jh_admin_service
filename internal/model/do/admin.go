// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Admin is the golang structure of table admin for DAO operations like Where/Data.
type Admin struct {
	g.Meta             `orm:"table:admin, do:true"`
	Id                 any         //
	SiteId             any         //
	Username           any         //
	Nickname           any         //
	Password           any         //
	AdminRoleId        any         //
	Status             any         //
	SwitchGoogle2Fa    any         // 二次验证开关。1=打开；0=关闭
	Google2FaSecret    any         // 二次验证密钥
	TransferAuditSound any         // 转账.审核提示声音控制。0=关闭 1=播放一次；2=循环播放
	SoundLoopTime      any         // 声音循环时间 单位秒
	PaymentSound       any         // 第三方支付提示声音控制。0=关闭 1=播放一次
	LastLoginIp        any         //
	LastLoginTime      *gtime.Time //
	CreatedAt          *gtime.Time //
	UpdatedAt          *gtime.Time //
	DeleteAt           *gtime.Time //
}
