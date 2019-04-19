package database

import (
	"github.com/johnnyeven/libtools/sqlx/presets"
	"github.com/johnnyeven/libtools/timelib"
)

//go:generate libtools gen model Token --database DBIn2Book --table-name t_token --with-comments
// @def primary ID
// @def unique_index U_token_id TokenID
// @def unique_index U_token AccessToken
// @def unique_index U_user_channel UserID ChannelID
type Token struct {
	presets.PrimaryID

	// 业务ID
	TokenID uint64 `json:"tokenID,string" db:"F_token_id" sql:"bigint(64) unsigned NOT NULL"`

	// 用户ID
	UserID uint64 `json:"userID,string" db:"F_user_id" sql:"bigint(64) unsigned NOT NULL"`

	// 通道ID
	ChannelID uint64 `json:"channelID,string" db:"F_channel_id" sql:"bigint(64) unsigned NOT NULL"`

	// AccessToken is the token that authorizes and authenticates the requests.
	AccessToken string `json:"accessToken" db:"F_access_token" sql:"varchar(64) NOT NULL"`

	// TokenType is the type of token. The Type method returns either this or "Bearer", the default.
	TokenType string `json:"tokenType" db:"F_tokenType" sql:"varchar(16) DEFAULT NULL"`

	// RefreshToken is a token that's used by the application(as opposed to the user) to refresh the access token if it expires.
	RefreshToken string `json:"refreshToken" db:"F_refresh_token" sql:"varchar(64) DEFAULT NULL"`

	// Expiry is the optional expiration time of the access token.
	ExpiryTime timelib.MySQLTimestamp `json:"expiry" db:"F_expiry_time" sql:"bigint(64) NOT NULL DEFAULT '0'"`

	presets.OperateTime
	presets.SoftDelete
}
