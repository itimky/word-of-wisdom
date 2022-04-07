package tcp

import (
	"net"

	"github.com/panjf2000/gnet/v2"
)

func GetClientIP(conn gnet.Conn) string {
	return conn.RemoteAddr().(*net.TCPAddr).IP.String()
}
