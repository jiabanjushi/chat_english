package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type ChatRoom struct {
	sync.RWMutex
	members map[string][]*User
}

// NewRoom 创建用户集合
func NewRoom() *ChatRoom {
	r := &ChatRoom{
		members: make(map[string][]*User),
	}
	r.cleanRoom()
	return r
}

//清理用户集合
func (r *ChatRoom) cleanRoom() {
	go func() {
		for {
			log.Println("cleanRoom start...")
			r.members = nil
			r.members = make(map[string][]*User)
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

// SendMessageToRoom 发送消息给集合
func (r *ChatRoom) SendMessageToRoom(roomId string, msg []byte) {
	members, _ := r.GetMembers(roomId)
	newMembers := make([]*User, 0)
	for _, member := range members {
		member.Mux.Lock()
		err := member.Conn.WriteMessage(websocket.TextMessage, msg)
		if err == nil {
			newMembers = append(newMembers, member)
		}
		member.Mux.Unlock()
	}
	r.SetMembers(roomId, newMembers)
}

// GetMembers 获取用户集合
func (r *ChatRoom) GetMembers(key string) ([]*User, bool) {
	r.RLock()
	value, ok := r.members[key]
	r.RUnlock()
	return value, ok
}

// SetMembers 设置用户集合
func (r *ChatRoom) SetMembers(key string, value []*User) {
	r.Lock()
	r.members[key] = value
	r.Unlock()
}

//添加用户
func (r *ChatRoom) addMember(key string, user *User) {
	members, _ := r.GetMembers(key)
	members = append(members, user)
	r.SetMembers(key, members)
}

//移除用户
func (r *ChatRoom) removeMember(key string, userId string) {
	members, _ := r.GetMembers(key)
	newMembers := make([]*User, 0)
	for _, member := range members {
		if member.Id != userId {
			newMembers = append(newMembers, member)
		}
	}
	r.SetMembers(key, newMembers)
}
