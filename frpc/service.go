package frpc

import "log"

type Service struct {
}
type Message struct {
	VisitorId, Content string
}
type Reply struct {
	Code, Msg string
}

func (that *Service) SendToVisitor(msg Message, reply *Reply) error {
	log.Println("SendToVisitor:", msg.VisitorId, msg.Content)
	reply.Msg = "ok"
	reply.Code = "200"
	return nil
}
