package tokens

import (
	"context"
	"github.com/in2store/service-in2-auth/database"
	"github.com/in2store/service-in2-auth/global"
	"github.com/in2store/service-in2-auth/modules"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/httplib"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.Register(courier.NewRouter(GetTokens{}))
}

// 根据用户ID和通道ID获取token
type GetTokens struct {
	httpx.MethodGet
	// 用户ID
	UserID uint64 `name:"userID" in:"query"`
	// 通道ID
	ChannelID uint64 `name:"channelID" in:"query"`
	httplib.Pager
}

func (req GetTokens) Path() string {
	return ""
}

type GetTokensResult struct {
	Data  database.TokenList `json:"data"`
	Total int32              `json:"total"`
}

func (req GetTokens) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.SlaveDB.Get()
	request := modules.GetTokensParams{
		UserID:    req.UserID,
		ChannelID: req.ChannelID,
		Pager: httplib.Pager{
			Size:   req.Size,
			Offset: req.Offset,
		},
	}
	tokens, count, err := modules.GetTokens(request, db)
	if err != nil {
		logrus.Errorf("[GetTokens] modules.GetTokens err: %v, request: %+v", err, req)
		return nil, err
	}

	return GetTokensResult{
		Data:  tokens,
		Total: count,
	}, nil
}
