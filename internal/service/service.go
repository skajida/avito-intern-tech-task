package service

import (
	"context"
	"fmt"
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
func (this *SegmentsService) ModifyBelonging(ctx context.Context, id int, add, delete []string) error {
	add = unique(add)
	delete = diff(unique(delete), add)
	if err := this.internal.UpdateBelonging(ctx, id, add, delete); err != nil {
		return fmt.Errorf(errorTemplate, "modify", err)
	}
	return nil
}
