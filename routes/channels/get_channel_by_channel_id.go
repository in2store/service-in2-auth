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
	Router.Register(courier.NewRouter(GetChannelByChannelID{}))
}

// 根据通道ID获取通道信息
type GetChannelByChannelID struct {
	httpx.MethodGet
	// 通道ID
	ChannelID uint64 `name:"channelID,string" in:"path"`
}

func (req GetChannelByChannelID) Path() string {
	return "/:channelID"
}

func (req GetChannelByChannelID) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.SlaveDB.Get()
	result, err = modules.GetChannelByChannelID(req.ChannelID, db, false)
	if err != nil {
		logrus.Errorf("[GetChannelByChannelID] modules.GetChannelByChannelID err: %v, request: %d", err, req.ChannelID)
	}
	return
}
