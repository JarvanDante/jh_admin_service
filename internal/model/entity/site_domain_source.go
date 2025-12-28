// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteDomainSource is the golang structure for table site_domain_source.
type SiteDomainSource struct {
	Id        int         `json:"id"        orm:"id"         description:"ID"`
	SiteId    int         `json:"siteId"    orm:"site_id"    description:"商户id"`
	Domain    string      `json:"domain"    orm:"domain"     description:"域名"`
	Status    int         `json:"status"    orm:"status"     description:"状态1开0关"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""`
}
