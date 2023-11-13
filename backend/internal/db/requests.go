package db

import (
	"context"

	"github.com/ssleert/tzproj/internal/vars"
	"github.com/ssleert/tzproj/pkg/genderize"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	PageSize = 10
)

type All struct {
	P People            `json:"people"`
	G Gender            `json:"gender"`
	N []Nationalization `json:"nationalizations"`
}

type People struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
	Age        int     `json:"age"`
}

type Gender struct {
	id          *int             `json:"id,omitempty"`
	Gender      genderize.Gender `json:"gender"`
	Probability float64          `json:"probability"`
}

type Nationalization struct {
	id          *int    `json:"id,omitempty"`
	CountryCode string  `json:"country_code"`
	Probability float64 `json:"probability"`
}

func GetPeopleId(conn *pgxpool.Conn, p People) (int, error) {
	var err error
	if p.Patronymic != nil {
		err = conn.QueryRow(context.TODO(),
			`SELECT people_id FROM peoples 
			 WHERE name = $1 AND surname = $2 AND age = $3 AND patronymic = $4`,
			p.Name, p.Surname, p.Age, *p.Patronymic,
		).Scan(&p.Id)
	} else {
		err = conn.QueryRow(context.TODO(),
			`SELECT people_id FROM peoples 
			 WHERE name = $1 AND surname = $2 AND age = $3`,
			p.Name, p.Surname, p.Age,
		).Scan(&p.Id)
	}
	if err != nil {
		return 0, err
	}
	return p.Id, nil
}

func CheckPeople(conn *pgxpool.Conn, p People) (bool, error) {
	var (
		exists bool
		err    error
	)
	if p.Patronymic != nil {
		err = conn.QueryRow(context.TODO(),
			`SELECT EXISTS(SELECT 1 FROM peoples 
			 WHERE name = $1 AND surname = $2 AND age = $3 AND patronymic = $4)`,
			p.Name, p.Surname, p.Age, *p.Patronymic,
		).Scan(&exists)
	} else {
		err = conn.QueryRow(context.TODO(),
			`SELECT EXISTS(SELECT 1 FROM peoples 
			 WHERE name = $1 AND surname = $2 AND age = $3)`,
			p.Name, p.Surname, p.Age,
		).Scan(&exists)
	}
	if err != nil {
		return false, err
	}
	return exists, nil
}

func InsertPeople(conn *pgxpool.Conn, p People) (int, error) {
	var (
		id  int
		err error
	)
	if p.Id <= 0 {
		err = conn.QueryRow(context.TODO(),
			`INSERT INTO peoples (name, surname, age, patronymic)
			 VALUES ($1, $2, $3, $4)
			 RETURNING people_id`,
			p.Name, p.Surname, p.Age, p.Patronymic,
		).Scan(&id)
	} else {
		err = conn.QueryRow(context.TODO(),
			`INSERT INTO peoples (people_id, name, surname, age, patronymic)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING people_id`,
			p.Id, p.Name, p.Surname, p.Age, p.Patronymic,
		).Scan(&id)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ExistPeopleById(conn *pgxpool.Conn, id int) (bool, error) {
	var (
		exists bool
		err    error
	)
	err = conn.QueryRow(context.TODO(),
		`SELECT EXISTS(SELECT 1 FROM peoples 
		 WHERE people_id = $1)`, id,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func InsertNationalizationsById(conn *pgxpool.Conn, id int, ns []Nationalization) error {
	nationalizationIds := make([]int, len(ns))
	for i, e := range ns {
		err := conn.QueryRow(context.TODO(),
			`INSERT INTO nationalizations (country_code, probability)
			 VALUES ($1, $2)
			 RETURNING nationalization_id`,
			e.CountryCode, e.Probability,
		).Scan(&nationalizationIds[i])
		if err != nil {
			return err
		}
		_, err = conn.Exec(context.TODO(),
			`INSERT INTO people_nationalizations (people_id, nationalization_id)
			 VALUES ($1, $2)`,
			id, nationalizationIds[i],
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetNationalizationsById(conn *pgxpool.Conn, id int) ([]Nationalization, error) {
	nationalizationIds := make([]int, 0, PageSize)
	nationalizations := make([]Nationalization, 0, PageSize)

	rows, err := conn.Query(context.TODO(),
		`SELECT nationalization_id
		 FROM people_nationalizations
		 WHERE people_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nId int
	for rows.Next() {
		err = rows.Scan(&nId)
		if err != nil {
			return nil, err
		}
		nationalizationIds = append(
			nationalizationIds,
			nId,
		)
	}

	for _, e := range nationalizationIds {
		rows, err = conn.Query(context.TODO(),
			`SELECT country_code,
							probability
			 FROM nationalizations
			 WHERE nationalization_id = $1`, e,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var (
			countryCode string
			probability float64
		)
		for rows.Next() {
			err = rows.Scan(
				&countryCode,
				&probability,
			)
			if err != nil {
				return nil, err
			}
			nationalizations = append(
				nationalizations,
				Nationalization{
					CountryCode: countryCode,
					Probability: probability,
				},
			)
		}
	}
	return nationalizations, nil
}

func InsertGenderById(conn *pgxpool.Conn, id int, g Gender) error {
	_, err := conn.Exec(context.TODO(),
		`INSERT INTO genders (people_id, gender, probability)
		 VALUES ($1, $2, $3)`,
		id, g.Gender, g.Probability,
	)
	if err != nil {
		return err
	}
	return nil
}

func InsertAll(conn *pgxpool.Conn, a All) (int, error) {
	exists, err := CheckPeople(conn, a.P)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, vars.ErrAlreadyInDb
	}
	peopleId, err := InsertPeople(conn, a.P)
	if err != nil {
		return 0, err
	}
	err = InsertGenderById(conn, peopleId, a.G)
	if err != nil {
		return 0, err
	}
	err = InsertNationalizationsById(conn, peopleId, a.N)
	if err != nil {
		return 0, err
	}
	return peopleId, nil
}

func DeletePeopleById(conn *pgxpool.Conn, id int) error {
	_, err := conn.Exec(context.TODO(),
		`DELETE FROM peoples WHERE people_id = $1`, id,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteGenderByid(conn *pgxpool.Conn, id int) error {
	_, err := conn.Exec(context.TODO(),
		`DELETE FROM genders WHERE people_id = $1`, id,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNationalizationsById(conn *pgxpool.Conn, id int) error {
	var (
		nationalizationId  int
		nationalizationIds []int
	)
	rows, err := conn.Query(context.TODO(),
		`DELETE FROM people_nationalizations 
		 WHERE people_id = $1
		 RETURNING nationalization_id`, id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&nationalizationId)
		if err != nil {
			return err
		}

		nationalizationIds = append(
			nationalizationIds,
			nationalizationId,
		)
	}
	for _, nId := range nationalizationIds {
		_, err = conn.Exec(context.TODO(),
			`DELETE FROM nationalizations
			 WHERE nationalization_id = $1`, nId,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteAllById(conn *pgxpool.Conn, id int) error {
	exists, err := ExistPeopleById(conn, id)
	if err != nil {
		return err
	}
	if !exists {
		return vars.ErrNotInDb
	}
	err = DeleteNationalizationsById(conn, id)
	if err != nil {
		return err
	}
	err = DeleteGenderByid(conn, id)
	if err != nil {
		return err
	}
	err = DeletePeopleById(conn, id)
	if err != nil {
		return err
	}
	return nil
}

func ReplaceAll(conn *pgxpool.Conn, a All) error {
	exists, err := ExistPeopleById(conn, a.P.Id)
	if err != nil {
		return err
	}
	if !exists {
		return vars.ErrNotInDb
	}
	_, err = conn.Exec(context.TODO(),
		`UPDATE peoples
		 SET name = $2,
		     surname = $3,
				 patronymic = $4,
				 age = $5
		 WHERE people_id = $1`, a.P.Id,
		a.P.Name, a.P.Surname,
		a.P.Patronymic, a.P.Age,
	)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.TODO(),
		`UPDATE genders
		 SET gender = $2,
		     probability = $3
		 WHERE people_id = $1`, a.P.Id,
		a.G.Gender, a.G.Probability,
	)
	if err != nil {
		return err
	}
	err = DeleteNationalizationsById(conn, a.P.Id)
	if err != nil {
		return err
	}
	err = InsertNationalizationsById(conn, a.P.Id, a.N)
	if err != nil {
		return err
	}
	return nil
}

func GetPeoples(conn *pgxpool.Conn, offset, limit int, fs []FilterOperation) ([]All, error) {
	allData := make([]All, 0, limit)
	prependArgs := []any{limit, offset}
	where, args := FilterOperationsToSql(len(prependArgs), fs)
	args = append(
		prependArgs,
		args...,
	)

	rows, err := conn.Query(context.TODO(),
		`SELECT people_id,
						name,
		        surname,
						patronymic,
						age
		 FROM peoples
		`+where+`
		 LIMIT  $1
		 OFFSET $2`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people People
	for rows.Next() {
		err = rows.Scan(
			&people.Id, &people.Name, &people.Surname,
			&people.Patronymic, &people.Age,
		)
		if err != nil {
			return nil, err
		}

		allData = append(allData,
			All{
				P: people,
			},
		)
	}
	if len(allData) < 1 {
		return nil, vars.ErrNotInDb
	}

	for i := range allData {
		err = conn.QueryRow(context.TODO(),
			`SELECT gender,
							probability
			 FROM genders
			 WHERE people_id = $1`,
			allData[i].P.Id,
		).Scan(
			&allData[i].G.Gender,
			&allData[i].G.Probability,
		)
		if err != nil {
			return nil, err
		}
	}

	for i, e := range allData {
		ns, err := GetNationalizationsById(conn, e.P.Id)
		if err != nil {
			return nil, err
		}

		allData[i].N = ns
	}

	return allData, nil
}
