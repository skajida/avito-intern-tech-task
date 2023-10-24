package repository

import (
	"context"
	"io"
	"os"
	mod "service-segs/internal/model"
	c "service-segs/internal/model/constants"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

type CsvRepository struct {
	volumePath string
}

func NewCsvRepository(volPath string) *CsvRepository {
	return &CsvRepository{volumePath: volPath}
}

func (cr CsvRepository) CreateHistoryFile(ctx context.Context, history mod.HistoryCollection) (mod.Filename, error) {
	filename := uuid.New().String() + ".csv"
	clientsFile, err := os.OpenFile(cr.volumePath+filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", c.InternalError
	}
	defer clientsFile.Close()
	if err = gocsv.MarshalFile(&history, clientsFile); err != nil {
		return "", c.InternalError
	}
	return mod.Filename(filename), nil
}

func (cr CsvRepository) DownloadHistoryFile(ctx context.Context, filename mod.Filename) (mod.RawData, error) {
	file, err := os.OpenFile(cr.volumePath+string(filename), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return mod.RawData{}, c.NotFound
	}
	defer file.Close()
	return io.ReadAll(file)
}
