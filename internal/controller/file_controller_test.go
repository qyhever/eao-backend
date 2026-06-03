package controller

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"eao/internal/service"

	"github.com/gin-gonic/gin"
)

func TestFileControllerListByDirRequiresDirName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	controller := NewFileController(service.NewFileService(&stubFileRepository{}))
	r.GET("/api/file/listByDir", controller.ListByDir)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/file/listByDir", nil)
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var body ResponseData
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Code != CodeInvalidParam {
		t.Fatalf("expected code %d, got %d", CodeInvalidParam, body.Code)
	}
}

func TestFileControllerUploadRequiresFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	controller := NewFileController(service.NewFileService(&stubFileRepository{}))
	r.POST("/api/file/upload", controller.Upload)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/file/upload", nil)
	r.ServeHTTP(resp, req)

	var body ResponseData
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Code != CodeInvalidParam {
		t.Fatalf("expected code %d, got %d", CodeInvalidParam, body.Code)
	}
}

type stubFileRepository struct{}

func (s *stubFileRepository) List(ctx context.Context) ([]byte, error) {
	return []byte(`{"code":1000,"message":"success","data":[]}`), nil
}

func (s *stubFileRepository) ListByDir(ctx context.Context, dirName string) ([]byte, error) {
	return []byte(`{"code":1000,"message":"success","data":[]}`), nil
}

func (s *stubFileRepository) Upload(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, dirName string) ([]byte, error) {
	return []byte(`{"code":1000,"message":"success","data":{}}`), nil
}
