package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"src/pkg/misc"
	"src/user_proto"
	"strconv"
	"strings"

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

func (p *PostgresDB) DbGetUserById(ctx context.Context, id int32) (*user_proto.User, error) {
	user := &user_proto.User{}
	err := p.Db.QueryRow("SELECT id, fname, city, phone, height, married FROM users WHERE id = $1", id).Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresDB) DbGetUsersByIds(ctx context.Context, ids []int32) ([]*user_proto.User, error) {
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

func (s *PostgresDB) DbSearchUsers(ctx context.Context, req *user_proto.SearchUsersRequest) ([]*user_proto.User, error) {
	query := "SELECT * FROM users WHERE 1=1"
	args := []interface{}{}
	var index int = 1

	// Handle filter criteria (e.g., city="Sydney|New York|Tokyo")
	for field, filterValue := range req.Filters {
		filterValues := strings.Split(filterValue, "|")
		for i, v := range filterValues {
			if i == 0 {
				query += fmt.Sprintf(" AND (%s = $%d", field, index)
			} else {
				query += fmt.Sprintf(" OR %s = $%d", field, index)
			}
			args = append(args, v)
			index++
		}
		query += ")"
	}

	// Handle range filters (e.g., low_height=5.0 high_height=6.0)
	for field, filterValue := range req.RangeFilters {
		filterValues := strings.Split(filterValue, " ")
		if len(filterValues) == 2 {
			// Assuming range filters are always "low_field high_field"
			lowValue, err := strconv.ParseFloat(filterValues[0], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid low value for field %s: %v", field, err)
			}
			highValue, err := strconv.ParseFloat(filterValues[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid high value for field %s: %v", field, err)
			}
			query += fmt.Sprintf(" AND %s BETWEEN $%d AND $%d", field, index, index+1)
			args = append(args, lowValue, highValue)
			index += 2
		} else {
			return nil, fmt.Errorf("invalid range filter for field %s: expected 'low_value high_value'", field)
		}
	}

	rows, err := s.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results into User objects
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

	return users, nil
}

func (p *PostgresDB) AddUser(user *user_proto.User) (*user_proto.User, error) {
	err := p.Db.QueryRow("INSERT INTO users (id, fname, city, phone, height, married) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Id, user.Fname, user.City, user.Phone, user.Height, user.Married).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
