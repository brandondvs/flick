package server

import (
	"encoding/json"
	"net/http"

	"github.com/brandondvs/flick/internal/feature"
	"github.com/brandondvs/flick/internal/store"
)

type Server struct {
	store *store.Memory
	mux   *http.ServeMux
}

type flagResponse struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type createRequest struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type updateRequest struct {
	Enabled bool `json:"enabled"`
}

func New(s *store.Memory) *Server {
	srv := &Server{
		store: s,
		mux:   http.NewServeMux(),
	}
	srv.routes()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/flags", s.handleFlags)
	s.mux.HandleFunc("/flags/", s.handleFlag)
}

// handleFlags routes GET /flags and POST /flags
func (s *Server) handleFlags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listFlags(w, r)
	case http.MethodPost:
		s.createFlag(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleFlag routes GET, PUT, DELETE on /flags/{name}
func (s *Server) handleFlag(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/flags/"):]
	if name == "" {
		http.Error(w, "flag name required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getFlag(w, r, name)
	case http.MethodPut:
		s.updateFlag(w, r, name)
	case http.MethodDelete:
		s.deleteFlag(w, r, name)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// POST /flags
func (s *Server) createFlag(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if existing := s.store.Get(req.Name); existing != nil {
		http.Error(w, "flag already exists", http.StatusConflict)
		return
	}

	flag := feature.Create(req.Name)
	flag.Set(req.Enabled)
	s.store.Store(req.Name, flag)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flagResponse{
		Name:    flag.Name(),
		Enabled: flag.IsEnabled(),
	})
}

// GET /flags
func (s *Server) listFlags(w http.ResponseWriter, r *http.Request) {
	keys := s.store.AllKeys()
	flags := make([]flagResponse, 0, len(keys))

	for _, key := range keys {
		f := s.store.Get(key)
		if f != nil {
			flags = append(flags, flagResponse{
				Name:    f.Name(),
				Enabled: f.IsEnabled(),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flags)
}

// GET /flags/{name}
func (s *Server) getFlag(w http.ResponseWriter, r *http.Request, name string) {
	flag := s.store.Get(name)
	if flag == nil {
		http.Error(w, "flag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flagResponse{
		Name:    flag.Name(),
		Enabled: flag.IsEnabled(),
	})
}

// PUT /flags/{name}
func (s *Server) updateFlag(w http.ResponseWriter, r *http.Request, name string) {
	flag := s.store.Get(name)
	if flag == nil {
		http.Error(w, "flag not found", http.StatusNotFound)
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	flag.Set(req.Enabled)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flagResponse{
		Name:    flag.Name(),
		Enabled: flag.IsEnabled(),
	})
}

// DELETE /flags/{name}
func (s *Server) deleteFlag(w http.ResponseWriter, r *http.Request, name string) {
	flag := s.store.Get(name)
	if flag == nil {
		http.Error(w, "flag not found", http.StatusNotFound)
		return
	}

	s.store.Delete(name)
	w.WriteHeader(http.StatusNoContent)
}
