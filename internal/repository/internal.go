package repository

import (
	"context"
	"database/sql"
	"fmt"
	c "service-segs/internal/model/constants"
	"service-segs/internal/model/exchange"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type IRepository struct {
	database *sql.DB
}

func NewIRepository(db *sql.DB) *IRepository {
	return &IRepository{database: db}
}

const errorTemplate = "repo %s: %w"

func (this *IRepository) InsertSegment(ctx context.Context, segId string) error {
	const request = `INSERT INTO segments(tag) VALUES($1);`
	_, err := this.database.ExecContext(ctx, request, segId)
	if err != nil {
		return fmt.Errorf(errorTemplate, "create", err)
	}
	return nil
}

func (this *IRepository) DeleteSegment(ctx context.Context, segId string) error {
	const request = `DELETE FROM segments WHERE tag = $1;`
	res, err := this.database.ExecContext(ctx, request, segId)
	if err != nil {
		return fmt.Errorf(errorTemplate, "delete", err)
	} else if q, _ := res.RowsAffected(); q == 0 {
		return fmt.Errorf(errorTemplate, "delete", c.NotFound)
	}
	return nil
}

func (this *IRepository) SelectBelonging(ctx context.Context, userId int) ([]string, error) {
	const request = `
	SELECT tag
	FROM users_segments u
	INNER JOIN segments s
	ON
		user_id = $1 AND
		(remove_time IS NULL OR NOW() < remove_time)
		AND u.seg_id = s.seg_id;
	`

	rows, err := this.database.QueryContext(ctx, request, userId)
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, "select", err)
	}
	defer rows.Close()

	answer := make([]string, 0)
	var segId string
	for rows.Next() {
		rows.Scan(&segId)
		answer = append(answer, segId)
	}
	return answer, nil
}

const getSegIdReq = `SELECT seg_id FROM segments WHERE tag = $1;`

func (this *IRepository) validateSegments(ctx context.Context, segments []string) error {
	for _, seg := range segments {
		var exists int8
		if err := this.database.QueryRowContext(ctx, getSegIdReq, seg).Scan(&exists); err != nil {
			return c.InvalidSegment
		}
	}
	return nil
}

func (this *IRepository) getSegmentIds(ctx context.Context, segTags []string) (result []int) {
	for _, tag := range segTags {
		var segId int
		this.database.QueryRowContext(ctx, getSegIdReq, tag).Scan(&segId)
		result = append(result, segId)
	}
	return
}

const updReq = `
UPDATE users_segments u
SET remove_time = NOW()
FROM segments s
WHERE
	user_id = $1 AND
	(remove_time IS NULL OR NOW() < remove_time) AND
	u.seg_id = s.seg_id AND
	tag IN ($2);`

func (this *IRepository) UpdateBelonging(
	ctx context.Context,
	userId int,
	addTo, removeFrom []string,
) error {
	const insReq = `
	INSERT INTO users_segments (user_id, seg_id)
	VALUES ($1, $2);
	`

	if this.validateSegments(ctx, addTo) != nil || this.validateSegments(ctx, removeFrom) != nil {
		return c.InvalidSegment
	}
	this.database.ExecContext(ctx, updReq, userId, strings.Join(removeFrom, "', '"))

	segIds := this.getSegmentIds(ctx, addTo)
	for _, id := range segIds {
		this.database.ExecContext(ctx, insReq, userId, id)
	}

	return nil
}

func (this *IRepository) UpdateBelongingTimer(
	ctx context.Context,
	userId int,
	addTo, removeFrom []string,
	before time.Time,
) error {
	const insReq = `
	INSERT INTO users_segments (user_id, seg_id, remove_time)
	VALUES ($1, $2, $3);
	`

	if this.validateSegments(ctx, addTo) != nil || this.validateSegments(ctx, removeFrom) != nil {
		return c.InvalidSegment
	}
	this.database.ExecContext(ctx, updReq, userId, strings.Join(removeFrom, "', '"))

	segIds := this.getSegmentIds(ctx, addTo)
	for _, id := range segIds {
		this.database.ExecContext(ctx, insReq, userId, id, before)
	}

	return nil
}

func (this *IRepository) SelectHistory(
	ctx context.Context,
	userId int,
	from time.Time,
	to time.Time,
) ([]exchange.HistoryEntry, error) {
	const request = `
	SELECT user_id, tag, create_time, remove_time
	FROM segments s INNER JOIN users_segments u
	ON
		s.seg_id = u.seg_id AND
		u.user_id = $1 AND
			($2 < create_time AND create_time < $3 OR
			remove_time IS NOT NULL AND $2 < remove_time AND remove_time < $3);
	`
	rows, err := this.database.QueryContext(ctx, request, userId, from, to)
	if err != nil {
		return []exchange.HistoryEntry{}, fmt.Errorf(errorTemplate, "history", err)
	}
	defer rows.Close()

	result := make([]exchange.HistoryEntry, 0)
	item := &struct {
		userId     int
		segTag     string
		createTime time.Time
		removeTime sql.NullTime
	}{}
	for rows.Next() {
		rows.Scan(&item.userId, &item.segTag, &item.createTime, &item.removeTime)
		if item.createTime.After(from) && item.createTime.Before(to) {
			result = append(result, exchange.HistoryEntry{
				UserId:    item.userId,
				SegTag:    item.segTag,
				Operation: "add",
				Time:      item.createTime,
			})
		}
		if item.removeTime.Valid && item.removeTime.Time.After(from) &&
			item.removeTime.Time.Before(to) {
			result = append(result, exchange.HistoryEntry{
				UserId:    item.userId,
				SegTag:    item.segTag,
				Operation: "remove",
				Time:      item.removeTime.Time,
			})
		}
	}
	return result, nil
}

func (this *IRepository) AutoApply(ctx context.Context, tag string, userIds []int) error {
	const insertReq = `INSERT INTO users_segments(user_id, seg_id) VALUES ($1, $2)`
	var segId int
	this.database.QueryRowContext(ctx, getSegIdReq, tag).Scan(&segId)
	for _, id := range userIds {
		this.database.ExecContext(ctx, insertReq, id, segId)
	}
	return nil
}
