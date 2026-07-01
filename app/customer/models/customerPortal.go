package models

import commonModels "baize/app/common/commonModels"

type CustomerPortalMenuDQL struct {
	MenuName string `form:"menuName" db:"menu_name"`
	Visible  string `form:"visible" db:"visible"`
	Status   string `form:"status" db:"status"`
	commonModels.BaseEntityDQL
}

type CustomerPortalMenuDML struct {
	MenuId    int64  `json:"menuId,string" db:"menu_id"`
	ParentId  int64  `json:"parentId,string" db:"parent_id"`
	MenuName  string `json:"menuName" db:"menu_name"`
	OrderNum  string `json:"orderNum" db:"order_num"`
	Path      string `json:"path" db:"path"`
	Component string `json:"component" db:"component"`
	MenuType  string `json:"menuType" db:"menu_type"`
	Visible   string `json:"visible" db:"visible"`
	Status    string `json:"status" db:"status"`
	Perms     string `json:"perms" db:"perms"`
	Icon      string `json:"icon" db:"icon"`
	Remark    string `json:"remark" db:"remark"`
	commonModels.BaseEntityDML
}

type CustomerPortalMenuVo struct {
	MenuId    int64                 `json:"menuId,string" db:"menu_id"`
	ParentId  int64                 `json:"parentId" db:"parent_id"`
	MenuName  string                `json:"menuName" db:"menu_name"`
	OrderNum  string                `json:"orderNum" db:"order_num"`
	Path      string                `json:"path" db:"path"`
	Component string                `json:"component" db:"component"`
	MenuType  string                `json:"menuType" db:"menu_type"`
	Visible   string                `json:"visible" db:"visible"`
	Status    string                `json:"status" db:"status"`
	Perms     string                `json:"perms" db:"perms"`
	Icon      string                `json:"icon" db:"icon"`
	Remark    *string               `json:"remark" db:"remark"`
	Children  []*CustomerPortalMenuVo `json:"children,omitempty"`
	commonModels.BaseEntity
}

type CustomerPortalRoleDQL struct {
	RoleName  string `form:"roleName" db:"role_name"`
	RoleKey   string `form:"roleKey" db:"role_key"`
	Status    string `form:"status" db:"status"`
	BeginTime string `form:"beginTime" db:"begin_time"`
	EndTime   string `form:"endTime" db:"end_time"`
	commonModels.BaseEntityDQL
}

type CustomerPortalRoleDML struct {
	RoleId   int64    `json:"roleId,string" db:"role_id"`
	RoleName string   `json:"roleName" db:"role_name"`
	RoleKey  string   `json:"roleKey" db:"role_key"`
	RoleSort int      `json:"roleSort" db:"role_sort"`
	Status   string   `json:"status" db:"status"`
	Remark   string   `json:"remark" db:"remark"`
	MenuIds  []string `json:"menuIds"`
	commonModels.BaseEntityDML
}

type CustomerPortalRoleVo struct {
	RoleId   int64   `json:"roleId,string" db:"role_id"`
	RoleName string  `json:"roleName" db:"role_name"`
	RoleKey  string  `json:"roleKey" db:"role_key"`
	RoleSort int     `json:"roleSort" db:"role_sort"`
	Status   string  `json:"status" db:"status"`
	DelFlag  string  `json:"delFlag" db:"del_flag"`
	Remark   *string `json:"remark" db:"remark"`
	commonModels.BaseEntity
}

type CustomerPortalRoleOptionVo struct {
	RoleId   int64  `json:"roleId,string" db:"role_id"`
	RoleName string `json:"roleName" db:"role_name"`
}

type CustomerPortalRoleMenu struct {
	RoleId int64 `db:"role_id"`
	MenuId int64 `db:"menu_id"`
}

type CustomerPortalAccountRole struct {
	AccountId int64 `db:"account_id"`
	RoleId    int64 `db:"role_id"`
}

type CustomerPortalRouteMeta struct {
	Title  string `json:"title"`
	Icon   string `json:"icon,omitempty"`
	MenuId int64  `json:"menuId,string"`
}

type CustomerPortalRoute struct {
	Name      string                 `json:"name"`
	Path      string                 `json:"path"`
	Component string                 `json:"component,omitempty"`
	Hidden    bool                   `json:"hidden,omitempty"`
	Meta      CustomerPortalRouteMeta `json:"meta"`
	Children  []*CustomerPortalRoute `json:"children,omitempty"`
}

type CustomerPortalProfile struct {
	User        *CustomerAccountVo `json:"user"`
	Roles       []string           `json:"roles"`
	Permissions []string           `json:"permissions"`
}
