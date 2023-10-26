package usersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Doittikorn/go-e-commerce/modules/users"
	"github.com/jmoiron/sqlx"
)

type InsertUserImpl interface {
	Customer() (InsertUserImpl, error)
	Admin() (InsertUserImpl, error)
	Result() (*users.UserPassport, error)
}

// struct ตัวแม่
type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sqlx.DB
}

// โรงงานใหญ่
func InsertUser(db *sqlx.DB, req *users.UserRegisterReq, isAdmin bool) InsertUserImpl {
	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

// โรงงานเล็ก
type customer struct {
	*userReq
}

// โรงงานเล็ก
type admin struct {
	*userReq
}

func newCustomer(db *sqlx.DB, req *users.UserRegisterReq) InsertUserImpl {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func newAdmin(db *sqlx.DB, req *users.UserRegisterReq) InsertUserImpl {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (f *userReq) Customer() (InsertUserImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 1)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}

	return f, nil
}

func (f *userReq) Admin() (InsertUserImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 2)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	return f, nil

}

func (f *userReq) Result() (*users.UserPassport, error) {

	query := `
	SELECT
		json_build_object(
			'user',"t",
			'token', null
		)
	FROM (
		SELECT
			"u"."id",
			"u"."email",
			"u"."username",
			"u"."role_id"
		FROM "users" u
		WHERE "u"."id" = $1
	) as "t"`

	data := make([]byte, 0)
	if err := f.db.Get(&data, query, f.id); err != nil {
		return nil, fmt.Errorf("query user failed: %v", err)
	}

	user := new(users.UserPassport)

	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}

	return user, nil
}
