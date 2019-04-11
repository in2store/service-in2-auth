package authorize

import (
	"context"
	"github.com/google/uuid"
	"github.com/in2store/service-in2-auth/constants/errors"
	"github.com/in2store/service-in2-auth/database"
	"github.com/in2store/service-in2-auth/global"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func init() {
	Router.Register(courier.NewRouter(GetAuthURL{}))
}

// 获取认证界面URL
type GetAuthURL struct {
	httpx.MethodGet
	// ChannelID
	ChannelID uint64 `name:"channelId,string" in:"path"`
}

func (req GetAuthURL) Path() string {
	return "/:channelId"
}

type GetAuthURLResponse struct {
	URL string `json:"url"`
}

func (req GetAuthURL) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.MasterDB.Get()
	channel := &database.Channel{
		ChannelID: req.ChannelID,
	}
	err = channel.FetchByChannelID(db)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.ChannelNotFound
		}
		logrus.Errorf("channel.FetchByChannelID err: %v, request: %+v", err, req)
		return nil, errors.InternalError
	}

	state := uuid.New().String()
	stateMap := &database.StateMap{
		State:     state,
		ChannelID: channel.ChannelID,
	}
	err = stateMap.Create(db)
	if err != nil {
		logrus.Errorf("stateMap.Create err: %v, request: %+v", err, req)
		return nil, errors.InternalError
	}

	conf := oauth2.Config{
		ClientID:     channel.ClientID,
		ClientSecret: channel.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  channel.AuthURL,
			TokenURL: channel.TokenURL,
		},
		Scopes: []string{"repo"},
	}
	url := conf.AuthCodeURL(state)
	return &GetAuthURLResponse{
		URL: url,
	}, nil
}
