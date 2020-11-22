package mysql

import (
    "context"
    "database/sql"
    // "fmt"
    "log"

    "shellrean.com/airlines/domain"
)

type mysqlUserRepository struct {
    Conn *sql.DB
}

// NewMysqlUserRepository will create an object that represent the user.Repository interface
func NewMysqlUserRepository(Conn *sql.DB) domain.UserRepository {
    return &mysqlUserRepository{
        Conn,
    }
}

func (m *mysqlUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
    rows, err := m.Conn.QueryContext(ctx, query, args...)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    defer func() {
        errRow := rows.Close()
        if errRow != nil {
            log.Println(errRow)
        }
    }()

    result = make([]domain.User, 0)
    for rows.Next() {
        u := domain.User{}
        err = rows.Scan(
            &u.ID,
            &u.Email,
            &u.Name,
            &u.Password,
            &u.CreatedAt,
            &u.UpdatedAt,
        )

        if err != nil {
            log.Println(err)
            return nil, err
        }

        result = append(result, u)
    }

    return result, nil
}

func (m *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
    query := `SELECT id,email,name,password,created_at,updated_at
                    FROM users WHERE email = ?`

    list, err := m.fetch(ctx, query, email)
    if err != nil {
        return
    }

    if len(list) > 0 {
        res = list[0]
    } else {
        return res, domain.ErrNotFound
    }

    return
}