package jse

import (
	"github.com/google/uuid"
	"net"
	"strconv"
)

type Event struct {
	Version      uint
	AppType      AppType
	AppHost      string
	ClientID     uuid.UUID
	ClientIPv4   net.IP
	ProjectTitle string
	ProjectSize  uint32
	ParseError   *ParseError
}

func newEvent1() Event {
	return Event{
		Version: 1,
	}
}

func (e Event) record() []string {
	return []string{
		strconv.Itoa(int(e.AppType)),
		e.AppHost,
	}
}

func (e *Event) Decode(s string) {

}

// func New(b []byte) (e Event, err error) {
// 	err = json.Unmarshal(b, &e)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// func (e Event) ClientIP() net.IP {
// 	ip := net.ParseIP(e.ClientIPv4).To4()
// 	if ip != nil {
// 		return ip
// 	}
// 	return net.ParseIP("0.0.0.0")
// }
