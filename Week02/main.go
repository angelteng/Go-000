package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type User struct {
	Id   int
	Name string
}

var NOTFOUND = errors.New("not_found")

// dao
func getUser(name string) (*User, error) {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/monitoring")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(NOTFOUND)
		} else {
			return nil, errors.WithStack(err)
		}
	}
	defer db.Close()

	user := &User{}

	err = db.QueryRow("select id, name from users where name = ? ", name).Scan(&user.Id, &user.Name)
	if err != nil {
		return user, errors.WithStack(err)
	}
	return user, nil
}

//controller
func getUserByName(name string) (*User, error) {
	return getUser(name)
}

//api
func api() (int, *User) {
	user, err := getUserByName("abc")
	if err != nil {
		// logger.Error("api_fail, err="+err.Error())
		if errors.Is(err, NOTFOUND) {
			return 404, user
		}

		return 500, user
	}

	return 200, user
}

func main() {
	api()
}
