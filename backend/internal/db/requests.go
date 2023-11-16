package db

import (
	"context"

	"github.com/zpx64/supreme-octopus/internal/model"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/jackc/pgx/v5/pgxpool"
)

func IsUserExist(
	ctx  context.Context,
	conn *pgxpool.Conn,
	data *model.UserNCred,
) (bool, error) {
	exists := [2]bool{}

	rows, err := conn.Query(context.TODO(),
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
	ctx  context.Context,
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

