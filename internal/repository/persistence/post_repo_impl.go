package persistence

import (
	"eao/internal/model"
	"eao/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var postRepoMu sync.Mutex

type PostRepositoryImpl struct {
	postFilePath string
}

func NewPostRepository() repository.PostRepository {
	return &PostRepositoryImpl{postFilePath: filepath.Join(".", "public", "post.json")}
}

func (r *PostRepositoryImpl) GetPostList(query *model.PostListQuery) ([]model.Post, int, error) {
	postRepoMu.Lock()
	defer postRepoMu.Unlock()

	posts, err := r.readPosts()
	if err != nil {
		return nil, 0, err
	}

	filtered := filterPostsByKeyword(posts, query.Keyword)
	total := len(filtered)

	start := (query.PageNum - 1) * query.PageSize
	if start >= total {
		return []model.Post{}, total, nil
	}

	end := start + query.PageSize
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *PostRepositoryImpl) GetPostByID(id string) (*model.Post, error) {
	postRepoMu.Lock()
	defer postRepoMu.Unlock()

	posts, err := r.readPosts()
	if err != nil {
		return nil, err
	}

	for i := range posts {
		if posts[i].ID == id {
			post := posts[i]
			return &post, nil
		}
	}

	return nil, fmt.Errorf("post 不存在: %s", id)
}

func (r *PostRepositoryImpl) CreatePost(req *model.CreatePostRequest) (string, error) {
	postRepoMu.Lock()
	defer postRepoMu.Unlock()

	posts, err := r.readPosts()
	if err != nil {
		return "", err
	}

	id := generatePostID()
	posts = append(posts, model.Post{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	})

	if err = r.writePosts(posts); err != nil {
		return "", err
	}

	return id, nil
}

func (r *PostRepositoryImpl) UpdatePost(id string, req *model.UpdatePostRequest) error {
	postRepoMu.Lock()
	defer postRepoMu.Unlock()

	posts, err := r.readPosts()
	if err != nil {
		return err
	}

	found := false
	for i := range posts {
		if posts[i].ID == id {
			if req.Title != nil {
				posts[i].Title = *req.Title
			}
			if req.Content != nil {
				posts[i].Content = *req.Content
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("post 不存在: %s", id)
	}

	return r.writePosts(posts)
}

func (r *PostRepositoryImpl) DeletePost(id string) error {
	postRepoMu.Lock()
	defer postRepoMu.Unlock()

	posts, err := r.readPosts()
	if err != nil {
		return err
	}

	filtered := make([]model.Post, 0, len(posts))
	deleted := false
	for i := range posts {
		if posts[i].ID == id {
			deleted = true
			continue
		}
		filtered = append(filtered, posts[i])
	}

	if !deleted {
		return fmt.Errorf("post 不存在: %s", id)
	}

	return r.writePosts(filtered)
}

func (r *PostRepositoryImpl) readPosts() ([]model.Post, error) {
	data, err := os.ReadFile(r.postFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []model.Post{}, nil
		}
		return nil, fmt.Errorf("读取 post 文件失败: %w", err)
	}

	if len(data) == 0 {
		return []model.Post{}, nil
	}

	posts := make([]model.Post, 0)
	if err = json.Unmarshal(data, &posts); err != nil {
		return nil, fmt.Errorf("解析 post 文件失败: %w", err)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) writePosts(posts []model.Post) error {
	data, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 post 数据失败: %w", err)
	}

	if err = os.WriteFile(r.postFilePath, data, 0o644); err != nil {
		return fmt.Errorf("写入 post 文件失败: %w", err)
	}

	return nil
}

func generatePostID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

func filterPostsByKeyword(posts []model.Post, keyword string) []model.Post {
	kw := strings.ToLower(strings.TrimSpace(keyword))
	if kw == "" {
		return posts
	}

	filtered := make([]model.Post, 0, len(posts))
	for i := range posts {
		title := strings.ToLower(posts[i].Title)
		content := strings.ToLower(posts[i].Content)
		if strings.Contains(title, kw) || strings.Contains(content, kw) {
			filtered = append(filtered, posts[i])
		}
	}

	return filtered
}
