package goetna

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Info(format string, v ...any)
	Debug(format string, v ...any)
	Error(format string, v ...any)
	Fatal(format string, v ...any)
}

type StdLogger struct {
	infoLogger  *log.Logger
	debugLogger *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	name        string
}

func (l StdLogger) Info(format string, v ...any)  { l.infoLogger.Printf(format, v...) }
func (l StdLogger) Debug(format string, v ...any) { l.debugLogger.Printf(format, v...) }
func (l StdLogger) Error(format string, v ...any) { l.errorLogger.Printf(format, v...) }
func (l StdLogger) Fatal(format string, v ...any) { l.fatalLogger.Fatalf(format, v...) }

func ColouredLogger(name string) Logger {
	const colGreen = "\033[1;32m"
	const colTeal = "\033[1;36m"
	const colRed = "\033[1;31m"
	const colEnd = "\033[0m"

	return StdLogger{
		infoLogger: log.New(os.Stdout, fmt.Sprintf("%sINFO  %s%s ", colGreen, name, colEnd),
			log.Lmsgprefix|log.LstdFlags|log.Lmicroseconds),
		debugLogger: log.New(os.Stdout, fmt.Sprintf("%sDEBUG %s%s ", colTeal, name, colEnd),
			log.Lmsgprefix|log.LstdFlags|log.Lmicroseconds),
		errorLogger: log.New(os.Stderr, fmt.Sprintf("%sERROR %s%s ", colRed, name, colEnd),
			log.Lmsgprefix|log.LstdFlags|log.Lmicroseconds),
		fatalLogger: log.New(os.Stderr, fmt.Sprintf("%sFATAL %s%s ", colRed, name, colEnd),
			log.Lmsgprefix|log.LstdFlags|log.Lmicroseconds),
	}
}
