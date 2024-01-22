package db

// TODO: split in separate files

import (
	"context"
	"time"

	"github.com/zpx64/supreme-octopus/internal/model"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func IsUserExist(
	ctx context.Context,
	conn *pgxpool.Conn,
	data *model.UserNCred,
) (bool, error) {
	exists := [2]bool{}

	rows, err := conn.Query(ctx,
		`SELECT EXISTS(SELECT 1 FROM users
		 WHERE nickname = $1)
		 UNION
		 SELECT EXISTS(SELECT 1 FROM users_credentials
		 WHERE email = $2)`,
		data.User.Nickname, data.Credentials.Email,
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		err = rows.Scan(&exists[i])
		if err != nil {
			return false, err
		}
		i++
	}
	if rows.Err() != nil {
		return false, err
	}

	return exists[0] || exists[1], nil
}

// create new user with credentials
// return created user id & error
func CreateUser(
	ctx context.Context,
	conn *pgxpool.Conn,
	data *model.UserNCred,
) (int, error) {
	exist, err := IsUserExist(ctx, conn, data)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, vars.ErrAlreadyInDb
	}

	id := 0
	err = conn.QueryRow(ctx,
		`WITH new_user_id AS (
		   INSERT INTO users (nickname, avatar_img, creation_date, name, surname)
		   VALUES ($1, $2, $3, $4, $5)
		   RETURNING user_id 
		 )
		 INSERT INTO users_credentials (user_id, email, password, pow)
		 SELECT user_id, $6, $7, $8
		 FROM new_user_id
		 RETURNING user_id`,
		data.User.Nickname, data.User.AvatarImg,
		data.User.CreationDate,
		data.User.Name, data.User.Surname,
		data.Credentials.Email, data.Credentials.Password,
		data.Credentials.Pow,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetCredentialsByEmail(
	ctx context.Context,
	conn *pgxpool.Conn,
	email string,
) (model.UserCredentials, error) {
	var (
		credentials_id       int
		credentials_email    string
		credentials_password string
		credentials_pow      string
	)
	err := conn.QueryRow(ctx,
		`SELECT user_id, email,
		        password, pow
		 FROM users_credentials
		 WHERE email = $1`,
		email,
	).Scan(
		&credentials_id,
		&credentials_email,
		&credentials_password,
		&credentials_pow,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.UserCredentials{}, vars.ErrNotInDb
		}
		return model.UserCredentials{}, err
	}

	return model.UserCredentials{
		UserId:   credentials_id,
		Email:    credentials_email,
		Password: credentials_password,
		Pow:      credentials_pow,
	}, nil
}

func InsertNewToken(
	ctx context.Context,
	conn *pgxpool.Conn,
	token model.UserToken,
) (int, error) {
	var (
		id int
	)
	err := conn.QueryRow(ctx,
		`INSERT INTO users_tokens (
		   user_id, device_id, 
		   refresh_token, user_agent, 
		   token_date
		 )
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING token_id`,
		token.UserId, token.DeviceId,
		token.RefreshToken, token.UserAgent,
		token.TokenDate,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// TODO: rewrite with pass by pointer model.UserToken
func UpdateToken(
	ctx context.Context,
	conn *pgxpool.Conn,
	token model.UserToken,
) error {
	cmdTag, err := conn.Exec(ctx,
		`UPDATE users_tokens
		 SET refresh_token = $2,
		     token_date = $3
		 WHERE token_id = $1`,
		token.TokenId,
		token.RefreshToken,
		token.TokenDate,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != 1 {
		return vars.ErrNotInDb
	}

	return nil
}

// TODO: rewrite with pass by pointer model.UserPost
func InsertNewPost(
	ctx context.Context,
	conn *pgxpool.Conn,
	post model.UserPost,
) (int, error) {
	var (
		id int
	)
	err := conn.QueryRow(ctx,
		`INSERT INTO users_posts (
		   user_id, creation_date, post_type, 
		   body, attachments, votes_amount,
		   comments_amount, is_comments_disallowed
		 )
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING post_id`,
		post.UserId, post.CreationDate, post.PostType,
		post.Body, post.Attachments, post.VotesAmount,
		post.CommentsAmount, post.IsCommentsDisallowed,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ListPosts(
	ctx context.Context,
	conn *pgxpool.Conn,
	offset uint,
	limit uint,
) ([]model.UserNPost, error) {
	rows, err := conn.Query(ctx,
		`SELECT u.nickname, u.avatar_img,
		        up.post_id, up.creation_date, up.post_type,
		        up.body, up.attachments,
		        up.votes_amount, up.comments_amount,
		        up.is_comments_disallowed
		 FROM users_posts AS up
		 JOIN users AS u 
		 ON u.user_id = up.user_id
		 ORDER BY up.creation_date DESC
		 LIMIT $1
		 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, vars.ErrNotInDb
		}
		return nil, err
	}
	defer rows.Close()

	// TODO: make smart preallocation
	//       based on rows amount
	posts := make([]model.UserNPost, 0, 32)
	for rows.Next() {
		userPost := model.UserNPost{}
		err = rows.Scan(
			&userPost.User.Nickname,
			&userPost.User.AvatarImg,
			&userPost.Post.PostId,
			&userPost.Post.CreationDate,
			&userPost.Post.PostType,
			&userPost.Post.Body,
			&userPost.Post.Attachments,
			&userPost.Post.VotesAmount,
			&userPost.Post.CommentsAmount,
			&userPost.Post.IsCommentsDisallowed,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, userPost)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return posts, nil
}

func GetPost(
	ctx context.Context,
	conn *pgxpool.Conn,
	postId int,
) (model.UserNPost, error) {
	var userPost model.UserNPost
	err := conn.QueryRow(ctx,
		`SELECT u.nickname, u.avatar_img,
		        up.post_id, up.creation_date, up.post_type,
		        up.body, up.attachments,
		        up.votes_amount, up.comments_amount,
		        up.is_comments_disallowed
		 FROM users_posts AS up
		 JOIN users AS u 
		 ON u.user_id = up.user_id
		 WHERE up.post_id = $1`,
		postId,
	).Scan(
		&userPost.User.Nickname,
		&userPost.User.AvatarImg,
		&userPost.Post.PostId,
		&userPost.Post.CreationDate,
		&userPost.Post.PostType,
		&userPost.Post.Body,
		&userPost.Post.Attachments,
		&userPost.Post.VotesAmount,
		&userPost.Post.CommentsAmount,
		&userPost.Post.IsCommentsDisallowed,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.UserNPost{}, vars.ErrNotInDb
		}
		return model.UserNPost{}, err
	}

	return userPost, nil
}

func IsPostVoted(
	ctx context.Context,
	conn *pgxpool.Conn,
	userId int,
	postId int,
) (model.VoteAction, error) {
	var (
		voteType model.VoteAction
	)
	err := conn.QueryRow(ctx,
		`SELECT vote_type FROM users_likes
		 WHERE user_id = $1 AND post_id = $2`,
		userId, postId,
	).Scan(&voteType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return voteType, nil
}

func IsPostExists(
	ctx context.Context,
	conn *pgxpool.Conn,
	postId int,
) (bool, error) {
	var existsPost bool
	err := conn.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users_posts
		 WHERE post_id = $1)`,
		postId,
	).Scan(&existsPost)
	if err != nil {
		return false, err
	}

	return existsPost, nil
}

func VotePost(
	ctx context.Context,
	conn *pgxpool.Conn,
	userId int,
	postId int,
	voteType model.VoteAction,
	creationDate time.Time,
) (int, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	existsPost, err := IsPostExists(ctx, conn, postId)
	if err != nil {
		return 0, err
	}

	if !existsPost {
		return 0, vars.ErrNotInDb
	}

	var (
		voteTypeFromDb model.VoteAction
		likeId         int
		existsLike     bool
	)
	err = tx.QueryRow(ctx,
		`SELECT vote_type, like_id FROM users_likes
		 WHERE post_id = $1 AND user_id = $2`,
		postId, userId,
	).Scan(&voteTypeFromDb, &likeId)
	if err != nil {
		if err != pgx.ErrNoRows {
			return 0, err
		}
	} else {
		existsLike = true
	}

	if existsLike {
		if voteTypeFromDb == voteType {
			return likeId, nil
		}
	}

	var (
		id int
	)
	if !existsLike {
		err = tx.QueryRow(ctx,
			`INSERT INTO users_likes (user_id, post_id, vote_type, creation_date)
			 VALUES ($1, $2, $3, $4)
			 RETURNING like_id`,
			userId, postId, voteType, creationDate,
		).Scan(&id)
	} else {
		err = tx.QueryRow(ctx,
			`UPDATE users_likes
			 SET vote_type = $3,
			     creation_date = $4
			 WHERE user_id = $1 AND post_id = $2
			 RETURNING like_id`,
			userId, postId, voteType, creationDate,
		).Scan(&id)
	}
	if err != nil {
		return 0, err
	}

	var (
		votesAppend = 1
	)
	if existsLike {
		if voteType != voteTypeFromDb {
			votesAppend += 1
		}
	}
	if voteType == model.VoteDownvote {
		votesAppend = -votesAppend
	}

	_, err = tx.Exec(ctx,
		`UPDATE users_posts 
		 SET votes_amount = votes_amount + $1
		 WHERE post_id = $2`,
		votesAppend, postId,
	)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func RemovePostVote(
	ctx context.Context,
	conn *pgxpool.Conn,
	postId int,
	userId int,
) error {
	postExists, err := IsPostExists(ctx, conn, postId)
	if err != nil {
		return err
	}
	if !postExists {
		return vars.ErrNotInDb
	}

	var voteAction model.VoteAction
	err = conn.QueryRow(ctx,
		`SELECT vote_type FROM users_likes
		 WHERE post_id = $1 AND user_id = $2`,
		postId, userId,
	).Scan(&voteAction)
	if err != nil {
		if err == pgx.ErrNoRows {
			return vars.ErrNotInDb
		}
		return err
	}

	recoveredVotes := 0
	if voteAction == model.VoteUpvote {
		recoveredVotes = -1
	} else if voteAction == model.VoteDownvote {
		recoveredVotes = 1
	}

	_, err = conn.Exec(ctx,
		`UPDATE users_posts 
		 SET votes_amount = votes_amount + $1
		 WHERE post_id = $2`,
		recoveredVotes, postId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return vars.ErrNotInDb
		}
		return err
	}

	cmdTag, err := conn.Exec(ctx,
		`DELETE FROM users_likes 
		 WHERE post_id = $1 AND user_id = $2`,
		postId, userId,
	)
	if cmdTag.RowsAffected() == 0 {
		return vars.ErrNotInDb
	}
	if err != nil {
		return err
	}
	return nil
}

func IsCommentsAllowedForPost(
	ctx context.Context,
	conn *pgxpool.Conn,
	postId int,
) (bool, error) {
	var (
		result bool
	)
	err := conn.QueryRow(ctx,
		`SELECT is_comments_disallowed FROM users_posts
		 WHERE post_id = $1`,
		postId,
	).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}

func InsertNewComment(
	ctx context.Context,
	conn *pgxpool.Conn,
	comment model.UserComment,
) (int, error) {
	var (
		id int
	)
	err := conn.QueryRow(ctx,
		`INSERT INTO users_comments (
		   user_id, post_id, body, attachments,
		   creation_date, votes_amount, reply_id
		 )
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING comment_id`,
		comment.UserId, comment.PostId, comment.Body,
		comment.Attachments, comment.CreationDate,
		comment.VotesAmount, comment.ReplyId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetCommentsByPostId(
	ctx context.Context,
	conn *pgxpool.Conn,
	postId int,
) ([]model.CommentWithUser, error) {
	rows, err := conn.Query(ctx,
		`SELECT uc.comment_id, u.nickname, u.avatar_img, 
		        uc.body, uc.attachments, uc.creation_date,
		        uc.votes_amount, uc.reply_id 
		 FROM users_comments AS uc
		 JOIN users AS u
		 ON u.user_id = uc.user_id
		 WHERE post_id = $1
		 ORDER BY comment_id`,
		postId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, vars.ErrNotInDb
		}
		return nil, err
	}
	defer rows.Close()

	// TODO: make smart preallocation
	//       based on rows amount
	comments := make([]model.CommentWithUser, 0, 32)
	for rows.Next() {
		comment := model.CommentWithUser{}
		err = rows.Scan(
			&comment.CommentId,
			&comment.Nickname,
			&comment.AvatarImg,
			&comment.Body,
			&comment.Attachments,
			&comment.CreationDate,
			&comment.VotesAmount,
			&comment.ReplyId,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil
}

func IsCommentExists(
	ctx context.Context,
	conn *pgxpool.Conn,
	postId int,
) (bool, error) {
	var existsComment bool
	err := conn.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users_comments
		 WHERE comment_id = $1)`,
		postId,
	).Scan(&existsComment)
	if err != nil {
		return false, err
	}

	return existsComment, nil
}

func VoteComment(
	ctx context.Context,
	conn *pgxpool.Conn,
	userId int,
	commentId int,
	voteType model.VoteAction,
	creationDate time.Time,
) (int, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	existsComment, err := IsCommentExists(ctx, conn, commentId)
	if err != nil {
		return 0, err
	}

	if !existsComment {
		return 0, vars.ErrNotInDb
	}

	var (
		voteTypeFromDb model.VoteAction
		likeId         int
		existsLike     bool
	)
	err = tx.QueryRow(ctx,
		`SELECT vote_type, like_id FROM users_comments_likes
		 WHERE comment_id = $1 AND user_id = $2`,
		commentId, userId,
	).Scan(&voteTypeFromDb, &likeId)
	if err != nil {
		if err != pgx.ErrNoRows {
			return 0, err
		}
	} else {
		existsLike = true
	}

	if existsLike {
		if voteTypeFromDb == voteType {
			return likeId, nil
		}
	}

	var (
		id int
	)
	if !existsLike {
		err = tx.QueryRow(ctx,
			`INSERT INTO users_comments_likes (
         user_id, comment_id, 
         vote_type, creation_date
       )
			 VALUES ($1, $2, $3, $4)
			 RETURNING like_id`,
			userId, commentId, voteType, creationDate,
		).Scan(&id)
	} else {
		err = tx.QueryRow(ctx,
			`UPDATE users_comments_likes
			 SET vote_type = $3,
			     creation_date = $4
			 WHERE user_id = $1 AND comment_id = $2
			 RETURNING like_id`,
			userId, commentId,
			voteType, creationDate,
		).Scan(&id)
	}
	if err != nil {
		return 0, err
	}

	var (
		votesAppend = 1
	)
	if existsLike {
		if voteType != voteTypeFromDb {
			votesAppend += 1
		}
	}
	if voteType == model.VoteDownvote {
		votesAppend = -votesAppend
	}

	_, err = tx.Exec(ctx,
		`UPDATE users_comments
		 SET votes_amount = votes_amount + $1
		 WHERE comment_id = $2`,
		votesAppend, commentId,
	)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}
