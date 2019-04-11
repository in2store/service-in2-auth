package authorize

import (
	"context"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
)

func init() {
	Router.Register(courier.NewRouter(Authorize{}))
}

// 处理认证回调
type Authorize struct {
	httpx.MethodGet
	Code  string `name:"code" in:"query"`
	State string `name:"state" in:"query"`
}

func (req Authorize) Path() string {
	return ""
}

func (req Authorize) Output(ctx context.Context) (result interface{}, err error) {
	return
}
