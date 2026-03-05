package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brandondvs/flick/internal/feature"
	"github.com/brandondvs/flick/internal/store"
)

func newTestServer() *Server {
	return New(store.New())
}

func encodeBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatalf("failed to encode request body: %v", err)
	}
	return &buf
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) flagResponse {
	t.Helper()
	var resp flagResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

func decodeListResponse(t *testing.T, w *httptest.ResponseRecorder) []flagResponse {
	t.Helper()
	var resp []flagResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	return resp
}

func TestCreateFlag(t *testing.T) {
	srv := newTestServer()

	body := encodeBody(t, createRequest{Name: "dark-mode", Enabled: true})
	req := httptest.NewRequest(http.MethodPost, "/flags", body)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	resp := decodeResponse(t, w)
	if resp.Name != "dark-mode" {
		t.Errorf("expected name %q, got %q", "dark-mode", resp.Name)
	}
	if !resp.Enabled {
		t.Error("expected flag to be enabled")
	}
}

func TestCreateFlagDisabledByDefault(t *testing.T) {
	srv := newTestServer()

	body := encodeBody(t, createRequest{Name: "beta-ui"})
	req := httptest.NewRequest(http.MethodPost, "/flags", body)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	resp := decodeResponse(t, w)
	if resp.Enabled {
		t.Error("expected flag to be disabled by default")
	}
}

func TestCreateFlagMissingName(t *testing.T) {
	srv := newTestServer()

	body := encodeBody(t, createRequest{Enabled: true})
	req := httptest.NewRequest(http.MethodPost, "/flags", body)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateFlagInvalidBody(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodPost, "/flags", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateFlagDuplicate(t *testing.T) {
	srv := newTestServer()

	body := encodeBody(t, createRequest{Name: "cache"})
	req := httptest.NewRequest(http.MethodPost, "/flags", body)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	body = encodeBody(t, createRequest{Name: "cache"})
	req = httptest.NewRequest(http.MethodPost, "/flags", body)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestGetFlag(t *testing.T) {
	srv := newTestServer()
	flag := feature.Create("notifications")
	flag.Set(true)
	srv.store.Store("notifications", flag)

	req := httptest.NewRequest(http.MethodGet, "/flags/notifications", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	resp := decodeResponse(t, w)
	if resp.Name != "notifications" {
		t.Errorf("expected name %q, got %q", "notifications", resp.Name)
	}
	if !resp.Enabled {
		t.Error("expected flag to be enabled")
	}
}

func TestGetFlagNotFound(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/flags/nonexistent", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetFlagMissingName(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/flags/", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListFlagsEmpty(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	resp := decodeListResponse(t, w)
	if len(resp) != 0 {
		t.Errorf("expected 0 flags, got %d", len(resp))
	}
}

func TestListFlags(t *testing.T) {
	srv := newTestServer()
	srv.store.Store("alpha", feature.Create("alpha"))
	srv.store.Store("bravo", feature.Create("bravo"))

	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	resp := decodeListResponse(t, w)
	if len(resp) != 2 {
		t.Errorf("expected 2 flags, got %d", len(resp))
	}
}

func TestUpdateFlag(t *testing.T) {
	srv := newTestServer()
	srv.store.Store("dark-mode", feature.Create("dark-mode"))

	body := encodeBody(t, updateRequest{Enabled: true})
	req := httptest.NewRequest(http.MethodPut, "/flags/dark-mode", body)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	resp := decodeResponse(t, w)
	if !resp.Enabled {
		t.Error("expected flag to be enabled after update")
	}
}

func TestUpdateFlagDisable(t *testing.T) {
	srv := newTestServer()
	flag := feature.Create("cache")
	flag.Set(true)
	srv.store.Store("cache", flag)

	body := encodeBody(t, updateRequest{Enabled: false})
	req := httptest.NewRequest(http.MethodPut, "/flags/cache", body)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	resp := decodeResponse(t, w)
	if resp.Enabled {
		t.Error("expected flag to be disabled after update")
	}
}

func TestUpdateFlagNotFound(t *testing.T) {
	srv := newTestServer()

	body := encodeBody(t, updateRequest{Enabled: true})
	req := httptest.NewRequest(http.MethodPut, "/flags/nonexistent", body)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateFlagInvalidBody(t *testing.T) {
	srv := newTestServer()
	srv.store.Store("beta", feature.Create("beta"))

	req := httptest.NewRequest(http.MethodPut, "/flags/beta", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteFlag(t *testing.T) {
	srv := newTestServer()
	srv.store.Store("temp", feature.Create("temp"))

	req := httptest.NewRequest(http.MethodDelete, "/flags/temp", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// confirm it's gone
	req = httptest.NewRequest(http.MethodGet, "/flags/temp", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d after delete, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteFlagNotFound(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodDelete, "/flags/nonexistent", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestMethodNotAllowedFlags(t *testing.T) {
	srv := newTestServer()

	methods := []string{http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		req := httptest.NewRequest(method, "/flags", nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s /flags: expected status %d, got %d", method, http.StatusMethodNotAllowed, w.Code)
		}
	}
}

func TestMethodNotAllowedFlag(t *testing.T) {
	srv := newTestServer()
	srv.store.Store("test", feature.Create("test"))

	req := httptest.NewRequest(http.MethodPatch, "/flags/test", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestContentTypeJSON(t *testing.T) {
	srv := newTestServer()
	srv.store.Store("check", feature.Create("check"))

	endpoints := []struct {
		method string
		path   string
		body   any
	}{
		{http.MethodGet, "/flags", nil},
		{http.MethodGet, "/flags/check", nil},
		{http.MethodPost, "/flags", createRequest{Name: "new"}},
		{http.MethodPut, "/flags/check", updateRequest{Enabled: true}},
	}

	for _, ep := range endpoints {
		var req *http.Request
		if ep.body != nil {
			req = httptest.NewRequest(ep.method, ep.path, encodeBody(t, ep.body))
		} else {
			req = httptest.NewRequest(ep.method, ep.path, nil)
		}
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)

		ct := w.Header().Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("%s %s: expected Content-Type %q, got %q", ep.method, ep.path, "application/json", ct)
		}
	}
}

func TestFullCRUDLifecycle(t *testing.T) {
	srv := newTestServer()

	// create
	body := encodeBody(t, createRequest{Name: "lifecycle", Enabled: false})
	req := httptest.NewRequest(http.MethodPost, "/flags", body)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected %d, got %d", http.StatusCreated, w.Code)
	}

	// read
	req = httptest.NewRequest(http.MethodGet, "/flags/lifecycle", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("read: expected %d, got %d", http.StatusOK, w.Code)
	}
	resp := decodeResponse(t, w)
	if resp.Enabled {
		t.Error("read: expected flag to be disabled")
	}

	// update
	body = encodeBody(t, updateRequest{Enabled: true})
	req = httptest.NewRequest(http.MethodPut, "/flags/lifecycle", body)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update: expected %d, got %d", http.StatusOK, w.Code)
	}
	resp = decodeResponse(t, w)
	if !resp.Enabled {
		t.Error("update: expected flag to be enabled")
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/flags/lifecycle", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete: expected %d, got %d", http.StatusNoContent, w.Code)
	}

	// confirm gone
	req = httptest.NewRequest(http.MethodGet, "/flags/lifecycle", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("after delete: expected %d, got %d", http.StatusNotFound, w.Code)
	}
}
