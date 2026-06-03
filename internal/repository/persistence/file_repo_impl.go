package persistence

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"eao/internal/config"
	"eao/internal/repository"
)

const (
	fileListPath      = "/api/file/list"
	fileListByDirPath = "/api/file/listByDir"
	fileUploadPath    = "/api/file/upload"
)

type FileRepositoryImpl struct {
	baseURL string
	secret  string
	client  *http.Client
	now     func() time.Time
}

func NewFileRepository(cfg config.FileAPIConfig) repository.FileRepository {
	timeout := cfg.TimeoutSeconds
	if timeout <= 0 {
		timeout = 10
	}
	return &FileRepositoryImpl{
		baseURL: strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/"),
		secret:  strings.TrimSpace(cfg.Secret),
		client:  &http.Client{Timeout: time.Duration(timeout) * time.Second},
		now:     time.Now,
	}
}

func (r *FileRepositoryImpl) List(ctx context.Context) ([]byte, error) {
	return r.do(ctx, http.MethodGet, fileListPath, nil, nil, "")
}

func (r *FileRepositoryImpl) ListByDir(ctx context.Context, dirName string) ([]byte, error) {
	query := url.Values{}
	query.Set("dirName", dirName)
	return r.do(ctx, http.MethodGet, fileListByDirPath, query, nil, "")
}

func (r *FileRepositoryImpl) Upload(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, dirName string) ([]byte, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("dirName", dirName); err != nil {
		return nil, fmt.Errorf("构造上传表单失败: %w", err)
	}

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, fmt.Errorf("构造上传文件失败: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭上传表单失败: %w", err)
	}

	return r.do(ctx, http.MethodPost, fileUploadPath, nil, &body, writer.FormDataContentType())
}

func (r *FileRepositoryImpl) do(ctx context.Context, method, path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
	if r.baseURL == "" || r.secret == "" {
		return nil, fmt.Errorf("第三方文件接口配置缺失")
	}

	rawURL := r.baseURL + path
	if encodedQuery := query.Encode(); encodedQuery != "" {
		rawURL += "?" + encodedQuery
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建第三方请求失败: %w", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	timestamp := fmt.Sprintf("%d", r.now().Unix())
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Sign", BuildFileSign(method, path, query.Encode(), timestamp, r.secret))

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用第三方文件接口失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取第三方响应失败: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("第三方文件接口状态异常: %d", resp.StatusCode)
	}

	return respBody, nil
}

func BuildFileSign(method, path, canonicalQuery, timestamp, secret string) string {
	raw := strings.Join([]string{
		strings.ToUpper(method),
		path,
		canonicalQuery,
		timestamp,
		secret,
	}, "\n")

	sum := md5.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}
