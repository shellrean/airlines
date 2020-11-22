package mysql

import (
    "context"
    "database/sql"
    "log"
    "fmt"

    "shellrean.com/airlines/domain"
)

type mysqlPlaneRepository struct {
    Conn *sql.DB
}

// NewMysqlPlaneRepository will create an object that represent the plane.Repository interface
func NewMysqlPlaneRepository(Conn *sql.DB) domain.PlaneRepository {
    return &mysqlPlaneRepository{Conn}
}

func (m *mysqlPlaneRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Plane, err error) {
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

    result = make([]domain.Plane, 0)
    for rows.Next() {
        p := domain.Plane{}
        err = rows.Scan(
            &p.ID,
            &p.Name,
            &p.PlaneCode,
            &p.SeatSize,
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

func (m *mysqlPlaneRepository) Fetch(ctx context.Context, num int64) (res []domain.Plane, err error) {
    query := `SELECT id,name,plane_code,seat_size,created_at,updated_at
                    FROM planes ORDER BY created_at LIMIT ?`

    res, err = m.fetch(ctx, query, num)
    if err != nil {
        return nil, err
    }

    return
}

func (m *mysqlPlaneRepository) GetByID(ctx context.Context, id int64) (res domain.Plane, err error){
    query := `SELECT id,name,plane_code,seat_size,created_at,updated_at
                    FROM planes WHERE id=?`
    list, err := m.fetch(ctx, query, id)
    if err != nil {
        return domain.Plane{}, err
    }

    if len(list) > 0 {
        res = list[0]
    } else {
        return domain.Plane{}, domain.ErrNotFound
    }

    return
}

func (m *mysqlPlaneRepository) Store(ctx context.Context, p *domain.Plane) (err error) {
    query := `INSERT planes SET name=?, plane_code=?, seat_size=?, created_at=?, updated_at=?`
    stmt, err := m.Conn.PrepareContext(ctx, query)
    if err != nil {
        return
    }

    res, err := stmt.ExecContext(ctx, p.Name, p.PlaneCode, p.SeatSize, p.CreatedAt, p.UpdatedAt)
    if err != nil {
        return
    }

    lastID, err := res.LastInsertId()
    if err != nil {
        return
    }

    p.ID = lastID
    return
}

func (m *mysqlPlaneRepository) Update(ctx context.Context, p *domain.Plane) (err error) {
    query := `UPDATE planes SET name=?, plane_code=?, seat_size=?, updated_at=? WHERE id=?`
    stmt, err := m.Conn.PrepareContext(ctx, query)
    if err != nil {
        return
    }

    res, err := stmt.ExecContext(ctx, p.Name, p.PlaneCode, p.SeatSize, p.UpdatedAt, p.ID)
    if err != nil {
        return
    }

    affect, err := res.RowsAffected()
    if err != nil {
        return
    }
    if affect != 1 {
        err = fmt.Errorf("Weird Behavior. Total Affected: %d", affect)
        return
    }

    return
}

func (m *mysqlPlaneRepository) Delete(ctx context.Context, id int64) (err error) {
    query := `DELETE FROM planes WHERE id=?`
    stmt, err := m.Conn.PrepareContext(ctx, query)
    if err != nil {
        return
    }

    res, err := stmt.ExecContext(ctx, id)
    if err != nil {
        return
    }

    affect, err := res.RowsAffected()
    if err != nil {
        return
    }

    if affect != 1 {
        err = fmt.Errorf("Weird Behavior. Total Affected: %d", affect)
        return
    }

    return
}