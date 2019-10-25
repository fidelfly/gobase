package cronx

import "github.com/robfig/cron/v3"

type Option func(*Cronx)

func WithCronOption(opts ...cron.Option) Option {
	return func(cronx *Cronx) {
		if cronx.inCron != nil {
			panic(`"WithCronOption" can be used for one time!`)
		}
		cronx.inCron = cron.New(opts...)
	}
}

func WithMiddleware(m ...JobMiddleware) Option {
	return func(cronx *Cronx) {
		cronx.middlewares = append(cronx.middlewares, m...)
	}
}
