package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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
		return fmt.Errorf(errorTemplate, "delete", fmt.Errorf("not found")) // TODO spec custom errors
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

func (this *IRepository) validateSegments(ctx context.Context, segments []string) error {
	const request = `SELECT seg_id FROM segments WHERE tag = $1;`
	for _, seg := range segments {
		var exists int8
		if err := this.database.QueryRowContext(ctx, request, seg).Scan(&exists); err != nil {
			return fmt.Errorf("invalid segment")
		}
	}
	return nil
}

func (this *IRepository) getSegmentIds(ctx context.Context, segTags []string) (result []int) {
	const request = `SELECT seg_id FROM segments WHERE tag = $1;`
	for _, tag := range segTags {
		var segId int
		this.database.QueryRowContext(ctx, request, tag).Scan(&segId)
		result = append(result, segId)
	}
	return
}

func (this *IRepository) UpdateBelonging(
	ctx context.Context,
	userId int,
	addTo, removeFrom []string,
) error {
	const (
		updReq = `
	UPDATE users_segments u
	SET remove_time = NOW()
	FROM segments s
	WHERE
		user_id = $1 AND
		(remove_time IS NULL OR NOW() < remove_time) AND
		u.seg_id = s.seg_id AND
		tag IN ($2);
	`
		insReq = `
	INSERT INTO users_segments (user_id, seg_id)
	VALUES ($1, $2);
	`
	)

	if this.validateSegments(ctx, addTo) != nil || this.validateSegments(ctx, removeFrom) != nil {
		return fmt.Errorf("invalid segment")
	}
	this.database.ExecContext(ctx, updReq, strconv.Itoa(userId), strings.Join(removeFrom, "', '"))

	segIds := this.getSegmentIds(ctx, addTo)
	for _, id := range segIds {
		this.database.ExecContext(ctx, insReq, userId, id)
	}

	return nil
}
