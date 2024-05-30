package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"src/pkg/misc"
	"src/pkg/myerror"
	"src/user_proto"

	"github.com/lib/pq"
)

type PostgresDB struct {
	Db *sql.DB
}

func GetConnection() *sql.DB {
	dbuser := misc.GetEnv("POSTGRES_USER")
	dbpassword := misc.GetEnv("POSTGRES_PASSWORD")
	dbname := misc.GetEnv("POSTGRES_DB")
	dbsslmode := misc.GetEnv("POSTGRES_SSLMODE")
	host := misc.GetEnv("POSTGRES_HOST")
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", host, dbuser, dbpassword, dbname, dbsslmode)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to db")

	return db

}

func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{Db: db}
}

func (p *PostgresDB) DbGetUserById(ctx context.Context, id int32) (*user_proto.User, error) {
	user := &user_proto.User{}
	err := p.Db.QueryRow("SELECT id, fname, city, phone, height, married FROM users WHERE id = $1", id).Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
	if err != nil {
		return nil, myerror.WrapError(err, "Error in fetching user")
	}
	return user, nil
}

func (p *PostgresDB) DbGetUsersByIds(ctx context.Context, ids []int32) ([]*user_proto.User, error) {
	users := make([]*user_proto.User, 0)
	rows, err := p.Db.Query("SELECT id, fname, city, phone, height, married FROM users WHERE id = ANY($1)", pq.Array(ids))
	if err != nil {
		return nil, myerror.WrapError(err, "Error in fetching users")
	}
	defer rows.Close()
	for rows.Next() {
		user := &user_proto.User{}
		err := rows.Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
		if err != nil {
			return nil, myerror.WrapError(err, "Error in scanning users")
		}
		users = append(users, user)
	}
	fmt.Printf("%#v", users)
	return users, nil
}

func (s *PostgresDB) DbSearchUsers(ctx context.Context, req *user_proto.SearchUsersRequest) ([]*user_proto.User, error) {
	query := "SELECT * FROM users WHERE 1=1"
	args := []interface{}{}
	var index int = 1

	// Building WHERE clause dynamically
	if req.Id != 0 {
		query += fmt.Sprintf(" AND id = $%d", index)
		args = append(args, req.Id)
		index++
	}
	if req.Fname != "" {
		query += fmt.Sprintf(" AND fname = $%d", index)
		args = append(args, req.Fname)
		index++
	}
	if req.City != "" {
		query += fmt.Sprintf(" AND city = $%d", index)
		args = append(args, req.City)
		index++
	}
	if req.Phone != 0 {
		query += fmt.Sprintf(" AND phone = $%d", index)
		args = append(args, req.Phone)
		index++
	}

	if req.Married != nil {
		query += fmt.Sprintf(" AND married = $%d", index)
		args = append(args, req.Married.IsMarried)
		index++

	}

	if req.Height != nil {
		if req.Height.StartValue != 0.0 || req.Height.EndValue != 0.0 {
			query += fmt.Sprintf(" AND height BETWEEN $%d AND $%d", index, index+1)
			args = append(args, req.Height.StartValue, req.Height.EndValue)
			index += 2
		}
	}
	// fmt.Printf("req.Height.StartValue req.Height.EndValue %#v", req.Height)
	rows, err := s.Db.Query(query, args...)
	if err != nil {
		return nil, myerror.WrapError(err, "Error in searching users")
	}
	defer rows.Close()

	// Scan results into User objects
	var users []*user_proto.User
	for rows.Next() {
		var u user_proto.User
		if err := rows.Scan(&u.Id, &u.Fname, &u.City, &u.Phone, &u.Height, &u.Married); err != nil {
			return nil, myerror.WrapError(err, "Error in scanning users")
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, myerror.WrapError(err, "Error in fetching users")
	}

	return users, nil
}

func (p *PostgresDB) AddUser(user *user_proto.User) (*user_proto.User, error) {
	err := p.Db.QueryRow("INSERT INTO users (id, fname, city, phone, height, married) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Id, user.Fname, user.City, user.Phone, user.Height, user.Married).Scan(&user.Id)
	if err != nil {
		return nil, myerror.WrapError(err, "Error in adding user")
	}
	return user, nil
}
