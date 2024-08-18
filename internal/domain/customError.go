package domain

import (
	"log"
	"net/http"
	"runtime"
	"sync"
)

var (
	errorCounter = 0
	mu           sync.Mutex
)

type CustomError struct {
	code        int
	httpStatus  int
	messageUser string
	messageLog  string
	file        string
	line        int
	funcErr     string
}

type Option func(customError *CustomError)

func (e *CustomError) Error() string {
	return e.messageUser
}

func (e *CustomError) GetHttpStatus() int {
	return e.httpStatus
}

func (e *CustomError) GetCode() int {
	return e.code

}

func (e *CustomError) GetUserLog() string {
	return e.messageUser
}

func NewCustomError(opts ...Option) *CustomError {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	e := &CustomError{
		file:    file,
		line:    line,
		funcErr: funcName,
	}

	for _, opt := range opts {
		opt(e)
	}

	log.Printf("%d Error: %s, in %s, at %s:%d\n", e.code, e.messageLog, e.funcErr, e.file, e.line)

	return e
}

func InvalidInputError() Option {
	return func(customError *CustomError) {
		customError.messageUser = "невалидные данные ввода"
		customError.httpStatus = http.StatusBadRequest
	}
}

func AuthError() Option {
	return func(customError *CustomError) {
		customError.messageUser = "неавторизованный доступ"
		customError.httpStatus = http.StatusUnauthorized
	}
}

func NoFoundUserError() Option {
	return func(customError *CustomError) {
		customError.messageUser = "пользователь не найден"
		customError.httpStatus = http.StatusNotFound
	}
}

func InternalError(err error) Option {
	return func(customError *CustomError) {
		mu.Lock()
		errorCounter++
		customError.code = errorCounter
		mu.Unlock()
		customError.messageUser = "ошибка сервера"
		customError.messageLog = err.Error()
		customError.httpStatus = http.StatusInternalServerError
	}
}
