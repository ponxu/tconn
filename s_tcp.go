package tconn

import (
    "errors"
)

type TCPAdapter struct{}

func (a *TCPAdapter) Accept(s *Server) error {
    return errors.New("TCP not yet implements")
}

func (a *TCPAdapter) Receive(tc *TConn, v *[]byte) error {
    return nil
}

func (a *TCPAdapter) Send(tc *TConn, v []byte) error {
    return nil
}
