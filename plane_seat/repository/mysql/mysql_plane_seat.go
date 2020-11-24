package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"shellrean.com/airlines/domain"
)

type mysqlPlaneSeat struct {
	Conn *sql.DB
}

func NewMysqlPlaneSeatRepository(Conn *sql.DB) domain.PlaneSeatRepository {
	return &mysqlPlaneSeat {
		Conn,
	}
}

func (m *mysqlPlaneSeat) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.PlaneSeat, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()
	}()

	result := make([]domain.PlaneSeat, 0)
	for rows.Next() {
		p := domain.PlaneSeat{}
		planeId := int64(0)
		err = rows.Scan(
			&p.ID,
			&planeId,
			&p.SeatClass,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		p.Plane = domain.Plane{
			ID: planeId,
		}
		result = append(result, p)
	}

	return result, nil
}

func (m *mysqlPlaneSeat) Fetch(ctx context.Context, num int64) ([]domain.PlaneSeat, error) {
	query := `SELECT id,plane_id,seat_class,created_at,updated_at
					FROM plane_seats LIMIT ?`
	
	res, err := m.fetch(ctx, query, num)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mysqlPlaneSeat) GetByPlaneId(ctx context.Context, id int64) ([]domain.PlaneSeat, error) {
	query := `SELECT id,plane_id,seat_class,created_at,updated_at
					FROM plane_seats WHERE plane_id=?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (m *mysqlPlaneSeat) Store(ctx context.Context, seat *domain.PlaneSeat) (error) {
	query := `INSERT seat_class SET plane_id=?,seat_class=?,created_at=?,updated_at=?`
	
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, seat.Plane.ID, seat.SeatClass, seat.CreatedAt, seat.UpdatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	seat.ID = lastID

	return nil
}

func (m *mysqlPlaneSeat) Update(ctx context.Context, seat *domain.PlaneSeat) (error) {
	query := `UPDATE plane_seats SET plane_id=?, seat_class=?, updated_at=? WHERE id=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, seat.Plane.ID, seat.SeatClass, seat.UpdatedAt, seat.ID)
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

func (m *mysqlPlaneSeat) Delete(ctx context.Context, id int64) (error) {
	query := `DELETE FROM plane_seats WHERE id=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
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
		return fmt.Errorf("Weird Behavior. Total Affected: %d", affected)
	}

	return nil
}

