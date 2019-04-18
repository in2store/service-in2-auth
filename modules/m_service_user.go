package modules

import (
	"github.com/in2store/service-in2-auth/clients/client_in2_user"
	"github.com/sirupsen/logrus"
)

func GetUsers(req client_in2_user.GetUsersRequest, client *client_in2_user.ClientIn2User) (result client_in2_user.UserList, count int32, err error) {
	resp, err := client.GetUsers(req)
	if err != nil {
		logrus.Errorf("client.GetUsers err: %v, request: %+v", err, req)
		return nil, 0, err
	}
	return resp.Body.Data, resp.Body.Total, nil
}

func CreateUser(req client_in2_user.CreateUserRequest, client *client_in2_user.ClientIn2User) (result *client_in2_user.User, err error) {
	resp, err := client.CreateUser(req)
	if err != nil {
		logrus.Errorf("client.CreateUser err: %v, request: %+v", err, req)
		return nil, err
	}
	return &resp.Body, nil
}
