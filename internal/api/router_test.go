package router

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"eao/internal/config"
	jwtpkg "eao/internal/pkg/jwt"
)

type responseEnvelope struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []videoRecord `json:"data"`
}

type videoRecord struct {
	FileName  string `json:"fileName"`
	VideoName string `json:"videoName"`
}

func TestSetupRouterReturnsVideoList(t *testing.T) {
	config.GlobalConfig = testRouterConfig("")

	r := SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/video", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var body responseEnvelope
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}

	if body.Code != 1000 {
		t.Fatalf("expected response code 1000, got %d", body.Code)
	}

	if len(body.Data) != 49 {
		t.Fatalf("expected 49 videos, got %d", len(body.Data))
	}

	if !strings.HasPrefix(body.Data[0].FileName, "https://www.painorth.bbroot.com/videos/") {
		t.Fatalf("unexpected fileName prefix: %s", body.Data[0].FileName)
	}

	if body.Data[0].FileName != "https://www.painorth.bbroot.com/videos/24u7qivyunz.mp4" {
		t.Fatalf("unexpected first fileName: %s", body.Data[0].FileName)
	}

	if body.Data[0].VideoName != "爸爸带着女儿买烧鸡" {
		t.Fatalf("unexpected first videoName: %s", body.Data[0].VideoName)
	}
}

func TestSetupRouterProxiesFileList(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/file/list" {
			t.Fatalf("unexpected upstream path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Timestamp") == "" {
			t.Fatal("missing X-Timestamp")
		}
		if r.Header.Get("X-Sign") == "" {
			t.Fatal("missing X-Sign")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1000,"message":"success","data":["a.txt"]}`))
	}))
	defer upstream.Close()
	config.GlobalConfig = testRouterConfig(upstream.URL)

	r := SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/file/list", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if strings.TrimSpace(resp.Body.String()) != `{"code":1000,"message":"success","data":["a.txt"]}` {
		t.Fatalf("unexpected body: %s", resp.Body.String())
	}
}

func TestSetupRouterFileListByDirMissingDirName(t *testing.T) {
	config.GlobalConfig = testRouterConfig("")
	r := SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/file/listByDir", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var body struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Code != 1001 {
		t.Fatalf("expected code 1001, got %d", body.Code)
	}
}

func TestSetupRouterFileUploadMissingDirName(t *testing.T) {
	config.GlobalConfig = testRouterConfig("")
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", "a.txt")
	if err != nil {
		t.Fatalf("create form file failed: %v", err)
	}
	_, _ = part.Write([]byte("hello"))
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer failed: %v", err)
	}

	r := SetupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/file/upload", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var body struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Code != 1001 {
		t.Fatalf("expected code 1001, got %d", body.Code)
	}
}

func TestSetupRouterAdminRequiresToken(t *testing.T) {
	config.GlobalConfig = testRouterConfig("")
	r := SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/1", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestSetupRouterAdminInvalidID(t *testing.T) {
	config.GlobalConfig = testRouterConfig("")
	token := testAccessToken(t)
	r := SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/bad", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var body struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Code != 1001 {
		t.Fatalf("expected code 1001, got %d", body.Code)
	}
}

func TestSetupRouterAdminSuccessEnvelope(t *testing.T) {
	config.GlobalConfig = testRouterConfig("")
	token := testAccessToken(t)
	r := SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var body struct {
		Code int `json:"code"`
		Data struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
			Status   string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Code != 1000 {
		t.Fatalf("expected code 1000, got %d, body: %s", body.Code, resp.Body.String())
	}
	if body.Data.ID != 1 || body.Data.Username != "admin" || body.Data.Status != "active" {
		t.Fatalf("unexpected admin data: %+v", body.Data)
	}
}

func testRouterConfig(fileAPIBaseURL string) *config.Config {
	if fileAPIBaseURL == "" {
		fileAPIBaseURL = "http://localhost:6301"
	}
	return &config.Config{
		PublicBaseURL: "https://www.painorth.bbroot.com/videos/",
		JWT: config.JWTConfig{
			Secret:           "test-secret",
			AccessExpiresIn:  "1h",
			RefreshExpiresIn: "24h",
		},
		Auth: config.AuthConfig{
			Admin: config.AdminSeedConfig{
				Username: "admin",
				Password: "password",
				Name:     "管理员",
			},
		},
		ThirdParty: config.ThirdPartyConfig{
			FileAPI: config.FileAPIConfig{
				BaseURL:        fileAPIBaseURL,
				Secret:         "test-secret",
				TimeoutSeconds: 1,
			},
		},
	}
}

func testAccessToken(t *testing.T) string {
	t.Helper()
	token, _, err := jwtpkg.GenToken(1, false)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	return token
}
