package persistence

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"eao/internal/config"
)

func TestBuildFileSignWithDocumentExample(t *testing.T) {
	sign := BuildFileSign(http.MethodGet, fileListByDirPath, "dirName=avatars", "1762012800", "test-secret")
	if sign != "55a8172f810a2381fd751b11a0cb3564" {
		t.Fatalf("unexpected sign: %s", sign)
	}
}

func TestBuildFileSignUsesEncodedSortedQuery(t *testing.T) {
	query := url.Values{}
	query.Set("z", "last")
	query.Set("dirName", "avatars")
	query.Set("a", "first value")

	encoded := query.Encode()
	if encoded != "a=first+value&dirName=avatars&z=last" {
		t.Fatalf("unexpected encoded query: %s", encoded)
	}

	sign1 := BuildFileSign(http.MethodGet, fileListByDirPath, encoded, "1762012800", "test-secret")
	sign2 := BuildFileSign(http.MethodGet, fileListByDirPath, query.Encode(), "1762012800", "test-secret")
	if sign1 != sign2 {
		t.Fatalf("expected stable signs, got %s and %s", sign1, sign2)
	}
}

func TestBuildFileSignForUploadUsesEmptyQuery(t *testing.T) {
	sign := BuildFileSign(http.MethodPost, fileUploadPath, "", "1762012800", "test-secret")
	signWithDirName := BuildFileSign(http.MethodPost, fileUploadPath, "dirName=avatars", "1762012800", "test-secret")
	if sign == signWithDirName {
		t.Fatal("expected dirName query to change sign, upload must use empty query")
	}
}

func TestFileRepositoryListByDirProxiesSignedRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != fileListByDirPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "dirName=avatars" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		if r.Header.Get("X-Timestamp") != "1762012800" {
			t.Fatalf("unexpected timestamp: %s", r.Header.Get("X-Timestamp"))
		}
		if r.Header.Get("X-Sign") != "55a8172f810a2381fd751b11a0cb3564" {
			t.Fatalf("unexpected sign: %s", r.Header.Get("X-Sign"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1000,"message":"success","data":["a.png"]}`))
	}))
	defer server.Close()

	repo := newTestFileRepository(server.URL)
	body, err := repo.ListByDir(context.Background(), "avatars")
	if err != nil {
		t.Fatalf("ListByDir failed: %v", err)
	}
	if string(body) != `{"code":1000,"message":"success","data":["a.png"]}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestFileRepositoryUploadSignsWithoutDirNameQuery(t *testing.T) {
	expectedSign := BuildFileSign(http.MethodPost, fileUploadPath, "", "1762012800", "test-secret")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != fileUploadPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		if r.Header.Get("X-Sign") != expectedSign {
			t.Fatalf("unexpected sign: %s", r.Header.Get("X-Sign"))
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Fatalf("parse multipart failed: %v", err)
		}
		if r.FormValue("dirName") != "avatars" {
			t.Fatalf("unexpected dirName: %s", r.FormValue("dirName"))
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("missing file: %v", err)
		}
		_ = file.Close()
		_, _ = w.Write([]byte(`{"code":1000,"message":"success","data":{"url":"x"}}`))
	}))
	defer server.Close()

	repo := newTestFileRepository(server.URL)
	file := multipartFileFromString(t, "hello")
	defer file.Close()
	body, err := repo.Upload(context.Background(), file, &multipart.FileHeader{Filename: "a.txt"}, "avatars")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}
	if !strings.Contains(string(body), `"url":"x"`) {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestFileRepositoryMapsNon2xxToError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":1005}`))
	}))
	defer server.Close()

	repo := newTestFileRepository(server.URL)
	_, err := repo.List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("expected status in error, got: %v", err)
	}
}

func newTestFileRepository(baseURL string) *FileRepositoryImpl {
	repo := NewFileRepository(config.FileAPIConfig{
		BaseURL:        baseURL,
		Secret:         "test-secret",
		TimeoutSeconds: 1,
	}).(*FileRepositoryImpl)
	repo.now = func() time.Time { return time.Unix(1762012800, 0) }
	return repo
}

func multipartFileFromString(t *testing.T, content string) multipart.File {
	t.Helper()
	return readSeekCloser{Reader: strings.NewReader(content)}
}

type readSeekCloser struct {
	*strings.Reader
}

func (r readSeekCloser) Close() error {
	return nil
}
