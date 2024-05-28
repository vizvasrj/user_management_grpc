package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
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

func (r *PostgresDB) SearchUsers(ctx context.Context, criteria []*user_proto.SearchCriteria) ([]*user_proto.User, error) {
	query := "SELECT * FROM users WHERE 1=1" // Start with base query
	args := []interface{}{}
	var index int = 1

	// Build WHERE clause for all criteria except OR
	for _, criterion := range criteria {
		switch criterion.Operator {
		case user_proto.Operator_EQ:
			query += fmt.Sprintf(" AND %s = $%d", criterion.Field, index)
		case user_proto.Operator_GT:
			query += fmt.Sprintf(" AND %s > $%d", criterion.Field, index)
		case user_proto.Operator_LT:
			query += fmt.Sprintf(" AND %s < $%d", criterion.Field, index)
		case user_proto.Operator_GTE:
			query += fmt.Sprintf(" AND %s >= $%d", criterion.Field, index)
		case user_proto.Operator_LTE:
			query += fmt.Sprintf(" AND %s <= $%d", criterion.Field, index)
		case user_proto.Operator_BETWEEN:
			if criterion.RangeCriteria != nil {
				query += fmt.Sprintf(" AND %s BETWEEN $%d AND $%d",
					criterion.Field, index, index+1)
				args = append(args, criterion.RangeCriteria.StartValue, criterion.RangeCriteria.EndValue)
				index += 2
			}
		}
		index++
	}

	// Remove the initial "AND" if no other conditions were added
	if index == 2 {
		query = "SELECT * FROM users"
	}

	rows, err := r.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Execute the query and handle results
	var users []*user_proto.User
	for rows.Next() {
		var u user_proto.User
		if err := rows.Scan(&u.Id, &u.Fname, &u.City, &u.Phone, &u.Height, &u.Married); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Filter users based on OR conditions in Golang
	filteredUsers := []*user_proto.User{}
	for _, criterion := range criteria {
		if criterion.Operator == user_proto.Operator_OR {
			for _, user := range users {
				userValue := reflect.ValueOf(user).FieldByName(criterion.Field).Interface()
				if userValue == criterion.Value {
					filteredUsers = append(filteredUsers, user)
					break // No need to check other users for this OR condition
				}
			}
		}
	}

	return filteredUsers, nil
}

func (p *PostgresDB) AddUser(user *user_proto.User) (*user_proto.User, error) {
	err := p.Db.QueryRow("INSERT INTO users (id, fname, city, phone, height, married) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Id, user.Fname, user.City, user.Phone, user.Height, user.Married).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
