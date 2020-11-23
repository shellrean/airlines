package main

import (
	"context"

	"shellrean.com/airlines/domain"
)

type mysqlAirportRepository struct {
	Conn *sql.DB
}

// NewMysqlAirportRepository represent domain.AirportRepository interface
func NewMysqlAirportRepository(Conn *sql.DB) domain.AirportRepository {
	return &mysqlAirportRepository{Conn}
}

func (a *mysqlAirportRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Airport, err error) {
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

	result = make([]domain.Airport, 0)
	for rows.Next() {
		p := domain.Airport{}
		err = rows.Scan(
			&p.ID,
			&p.Code,
			&p.Name,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

func (a *mysqlAirportRepository) Fetch(ctx context.Context, num int64) ([]domain.Airport, error) {
	query := `SELECT id,code,name,created_at,updated_at
					FROM airports ORDER BY created_at LIMIT ?`

	res, err = m.fetch(ctx, query, num)
	if err != nil {
		return nil, err
	}

	return res, nil
}