package database

import (
	"github.com/johnnyeven/libtools/sqlx/presets"
)

//go:generate libtools gen model Channel --database DBIn2Book --table-name t_channel --with-comments
// @def primary ID
// @def unique_index U_channel_id ChannelID
type Channel struct {
	presets.PrimaryID
	// 业务ID
	ChannelID uint64 `json:"channelID,string" db:"F_channel_id" sql:"bigint(64) unsigned NOT NULL"`
	// 名称
	Name string `json:"name" db:"F_name" sql:"varchar(32) NOT NULL"`
	// ClientID
	ClientID string `json:"clientID" db:"F_client_id" sql:"varchar(255) NOT NULL"`
	// ClientSecret
	ClientSecret string `json:"clientSecret" db:"F_client_secret" sql:"varchar(255) NOT NULL"`
	// 认证URL
	AuthURL string `json:"authURL" db:"F_auth_url" sql:"varchar(255) DEFAULT NULL"`
	// 交换tokenURL
	TokenURL string `json:"tokenURL" db:"F_token_url" sql:"varchar(255) DEFAULT NULL"`
	// raw文件访问URL
	RawURL string `json:"rawURL" db:"F_raw_url" sql:"varchar(255) DEFAULT NULL"`

	presets.OperateTime
	presets.SoftDelete
}
