package cronx

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Cronx struct {
	inCron      *cron.Cron
	middlewares []JobMiddleware
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

type Metadata struct {
	meta map[string]string
}

func (md *Metadata) Get(key string) (string, bool) {
	if v, ok := md.meta[key]; ok {
		return v, true
	}
	return "", false
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
	cx := &Cronx{}
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

func (cx *Cronx) AddJob(spec string, job Job, mds ...map[string]string) (int, error) {
	if len(cx.middlewares) > 0 {
		job = AttachMiddleware(job, cx.middlewares...)
	}
	id, err := cx.inCron.AddJob(spec, newCronJob(job, mds...))
	return int(id), err
}

func (cx *Cronx) Remove(id int) {
	cx.inCron.Remove(cron.EntryID(id))
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
	job Job
	md  *Metadata
}

type ctxMetadataKey struct{}

func (cj *cronJob) Run() {
	ctx := context.WithValue(context.Background(), ctxMetadataKey{}, cj.md)
	cj.job.Run(ctx)
}

func newCronJob(job Job, mds ...map[string]string) cron.Job {
	return &cronJob{
		job: job,
		md:  NewMetadata(mds...),
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
