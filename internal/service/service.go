package service

import (
	"context"
	"fmt"
	"math/rand"
	mod "service-segs/internal/model"
	"time"
)

type SegmentsService struct {
	external  erepository
	internal  irepository
	csvVolume csvRepository
}

func NewSegmentsService(extRepo erepository, intRepo irepository, csvVol csvRepository) *SegmentsService {
	return &SegmentsService{external: extRepo, internal: intRepo, csvVolume: csvVol}
}

const errorTemplate = "svc %s: %w"

func (this *SegmentsService) AddSegment(ctx context.Context, tag string) error {
	if err := this.internal.InsertSegment(ctx, tag); err != nil {
		return fmt.Errorf(errorTemplate, "add", err)
	}
	return nil
}

func (this *SegmentsService) AddSegmentAutoApply(
	ctx context.Context,
	tag string,
	propotion float64,
) error {
	if err := this.internal.InsertSegment(ctx, tag); err != nil {
		return fmt.Errorf(errorTemplate, "add", err)
	}

	gen := rand.New(rand.NewSource(time.Now().Unix()))
	totalUsers := this.external.Count(ctx)
	userIds := make([]int, 0)
	for i := 0; i < int(propotion*float64(totalUsers)); i++ {
		id := gen.Int() % totalUsers
		if this.external.Exists(ctx, id) {
			userIds = append(userIds, id)
		}
	}
	if err := this.internal.AutoApply(ctx, tag, userIds); err != nil {
		return fmt.Errorf(errorTemplate, "auto", err)
	}
	return nil
}

func (this *SegmentsService) DeleteSegment(ctx context.Context, tag string) error {
	if err := this.internal.DeleteSegment(ctx, tag); err != nil {
		return fmt.Errorf(errorTemplate, "delete", err)
	}
	return nil
}

func (this *SegmentsService) ReadBelonging(ctx context.Context, id int) ([]string, error) {
	res, err := this.internal.SelectBelonging(ctx, id)
	if err != nil {
		return []string{}, fmt.Errorf(errorTemplate, "read", err)
	}
	return res, nil
}

func unique(slice []string) (result []string) {
	found := make(map[string]struct{}, len(slice))
	for _, str := range slice {
		if _, exists := found[str]; !exists {
			found[str] = struct{}{}
			result = append(result, str)
		}
	}
	return
}

func diff(minuend, subtrahend []string) (result []string) {
	sub := make(map[string]struct{}, len(subtrahend))
	for _, str := range subtrahend {
		sub[str] = struct{}{}
	}
	for _, str := range minuend {
		if _, exists := sub[str]; !exists {
			result = append(result, str)
		}
	}
	return
}

// deleting FIRST adding AFTER
func (this *SegmentsService) ModifyBelonging(
	ctx context.Context,
	id int,
	add, delete []string,
) error {
	add = unique(add)
	delete = diff(unique(delete), add)
	exist, selErr := this.internal.SelectBelonging(ctx, id)
	if selErr != nil {
		return fmt.Errorf(errorTemplate, "modify", selErr)
	}
	add = diff(add, exist)
	if err := this.internal.UpdateBelonging(ctx, id, add, delete); err != nil {
		return fmt.Errorf(errorTemplate, "modify", err)
	}
	return nil
}

func (this *SegmentsService) ModifyBelongingTimer(
	ctx context.Context,
	id int,
	add, delete []string,
	before time.Time,
) error {
	add = unique(add)
	delete = diff(unique(delete), add)
	exist, selErr := this.internal.SelectBelonging(ctx, id)
	if selErr != nil {
		return fmt.Errorf(errorTemplate, "modify timer", selErr)
	}
	add = diff(add, exist)
	if err := this.internal.UpdateBelongingTimer(ctx, id, add, delete, before); err != nil {
		return fmt.Errorf(errorTemplate, "modify", err)
	}
	return nil
}

func (this *SegmentsService) RequestHistory(
	ctx context.Context,
	id int,
	start time.Time,
) (mod.Filename, error) {
	end := time.Date(start.Year(), start.Month()+1, 1, 0, 0, 0, 0, time.Local)
	history, err := this.internal.SelectHistory(ctx, id, start, end)
	if err != nil {
		return "", fmt.Errorf(errorTemplate, "history", err)
	}
	return this.csvVolume.CreateHistoryFile(ctx, history)
}

func (this *SegmentsService) DownloadFile(ctx context.Context, filename mod.Filename) (mod.RawData, error) {
	return this.csvVolume.DownloadHistoryFile(ctx, filename)
}
