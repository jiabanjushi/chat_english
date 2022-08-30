package frpc

import (
	"log"
	"net/rpc/jsonrpc"
)

func ClientRpc() {
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		log.Fatal(err)
	}
	var msg = Message{"ssss", "hello"}
	var res Reply
	errCall := conn.Call("Service.SendToVisitor", msg, &res)
	log.Println(res)
	if errCall != nil {
		log.Fatal("Service.SendToVisitor error:", errCall)
	}
}
