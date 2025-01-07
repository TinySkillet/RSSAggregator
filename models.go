package main

import (
	"goapi/internal/database"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Created_At string    `json:"created_at"`
	Updated_At string    `json:"updated_at"`
	Name       string    `json:"name"`
}

type Feed struct {
	ID         uuid.UUID `json:"id"`
	Created_At string    `json:"created_at"`
	Updated_At string    `json:"updated_at"`
	Name       string    `json:"feed_name"`
	Url        string    `json:"url"`
	UserID     uuid.UUID `json:"user_id"`
}

type Post struct {
	ID           uuid.UUID `json:"id"`
	Created_At   string    `json:"created_at"`
	Updated_At   string    `json:"updated_at"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	Published_At string    `json:"published_at"`
	Url          string    `json:"url"`
	FeedID       uuid.UUID `json:"feed_id"`
}

type FollowedFeed struct {
	ID         uuid.UUID `json:"id"`
	Created_At string    `json:"created_at"`
	Updated_At string    `json:"updated_at"`
	UserID     uuid.UUID `json:"user_id"`
	FeedID     uuid.UUID `json:"feed_id"`
}

type PrivateUser struct {
	User
	ApiKey string `json:"api_key"`
}

type CreateUserParams struct {
	Name string
}

type CreateFeedParams struct {
	Name string
	Url  string
}

type CreateFeedFollowParams struct {
	FeedID string `json:"feed_id"`
}

func dbUserToUser(dbUser *database.User) *User {
	return &User{
		ID:         dbUser.ID,
		Created_At: dbUser.CreatedAt.Format("2006-01-02 15:04"),
		Updated_At: dbUser.UpdatedAt.Format("2006-01-02 15:04"),
		Name:       dbUser.Name,
	}
}

func dbFeedToFeed(dbFeed *database.Feed) *Feed {
	return &Feed{
		ID:         dbFeed.ID,
		Created_At: dbFeed.CreatedAt.Format("2006-01-02 15:04"),
		Updated_At: dbFeed.CreatedAt.Format("2006-01-02 15:04"),
		Name:       dbFeed.Name,
		Url:        dbFeed.Url,
		UserID:     dbFeed.UserID,
	}
}

func dbPostToPost(dbPost *database.Post) *Post {

	var desc *string
	if dbPost.Description.Valid {
		desc = &dbPost.Description.String
	}

	return &Post{
		ID:           dbPost.ID,
		Created_At:   dbPost.CreatedAt.Format("2006-01-02 15:04"),
		Updated_At:   dbPost.UpdatedAt.Format("2006-01-02 15:04"),
		Title:        dbPost.Title,
		Description:  desc,
		Published_At: dbPost.PublishedAt.Format("2006-01-02 15:04"),
		Url:          dbPost.Url,
		FeedID:       dbPost.FeedID,
	}
}

func dbPostsToPosts(dbPosts []database.Post) []*Post {
	posts := make([]*Post, 0, len(dbPosts))
	for _, dbPost := range dbPosts {
		posts = append(posts, dbPostToPost(&dbPost))
	}
	return posts
}

func dbUserToPrivateUser(dbUser *database.User) *PrivateUser {
	return &PrivateUser{
		User: User{
			ID:         dbUser.ID,
			Created_At: dbUser.CreatedAt.Format("2006-01-02 15:04"),
			Updated_At: dbUser.UpdatedAt.Format("2006-01-02 15:04"),
			Name:       dbUser.Name,
		},
		ApiKey: dbUser.ApiKey,
	}
}

func dbUsersToUsers(dbUsers []database.User) []*User {
	users := make([]*User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = dbUserToUser(&dbUser)
	}
	return users
}

func dbFeedsToFeeds(dbFeeds []database.Feed) []*Feed {
	feeds := make([]*Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		feed := dbFeedToFeed(&dbFeed)
		feeds[i] = feed
	}
	return feeds
}

func dbFollowedFeedToFollowedFeed(dbFollowedFeed *database.FeedFollow) *FollowedFeed {
	return &FollowedFeed{
		ID:         dbFollowedFeed.ID,
		Created_At: dbFollowedFeed.CreatedAt.Format("2006-01-02 15:04"),
		Updated_At: dbFollowedFeed.CreatedAt.Format("2006-01-02 15:04"),
		UserID:     dbFollowedFeed.UserID,
		FeedID:     dbFollowedFeed.FeedID,
	}
}

func dbFollowedFeedsToFollowedFeeds(dbFollowedFeeds []database.FeedFollow) []*FollowedFeed {
	followedFeeds := make([]*FollowedFeed, len(dbFollowedFeeds))
	for i, dbFollowedFeed := range dbFollowedFeeds {
		followedFeed := dbFollowedFeedToFollowedFeed(&dbFollowedFeed)
		followedFeeds[i] = followedFeed
	}
	return followedFeeds
}
