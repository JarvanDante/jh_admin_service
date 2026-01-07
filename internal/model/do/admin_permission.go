// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminPermission is the golang structure of table admin_permission for DAO operations like Where/Data.
type AdminPermission struct {
	g.Meta      `orm:"table:admin_permission, do:true"`
	Id          any         //
	ParentId    any         // 父级id
	Type        any         // 权限类型；1=菜单；2=操作权限
	Name        any         // 权限名称
	BackendUrl  any         // 后端url
	FrontendUrl any         // 前端url
	Status      any         // 状态。1=可用；0=禁用
	Sort        any         // 排序。值越小，越靠前
	Icon        any         // 菜单图标
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
