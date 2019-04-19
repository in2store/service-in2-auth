package database

import (
	"github.com/johnnyeven/libtools/sqlx/presets"
	"github.com/johnnyeven/libtools/timelib"
)

//go:generate libtools gen model Session --database DBIn2Book --table-name t_session --with-comments
// @def primary ID
// @def unique_index U_session_id SessionID
// @def index I_user_id UserID
type Session struct {
	presets.PrimaryID
	// 业务ID
	SessionID string `json:"sessionID" db:"F_session_id" sql:"varchar(255) NOT NULL"`
	// 用户ID
	UserID uint64 `json:"userID,string" db:"F_user_id" sql:"bigint(64) unsigned NOT NULL"`
	// 过期时间
	ExpireTime timelib.MySQLTimestamp `json:"expireTime" db:"F_expire_time" sql:"bigint(64) NOT NULL DEFAULT '0'"`

	presets.OperateTime
	presets.SoftDelete
}
