package httprxr

import (
	"errors"
)

type WsProgressHandler WsConnect

func (wph *WsProgressHandler) SendData(msg interface{}) error {
	conn := (*WsConnect)(wph)
	if conn.IsOpen() {
		conn.SendMessage(msg)
	} else {
		return errors.New("websocket connection is not open")
	}
	return nil
}
