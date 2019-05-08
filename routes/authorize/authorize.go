package authorize

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/in2store/service-in2-auth/clients/client_in2_user"
	"github.com/in2store/service-in2-auth/constants/errors"
	"github.com/in2store/service-in2-auth/database"
	"github.com/in2store/service-in2-auth/global"
	"github.com/in2store/service-in2-auth/modules"
	"github.com/johnnyeven/eden-library/libModule"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/johnnyeven/libtools/sqlx"
	"github.com/johnnyeven/libtools/timelib"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func init() {
	Router.Register(courier.NewRouter(&Authorize{}))
}

// 处理认证回调
type Authorize struct {
	httpx.MethodGet
	Code    string `name:"code" in:"query"`
	State   string `name:"state" in:"query"`
	cookies *http.Cookie
}

func (req Authorize) Cookies() *http.Cookie {
	return req.cookies
}

func (req Authorize) Path() string {
	return ""
}

func (req *Authorize) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.MasterDB.Get()
	stateMap := &database.StateMap{
		State: req.State,
	}
	err = stateMap.FetchByState(db)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.StateMapNotFound
		}
		logrus.Errorf("[Authorize] stateMap.FetchByState err: %v, request: %+v", err, req)
		return nil, errors.InternalError
	}

	// 获取channel
	channel := &database.Channel{
		ChannelID: stateMap.ChannelID,
	}
	err = channel.FetchByChannelID(db)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.ChannelNotFound
		}
		logrus.Errorf("[Authorize] channel.FetchByChannelID err: %v, request: %+v", err, req)
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
	token, err := conf.Exchange(context.Background(), req.Code)
	if err != nil {
		logrus.Errorf("[Authorize] conf.Exchange err: %v, conf: %+v", err, conf)
		return nil, errors.ExchangeAccessTokenError
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		logrus.Errorf("[Authorize] client.Users.Get err: %v", err)
		return nil, err
	}

	// 根据用户ID、channelID 查询是否存在用户
	var userID uint64
	userRequest := client_in2_user.GetUsersRequest{
		EntryID:   *user.Login,
		ChannelID: channel.ChannelID,
	}
	users, count, err := modules.GetUsers(userRequest, global.Config.ClientUser)
	if err != nil {
		logrus.Errorf("[Authorize] modules.GetUsers err: %v, request: %+v", err, userRequest)
		return nil, err
	}
	if count == 0 {
		// 如果不存在则创建
		request := client_in2_user.CreateUserRequest{
			Body: client_in2_user.CreateUserParams{
				Name: userRequest.EntryID,
				Entries: []client_in2_user.CreateUserParamsEntry{
					{
						ChannelID: userRequest.ChannelID,
						EntryID:   userRequest.EntryID,
					},
				},
			},
		}
		user, err := modules.CreateUser(request, global.Config.ClientUser)
		if err != nil {
			logrus.Errorf("[Authorize] modules.CreateUser err: %v, request: %+v", err, request)
			return nil, err
		}
		userID = user.UserID
	} else {
		userID = users[0].UserID
	}

	// 根据用户ID、channelID 查询是否存在授权信息
	tokenRequest := modules.GetTokensParams{
		UserID:    userID,
		ChannelID: channel.ChannelID,
	}
	_, count, err = modules.GetTokens(tokenRequest, db)
	if err != nil {
		logrus.Errorf("[Authorize] modules.GetTokens err: %v, request: %+v", err, tokenRequest)
		return nil, err
	}
	if count == 0 {
		// 如果不存在则绑定
		tokenID, err := libModule.NewUniqueID(global.Config.ClientID)
		if err != nil {
			logrus.Errorf("[Authorize] libModule.NewUniqueID err: %v", err)
			return nil, err
		}
		request := modules.CreateTokenParams{
			TokenID:      tokenID,
			UserID:       userID,
			ChannelID:    channel.ChannelID,
			AccessToken:  token.AccessToken,
			TokenType:    token.TokenType,
			RefreshToken: token.RefreshToken,
			ExpiryTime:   timelib.MySQLTimestamp(token.Expiry),
		}
		_, err = modules.CreateToken(request, db)
		if err != nil {
			logrus.Errorf("[Authorize] modules.CreateToken err: %v, request: %+v", err, request)
			return nil, err
		}
	}

	// 注册session
	session, err := modules.GetSessionByUserID(userID, db)
	if err != nil {
		logrus.Errorf("[Authorize] modules.GetSessionByUserID err: %v, request: userID=%d", err, userID)
		return nil, err
	}
	if session != nil {
		// 刷新已有的session
		session, err = modules.RefreshSessionID(session.SessionID, db)
		if err != nil {
			logrus.Errorf("[Authorize] modules.RefreshSessionID err: %v, request: sessionID=%d", err, session.SessionID)
			return nil, err
		}
	} else {
		// 创建session
		session, err = modules.CreateSession(userID, db)
		if err != nil {
			logrus.Errorf("[Authorize] modules.CreateSession err: %v, request: userID=%d", err, userID)
			return nil, err
		}
	}

	req.cookies = &http.Cookie{
		Name:   "in2store_auth_token",
		Value:  "INNER:" + session.SessionID,
		Path:   "/",
		Domain: global.Config.AuthRedirectURL,
	}

	return httpx.RedirectWithStatusMovedPermanently(global.Config.AuthRedirectURL), nil
}
