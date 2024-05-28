package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"src/pkg/misc"
	"src/user_proto"

	"github.com/lib/pq"
)

type PostgresDB struct {
	Db *sql.DB
}

func GetConnection() *sql.DB {
	dbuser := misc.GetEnv("dbuser")
	dbpassword := misc.GetEnv("dbpassword")
	dbname := misc.GetEnv("dbname")
	dbsslmode := misc.GetEnv("dbsslmode")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbuser, dbpassword, dbname, dbsslmode)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close() //! remove
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

func (p *PostgresDB) GetUserById(ctx context.Context, id int32) (*user_proto.User, error) {
	user := &user_proto.User{}
	err := p.Db.QueryRow("SELECT id, fname, city, phone, height, married FROM users WHERE id = $1", id).Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresDB) GetUsersByIds(ctx context.Context, ids []int32) ([]*user_proto.User, error) {
	users := make([]*user_proto.User, 0)
	rows, err := p.Db.Query("SELECT id, fname, city, phone, height, married FROM users WHERE id = ANY($1)", pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := &user_proto.User{}
		err := rows.Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (p *PostgresDB) SearchUsers(ctx context.Context, req *user_proto.SearchUsersRequest) ([]*user_proto.User, error) {
	users := make([]*user_proto.User, 0)
	query := "SELECT id, fname, city, phone, height, married FROM users WHERE "
	args := make([]interface{}, 0)
	for i, criterion := range req.Criteria {
		if i > 0 {
			query += " AND "
		}
		switch criterion.Field {
		case "city":
			query += fmt.Sprintf("city = $%d", len(args)+1)
			args = append(args, criterion.Value)
		case "married":
			query += fmt.Sprintf("married = $%d", len(args)+1)
			args = append(args, criterion.Value)
		case "height":
			query += fmt.Sprintf("height = $%d", len(args)+1)
			args = append(args, criterion.Value)
		case "phone":
			query += fmt.Sprintf("phone = $%d", len(args)+1)
			args = append(args, criterion.Value)
		case "fname":
			query += fmt.Sprintf("fname = $%d", len(args)+1)
			args = append(args, criterion.Value)
		}
	}
	rows, err := p.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := &user_proto.User{}
		err := rows.Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (p *PostgresDB) AddUser(user *user_proto.User) (*user_proto.User, error) {
	err := p.Db.QueryRow("INSERT INTO users (id, fname, city, phone, height, married) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Id, user.Fname, user.City, user.Phone, user.Height, user.Married).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
