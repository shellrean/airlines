package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
	rows, err := a.Conn.QueryContext(ctx, query, args...)
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

	res, err := a.fetch(ctx, query, num)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *mysqlAirportRepository) GetByID(ctx context.Context, id int64) (domain.Airport, error) {
	query := `SELECT id,code,name,created_at,updated_at
					FROM airports WHERE id=?`
	
	list, err := a.fetch(ctx, query, id)
	if err != nil {
		return domain.Airport{}, err
	}

	var res domain.Airport

	if len(list) > 0 {
		res = list[0]
	} else {
		return domain.Airport{}, err
	}

	return res, nil
}

func (a *mysqlAirportRepository) Store(ctx context.Context, airport *domain.Airport) (error) {
	query := `INSERT airports SET code=?, name=?, created_at=?, updated_at=?`
	
	stmt, err := a.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, airport.Code, airport.Name, airport.CreatedAt, airport.UpdatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	airport.ID = lastID
	
	return nil
}

func (a *mysqlAirportRepository) Update(ctx context.Context, airport *domain.Airport) (error) {
	query := `UPDATE airports SET code=?, name=?, updated_at=? WHERE id=?`

	stmt, err := a.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, airport.Code, airport.Name, airport.UpdatedAt, airport.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return fmt.Errorf("Weird Behavior. Total Affected: %d", affected)
	}

	return nil
}

func (a *mysqlAirportRepository) Delete(ctx context.Context, id int64) (error) {
	query := `DELETE FROM airports WHERE id=?`

	stmt, err := a.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return fmt.Errorf("Weird Behavior rows affected: %d", affected)
	}

	return nil
}