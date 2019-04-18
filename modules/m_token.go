package modules

import (
	"github.com/in2store/service-in2-auth/database"
	"github.com/johnnyeven/libtools/httplib"
	"github.com/johnnyeven/libtools/sqlx"
	"github.com/johnnyeven/libtools/sqlx/builder"
	"github.com/johnnyeven/libtools/timelib"
)

type GetTokensParams struct {
	// 业务ID
	TokenIDs httplib.Uint64List `json:"tokenIDs"`
	// 用户ID
	UserID uint64 `json:"userID,string"`
	// 通道ID
	ChannelID uint64 `json:"channelID,string"`
	httplib.Pager
}

func GetTokens(req GetTokensParams, db *sqlx.DB) (result database.TokenList, count int32, err error) {
	token := &database.Token{}
	t := token.T()

	var condition *builder.Condition
	if req.TokenIDs != nil && len(req.TokenIDs) > 0 {
		condition = builder.And(condition, t.F("TokenID").In(req.TokenIDs))
	}
	if req.UserID != 0 {
		condition = builder.And(condition, t.F("UserID").Eq(req.UserID))
	}
	if req.ChannelID != 0 {
		condition = builder.And(condition, t.F("ChannelID").Eq(req.ChannelID))
	}

	result, count, err = token.FetchList(db, req.Size, req.Offset, condition)
	return
}

type CreateTokenParams struct {
	// TokenID
	TokenID uint64 `json:"tokenID,string"`
	// 用户ID
	UserID uint64 `json:"userID,string"`
	// 通道ID
	ChannelID uint64 `json:"channelID,string"`
	// AccessToken is the token that authorizes and authenticates the requests.
	AccessToken string `json:"accessToken"`
	// TokenType is the type of token. The Type method returns either this or "Bearer", the default.
	TokenType string `json:"tokenType"`
	// RefreshToken is a token that's used by the application(as opposed to the user) to refresh the access token if it expires.
	RefreshToken string `json:"refreshToken"`
	// Expiry is the optional expiration time of the access token.
	ExpiryTime timelib.MySQLTimestamp `json:"expiry"`
}

func CreateToken(req CreateTokenParams, db *sqlx.DB) (*database.Token, error) {
	token := &database.Token{
		TokenID:      req.TokenID,
		UserID:       req.UserID,
		ChannelID:    req.ChannelID,
		AccessToken:  req.AccessToken,
		TokenType:    req.TokenType,
		RefreshToken: req.RefreshToken,
		ExpiryTime:   req.ExpiryTime,
	}
	err := token.Create(db)
	if err != nil {
		return nil, err
	}

	return token, nil
}
