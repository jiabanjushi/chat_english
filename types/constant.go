package types

type constant struct {
	AccountForbidden int64
	SuperAdminRoleId uint
	EntRoleId        uint
	AgentRoleId      uint
}

var Constant = &constant{
	AccountForbidden: 1,
	SuperAdminRoleId: 1,
	EntRoleId:        2,
	AgentRoleId:      3,
}
