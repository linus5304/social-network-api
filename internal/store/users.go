package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exist")
	ErrDuplicateUsername = errors.New("a user with that username already exist")
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	IsActive  bool     `json:"is_active"`
	RoleId    int64    `json:"role_id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Role      Role     `json:"role"`
}

type UserWithToken struct {
	*User
	Token string `json:"token"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash
	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		insert into users (username, password, email, role_id)
		values ($1, $2, $3, (select id from roles where name = $4))
		returning id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	err := tx.QueryRowContext(ctx, query, user.Username, user.Password.hash, user.Email, role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value voilates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		select u.id, username, email, password, created_at, updated_at, r.id, r.name, r.level, r.description 
		from users u
		join roles r on u.role_id = r.id
		where u.id = $1 and u.is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		select id, username, email, password, created_at, updated_at from users
		where email = $1 and is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		// create the user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 1. find user that token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		//2. update user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		//3. clean the invitations
		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Delete(ctx context.Context, userId int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userId); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userId); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId int64) error {
	query := `insert into user_invitations (token, user_id, expiry) values ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(exp))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		select u.id, u.username, u.email, u.created_at, u.is_active
		from users u
		join user_invitations ui on u.id = ui.user_id
		where ui.token = $1 and ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `update users set username = $1, email = $2, is_active = $3 where id = $4`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `delete from user_invitations where user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `delete from users where id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}
