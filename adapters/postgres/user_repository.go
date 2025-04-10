package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
	"github.com/chrikar/chatheon/internal/config"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)",
		user.ID, user.Username, user.PasswordHash)
	return err
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func Connect(cfg config.Config) (*sql.DB, error) {
	connStr := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " sslmode=disable"
	return sql.Open("postgres", connStr)
}
