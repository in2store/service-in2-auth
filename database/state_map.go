package database

import (
	"github.com/johnnyeven/libtools/sqlx/presets"
)

//go:generate libtools gen model StateMap --database DBIn2Book --table-name t_state_map --with-comments
// @def primary ID
// @def unique_index U_state State
type StateMap struct {
	presets.PrimaryID
	// 业务ID
	State string `json:"state" db:"F_state" sql:"char(36) NOT NULL"`
	// ChannelID
	ChannelID uint64 `json:"channelID,string" db:"F_channel_id" sql:"bigint(64) unsigned NOT NULL"`

	presets.OperateTime
	presets.SoftDelete
}
