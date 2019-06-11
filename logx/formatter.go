package logx

import "github.com/sirupsen/logrus"

type Formatter interface {
	Format(*Entry) ([]byte, error)
}

type logrusFormatter struct {
	Formatter
}

func (lf *logrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return lf.Formatter.Format(&Entry{entry})
}

type formatLogrus struct {
	logrus.Formatter
}

func (fl *formatLogrus) Format(entry *Entry) ([]byte, error) {
	return fl.Formatter.Format(entry.Entry)
}

//export
func ColorFormatter(timeFormat string) Formatter {
	return &formatLogrus{&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: timeFormat,
		ForceColors:     true,
	}}
}

func PlainFormatter(timeFormat string) Formatter {
	return &PlainTextFormatter{
		TimestampFormat: timeFormat,
	}
}

func JSONFormatter(timeFormat string) Formatter {
	return &formatLogrus{&logrus.JSONFormatter{
		TimestampFormat: timeFormat,
	}}
}
