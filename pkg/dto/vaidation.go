package dto

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}
