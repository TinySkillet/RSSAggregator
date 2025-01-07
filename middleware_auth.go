package main

import (
	"goapi/auth"
	"goapi/internal/database"
	"net/http"
)

type authedHandler func(http.ResponseWriter, *http.Request, *database.User) *APIError

func (a *APIServer) middlewareAuth(handler authedHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		queries := database.New(a.store.db)

		apiKey, err := auth.GetApiKeyFromHeader(r.Header)
		if err != nil {
			return &APIError{
				Message: err.Error(),
				status:  401,
			}
		}
		user, err := queries.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			return &APIError{
				Message: "User not found!",
				status:  401,
			}
		}
		apiErr := handler(w, r, &user)
		return apiErr
	}
}
