package modules

import (
	"github.com/in2store/service-in2-auth/database"
	"github.com/johnnyeven/libtools/sqlx"
)

type CreateChannelParams struct {
	// 名称
	Name string `json:"name"`
	// ClientID
	ClientID string `json:"clientId"`
	// ClientSecret
	ClientSecret string `json:"clientSecret"`
	// 认证URL
	AuthURL string `json:"authURL"`
	// 交换tokenURL
	TokenURL string `json:"tokenURL"`
}

func CreateChannel(channelID uint64, req CreateChannelParams, db *sqlx.DB) (channel *database.Channel, err error) {
	channel = &database.Channel{
		ChannelID:    channelID,
		Name:         req.Name,
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		AuthURL:      req.AuthURL,
		TokenURL:     req.TokenURL,
	}
	err = channel.Create(db)
	if err != nil {
		return nil, err
	}
	return
}

func GetChannels(db *sqlx.DB) (channelList database.ChannelList, err error) {
	channel := &database.Channel{}
	channelList, _, err = channel.FetchList(db, -1, 0)
	if err != nil {
		return nil, err
	}
	return
}

func GetChannelByChannelID(channelID uint64, db *sqlx.DB, withLock bool) (channel *database.Channel, err error) {
	channel = &database.Channel{
		ChannelID: channelID,
	}
	if withLock {
		err = channel.FetchByChannelIDForUpdate(db)
	} else {
		err = channel.FetchByChannelID(db)
	}
	if err != nil {
		return nil, err
	}
	return
}
