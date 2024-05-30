package myerror

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"

	pg "github.com/lib/pq"
)

type MyError struct {
	Inner            error
	Message          string
	SingleStacktrace string
	StackTrace       []string
	Misc             map[string]any
	StatusCode       int
	Custom           bool
}

func (e MyError) Status() int {
	if e.StatusCode == 0 {
		return 500
	}
	return e.StatusCode
}

func WrapError(err error, messagef string, msgArgs ...any) MyError {
	_, currentFile, currentLine, _ := runtime.Caller(1)
	stackTrace := fmt.Sprintf(">>: %s:%d\n", currentFile, currentLine)

	switch err2 := err.(type) {
	case MyError:
		myerror := err2
		message := fmt.Sprintf("%s >> %s", myerror.Message, fmt.Sprintf(messagef, msgArgs...))
		myerror.Message = message
		myerror.SingleStacktrace = stackTrace
		// myerror.StackTrace = append(myerror.StackTrace, stackTrace)
		PrintError(myerror)

		return myerror

	case *pg.Error:
		var myerror MyError
		pgerr := err.(*pg.Error)
		if pgerr.Code == "23505" {
			myerror.Inner = err2
			myerror.Message = pgerr.Message
			// myerror.StackTrace = []string{stackTrace}
			myerror.SingleStacktrace = stackTrace
			PrintError(myerror)
			return myerror
		} else {
			myerror.Inner = err2
			myerror.Message = pgerr.Message
			// myerror.StackTrace = []string{stackTrace}
			myerror.SingleStacktrace = stackTrace
			PrintError(myerror)
			return myerror
		}

	default:
		myerror := MyError{
			Inner:   err,
			Message: fmt.Sprintf(messagef, msgArgs...),
			// StackTrace: []string{stackTrace},
			SingleStacktrace: stackTrace,
			Misc:             make(map[string]any),
		}
		PrintError(myerror)
		return myerror
	}
}

func (err MyError) Error() string {
	return err.Inner.Error()
}

// func Is(err error, )
func New(status_code int, text string, any ...interface{}) error {
	return MyError{
		StatusCode: status_code,
		Inner:      fmt.Errorf(text, any...),
		StackTrace: []string{string(debug.Stack())},
		Custom:     true,
	}
}

// incomingError if you want to match then it should be of
// MyError instance and matchError is simple myerror.New()
// instance what could be any
func Is(incomingError error, matchError error) bool {
	ierr, ok := incomingError.(MyError)
	if !ok {
		// color.Black("ierror error")
		return false
	}

	return ierr.Inner == matchError
}

func PrintError(err MyError) {
	if !err.Custom {
		log.Printf(
			"[[Error:: %s]], [[Message:: %s]], [[Misc:: %#v]], [[StackTrace:: %#v]]",
			err.Error(),
			err.Message,
			err.Misc,
			err.SingleStacktrace,
		)
	}
}
