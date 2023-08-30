package segment

import (
	"context"
	"database/sql"
	"log"
	"maze.io/x/duration"
	"strconv"
	"strings"
	"time"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetSegmentsByUserId(context context.Context, userId uint64) ([]string, error) {
	query := "SELECT s.slug FROM users_segments AS us JOIN segments s on us.segment_id = s.id WHERE us.user_id = $1;"
	slugs := make([]string, 0)
	rows, err := r.db.QueryContext(context, query, userId)
	if err != nil {
		return slugs, err
	}
	for rows.Next() {
		slug := ""
		if err := rows.Scan(&slug); err != nil {
			return slugs, err
		}
		slugs = append(slugs, slug)
	}
	return slugs, err
}

func (r *repository) CreateSegment(context context.Context, slug string) (*Segment, error) {
	var segmentId uint64
	query := "INSERT INTO \"segments\" (slug) VALUES ($1) RETURNING id"
	err := r.db.QueryRowContext(context, query, slug).Scan(&segmentId)
	if err != nil {
		return nil, err
	}
	segment := Segment{
		segmentId,
		slug,
	}
	return &segment, nil
}

func (r *repository) DeleteSegment(context context.Context, slug string) error {
	query := "DELETE FROM \"segments\" WHERE slug = $1"
	_, err := r.db.ExecContext(context, query, slug)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateSegments(context context.Context, req UpdateSegmentsReq) error {
	err := r.AddUserToSegments(context, req)
	if err != nil {
		return err
	}
	err = r.DeleteUserFromSegments(context, req)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) AddUserToSegments(context context.Context, req UpdateSegmentsReq) error {
	if len(req.AddSegments) == 0 {
		return nil
	}
	query := "INSERT INTO users_segments(user_id, segment_id, created_at, ttl) VALUES "
	vals := []interface{}{}
	now := time.Now()
	var ttl *time.Time
	ttl = nil
	if req.Ttl != nil {
		dur, err := duration.ParseDuration(*req.Ttl)
		if err != nil {
			return err
		}
		temp := now.Add(time.Duration(dur))
		ttl = &temp

	}
	for _, slug := range req.AddSegments {
		query += "(?, (SELECT id FROM segments WHERE slug = ?), ?, ?),"
		vals = append(vals, req.UserId, slug, now, ttl)
	}
	query = query[0 : len(query)-1]
	query = ReplaceSQL(query, "?")
	stmt, err := r.db.PrepareContext(context, query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vals...)
	return err
}

func (r *repository) DeleteUserFromSegments(context context.Context, req UpdateSegmentsReq) error {
	if len(req.DeleteSegments) == 0 {
		return nil
	}
	query := "DELETE FROM users_segments WHERE user_id = ? AND segment_id IN (SELECT us.segment_id FROM users_segments us JOIN segments s on s.id = us.segment_id WHERE "
	vals := []interface{}{req.UserId}
	for _, slug := range req.DeleteSegments {
		query += "s.slug = ? OR "
		vals = append(vals, slug)
	}
	query = strings.TrimSuffix(query, "OR ") + ");"
	query = ReplaceSQL(query, "?")
	stmt, err := r.db.PrepareContext(context, query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vals...)
	log.Print(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) ClearExpiredRows(context context.Context) error {
	_, err := r.db.ExecContext(context, "DELETE FROM users_segments WHERE ttl <= now()")
	return err
}

func (r *repository) GetReports(context context.Context, date string) ([]Report, error) {
	query := "SELECT o.user_id, s.slug, o.operation, o.time FROM operations o JOIN segments s on o.segment_id = s.id WHERE to_char(o.time, 'YYYY-MM') = $1;"
	rows, err := r.db.QueryContext(context, query, date)
	reports := make([]Report, 0)
	if err != nil {
		return reports, err
	}
	for rows.Next() {
		var report Report
		if err = rows.Scan(&report.UserId, &report.Segment, &report.Operation, &report.Time); err != nil {
			return reports, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func (r *repository) AddUsersToSegment(context context.Context, usersIds []uint64, segmentId uint64, ttl *duration.Duration) error {
	query := "INSERT INTO users_segments(user_id, segment_id, created_at, ttl) VALUES "
	vals := []interface{}{}
	now := time.Now()
	var deadline *time.Time
	deadline = nil
	if ttl != nil {
		temp := now.Add(time.Duration(*ttl))
		deadline = &temp

	}
	for _, id := range usersIds {
		query += "(?, ?, ?, ?),"
		vals = append(vals, id, segmentId, now, deadline)
	}
	query = query[0 : len(query)-1]
	query = ReplaceSQL(query, "?")
	stmt, err := r.db.PrepareContext(context, query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vals...)
	return err
}

func (r *repository) GetUserIds(context context.Context) ([]uint64, error) {
	query := "SELECT id FROM users;"
	userIds := make([]uint64, 0)
	rows, err := r.db.QueryContext(context, query)
	if err != nil {
		return userIds, err
	}
	for rows.Next() {
		var userId uint64
		if err = rows.Scan(&userId); err != nil {
			return userIds, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, nil
}
