package db

// TODO: split in separate files

import (
	"context"

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
		   INSERT INTO users (nickname, creation_date, name, surname)
		   VALUES ($1, $2, $3, $4)
		   RETURNING user_id 
		 )
		 INSERT INTO users_credentials (user_id, email, password, pow)
		 SELECT user_id, $5, $6, $7
		 FROM new_user_id
		 RETURNING user_id`,
		data.User.Nickname, data.User.CreationDate,
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
			 comments_amount
		 )
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING post_id`,
		post.UserId, post.CreationDate, post.PostType, 
		post.Body, post.Attachments, post.VotesAmount,
		post.CommentsAmount,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
