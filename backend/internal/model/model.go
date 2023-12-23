package model

import (
	"time"
)

type Post int

const (
	PostArticle Post = iota + 1
	PostThought
)

type VoteAction int

const (
	VoteUpvote VoteAction = iota + 1
	VoteDownvote
)

type UserNCred struct {
	User        User            `json:"user"`
	Credentials UserCredentials `json:"credentials"`
}

type UserNPost struct {
	User User     `json:"user"`
	Post UserPost `json:"post"`
}

type CommentWithUser struct {
	CommentId    int       `json:"comment_id"`
	Nickname     string    `json:"nickname"`
	AvatarImg    string    `json:"avatar_img"`
	Body         string    `json:"body"`
	Attachments  []string  `json:"attachments"`
	CreationDate time.Time `json:"creation_date"`
	VotesAmount  int       `json:"votes_amount"`
	ReplyId      *int      `json:"reply_id"`
}

type User struct {
	UserId       int       `json:"user_id"`
	CreationDate time.Time `json:"creation_date"`
	Nickname     string    `json:"nickname"`
	AvatarImg    string    `json:"avatar_img"`
	Name         *string   `json:"name,omitempty"`
	Surname      *string   `json:"surname,omitempty"`
}

type UserCredentials struct {
	UserId   int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Pow      string `json:"pow"`
	// i really hate local pows
	// but i think we need it(
}

type UserToken struct {
	TokenId      int    `json:"token_id"`
	UserId       int    `json:"user_id"`
	DeviceId     string `json:"device_id"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	TokenDate    int64  `json:"token_date"`
}

type UserPost struct {
	PostId         int       `json:"post_id"`
	UserId         int       `json:"user_id"`
	CreationDate   time.Time `json:"creation_date"`
	PostType       Post      `json:"post_type"`
	Body           string    `json:"body"`
	Attachments    []string  `json:"attachments"`
	VotesAmount    int       `json:"votes_amount"`
	CommentsAmount int       `json:"comments_amount"`
}

type UserLike struct {
	LikeId       int        `json:"like_id"`
	UserId       int        `json:"user_id"`
	PostId       int        `json:"post_id"`
	VoteType     VoteAction `json:"vote_type"`
	CreationDate time.Time  `json:"creation_date"`
}

type UserComment struct {
	CommentId    int       `json:"comment_id"`
	UserId       int       `json:"user_id"`
	PostId       int       `json:"post_id"`
	Body         string    `json:"body"`
	Attachments  []string  `json:"attachments"`
	CreationDate time.Time `json:"creation_date"`
	VotesAmount  int       `json:"votes_amount"`
	ReplyId      *int      `json:"reply_id"`
}
