// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteDomainSource is the golang structure of table site_domain_source for DAO operations like Where/Data.
type SiteDomainSource struct {
	g.Meta    `orm:"table:site_domain_source, do:true"`
	Id        any         // ID
	SiteId    any         // 商户id
	Domain    any         // 域名
	Status    any         // 状态1开0关
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
