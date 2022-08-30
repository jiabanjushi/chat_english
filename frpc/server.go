package frpc

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func NewRpcServer(address string) {
	rpc.Register(new(Service))
	sock, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("rpc listen error:", err)
	}
	log.Println("start rpc server " + address)
	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Println(" rpc server error", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
