package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"eao/internal/model"
	"eao/internal/repository"
)

type DiscRepositoryImpl struct {
	discFilePath string
}

func NewDiscRepository() repository.DiscRepository {
	return &DiscRepositoryImpl{discFilePath: resolveDiscFilePath()}
}

func resolveDiscFilePath() string {
	exePath, err := os.Executable()
	if err == nil {
		path := filepath.Join(filepath.Dir(exePath), "data", "disc.json")
		if fileExists(path) {
			return path
		}
	}

	wd, err := os.Getwd()
	if err == nil {
		for _, path := range discFilePathCandidates(wd) {
			if fileExists(path) {
				return path
			}
		}
	}

	return filepath.Join("data", "disc.json")
}

func discFilePathCandidates(startDir string) []string {
	candidates := make([]string, 0)
	for dir := startDir; ; dir = filepath.Dir(dir) {
		candidates = append(candidates,
			filepath.Join(dir, "internal", "data", "disc.json"),
			filepath.Join(dir, "data", "disc.json"),
		)

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return candidates
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
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
