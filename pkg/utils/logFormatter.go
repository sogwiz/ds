package utils

import (
	"bytes"
	"ds/pkg/color"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogFormatter ...
type LogFormatter struct{}

// Format ...
func (f LogFormatter) Format(e *logrus.Entry) ([]byte, error) {
	skip := 6
	var fn string
	var line int
	for {
		_, fn, line, _ = runtime.Caller(skip)
		if !strings.Contains(fn, "/logrus/") || skip >= 10 {
			break
		}
		skip++
	}
	var buffer bytes.Buffer
	var level string
	switch e.Level {
	case logrus.DebugLevel:
		level = color.Magenta("DEBU")
	case logrus.InfoLevel:
		level = color.Cyan("INFO")
	case logrus.WarnLevel:
		level = color.Yellow("WARN")
	case logrus.ErrorLevel:
		level = color.Red("ERRO")
	case logrus.FatalLevel:
		level = color.Red("FATA")
	case logrus.PanicLevel:
		level = color.Red("PANI")
	}
	repopath := fmt.Sprintf("%s/src/ds", os.Getenv("GOPATH"))
	filename := strings.Replace(fn, repopath, "", -1)
	buffer.WriteString(e.Time.Format("15:04:05"))
	buffer.WriteString(" ")
	buffer.WriteString(level)
	buffer.WriteString(" ")
	buffer.WriteString(color.Magenta("[" + filename + ":" + strconv.Itoa(line) + "]"))
	buffer.WriteString(" ")
	buffer.WriteString(e.Message)
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}
