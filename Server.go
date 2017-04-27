package iMQ

import (
	"errors"
	"sync"
)

type Client struct {
	Id   int
	Name string       //can be ip addr
	Ccb  func([]byte,interface{}) //Client call back
	Para interface{}
}

type Topic struct {
	Name string
	Msg  []byte
	Tcb  func([]byte) //topic call back
}

type Server struct {
	Dict map[string]*Channel //map[Channel.Name]*Channel
	sync.RWMutex
}

var Imqsrv Server

func Init() {
	Imqsrv.Dict = make(map[string]*Channel) //所有channel
}

func NewServer() *Server {
	s := &Server{}
	s.Dict = make(map[string]*Channel) //所有channel
	return s
}

//订阅
func (srv *Server) Subscribe(client *Client, channelName string) {

	// 客户是否在Channel的客户列表中
	srv.RLock()
	ch, found := srv.Dict[channelName]
	srv.RUnlock()

	if !found {
		ch = NewChannel(channelName)
		ch.AddClient(client)
		srv.Lock()
		srv.Dict[channelName] = ch
		srv.Unlock()

	} else {
		ch.AddClient(client)
	}

}

//取消订阅
func (srv *Server) Unsubscribe(client *Client, channelName string) {
	srv.RLock()
	ch, found := srv.Dict[channelName]
	srv.RUnlock()
	if found {
		if ch.DeleteClient(client) == 0 {
			ch.Exit()
			srv.Lock()
			delete(srv.Dict, channelName)
			srv.Unlock()
		}
	}
}

//发布消息
func (srv *Server) PublishMessage(channelName string, message []byte) (bool, error) {
	srv.RLock()
	ch, found := srv.Dict[channelName]
	if !found {
		srv.RUnlock()
		return false, errors.New("channelName不存在!")
	}
	srv.RUnlock()

	ch.Notify(message)
	ch.Wait()
	return true, nil
}
