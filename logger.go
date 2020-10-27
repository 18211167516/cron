package cron

import (
	"os"
	"fmt"
	"io"
	"time"
)

var defaultFormat = func(Values ...interface{}) string {
	var formattedArgs []interface{}
	for _, arg := range Values {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format("2006-01-02 15:04:05")
		}
		formattedArgs = append(formattedArgs, arg," ")
	}
	return fmt.Sprint(formattedArgs...)
}

// DefaultLogger 默认Logger
var DefaultLogger Logger = PrintfLogger(os.Stdout,defaultFormat)
// LogFormatter log格式化func
type LogFormatter func(Values ...interface{}) string
// PrintLogger ...
type PrintLogger struct {
	out io.Writer
	Format LogFormatter
}

func PrintfLogger(o io.Writer,f LogFormatter) Logger {
	return PrintLogger{o, f}
}

func (l PrintLogger) Info(values ...interface{}) {
	fmt.Fprintln(l.out,"[Cron]",l.Format(values...))
}

func (l PrintLogger) Error(err error,values ...interface{}) {
	fmt.Fprintln(l.out,"[Error]",err, l.Format(values...))
}

// Logger is Log interface
type Logger interface {
	Info(values ...interface{})
	Error(err error,value ...interface{})
}


