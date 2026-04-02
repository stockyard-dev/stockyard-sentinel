package store
import ("database/sql";"fmt";"os";"path/filepath";"strings";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type AlertRule struct{ID string `json:"id"`;Name string `json:"name"`;Severity string `json:"severity"`;Description string `json:"description,omitempty"`;Channel string `json:"channel,omitempty"`;Enabled bool `json:"enabled"`;CreatedAt string `json:"created_at"`;FireCount int `json:"fire_count"`;LastFired string `json:"last_fired,omitempty"`}
type Alert struct{ID string `json:"id"`;RuleID string `json:"rule_id"`;RuleName string `json:"rule_name,omitempty"`;Severity string `json:"severity"`;Status string `json:"status"`;Message string `json:"message,omitempty"`;Source string `json:"source,omitempty"`;FiredAt string `json:"fired_at"`;AckedAt string `json:"acked_at,omitempty"`;ResolvedAt string `json:"resolved_at,omitempty"`;AckedBy string `json:"acked_by,omitempty"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"sentinel.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
for _,q:=range[]string{
`CREATE TABLE IF NOT EXISTS rules(id TEXT PRIMARY KEY,name TEXT NOT NULL,severity TEXT DEFAULT 'warning',description TEXT DEFAULT '',channel TEXT DEFAULT '',enabled INTEGER DEFAULT 1,created_at TEXT DEFAULT(datetime('now')))`,
`CREATE TABLE IF NOT EXISTS alerts(id TEXT PRIMARY KEY,rule_id TEXT DEFAULT '',rule_name TEXT DEFAULT '',severity TEXT DEFAULT 'warning',status TEXT DEFAULT 'firing',message TEXT DEFAULT '',source TEXT DEFAULT '',fired_at TEXT DEFAULT(datetime('now')),acked_at TEXT DEFAULT '',resolved_at TEXT DEFAULT '',acked_by TEXT DEFAULT '')`,
`CREATE INDEX IF NOT EXISTS idx_alerts_rule ON alerts(rule_id)`,
`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status)`,
}{if _,err:=db.Exec(q);err!=nil{return nil,fmt.Errorf("migrate: %w",err)}};return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)CreateRule(r *AlertRule)error{r.ID=genID();r.CreatedAt=now();if r.Severity==""{r.Severity="warning"};en:=1;if!r.Enabled{en=0}
_,err:=d.db.Exec(`INSERT INTO rules(id,name,severity,description,channel,enabled,created_at)VALUES(?,?,?,?,?,?,?)`,r.ID,r.Name,r.Severity,r.Description,r.Channel,en,r.CreatedAt);return err}
func(d *DB)GetRule(id string)*AlertRule{var r AlertRule;var en int;if d.db.QueryRow(`SELECT id,name,severity,description,channel,enabled,created_at FROM rules WHERE id=?`,id).Scan(&r.ID,&r.Name,&r.Severity,&r.Description,&r.Channel,&en,&r.CreatedAt)!=nil{return nil};r.Enabled=en==1
d.db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE rule_id=?`,r.ID).Scan(&r.FireCount);d.db.QueryRow(`SELECT fired_at FROM alerts WHERE rule_id=? ORDER BY fired_at DESC LIMIT 1`,r.ID).Scan(&r.LastFired);return &r}
func(d *DB)ListRules()[]AlertRule{rows,_:=d.db.Query(`SELECT id,name,severity,description,channel,enabled,created_at FROM rules ORDER BY name`);if rows==nil{return nil};defer rows.Close()
var o []AlertRule;for rows.Next(){var r AlertRule;var en int;rows.Scan(&r.ID,&r.Name,&r.Severity,&r.Description,&r.Channel,&en,&r.CreatedAt);r.Enabled=en==1
d.db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE rule_id=?`,r.ID).Scan(&r.FireCount);d.db.QueryRow(`SELECT fired_at FROM alerts WHERE rule_id=? ORDER BY fired_at DESC LIMIT 1`,r.ID).Scan(&r.LastFired);o=append(o,r)};return o}
func(d *DB)DeleteRule(id string)error{_,err:=d.db.Exec(`DELETE FROM rules WHERE id=?`,id);return err}
func(d *DB)ToggleRule(id string)error{_,err:=d.db.Exec(`UPDATE rules SET enabled=1-enabled WHERE id=?`,id);return err}
func(d *DB)Fire(ruleID,message,source string)*Alert{r:=d.GetRule(ruleID);sev:="warning";name:=""
if r!=nil{sev=r.Severity;name=r.Name}
a:=&Alert{ID:genID(),RuleID:ruleID,RuleName:name,Severity:sev,Status:"firing",Message:message,Source:source,FiredAt:now()}
d.db.Exec(`INSERT INTO alerts(id,rule_id,rule_name,severity,status,message,source,fired_at)VALUES(?,?,?,?,?,?,?,?)`,a.ID,a.RuleID,a.RuleName,a.Severity,a.Status,a.Message,a.Source,a.FiredAt);return a}
func(d *DB)Ack(id,by string)error{_,err:=d.db.Exec(`UPDATE alerts SET status='acked',acked_at=?,acked_by=? WHERE id=?`,now(),by,id);return err}
func(d *DB)Resolve(id string)error{_,err:=d.db.Exec(`UPDATE alerts SET status='resolved',resolved_at=? WHERE id=?`,now(),id);return err}
func(d *DB)ListAlerts(status string,limit int)[]Alert{if limit<=0{limit=100};q:=`SELECT id,rule_id,rule_name,severity,status,message,source,fired_at,acked_at,resolved_at,acked_by FROM alerts`;args:=[]any{}
if status!=""&&status!="all"{q+=` WHERE status=?`;args=append(args,status)};q+=` ORDER BY fired_at DESC LIMIT ?`;args=append(args,limit)
rows,_:=d.db.Query(q,args...);if rows==nil{return nil};defer rows.Close()
var o []Alert;for rows.Next(){var a Alert;rows.Scan(&a.ID,&a.RuleID,&a.RuleName,&a.Severity,&a.Status,&a.Message,&a.Source,&a.FiredAt,&a.AckedAt,&a.ResolvedAt,&a.AckedBy);o=append(o,a)};return o}
type Stats struct{Rules int `json:"rules"`;Firing int `json:"firing"`;Acked int `json:"acked"`;Resolved int `json:"resolved"`;Total int `json:"total"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM rules`).Scan(&s.Rules);d.db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE status='firing'`).Scan(&s.Firing);d.db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE status='acked'`).Scan(&s.Acked);d.db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE status='resolved'`).Scan(&s.Resolved);d.db.QueryRow(`SELECT COUNT(*) FROM alerts`).Scan(&s.Total);return s}
var _=strings.Join
