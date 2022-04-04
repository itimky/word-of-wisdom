package testutils

import "net"

//go:generate mockery --name Conn --with-expecter=true

type Conn interface {
	net.Conn
}
