package service

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	c "service-segs/internal/model/constants"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

type SegmentsService struct {
	external erepository
	internal irepository
}

func NewSegmentsService(extRepo erepository, intRepo irepository) *SegmentsService {
	return &SegmentsService{external: extRepo, internal: intRepo}
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
	fmt.Println(add, delete)
	if err := this.internal.UpdateBelonging(ctx, id, add, delete); err != nil {
		fmt.Println(fmt.Errorf(errorTemplate, "modify", err))
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
) (string, error) {
	end := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
	history, err := this.internal.SelectHistory(ctx, id, start, end)
	if err != nil {
		return "", fmt.Errorf(errorTemplate, "history", err)
	}
	filename := uuid.New().String() + ".csv"
	clientsFile, err := os.OpenFile(c.VolumePath+filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", c.InternalError
	}
	defer clientsFile.Close()
	if err = gocsv.MarshalFile(&history, clientsFile); err != nil {
		return "", c.InternalError
	}
	return filename, nil
}

func (this *SegmentsService) DownloadFile(ctx context.Context, filename string) ([]byte, error) {
	file, err := os.OpenFile(c.VolumePath+filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return []byte{}, c.NotFound
	}
	defer file.Close()
	return io.ReadAll(file)
}
