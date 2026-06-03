package repository

import (
	"context"
	"mime/multipart"
)

type FileRepository interface {
	List(ctx context.Context) ([]byte, error)
	ListByDir(ctx context.Context, dirName string) ([]byte, error)
	Upload(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, dirName string) ([]byte, error)
}
