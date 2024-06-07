package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// var stdout *log.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
// var stderr *log.Logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
var stdout *log.Logger = log.New(os.Stdout, "", log.LstdFlags)
var stderr *log.Logger = log.New(os.Stderr, "", log.LstdFlags)

func Info(text string) {
	stdout.Println(getAdditionalInfo() + text)
}

func Error(text string) {
	stderr.Println(getAdditionalInfo() + text)
}

func getAdditionalInfo() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	// Short file name
	fileParts := strings.Split(file, "/")
	// Short function name
	funcParts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	return fmt.Sprintf("%s:%s:%d: ", fileParts[len(fileParts) - 1], funcParts[len(funcParts) - 1], line)
}
