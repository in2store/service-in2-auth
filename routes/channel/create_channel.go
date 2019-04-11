package channel

import (
	"context"
	"github.com/in2store/service-in2-auth/database"
	"github.com/in2store/service-in2-auth/global"
	"github.com/johnnyeven/eden-library/modules"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/service-account/constants/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.Register(courier.NewRouter(CreateChannel{}))
}

// 创建通道
type CreateChannel struct {
	httpx.MethodPost
	Body CreateChannelBody `name:"body" in:"body"`
}

type CreateChannelBody struct {
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

func (req CreateChannel) Path() string {
	return ""
}

func (req CreateChannel) Output(ctx context.Context) (result interface{}, err error) {
	id, err := modules.NewUniqueID(global.Config.ClientID)
	if err != nil {
		logrus.Errorf("modules.NewUniqueID err: %v", err)
		return nil, errors.InternalError.StatusError().WithMsg("底层服务异常，请稍后再试").WithErrTalk()
	}

	db := global.Config.MasterDB.Get()
	channel := &database.Channel{
		ChannelID:    id,
		Name:         req.Body.Name,
		ClientID:     req.Body.ClientID,
		ClientSecret: req.Body.ClientSecret,
		AuthURL:      req.Body.AuthURL,
		TokenURL:     req.Body.TokenURL,
	}
	err = channel.Create(db)
	if err != nil {
		logrus.Errorf("channel.Create err: %v, request: %+v", err, req.Body)
		return nil, errors.InternalError
	}

	return channel, nil
}
