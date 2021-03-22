package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

// Why do the set this way? There is no easy way to list keys. You can check for
// them, so by listing all attributes as a constant and only adding the ones
// that exist.
// https://stackoverflow.com/questions/54926712/is-there-a-way-to-list-keys-in-context-context
var CtxRequestID = "requestId"
var CtxDomain = "domain"
var CtxHandlerMethod = "handlerMethod"
var CtxServiceMethod = "serviceMethod"
var CtxHelpMethods = "helperMethods"
var CtxEmail = "email"
var CtxClientIP = "clientIP"
var loggingAttributes = []string{
	CtxRequestID,
	CtxDomain,
	CtxHandlerMethod,
	CtxServiceMethod,
	CtxEmail,
	CtxClientIP,
	CtxHelpMethods,
}

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func NewLogger() Logger {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "logs.txt"
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(file, os.Stdout)

	infoLogger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger := log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return Logger{
		infoLogger,
		warningLogger,
		errorLogger,
	}
}

func (l *Logger) Info(ctx context.Context, message string) {
	l.infoLogger.Println(l.buildPayload(ctx, message, ""))
}

func (l *Logger) Warning(ctx context.Context, message string, error error) {
	l.warningLogger.Println(l.buildPayload(ctx, message, error.Error()))
}

func (l *Logger) Error(ctx context.Context, message string, error error) {
	l.errorLogger.Println(l.buildPayload(ctx, message, error.Error()))
}

func (l *Logger) buildPayload(ctx context.Context, message string, error string) string {
	payload := fmt.Sprintf("\"message\": \"%s\"", message)
	if error != "" {
		payload += fmt.Sprintf(", \"error\": \"" + error + "\"")
	}
	for _, loggingAttribute := range loggingAttributes {
		value, ok := ctx.Value(loggingAttribute).(string)
		if ok {
			payload += fmt.Sprintf(", \"%s\": \"%s\"", loggingAttribute, value)
		}
	}
	return "{" + payload + "}"
}

// AddToHelperMethods pulls the string list of help methods in the context and appends another value
// to it. This method helps minimize writing this code in multiple spots throughout.
func AddToHelperMethods(ctx context.Context, newMethod string) string {
	helperMethods := ""
	if  ctx.Value(CtxHelpMethods) != nil {
		helperMethods = ctx.Value(CtxHelpMethods).(string)
	}
	if helperMethods != "" {
		helperMethods = fmt.Sprintf("%s, %s", helperMethods, newMethod)
	}
	return helperMethods
}
