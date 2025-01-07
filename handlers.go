package main

import (
	"encoding/json"
	"goapi/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type apiHandler func(http.ResponseWriter, *http.Request) *APIError

func makeHTTPHandler(a apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiErr := a(w, r); apiErr != nil {
			if apiErr.status > 499 {
				// log the reason behind the 5XX error, but don't write it to w
				log.Printf("Internal Error: %v", apiErr.Message)
				http.Error(w, http.StatusText(apiErr.status), apiErr.status)
			}
			WriteJSON(w, apiErr.status, apiErr)
		}
	}
}

func (a *APIServer) handleReadiness(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, 200, struct{}{})
}

func (a *APIServer) handleError(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, 400, struct{}{})
}

func (a *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) *APIError {
	queries := database.New(a.store.db)
	users, err := queries.GetUsers(r.Context())
	if err != nil {
		return &APIError{
			Message: "Unable to fetch users - Error: " + err.Error(),
			status:  500,
		}
	}
	WriteJSON(w, 200, dbUsersToUsers(users))
	return nil
}

func (a *APIServer) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request, user *database.User) *APIError {
	WriteJSON(w, 200, dbUserToPrivateUser(user))
	return nil
}

func (a *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) *APIError {
	decoder := json.NewDecoder(r.Body)
	var params CreateUserParams
	if err := decoder.Decode(&params); err != nil {
		return &APIError{
			Message: "Invalid JSON: unable to parse request body!",
			status:  400,
		}
	}
	if params.Name == "" {
		return &APIError{
			Message: "Missing required field: name!",
			status:  400,
		}
	}
	queries := database.New(a.store.db)
	user, err := queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		return &APIError{
			Message: "Unable to create user: " + err.Error(),
			status:  400,
		}
	}
	WriteJSON(w, 201, dbUserToPrivateUser(&user))
	return nil
}

func (a *APIServer) handleCreateFeed(w http.ResponseWriter, r *http.Request, user *database.User) *APIError {
	decoder := json.NewDecoder(r.Body)
	var params CreateFeedParams
	if err := decoder.Decode(&params); err != nil {
		return &APIError{
			Message: "Invalid JSON: unable to parse request body!",
			status:  400,
		}
	}

	if params.Name == "" || params.Url == "" {
		return &APIError{
			Message: "name and url are required fields!",
			status:  400,
		}
	}

	queries := database.New(a.store.db)
	feed, err := queries.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		return &APIError{
			Message: "Unable to create feed: " + err.Error(),
			status:  400,
		}
	}
	WriteJSON(w, 201, dbFeedToFeed(&feed))
	return nil
}

func (a *APIServer) handleGetFeeds(w http.ResponseWriter, r *http.Request) *APIError {
	queries := database.New(a.store.db)
	feeds, err := queries.GetFeeds(r.Context())
	if err != nil {
		return &APIError{
			Message: "Unable to fetch feeds: " + err.Error(),
			status:  400,
		}
	}
	WriteJSON(w, 200, dbFeedsToFeeds(feeds))
	return nil
}

func (a *APIServer) handleGetPostsForUser(w http.ResponseWriter, r *http.Request, user *database.User) *APIError {
	queries := database.New(a.store.db)
	posts, err := queries.GetPostsForUsers(r.Context(), database.GetPostsForUsersParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		return &APIError{
			Message: "Unable to fetch posts for user" + err.Error(),
			status:  400,
		}
	}
	WriteJSON(w, 200, dbPostsToPosts(posts))
	return nil
}

func (a *APIServer) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user *database.User) *APIError {
	queries := database.New(a.store.db)
	feedFollows, err := queries.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		return &APIError{
			Message: "Unable to fetch feed follows: " + err.Error(),
			status:  400,
		}
	}
	WriteJSON(w, 200, dbFollowedFeedsToFollowedFeeds(feedFollows))
	return nil
}

func (a *APIServer) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user *database.User) *APIError {
	var params CreateFeedFollowParams
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return &APIError{
			Message: "Invalid JSON: unable to parse request body!",
			status:  400,
		}
	}
	if params.FeedID == "" {
		return &APIError{
			Message: "feed_id is a required field!",
			status:  400,
		}
	}
	// check if feed_id is actually a uuid
	feedID, err := uuid.Parse(params.FeedID)
	if err != nil {
		return &APIError{
			Message: "Invalid feed_id!",
			status:  400,
		}
	}

	queries := database.New(a.store.db)
	followedFeed, err := queries.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedID,
	})
	if err != nil {
		return &APIError{
			Message: "Unable to follow feed!: " + err.Error(),
			status:  400,
		}
	}
	WriteJSON(w, 201, dbFollowedFeedToFollowedFeed(&followedFeed))
	return nil
}

func (a *APIServer) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user *database.User) *APIError {

	feedId := chi.URLParam(r, "feedID")
	if feedId == "" {
		return &APIError{
			Message: "feedID is a required URL parameter! Required: /feedfollow/{feedID}",
			status:  400,
		}
	}

	feedID, err := uuid.Parse(feedId)
	if err != nil {
		return &APIError{
			Message: "Invalid feed_id!",
			status:  400,
		}
	}

	queries := database.New(a.store.db)
	err = queries.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedID,
	})
	if err != nil {
		return &APIError{
			Message: "Unable to unfollow feed!: " + err.Error(),
			status:  400,
		}
	}
	WriteJSON(w, 200, "Unfollowed succesfully!")
	return nil
}
