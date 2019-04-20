package channels

import (
	"context"
	"github.com/in2store/service-in2-auth/global"
	"github.com/in2store/service-in2-auth/modules"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.Register(courier.NewRouter(GetChannels{}))
}

// 获取通道列表
type GetChannels struct {
	httpx.MethodGet
}

func (req GetChannels) Path() string {
	return ""
}

func (req GetChannels) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.SlaveDB.Get()
	result, err = modules.GetChannels(db)
	if err != nil {
		logrus.Errorf("[GetChannels] modules.GetChannels err: %v", err)
	}
	return
}
