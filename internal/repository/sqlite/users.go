package sqlite

import (
	"context"
	"errors"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/jackc/pgx/v4"
)

type UsersRepository struct {
	db *SQLite
}

var _ repository.Users = (*UsersRepository)(nil)

func NewUsersRepository(db *SQLite) repository.Users {
	return &UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	return r.find(ctx, "login", login)
}

func (r *UsersRepository) Create(ctx context.Context, login, password, aesSecret, privateKey string) error {
	sql := `INSERT INTO users (login, password, aes_secret, private_key) VALUES ($1, $2, $3, $4);`

	if _, err := r.db.ExecContext(ctx, sql, login, password, aesSecret, privateKey); err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) find(ctx context.Context, column string, value interface{}) (*models.User, error) {
	sql := `SELECT id, login, password, aes_secret, private_key FROM users WHERE %s = $1;`

	user := new(models.User)

	row := r.db.QueryRowContext(ctx, fmt.Sprintf(sql, column), value)
	if err := r.scanUser(row, user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UsersRepository) scanUser(row pgx.Row, user *models.User) error {
	if err := row.Scan(&user.ID, &user.Login, &user.Password, &user.AesSecret, &user.PrivateKey); err != nil {
		return err
	}

	return nil
}
