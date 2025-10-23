package models

import "fmt"

type ErrorUnprocessableEntity struct {
	Message string `json:"message"`
}

func (e ErrorUnprocessableEntity) Error() string {
	return e.Message
}

type ErrorNotFound struct {
	Message string `json:"message"`
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("Not found - %s", e.Message)
}

type ErrorGeneric struct {
	Message string `json:"message"`
}

func (e ErrorGeneric) Error() string {
	return e.Message
}

type ErrorBadRequest struct {
	Message string `json:"message"`
}

func (e ErrorBadRequest) Error() string {
	return e.Message
}
