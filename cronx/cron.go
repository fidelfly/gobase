package cronx

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/fidelfly/gox/pkg/randx"
)

type Cronx struct {
	inCron      *cron.Cron
	middlewares []JobMiddleware
	keyMap      map[string]int
}

type Job interface {
	Run(ctx context.Context)
}

type FuncJob func(ctx context.Context)

func (fj FuncJob) Run(ctx context.Context) {
	fj(ctx)
}

type JobMiddleware func(Job) Job

type Entry cron.Entry

func (e Entry) Valid() bool { return e.ID != 0 }
func (e Entry) Meta() *Metadata {
	if e.Job != nil {
		if cj, ok := e.Job.(*cronJob); ok {
			return cj.md
		}
	}
	return NewMetadata()
}
func (e Entry) Key() string {
	md := e.Meta()
	if md != nil {
		return md.GetJobKey()
	}
	return ""
}

type Metadata struct {
	meta map[string]string
}

func (md *Metadata) Get(key string) (string, bool) {
	if v, ok := md.meta[key]; ok {
		return v, true
	}
	return "", false
}

func (md *Metadata) GetJobKey() string {
	key, _ := md.Get(metaJobKey)
	return key
}

func NewMetadata(datas ...map[string]string) *Metadata {
	md := make(map[string]string)
	for _, mp := range datas {
		if len(mp) > 0 {
			for k, v := range mp {
				md[k] = v
			}
		}
	}
	return &Metadata{md}
}

func New(opts ...Option) *Cronx {
	cx := &Cronx{keyMap: make(map[string]int)}
	for _, opt := range opts {
		opt(cx)
	}

	if cx.inCron == nil {
		cx.inCron = cron.New()
	}
	return cx
}

func (cx *Cronx) AddFunc(spec string, cmd func(context.Context), mds ...map[string]string) (int, error) {
	return cx.AddJob(spec, FuncJob(cmd), mds...)
}

const uuidSeed = "job.uuid"

func (cx *Cronx) AddJob(spec string, job Job, mds ...map[string]string) (int, error) {
	if len(cx.middlewares) > 0 {
		job = AttachMiddleware(job, cx.middlewares...)
	}
	jobKey := randx.GenUUID(uuidSeed)
	id, err := cx.inCron.AddJob(spec, newCronJob(jobKey, job, mds...))
	cx.keyMap[jobKey] = int(id)
	return int(id), err
}

func (cx *Cronx) removeTimerJob(job Job) Job {
	return FuncJob(func(ctx context.Context) {
		md := GetMetadata(ctx)
		job.Run(ctx)
		if jobKey := md.GetJobKey(); len(jobKey) > 0 {
			if id, ok := cx.keyMap[jobKey]; ok {
				go cx.Remove(id)
			}
		}
	})
}

func (cx *Cronx) AddTimerFunc(t time.Time, cmd func(context.Context), mds ...map[string]string) int {
	return cx.AddTimerJob(t, FuncJob(cmd), mds...)
}

func (cx *Cronx) AddTimerJob(t time.Time, job Job, mds ...map[string]string) int {
	if len(cx.middlewares) > 0 {
		job = cx.removeTimerJob(AttachMiddleware(job, cx.middlewares...))
	}
	jobKey := randx.GenUUID(uuidSeed)
	schedule := NewTimerSchedule(t)
	id := cx.inCron.Schedule(schedule, newCronJob(jobKey, job, mds...))
	cx.keyMap[jobKey] = int(id)
	return int(id)
}

func (cx *Cronx) Remove(id int) {
	cx.inCron.Remove(cron.EntryID(id))
	for k, v := range cx.keyMap {
		if v == id {
			delete(cx.keyMap, k)
			return
		}
	}
}

func (cx *Cronx) Start() {
	cx.inCron.Start()
}

func (cx *Cronx) Stop() context.Context {
	return cx.inCron.Stop()
}

func (cx *Cronx) Entry(id int) Entry {
	return Entry(cx.inCron.Entry(cron.EntryID(id)))
}

func (cx *Cronx) Entries() []Entry {
	list := cx.inCron.Entries()
	entries := make([]Entry, len(list))
	for i := range list {
		entries[i] = Entry(list[i])
	}
	return entries
}

func AttachMiddleware(job Job, ms ...JobMiddleware) Job {
	for i := len(ms) - 1; i >= 0; i-- {
		job = ms[i](job)
	}
	return job
}

type cronJob struct {
	key string
	job Job
	md  *Metadata
}

type ctxMetadataKey struct{}

func (cj *cronJob) Run() {
	ctx := context.WithValue(context.Background(), ctxMetadataKey{}, cj.md)
	cj.job.Run(ctx)
}

const metaJobKey = "job.meta.key"

func newCronJob(jobKey string, job Job, mds ...map[string]string) cron.Job {
	jobMD := NewMetadata(mds...)
	jobMD.meta[metaJobKey] = jobKey
	return &cronJob{
		key: jobKey,
		job: job,
		md:  jobMD,
	}
}

func GetMetadata(ctx context.Context) *Metadata {
	if v := ctx.Value(ctxMetadataKey{}); v != nil {
		if md, ok := v.(*Metadata); ok {
			return md
		}
	}
	return nil
}

type TimerSchedule struct {
	t time.Time
}

func (ts TimerSchedule) Next(now time.Time) time.Time {
	if ts.t.Before(now) {
		return time.Time{}
	}
	return ts.t
}

func NewTimerSchedule(t time.Time) *TimerSchedule {
	return &TimerSchedule{t: t}
}
