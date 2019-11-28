package httprxr

import (
	"net/http"
	"sync"
	"time"

	gws "github.com/gorilla/websocket"

	"github.com/fidelfly/gox/logx"
)

const (
	//STANDBY = iota
	OPENED = iota + 1
	CLOSED
)

type WsConnect struct {
	Code             string
	Decoder          WsDecoder
	Encoder          WsEncoder
	Status           uint
	Conn             *gws.Conn
	Duration         time.Duration
	receivers        []WsReceiver
	closeHandlers    []WsCloseHandler
	receiveLock      *sync.Mutex
	closeHandlerLock *sync.Mutex
	writerChan       chan interface{}
}

type WsDecoder func([]byte) (interface{}, error)
type WsEncoder func(interface{}) (int, []byte, error)
type WsCloseHandler func(int, string) error

type WsWriter interface {
	EncodeWsMessage() (int, []byte, error)
}

type WsReceiver func(interface{})

func (wsc *WsConnect) AddReceiver(receiver WsReceiver) {
	wsc.receiveLock.Lock()
	defer wsc.receiveLock.Unlock()
	if wsc.receivers == nil {
		wsc.receivers = make([]WsReceiver, 1)
	}
	wsc.receivers = append(wsc.receivers, receiver)
}

func (wsc *WsConnect) AddCloseHandler(handler WsCloseHandler) {
	wsc.closeHandlerLock.Lock()
	defer wsc.closeHandlerLock.Unlock()
	if wsc.closeHandlers == nil {
		wsc.closeHandlers = make([]WsCloseHandler, 1)
	}
	wsc.closeHandlers = append(wsc.closeHandlers, handler)
}

func (wsc *WsConnect) SendMessage(message interface{}) {
	if wsc.Status == OPENED {
		wsc.writerChan <- message
	}
}

func (wsc *WsConnect) SetupConnection(ws *gws.Conn) {
	wsc.Conn = ws
	wsc.Status = OPENED
	wsc.writerChan = make(chan interface{}, 100)
	wsc.Conn.SetCloseHandler(wsc.onClose)
}

func (wsc *WsConnect) GetStatus() uint {
	return wsc.Status
}

func (wsc *WsConnect) IsOpen() bool {
	return wsc.Status == OPENED
}

func (wsc *WsConnect) ListenAndServe() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		wsc.startReader()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		wsc.startWriter()
	}()

	wg.Wait()
}

func (wsc *WsConnect) onClose(code int, text string) error {
	wsc.Status = CLOSED
	if len(wsc.closeHandlers) > 0 {
		for _, handler := range wsc.closeHandlers {
			err := handler(code, text)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (wsc *WsConnect) notifyReceiver(message interface{}) {
	if len(wsc.receivers) > 0 {
		for _, receiver := range wsc.receivers {
			receiver(message)
		}
	}
}

func (wsc *WsConnect) startReader() {
	if wsc.Status == OPENED {
		for {
			if wsc.Status == OPENED {
				_, p, err := wsc.Conn.ReadMessage()
				if err != nil {
					if gws.IsUnexpectedCloseError(err, gws.CloseGoingAway, gws.CloseAbnormalClosure) {
						logx.Errorf("error: %v", err)
					}
					break
				}
				var message interface{}
				if wsc.Decoder != nil {
					message, err = wsc.Decoder(p)
					if err != nil {
						message = p
					}
				} else {
					message = p
				}

				wsc.notifyReceiver(message)
			} else {
				break
			}
		}
	}
}

func (wsc *WsConnect) sendToReceiver(message interface{}) {
	if wsc.Encoder != nil {
		msgType, data, err := wsc.Encoder(message)
		if err != nil {
			logx.CaptureError(wsc.Conn.WriteMessage(msgType, data))
			return
		}
	} else {
		if encoder, ok := message.(WsWriter); ok {
			msgType, data, err := encoder.EncodeWsMessage()
			if err == nil {
				logx.CaptureError(wsc.Conn.WriteMessage(msgType, data))
				return
			}
		}
		if text, ok := message.(string); ok {
			logx.CaptureError(wsc.Conn.WriteMessage(gws.TextMessage, []byte(text)))
			return
		}
		logx.CaptureError(wsc.Conn.WriteJSON(message))
	}
}

// nolint:gocyclo
func (wsc *WsConnect) startWriter() {
	wsTicker := time.NewTicker(30 * time.Second)
	defer wsTicker.Stop()
	var message interface{}
	if wsc.Duration <= 0 {
		for {
			select {
			case message = <-wsc.writerChan:
				wsc.sendToReceiver(message)
				break
			case <-wsTicker.C: //check ws connection status for every 30 seconds
				if wsc.Status != OPENED {
					return
				}
			}
		}
	} else {
		ticker := time.NewTicker(wsc.Duration)
		defer ticker.Stop()
		for {
			/*			select {
						case message = <-wsc.writerChan:
							break
						default:
							if wsc.Status != OPENED {
								return
							}
						}
			*/
			select {
			case message = <-wsc.writerChan:
				//do nothing
			case <-ticker.C:
				if message != nil {
					wsc.sendToReceiver(message)
					message = nil
				}
			case <-wsTicker.C: //check ws connection status for every 30 seconds
				if wsc.Status != OPENED {
					return
				}
			}
		}
	}

}

//export
func SetupWebsocket(wsc *WsConnect, w http.ResponseWriter, r *http.Request, headers ...map[string]string) (err error) {
	respHeader := http.Header{}
	if len(headers) > 0 {
		for _, header := range headers {
			for key, value := range header {
				respHeader.Set(key, value)
			}
		}
	}

	webSocket, err := upgrader.Upgrade(w, r, respHeader)
	if err != nil {
		return
	}

	wsc.SetupConnection(webSocket)
	return
}

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
