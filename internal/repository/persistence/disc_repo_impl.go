package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"eao/internal/model"
	"eao/internal/repository"
)

type DiscRepositoryImpl struct {
	discFilePath string
}

func NewDiscRepository() repository.DiscRepository {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return &DiscRepositoryImpl{discFilePath: filepath.Join(".", "internal", "data", "disc.json")}
	}

	dataPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "data", "disc.json")
	return &DiscRepositoryImpl{discFilePath: dataPath}
}

func (r *DiscRepositoryImpl) GetDiscList(query *model.DiscListQuery) ([]model.Disc, int, error) {
	discs, err := r.readDiscs()
	if err != nil {
		return nil, 0, err
	}

	total := len(discs)
	start := (query.PageNum - 1) * query.PageSize
	if start >= total {
		return []model.Disc{}, total, nil
	}

	end := start + query.PageSize
	if end > total {
		end = total
	}

	return discs[start:end], total, nil
}

func (r *DiscRepositoryImpl) readDiscs() ([]model.Disc, error) {
	data, err := os.ReadFile(r.discFilePath)
	if err != nil {
		return nil, fmt.Errorf("读取 disc 文件失败: %w", err)
	}

	discs := make([]model.Disc, 0)
	if err = json.Unmarshal(data, &discs); err != nil {
		return nil, fmt.Errorf("解析 disc 文件失败: %w", err)
	}

	return discs, nil
}
