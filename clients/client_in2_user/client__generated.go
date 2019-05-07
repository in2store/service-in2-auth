package client_in2_user

import (
	"fmt"

	github_com_johnnyeven_libtools_courier "github.com/johnnyeven/libtools/courier"
	github_com_johnnyeven_libtools_courier_client "github.com/johnnyeven/libtools/courier/client"
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
)

type ClientIn2UserInterface interface {
	CreateUser(req CreateUserRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *CreateUserResponse, err error)
	GetEntries(req GetEntriesRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetEntriesResponse, err error)
	GetUserByUserID(req GetUserByUserIDRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetUserByUserIDResponse, err error)
	GetUsers(req GetUsersRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetUsersResponse, err error)
}

type ClientIn2User struct {
	github_com_johnnyeven_libtools_courier_client.Client
}

func (ClientIn2User) MarshalDefaults(v interface{}) {
	if cl, ok := v.(*ClientIn2User); ok {
		cl.Name = "in2-user"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c ClientIn2User) Init() {
	c.CheckService()
}

func (c ClientIn2User) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := github_com_johnnyeven_libtools_courier_status_error.FromError(err)
	if statusErr.Code == int64(github_com_johnnyeven_libtools_courier_status_error.RequestTimeout) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}

type CreateUserRequest struct {
	//
	Body CreateUserParams `fmt:"json" in:"body"`
}

func (c ClientIn2User) CreateUser(req CreateUserRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *CreateUserResponse, err error) {
	resp = &CreateUserResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".CreateUser", "POST", "/in2-user/v0/users", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type CreateUserResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body User
}

type GetEntriesRequest struct {
	// 用户ID
	UserID uint64 `in:"path" name:"userID"`
}

func (c ClientIn2User) GetEntries(req GetEntriesRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetEntriesResponse, err error) {
	resp = &GetEntriesResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".GetEntries", "GET", "/in2-user/v0/users/:userID/entries", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetEntriesResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body UserEntryList
}

type GetUserByUserIDRequest struct {
	// 用户ID
	UserID uint64 `in:"path" name:"userID"`
}

func (c ClientIn2User) GetUserByUserID(req GetUserByUserIDRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetUserByUserIDResponse, err error) {
	resp = &GetUserByUserIDResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".GetUserByUserID", "GET", "/in2-user/v0/users/:userID", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetUserByUserIDResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body User
}

type GetUsersRequest struct {
	// 入口ID
	EntryID string `in:"query" name:"entryID"`
	// 通道ID
	ChannelID uint64 `in:"query" name:"channelID"`
	// 分页
	Size int32 `default:"10" in:"query" name:"size,omitempty"`
	// 偏移量
	Offset int32 `default:"0" in:"query" name:"offset,omitempty"`
}

func (c ClientIn2User) GetUsers(req GetUsersRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *GetUsersResponse, err error) {
	resp = &GetUsersResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".GetUsers", "GET", "/in2-user/v0/users", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetUsersResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body GetUsersResult
}

func (c ClientIn2User) Swagger(metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *SwaggerResponse, err error) {
	resp = &SwaggerResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".Swagger", "GET", "/in2-user", nil, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type SwaggerResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body JSONBytes
}
