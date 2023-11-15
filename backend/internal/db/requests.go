package db

import (
	"context"
	"github.com/zpx64/supreme-octopus/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// create new user with credentials
// return created user id & error
func CreateUser(
	ctx  context.Context,
	conn *pgxpool.Conn,
	data *model.UserNCred,
) (int, error) {
	var (
		id  int
		err error
	)

	err = conn.QueryRow(ctx,
		`INSERT INTO users (nickname, creation_date, name, surname)
		 VALUES ($1, $2, $3, $4)
		 RETURNING user_id`,
		data.User.Nickname, data.User.CreationDate,
		data.User.Name, data.User.Surname,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	_, err = conn.Exec(ctx,
		`INSERT INTO users_credentials (user_id, email, password, pow)
		 VALUES ($1, $2, $3, $4)`,
		id, data.Credentials.Email,
		data.Credentials.Password, data.Credentials.Pow,
	)
	if err != nil {
		return 0, err
	}

	return id, nil
}

