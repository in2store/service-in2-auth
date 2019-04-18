package modules

import (
	"github.com/johnnyeven/eden-library/clients/client_id"
	"github.com/sirupsen/logrus"
)

func NewUniqueID(client *client_id.ClientID) (uint64, error) {
	resp, err := client.GetNewId()
	if err != nil {
		logrus.Errorf("ClientID.GetNewId err: %v", err)
		return 0, err
	}

	return resp.Body.ID, nil
}
