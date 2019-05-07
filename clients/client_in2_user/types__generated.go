package client_in2_user

import (
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
	github_com_johnnyeven_libtools_courier_swagger "github.com/johnnyeven/libtools/courier/swagger"
	github_com_johnnyeven_libtools_sqlx_presets "github.com/johnnyeven/libtools/sqlx/presets"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

type CreateUserParams struct {
	//
	Entries []CreateUserParamsEntry `json:"entries"`
	//
	Name string `json:"name"`
}

type CreateUserParamsEntry struct {
	// 入口系统通道ID
	ChannelID uint64 `json:"channelID,string"`
	// 入口系统唯一标识
	EntryID string `json:"entryID"`
}

type ErrorField = github_com_johnnyeven_libtools_courier_status_error.ErrorField

type ErrorFields = github_com_johnnyeven_libtools_courier_status_error.ErrorFields

type GetUsersResult struct {
	//
	Data UserList `json:"data"`
	//
	Total int32 `json:"total"`
}

type JSONBytes = github_com_johnnyeven_libtools_courier_swagger.JSONBytes

type MySQLTimestamp = github_com_johnnyeven_libtools_timelib.MySQLTimestamp

type OperateTime = github_com_johnnyeven_libtools_sqlx_presets.OperateTime

type PrimaryID = github_com_johnnyeven_libtools_sqlx_presets.PrimaryID

type SoftDelete = github_com_johnnyeven_libtools_sqlx_presets.SoftDelete

type StatusError = github_com_johnnyeven_libtools_courier_status_error.StatusError

type User struct {
	//
	PrimaryID
	//
	OperateTime
	//
	SoftDelete
	// 用户名
	Name string `json:"name"`
	// 业务ID
	UserID uint64 `json:"userID,string"`
}

type UserEntry struct {
	//
	PrimaryID
	//
	OperateTime
	//
	SoftDelete
	// 入口系统通道ID
	ChannelID uint64 `json:"channelID,string"`
	// 入口系统唯一标识
	EntryID string `json:"entryID"`
	// 业务ID
	UserEntryID uint64 `json:"userEntryID,string"`
	// UserID
	UserID uint64 `json:"userID,string"`
}

type UserEntryList []UserEntry

type UserList []User
