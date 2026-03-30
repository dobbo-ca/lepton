// Package api provides the HTTP REST API server.
package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dobbo-ca/lepton/internal/db"
	"github.com/dobbo-ca/lepton/internal/domain"
)

// Server holds application dependencies and the HTTP router.
type Server struct {
	router *chi.Mux
	db     *db.DB
}

// New creates a new Server wired up with all routes.
func New(database *db.DB) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     database,
	}
	s.routes()
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	r := s.router
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		// Agents
		r.Get("/agents", s.listAgents)
		r.Post("/agents", s.createAgent)
		r.Patch("/agents/{id}", s.updateAgent)
		r.Delete("/agents/{id}", s.deleteAgent)

		// Issues
		r.Get("/issues", s.listIssues)
		r.Post("/issues", s.createIssue)
		r.Patch("/issues/{id}", s.updateIssue)

		// Runs
		r.Get("/runs", s.listRuns)
		r.Get("/runs/{id}", s.getRun)

		// Routines
		r.Get("/routines", s.listRoutines)
		r.Post("/routines", s.createRoutine)

		// Projects
		r.Get("/projects", s.listProjects)

		// Goals
		r.Get("/goals", s.listGoals)
	})
}

// ---- helpers ----

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// companyID returns a placeholder company ID. In production this would come
// from auth middleware extracting the tenant from the request.
func companyID(_ *http.Request) string {
	return "default"
}

// ---- Agents ----

func (s *Server) listAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := s.db.Agents().List(r.Context(), companyID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, agents)
}

func (s *Server) createAgent(w http.ResponseWriter, r *http.Request) {
	var a domain.Agent
	if err := decodeJSON(r, &a); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.ID = newID()
	a.CompanyID = companyID(r)
	a.CreatedAt = time.Now()
	a.UpdatedAt = a.CreatedAt
	if err := s.db.Agents().Create(r.Context(), &a); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (s *Server) updateAgent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := s.db.Agents().Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "agent not found")
		return
	}
	if err := decodeJSON(r, a); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.ID = id
	a.UpdatedAt = time.Now()
	if err := s.db.Agents().Update(r.Context(), a); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) deleteAgent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.db.Agents().Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- Issues ----

func (s *Server) listIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := s.db.Issues().List(r.Context(), companyID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, issues)
}

func (s *Server) createIssue(w http.ResponseWriter, r *http.Request) {
	var i domain.Issue
	if err := decodeJSON(r, &i); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	i.ID = newID()
	i.CompanyID = companyID(r)
	if i.Status == "" {
		i.Status = domain.IssueStatusTodo
	}
	if i.Priority == "" {
		i.Priority = domain.IssuePriorityMedium
	}
	i.CreatedAt = time.Now()
	i.UpdatedAt = i.CreatedAt
	if err := s.db.Issues().Create(r.Context(), &i); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, i)
}

func (s *Server) updateIssue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := s.db.Issues().Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "issue not found")
		return
	}
	var update domain.IssueUpdate
	if err := decodeJSON(r, &update); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	applyIssueUpdate(i, update)
	i.UpdatedAt = time.Now()
	if err := s.db.Issues().Update(r.Context(), i); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, i)
}

func applyIssueUpdate(i *domain.Issue, u domain.IssueUpdate) {
	if u.Title != nil {
		i.Title = *u.Title
	}
	if u.Description != nil {
		i.Description = *u.Description
	}
	if u.Status != nil {
		i.Status = *u.Status
	}
	if u.Priority != nil {
		i.Priority = *u.Priority
	}
	if u.AssigneeAgentID != nil {
		i.AssigneeAgentID = *u.AssigneeAgentID
	}
}

// ---- Runs ----

func (s *Server) listRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := s.db.Runs().List(r.Context(), companyID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, runs)
}

func (s *Server) getRun(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	run, err := s.db.Runs().Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}
	writeJSON(w, http.StatusOK, run)
}

// ---- Routines ----

func (s *Server) listRoutines(w http.ResponseWriter, r *http.Request) {
	routines, err := s.db.Routines().List(r.Context(), companyID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, routines)
}

func (s *Server) createRoutine(w http.ResponseWriter, r *http.Request) {
	var rt domain.Routine
	if err := decodeJSON(r, &rt); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rt.ID = newID()
	rt.CompanyID = companyID(r)
	rt.CreatedAt = time.Now()
	rt.UpdatedAt = rt.CreatedAt
	if err := s.db.Routines().Create(r.Context(), &rt); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rt)
}

// ---- Projects ----

func (s *Server) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.db.Projects().List(r.Context(), companyID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

// ---- Goals ----

func (s *Server) listGoals(w http.ResponseWriter, r *http.Request) {
	goals, err := s.db.Goals().List(r.Context(), companyID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, goals)
}

// newID generates a unique ID using a context-free approach.
func newID() string {
	return generateID(context.Background())
}

func generateID(_ context.Context) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s-%d", hex.EncodeToString(b), time.Now().UnixNano())
}
