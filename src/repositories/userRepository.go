package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kylerequez/go-sample-dashboard/src/models"
)

type UserRepository struct {
	Db    *pgx.Conn
	Table string
}

type userRepository interface {
	GetAllUsers(context.Context)
}

func NewUserRepository(db *pgx.Conn, table string) *UserRepository {
	return &UserRepository{
		Db:    db,
		Table: table,
	}
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	sql := fmt.Sprintf(`
			SELECT
				id,
				name,
				email,
				created_at,
				updated_at
			FROM
				%s
			LIMIT
				50;
		`, ur.Table)

	_, err := ur.Db.Prepare(ctx, sql, sql)
	if err != nil {
		return []models.User{}, err
	}

	rows, err := ur.Db.Query(ctx, sql)
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return []models.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := fmt.Sprintf(`
			SELECT
				id,
				password
			FROM
				%s
			WHERE
				email = $1
			LIMIT
				1;
		`, ur.Table)

	_, err := ur.Db.Prepare(ctx, sql, sql)
	if err != nil {
		return nil, err
	}

	user := models.User{}
	row := ur.Db.QueryRow(ctx, sql,
		email,
	)
	if err := row.Scan(
		&user.ID,
		&user.Password,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user does not exists")
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	sql := fmt.Sprintf(`
			INSERT INTO
				%s 
			(	name,
				email,
				password
			) VALUES (
				$1,
				$2,
				$3
			);
		`, ur.Table)

	_, err := ur.Db.Prepare(ctx, sql, sql)
	if err != nil {
		return err
	}

	res, err := ur.Db.Exec(ctx, sql,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return err
	}

	if count := res.RowsAffected(); count < 1 {
		return errors.New("user was not created")
	}

	return nil
}
