package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	*pgxpool.Pool
}

func NewPostgres(ctx context.Context, address string) (*Postgres, error) {
	pool, err := pgxpool.Connect(ctx, address)
	if err != nil {
		return nil, err
	}

	db := &Postgres{
		pool,
	}

	go func() {
		<-ctx.Done()
		db.Close()
	}()

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err = db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

func (p *Postgres) migrate() error {
	if err := p.createSettingsTable(); err != nil {
		return err
	}

	if err := p.createUsersTable(); err != nil {
		return err
	}

	if err := p.createCredsSecretsTable(); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) createSettingsTable() error {
	query := `CREATE TABLE IF NOT EXISTS settings (
		key   TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL
	)`
	if _, err := p.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) createUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id          SERIAL PRIMARY KEY,
		login       TEXT NOT NULL UNIQUE,
		password    TEXT NOT NULL,
		aes_secret  TEXT NOT NULL,
		private_key TEXT NOT NULL
	)`
	if _, err := p.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) createCredsSecretsTable() error {
	query := `CREATE TABLE IF NOT EXISTS creds_secrets (
		id              SERIAL PRIMARY KEY,
		uid				BIGINT NOT NULL,
		website         TEXT NOT NULL,
		login           TEXT NOT NULL,
		enc_password    TEXT NOT NULL,
		additional_data TEXT NOT NULL,
		user_id			INTEGER NOT NULL,
		UNIQUE (uid),
		UNIQUE (website, login, user_id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`
	if _, err := p.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}
