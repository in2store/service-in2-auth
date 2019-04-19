package channels

import "github.com/johnnyeven/libtools/courier"

var Router = courier.NewRouter(ChannelGroup{})

type ChannelGroup struct {
	courier.EmptyOperator
}

func (ChannelGroup) Path() string {
	return "/channels"
}
