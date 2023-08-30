package segment

import (
	"context"
	"math/rand"
	"maze.io/x/duration"
	"os"
	"strconv"
	"time"
)

type service struct {
	Repository
	TtlCheckPeriod time.Duration
}

func NewService(repository Repository) Service {
	ttlCheckPeriodString := os.Getenv("TTL_CHECK_PERIOD")
	ttlCheckPeriod, err := time.ParseDuration(ttlCheckPeriodString)
	if err != nil {
		ttlCheckPeriod = time.Hour
	}
	return &service{
		repository,
		ttlCheckPeriod,
	}
}

func (s *service) CreateSegment(context context.Context, slug string) (*Segment, error) {
	segment, err := s.Repository.CreateSegment(context, slug)
	return segment, err
}

func (s *service) CreateSegmentWithUsers(context context.Context, req CreateSegmentWithUsersReq) error {
	segment, err := s.CreateSegment(context, req.Slug)
	if err != nil {
		return err
	}
	segmentId := segment.Id
	userIds, err := s.Repository.GetUserIds(context)
	if err != nil {
		return err
	}
	usersCount := len(userIds)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userIds), func(i, j int) { userIds[i], userIds[j] = userIds[j], userIds[i] })
	upperBound := usersCount * req.Percent / 100
	if req.Percent == 0 {
		upperBound = 0
	} else if req.Percent != 100 {
		upperBound++
	}
	userIds = userIds[0 : (usersCount*req.Percent/100)+1]
	var ttl *duration.Duration
	ttl = nil
	if req.Ttl != nil {
		temp, err := duration.ParseDuration(*req.Ttl)
		if err != nil {
			return err
		}
		ttl = &temp
	}
	err = s.Repository.AddUsersToSegment(context, userIds, segmentId, ttl)
	return err
}

func (s *service) DeleteSegment(context context.Context, slug string) error {
	err := s.Repository.DeleteSegment(context, slug)
	return err
}

func (s *service) UpdateSegments(context context.Context, req UpdateSegmentsReq) error {
	err := s.Repository.UpdateSegments(context, req)
	return err
}

func (s *service) TtlService(context context.Context) error {
	go func() {
		err := func() error {
			for {
				select {
				case <-context.Done():
					return nil
				default:
					err := s.Repository.ClearExpiredRows(context)
					if err != nil {
						return err
					}
				}
				time.Sleep(s.TtlCheckPeriod)
			}
		}()
		if err != nil {
		}
	}()

	return nil
}

func (s *service) GetReports(context context.Context, date string) (string, error) {
	reports, err := s.Repository.GetReports(context, date)
	if err != nil {
		return "", err
	}
	rand.Seed(time.Now().Unix())
	var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	filename := string(b)
	f, err := os.Create(filename + ".csv")
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	for _, report := range reports {
		_, err := f.WriteString(strconv.FormatUint(report.UserId, 10) + ";" + report.Segment + ";" + report.Operation + ";" + report.Time.String() + "\n")
		if err != nil {
			return "", err
		}
	}
	return "0.0.0.0:8081/download/" + filename, nil
}
