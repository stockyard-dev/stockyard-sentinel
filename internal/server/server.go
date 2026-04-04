package server

import (
	"encoding/json"
	"net/http"

	"github.com/stockyard-dev/stockyard-sentinel/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	limits Limits
}

func New(db *store.DB, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), limits: limits}

	s.mux.HandleFunc("GET /api/rules", s.listRules)
	s.mux.HandleFunc("POST /api/rules", s.createRule)
	s.mux.HandleFunc("GET /api/rules/{id}", s.getRule)
	s.mux.HandleFunc("DELETE /api/rules/{id}", s.deleteRule)
	s.mux.HandleFunc("POST /api/rules/{id}/toggle", s.toggleRule)

	s.mux.HandleFunc("GET /api/alerts", s.listAlerts)
	s.mux.HandleFunc("POST /api/alerts/fire", s.fireAlert)
	s.mux.HandleFunc("POST /api/alerts/{id}/ack", s.ackAlert)
	s.mux.HandleFunc("POST /api/alerts/{id}/resolve", s.resolveAlert)

	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{"tier": s.limits.Tier, "upgrade_url": "https://stockyard.dev/sentinel/"})
	})
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }
func wj(w http.ResponseWriter, c int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(c); json.NewEncoder(w).Encode(v) }
func we(w http.ResponseWriter, c int, m string) { wj(w, c, map[string]string{"error": m}) }
func (s *Server) root(w http.ResponseWriter, r *http.Request) { if r.URL.Path != "/" { http.NotFound(w, r); return }; http.Redirect(w, r, "/ui", 302) }

func oe(a []store.AlertRule) []store.AlertRule { if a == nil { return []store.AlertRule{} }; return a }
func oa(a []store.Alert) []store.Alert { if a == nil { return []store.Alert{} }; return a }

func (s *Server) listRules(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"rules": oe(s.db.ListRules())}) }
func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 && len(s.db.ListRules()) >= s.limits.MaxItems { we(w, 402, "Free tier limit reached"); return }
	var rule store.AlertRule; json.NewDecoder(r.Body).Decode(&rule)
	if rule.Name == "" { we(w, 400, "name required"); return }
	rule.Enabled = true
	s.db.CreateRule(&rule); wj(w, 201, s.db.GetRule(rule.ID))
}
func (s *Server) getRule(w http.ResponseWriter, r *http.Request) {
	rule := s.db.GetRule(r.PathValue("id")); if rule == nil { we(w, 404, "not found"); return }; wj(w, 200, rule)
}
func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) { s.db.DeleteRule(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }
func (s *Server) toggleRule(w http.ResponseWriter, r *http.Request) { s.db.ToggleRule(r.PathValue("id")); wj(w, 200, s.db.GetRule(r.PathValue("id"))) }

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status"); wj(w, 200, map[string]any{"alerts": oa(s.db.ListAlerts(status, 100))})
}
func (s *Server) fireAlert(w http.ResponseWriter, r *http.Request) {
	var body struct { RuleID string `json:"rule_id"`; Message string `json:"message"`; Source string `json:"source"` }
	json.NewDecoder(r.Body).Decode(&body)
	a := s.db.Fire(body.RuleID, body.Message, body.Source); wj(w, 201, a)
}
func (s *Server) ackAlert(w http.ResponseWriter, r *http.Request) {
	var body struct { By string `json:"by"` }; json.NewDecoder(r.Body).Decode(&body)
	s.db.Ack(r.PathValue("id"), body.By); wj(w, 200, map[string]string{"status": "acknowledged"})
}
func (s *Server) resolveAlert(w http.ResponseWriter, r *http.Request) {
	s.db.Resolve(r.PathValue("id")); wj(w, 200, map[string]string{"status": "resolved"})
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) }
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats(); wj(w, 200, map[string]any{"service": "sentinel", "status": "ok", "firing": st.Firing, "rules": st.Rules})
}
