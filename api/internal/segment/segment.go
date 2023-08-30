package segment

import (
	"context"
	"maze.io/x/duration"
	"time"
)

type Segment struct {
	Id   uint64 `json:"id"`
	Slug string `json:"slug"`
}

type UpdateSegmentsReq struct {
	AddSegments    []string `json:"add_segments"`
	DeleteSegments []string `json:"delete_segments"`
	UserId         uint64   `json:"user_id"`
	Ttl            *string  `json:"ttl"`
}

type CreateSegmentWithUsersReq struct {
	Slug    string  `json:"slug"`
	Percent int     `json:"percent"`
	Ttl     *string `json:"ttl"`
}

type Report struct {
	UserId    uint64
	Segment   string
	Operation string
	Time      time.Time
}

type Repository interface {
	CreateSegment(context context.Context, slug string) (*Segment, error)
	DeleteSegment(context context.Context, slug string) error
	UpdateSegments(context context.Context, req UpdateSegmentsReq) error
	ClearExpiredRows(context context.Context) error
	GetSegmentsByUserId(context context.Context, userId uint64) ([]string, error)
	GetReports(context context.Context, date string) ([]Report, error)
	GetUserIds(context context.Context) ([]uint64, error)
	AddUsersToSegment(context context.Context, usersIds []uint64, segmentId uint64, ttl *duration.Duration) error
}

type Service interface {
	CreateSegment(context context.Context, slug string) (*Segment, error)
	CreateSegmentWithUsers(context context.Context, req CreateSegmentWithUsersReq) error
	DeleteSegment(context context.Context, slug string) error
	UpdateSegments(context context.Context, req UpdateSegmentsReq) error
	TtlService(context context.Context) error
	GetSegmentsByUserId(context context.Context, userId uint64) ([]string, error)
	GetReports(context context.Context, date string) (string, error)
}
