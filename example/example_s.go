package main

import (
    "fmt"
    . "github.com/ponxu/tconn"
)

type MyHandler struct{}

func (handler *MyHandler) OnConnected(tc *TConn) {
    tc.Tag = "xxx"
    fmt.Println("OnConnected:", tc.Tag, tc)
}
func (handler *MyHandler) OnClosed(tc *TConn) {
    fmt.Println("OnClosed:", tc.Tag, tc)
}
func (handler *MyHandler) OnReceived(tc *TConn, data []byte) {
    fmt.Println("OnReceived:", tc.Tag, tc, string(data))

    if "STOP" == string(data) {
        tc.Close()
    }
}

func main() {
    s := &Server{
        Protocol: ProtocolWebsocket,
        Addr:     ":8080",
        Handler:  &MyHandler{},
    }
    s.Start()
}
