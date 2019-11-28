package progx

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/fidelfly/gox/logx"
)

const (
	ProgressActive    = "active"
	ProgressException = "exception"
	ProgressSuccess   = "success"
)

type ProgressGetter interface {
	GetPercent() int
	GetStatus() string
	GetMessage() interface{}
}

type ProgressSetter interface {
	ProgressGetter
	Set(percent int, status string, message ...interface{})
	update(percent int, status string, message ...interface{})
}

type ProgressSubscriber interface {
	ProgressSet(percent int, status string, messages ...interface{})
}

type ProgressSuperior interface {
	ProgressChanged(subProgress *SubProgress)
}

type SubProgress struct {
	superior    ProgressSuperior
	Proportion  int
	Code        string
	Message     interface{}
	Percent     int
	Status      string
	Propagation bool
}

func (sp *SubProgress) ProgressSet(percent int, status string, message ...interface{}) {
	sp.Percent = percent
	sp.Status = status
	if len(message) > 0 {
		sp.Message = message[0]
	}

	sp.superior.ProgressChanged(sp)
}

func (sp *SubProgress) IsDone() bool {
	return sp.Percent >= 100
}

//export
func NewProgressDispatcher(code string, subscriber ...ProgressSubscriber) *ProgressDispatcher {
	return &ProgressDispatcher{Code: code, Subscribers: subscriber}
}

// Progress Dispatcher
type ProgressDispatcher struct {
	Code        string
	Subscribers []ProgressSubscriber
	message     interface{}
	percent     int
	status      string
	auto        *AutoProgress
	sub         []*SubProgress
	mux         sync.Mutex
	data        map[string]interface{}
	notifyLock  sync.Mutex
}

func (pd *ProgressDispatcher) GetPercent() int {
	return pd.percent
}
func (pd *ProgressDispatcher) GetStatus() string {
	return pd.status
}
func (pd *ProgressDispatcher) GetMessage() interface{} {
	return pd.message
}

func (pd *ProgressDispatcher) SetStatus(status string, message ...interface{}) {
	pd.Set(pd.percent, status, message...)
}

func (pd *ProgressDispatcher) Exception(percent int, message ...interface{}) {
	pd.Set(percent, ProgressException, message...)
}

func (pd *ProgressDispatcher) Active(percent int, message ...interface{}) {
	pd.Set(percent, ProgressActive, message...)
}

func (pd *ProgressDispatcher) Success(message ...interface{}) {
	if len(message) == 0 {
		pd.Set(100, ProgressSuccess, "")
	} else {
		pd.Set(100, ProgressSuccess, message...)
	}
}

func (pd *ProgressDispatcher) Done(status string, message ...interface{}) {
	if len(status) == 0 {
		if pd.status == ProgressActive {
			status = ProgressSuccess
		} else {
			status = pd.status
		}
	}
	if len(message) == 0 {
		pd.Set(100, status, "")
	} else {
		pd.Set(100, status, message...)
	}
}

func (pd *ProgressDispatcher) notifySubscriber() {
	pd.notify(pd.percent, pd.status, pd.message)
}

func (pd *ProgressDispatcher) updateData(percent int, status string, message interface{}) bool {
	dataChange := false
	if pd.data == nil {
		pd.data = make(map[string]interface{})
	}

	if pd.data["percent"] != percent {
		dataChange = true
		pd.data["percent"] = percent
	}

	if pd.data["status"] != status {
		pd.data["status"] = status
	}

	if pd.data["message"] != message {
		pd.data["message"] = message
	}

	return dataChange
}
func (pd *ProgressDispatcher) notify(percent int, status string, message interface{}) {
	pd.notifyLock.Lock()
	defer pd.notifyLock.Unlock()
	if percent < 0 {
		percent = pd.percent
	}
	if len(status) == 0 {
		status = pd.status
	}
	if message == nil {
		message = pd.message
	}

	if !pd.updateData(percent, status, message) {
		return
	}
	if len(pd.Subscribers) > 0 {
		for _, subscriber := range pd.Subscribers {
			subscriber.ProgressSet(percent, status, message)
		}
	} else {
		msg := ""
		if message != nil {
			if msgText, ok := message.(string); ok {
				msg = msgText
			} else if msgData, err := json.Marshal(message); err == nil {
				msg = string(msgData)
			}
		}
		logx.Infof("Progress(%s) : percent = %d%%, status = %s, message = %s", pd.Code, percent, status, msg)
	}
}
func (pd *ProgressDispatcher) update(percent int, status string, message ...interface{}) {
	pd.percent = percent
	pd.status = status
	if len(message) > 0 {
		pd.message = message[0]
	}
	pd.notifySubscriber()
}
func (pd *ProgressDispatcher) Set(percent int, status string, message ...interface{}) {
	if pd.auto != nil {
		pd.auto.Stop()
		pd.auto = nil
	}
	pd.update(percent, status, message...)
}

func (pd *ProgressDispatcher) AutoProgress(stepValue int, duration time.Duration, maxValue int, message ...interface{}) {
	if len(message) > 0 {
		pd.message = message[0]
		pd.notifySubscriber()
	}

	if pd.auto != nil {
		pd.auto.Stop()
		pd.auto = nil
	}
	pd.auto = newAutoProgress(pd, stepValue, duration, maxValue)
	pd.auto.Start()
}

func (pd *ProgressDispatcher) Step(stepValue int, message ...interface{}) {
	pd.Set(pd.percent+stepValue, ProgressActive, message...)
}

func (pd *ProgressDispatcher) NewSubProgress(proportion int) *SubProgress {
	pd.mux.Lock()
	defer pd.mux.Unlock()
	sp := &SubProgress{superior: pd, Proportion: proportion}
	pd.sub = append(pd.sub, sp)
	return sp
}

// nolint[:gocyclo,dupl]
func (pd *ProgressDispatcher) ProgressChanged(subProgress *SubProgress) {
	pd.mux.Lock()
	defer pd.mux.Unlock()
	if pd.sub != nil && len(pd.sub) > 0 {
		subValue := 0
		index := -1
		for i, sp := range pd.sub {
			if sp == subProgress && sp.IsDone() {
				index = i
				pd.percent += sp.Proportion
			} else {
				value := 0
				if sp.IsDone() {
					value = sp.Proportion
				} else {
					value = sp.Proportion * sp.Percent / 100
				}
				subValue += value
			}
		}

		if index >= 0 {
			newSub := make([]*SubProgress, 0)
			if index > 0 {
				newSub = append(newSub, pd.sub[:index]...)
			}
			if index < len(pd.sub)-1 {
				newSub = append(newSub, pd.sub[index+1:]...)
			}
			pd.sub = newSub
		}

		percent := pd.percent + subValue
		msg := pd.message
		if subProgress.Message != nil {
			msg = subProgress.Message
		}

		if subProgress.Propagation && subProgress.Status == ProgressException {
			pd.status = subProgress.Status
		}
		pd.notify(percent, pd.status, msg)
	}

}

//Struct AutoProgress
type AutoProgress struct {
	progress  ProgressSetter
	stepValue int
	maxValue  int
	duration  time.Duration
	ticker    *time.Ticker
}

func (ap *AutoProgress) Start() {
	ap.ticker = time.NewTicker(ap.duration)
	go func() {
		defer ap.ticker.Stop()
		for range ap.ticker.C {
			percent := ap.progress.GetPercent() + ap.stepValue
			if percent > ap.maxValue {
				percent = ap.maxValue
			}
			ap.progress.update(percent, ap.progress.GetStatus(), ap.progress.GetMessage())
			if percent >= ap.maxValue {
				return
			}
		}
	}()
}

func (ap *AutoProgress) Stop() {
	if ap.ticker != nil {
		ap.ticker.Stop()
	}
	if ap.progress.GetPercent() < ap.maxValue {
		ap.progress.update(ap.maxValue, ap.progress.GetStatus(), ap.progress.GetMessage())
	}
}

func newAutoProgress(progress ProgressSetter, stepValue int, duration time.Duration, maxValue int) *AutoProgress {
	return &AutoProgress{progress: progress, stepValue: stepValue, duration: duration, maxValue: maxValue}
}

//Core Struct : Progress
type Progress struct {
	//ws           *httprxr.WsConnect
	handler      ProgressHandler
	Code         string
	Message      interface{}
	Percent      int
	Status       string
	auto         *AutoProgress
	sub          []*SubProgress
	mux          sync.Mutex
	data         map[string]interface{}
	delayed      bool
	delayMessage bool
	senderLock   sync.Mutex
}

type ProgressHandler interface {
	SendData(msg interface{}) error
}

func NewProgress(handler ProgressHandler, code string) *Progress {
	return &Progress{handler: handler, Code: code}
}

/*func NewWsProgress(ws *httprxr.WsConnect, code string) *Progress {
	return &Progress{ws: ws, Code: code}
}*/

func (p *Progress) GetPercent() int {
	return p.Percent
}

func (p *Progress) GetStatus() string {
	return p.Status
}
func (p *Progress) GetMessage() interface{} {
	return p.Message
}

func (p *Progress) SetStatus(status string, message ...interface{}) {
	p.Set(p.Percent, status, message...)
}

func (p *Progress) Exception(percent int, message ...interface{}) {
	p.Set(percent, ProgressException, message...)
}

func (p *Progress) Active(percent int, message ...interface{}) {
	p.Set(percent, ProgressActive, message...)
}

func (p *Progress) Success(message ...interface{}) {
	p.Set(100, ProgressSuccess, message...)
}

func (p *Progress) Done(status string, message ...interface{}) {
	p.Set(100, status, message...)
}

func (p *Progress) update(percent int, status string, message ...interface{}) {
	p.Percent = percent
	p.Status = status
	if len(message) > 0 {
		p.Message = message[0]
	}

	p.SendMsg()
}

func (p *Progress) Set(percent int, status string, message ...interface{}) {
	if p.auto != nil {
		p.auto.Stop()
		p.auto = nil
	}
	p.update(percent, status, message...)
}

func (p *Progress) updateData(percent int, status string, message interface{}) (dataChange bool) {
	newData := map[string]interface{}{
		"code":    p.Code,
		"percent": percent,
		"status":  status,
		"message": message,
	}

	if p.data == nil {
		dataChange = true
	} else {
		dataChange = newData["percent"] != p.data["percent"] ||
			newData["status"] != p.data["status"] ||
			newData["message"] != p.data["message"]
	}

	if dataChange {
		p.data = newData
	}

	return dataChange
}

func (p *Progress) delaySend(timer *time.Timer) {
	p.delayed = true
	go func() {
		defer timer.Stop()
		for _ = range timer.C {
			p.senderLock.Lock()
			if p.delayMessage {
				if err := p.handler.SendData(p.data); err == nil {
					timer.Reset(100 * time.Millisecond)
					p.delayMessage = false
					p.senderLock.Unlock()
				} else {
					p.delayed = false
					p.delayMessage = false
					p.senderLock.Unlock()
					return
				}
			} else {
				p.delayed = false
				p.delayMessage = false
				p.senderLock.Unlock()
				return
			}
		}
		/*		for {
				select {
				case <-timer.C:
					p.senderLock.Lock()
					if p.delayMessage {
						if err := p.handler.SendData(p.data); err == nil {
							timer.Reset(100 * time.Millisecond)
							p.delayMessage = false
							p.senderLock.Unlock()
						} else {
							p.delayed = false
							p.delayMessage = false
							p.senderLock.Unlock()
							return
						}
					} else {
						p.delayed = false
						p.delayMessage = false
						p.senderLock.Unlock()
						return
					}
				}
			}*/

	}()
}

func (p *Progress) Send(percent int, status string, message interface{}) {
	p.senderLock.Lock()
	defer p.senderLock.Unlock()
	if !p.updateData(percent, status, message) {
		return
	}
	if !p.delayed {
		if err := p.handler.SendData(p.data); err == nil {
			p.delaySend(time.NewTimer(100 * time.Millisecond))
		} else {
			msg := ""
			if message != nil {
				if msgText, ok := message.(string); ok {
					msg = msgText
				} else if msgData, err := json.Marshal(message); err == nil {
					msg = string(msgData)
				}
			}
			logx.Infof("Progress(%s) : percent = %d%%, status = %s, message = %s", p.Code, percent, status, msg)
		}
	} else {
		p.delayMessage = true
	}
}

func (p *Progress) SendMsg() {
	p.Send(p.Percent, p.Status, p.Message)
	/*if p.ws != nil && p.ws.IsOpen() {
		p.ws.SendMessage(p)
	} else {
		msg := ""
		if p.Message != nil {
			if msgText, ok := p.Message.(string); ok {
				msg = msgText
			} else {
				if msgData, err := json.Marshal(p.Message); err == nil {
					msg = string(msgData)
				}
			}
		}
		logrus.Infof("Progress(%s) : percent = %d%%, status = %s, message = %s", p.Code, p.Percent, p.Status, msg)
	}*/
}

func (p *Progress) AutoProgress(stepValue int, duration time.Duration, maxValue int, message ...interface{}) {
	if len(message) > 0 {
		p.Message = message[0]
		p.SendMsg()
	}

	if p.auto != nil {
		p.auto.Stop()
		p.auto = nil
	}

	p.auto = newAutoProgress(p, stepValue, duration, maxValue)
	p.auto.Start()
}

func (p *Progress) Step(stepValue int, message ...interface{}) {
	p.Set(p.Percent+stepValue, ProgressActive, message...)
}

func (p *Progress) NewSubProgress(proportion int) *SubProgress {
	p.mux.Lock()
	defer p.mux.Unlock()
	sp := &SubProgress{superior: p, Proportion: proportion}
	if p.sub == nil {
		p.sub = make([]*SubProgress, 0)
	}
	p.sub = append(p.sub, sp)
	return sp
}

// nolint[:gocyclo,dupl]
func (p *Progress) ProgressChanged(subProgress *SubProgress) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.sub != nil && len(p.sub) > 0 {
		subValue := 0
		index := -1
		for i, sp := range p.sub {
			if sp == subProgress && sp.IsDone() {
				index = i
				p.Percent += sp.Proportion
			} else {
				value := 0
				if sp.IsDone() {
					value = sp.Proportion
				} else {
					value = sp.Proportion * sp.Percent / 100
				}
				subValue += value
			}
		}

		if index >= 0 {
			newSub := make([]*SubProgress, 0)
			if index > 0 {
				newSub = append(newSub, p.sub[:index]...)
			}
			if index < len(p.sub)-1 {
				newSub = append(newSub, p.sub[index+1:]...)
			}
			p.sub = newSub
		}

		percent := p.Percent + subValue
		msg := p.Message
		if subProgress.Message != nil {
			msg = subProgress.Message
		}

		if subProgress.Propagation && subProgress.Status == ProgressException {
			p.Status = subProgress.Status
		}
		p.Send(percent, p.Status, msg)
	}

}
