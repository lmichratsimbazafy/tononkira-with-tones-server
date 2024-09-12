package helpers

import (
	"fmt"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ErrorWrapper(err error) error {
	err = errors.Wrap(err, "ERROR")
	if err != nil {
		if err, ok := err.(stackTracer); ok {
			fmt.Printf("%v", err)
			fmt.Println()
			for _, f := range err.StackTrace() {
				fmt.Printf("%+s:%d\n", f, f)
			}
		}
	}
	return err
}
