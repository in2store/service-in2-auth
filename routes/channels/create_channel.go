package channels

import (
	"context"
	"github.com/in2store/service-in2-auth/global"
	"github.com/in2store/service-in2-auth/modules"
	"github.com/johnnyeven/eden-library/libModule"
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
	Body modules.CreateChannelParams `name:"body" in:"body"`
}

func (req CreateChannel) Path() string {
	return ""
}

func (req CreateChannel) Output(ctx context.Context) (result interface{}, err error) {
	id, err := libModule.NewUniqueID(global.Config.ClientID)
	if err != nil {
		logrus.Errorf("modules.NewUniqueID err: %v", err)
		return nil, errors.InternalError.StatusError().WithMsg("底层服务异常，请稍后再试").WithErrTalk()
	}

	db := global.Config.MasterDB.Get()
	channel, err := modules.CreateChannel(id, req.Body, db)
	if err != nil {
		logrus.Errorf("channel.Create err: %v, request: %+v", err, req.Body)
		return nil, errors.InternalError
	}

	return channel, nil
}
