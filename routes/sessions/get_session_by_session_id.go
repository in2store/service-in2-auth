package sessions

import (
	"context"
	"github.com/in2store/service-in2-auth/global"
	"github.com/in2store/service-in2-auth/modules"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.Register(courier.NewRouter(GetSessionBySessionID{}))
}

// 根据SessionID获取session
type GetSessionBySessionID struct {
	httpx.MethodGet
	// SessionID
	SessionID string `name:"sessionID" in:"path"`
}

func (req GetSessionBySessionID) Path() string {
	return "/:sessionID"
}

func (req GetSessionBySessionID) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.SlaveDB.Get()
	session, err := modules.GetSessionBySessionID(req.SessionID, db)
	if err != nil {
		logrus.Errorf("[GetSessionBySessionID] modules.GetSessionBySessionID err: %v, request: %+v", err, req)
		return nil, err
	}
	return session, nil
}
