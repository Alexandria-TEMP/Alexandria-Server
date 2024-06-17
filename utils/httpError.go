package utils

type HTTPError struct {
	Error string `json:"error" example:"failed to something: reason"`
}
