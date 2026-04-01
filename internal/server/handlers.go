package server
import("encoding/json";"net/http";"strconv";"time";"github.com/stockyard-dev/stockyard-sentinel/internal/store")
func(s *Server)handleListTeams(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListTeams();if list==nil{list=[]store.Team{}};writeJSON(w,200,list)}
func(s *Server)handleCreateTeam(w http.ResponseWriter,r *http.Request){
    var t store.Team;json.NewDecoder(r.Body).Decode(&t)
    if t.Name==""{writeError(w,400,"name required");return}
    if err:=s.db.CreateTeam(&t);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,t)}
func(s *Server)handleDeleteTeam(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteTeam(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListMembers(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListMembers(id);if list==nil{list=[]store.Member{}};writeJSON(w,200,list)}
func(s *Server)handleCreateMember(w http.ResponseWriter,r *http.Request){
    id,_:=strconv.ParseInt(r.PathValue("id"),10,64)
    var m store.Member;json.NewDecoder(r.Body).Decode(&m);m.TeamID=id
    if m.Name==""{writeError(w,400,"name required");return}
    if m.TZ==""{m.TZ="UTC"}
    if err:=s.db.CreateMember(&m);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,m)}
func(s *Server)handleDeleteMember(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteMember(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListShifts(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListShifts(id);if list==nil{list=[]store.Shift{}};writeJSON(w,200,list)}
func(s *Server)handleCreateShift(w http.ResponseWriter,r *http.Request){
    id,_:=strconv.ParseInt(r.PathValue("id"),10,64)
    var req struct{MemberID int64 `json:"member_id"`;StartsAt string `json:"starts_at"`;EndsAt string `json:"ends_at"`;Notes string `json:"notes"`}
    json.NewDecoder(r.Body).Decode(&req)
    st,err:=time.Parse("2006-01-02T15:04",req.StartsAt);if err!=nil{writeError(w,400,"invalid starts_at (use 2006-01-02T15:04)");return}
    et,err:=time.Parse("2006-01-02T15:04",req.EndsAt);if err!=nil{writeError(w,400,"invalid ends_at");return}
    sh:=&store.Shift{TeamID:id,MemberID:req.MemberID,StartsAt:st,EndsAt:et,Notes:req.Notes}
    if err:=s.db.CreateShift(sh);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,sh)}
func(s *Server)handleCurrentOnCall(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);shift,_:=s.db.CurrentOnCall(id);if shift==nil{writeJSON(w,200,map[string]string{"status":"no_oncall"});return};writeJSON(w,200,shift)}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){t,_:=s.db.CountTeams();sh,_:=s.db.CountShifts();writeJSON(w,200,map[string]interface{}{"teams":t,"shifts":sh})}
