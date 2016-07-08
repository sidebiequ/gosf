package gosf

import "fmt"

// Logger inteface
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

var defaultLogger Logger = &logger{}

type logger struct{}

func (l *logger) Print(v ...interface{}) {
	fmt.Println(v...)
}

func (l *logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}
