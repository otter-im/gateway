package handler

import (
	"fmt"
	"github.com/go-oauth2/oauth2/v4/errors"
)

func InternalErrorHandler(err error) (re *errors.Response) {
	fmt.Println("Internal error:", err.Error())
	return
}

func ResponseErrorHandler(re *errors.Response) {
	fmt.Println("Response error:", re.Error.Error())
}
