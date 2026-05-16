package router

import (
	"eao/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	config.GlobalConfig = &config.Config{PublicBaseURL: "https://www.painorth.bbroot.com/videos/"}

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

	if len(body.Data) != 37 {
		t.Fatalf("expected 37 videos, got %d", len(body.Data))
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
