package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Sentinel</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#4a7ec9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.stats-bar{font-family:var(--mono);font-size:.7rem;display:flex;gap:1rem}
.tabs{display:flex;border-bottom:1px solid var(--bg3);padding:0 1.5rem;font-family:var(--mono);font-size:.75rem}
.tab{padding:.7rem 1.2rem;cursor:pointer;color:var(--cm);border-bottom:2px solid transparent}.tab:hover{color:var(--cream)}.tab.active{color:var(--rust);border-color:var(--rust)}
.wrap{padding:1.5rem;max-width:900px;margin:0 auto}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.5rem;padding:.8rem 1rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.card-name{font-family:var(--mono);font-size:.82rem}
.badge{font-family:var(--mono);font-size:.55rem;padding:.12rem .4rem;text-transform:uppercase;letter-spacing:1px}
.b-critical{background:#c9444433;color:#ff6b6b;border:1px solid #c9444455}
.b-warning{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.b-info{background:#4a7ec922;color:var(--blue);border:1px solid #4a7ec944}
.b-firing{background:#c9444433;color:var(--red)}.b-acked{background:#d4843a22;color:var(--orange)}.b-resolved{background:#4a9e5c22;color:var(--green)}
.meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.3rem}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.toggle{position:relative;width:36px;height:18px;cursor:pointer;display:inline-block;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}
.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem}
.actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic}
.fire-dot{display:inline-block;width:6px;height:6px;border-radius:50%;background:var(--red);margin-right:.3rem;animation:pulse 1.5s infinite}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.3}}
</style></head><body>
<div class="hdr"><h1>SENTINEL</h1><div class="stats-bar" id="stats"></div></div>
<div class="tabs">
<div class="tab active" onclick="showTab('alerts')">Alerts</div>
<div class="tab" onclick="showTab('rules')">Rules</div>
</div>
<div class="wrap" id="main"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let tab='alerts',rules=[],alerts=[];
async function load(){const[r,a,s]=await Promise.all([fetch(A+'/rules').then(r=>r.json()),fetch(A+'/alerts?status=all&limit=100').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
rules=r.rules||[];alerts=a.alerts||[];
document.getElementById('stats').innerHTML=(s.firing>0?'<span><span class="fire-dot"></span>'+s.firing+' firing</span>':'<span style="color:var(--green)">0 firing</span>')+'<span>'+s.acked+' acked</span><span>'+s.resolved+' resolved</span>';
render();}
function showTab(t){tab=t;document.querySelectorAll('.tab').forEach((e,i)=>e.classList.toggle('active',['alerts','rules'][i]===t));render();}
function render(){const m=document.getElementById('main');if(tab==='rules')renderRules(m);else renderAlerts(m);}
function renderAlerts(m){let h='<div style="display:flex;justify-content:space-between;margin-bottom:1rem"><h2 style="font-family:var(--mono);font-size:.75rem;color:var(--leather)">ALERTS</h2><div style="display:flex;gap:.3rem"><button class="btn" onclick="filterAlerts(\'all\')">All</button><button class="btn" onclick="filterAlerts(\'firing\')">Firing</button><button class="btn" onclick="filterAlerts(\'acked\')">Acked</button></div></div>';
if(!alerts||!alerts.length){h+='<div class="empty">No alerts. All quiet.</div>';}
else{const firing=alerts.filter(a=>a.status==='firing');const acked=alerts.filter(a=>a.status==='acked');const resolved=alerts.filter(a=>a.status==='resolved');
if(firing.length){h+='<div style="font-family:var(--mono);font-size:.6rem;color:var(--red);margin-bottom:.4rem;text-transform:uppercase;letter-spacing:1px">FIRING ('+firing.length+')</div>';firing.forEach(a=>{h+=alertCard(a)});}
if(acked.length){h+='<div style="font-family:var(--mono);font-size:.6rem;color:var(--orange);margin:1rem 0 .4rem;text-transform:uppercase;letter-spacing:1px">ACKNOWLEDGED ('+acked.length+')</div>';acked.forEach(a=>{h+=alertCard(a)});}
if(resolved.length){h+='<div style="font-family:var(--mono);font-size:.6rem;color:var(--green);margin:1rem 0 .4rem;text-transform:uppercase;letter-spacing:1px">RESOLVED ('+resolved.length+')</div>';resolved.slice(0,20).forEach(a=>{h+=alertCard(a)});}}
m.innerHTML=h;}
function alertCard(a){let h='<div class="card"><div class="card-top"><div><span class="badge b-'+a.severity+'">'+a.severity+'</span> <span class="card-name">'+esc(a.rule_name||'Manual')+'</span></div><span class="badge b-'+a.status+'">'+(a.status==='firing'?'<span class="fire-dot"></span>':'')+a.status+'</span></div>';
if(a.message)h+='<div style="font-size:.78rem;color:var(--cd);margin-top:.3rem">'+esc(a.message)+'</div>';
h+='<div class="meta">Fired: '+ft(a.fired_at);if(a.source)h+=' · Source: '+esc(a.source);if(a.acked_by)h+=' · Acked by: '+esc(a.acked_by);h+='</div>';
h+='<div style="display:flex;gap:.3rem;margin-top:.4rem">';
if(a.status==='firing')h+='<button class="btn" onclick="ack(\''+a.id+'\')">Acknowledge</button>';
if(a.status!=='resolved')h+='<button class="btn" onclick="resolve(\''+a.id+'\')">Resolve</button>';
h+='</div></div>';return h;}
function renderRules(m){let h='<div style="display:flex;justify-content:space-between;margin-bottom:1rem"><h2 style="font-family:var(--mono);font-size:.75rem;color:var(--leather)">ALERT RULES</h2><button class="btn btn-primary" onclick="openRuleForm()">+ New Rule</button></div>';
if(!rules||!rules.length){h+='<div class="empty">No alert rules. Create one to start routing alerts.</div>';}
else{rules.forEach(r=>{
h+='<div class="card"><div class="card-top"><div><label class="toggle"><input type="checkbox" '+(r.enabled?'checked':'')+' onchange="toggleRule(\''+r.id+'\')"><span class="sl"></span></label> <span class="card-name">'+esc(r.name)+'</span> <span class="badge b-'+r.severity+'">'+r.severity+'</span></div><div style="display:flex;gap:.3rem"><button class="btn" style="font-size:.55rem" onclick="fire(\''+r.id+'\')">Fire Test</button><button class="btn" style="font-size:.55rem;color:var(--red)" onclick="delRule(\''+r.id+'\')">Delete</button></div></div>';
if(r.description)h+='<div style="font-size:.75rem;color:var(--cm);margin-top:.2rem">'+esc(r.description)+'</div>';
h+='<div class="meta">';if(r.channel)h+='Channel: '+esc(r.channel)+' · ';h+='Fired '+r.fire_count+'× ';if(r.last_fired)h+=' · Last: '+ft(r.last_fired);h+='</div></div>';});}
m.innerHTML=h;}
async function ack(id){await fetch(A+'/alerts/'+id+'/ack',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({by:'admin'})});load();}
async function resolve(id){await fetch(A+'/alerts/'+id+'/resolve',{method:'POST'});load();}
async function toggleRule(id){await fetch(A+'/rules/'+id+'/toggle',{method:'POST'});load();}
async function fire(id){await fetch(A+'/alerts/fire',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({rule_id:id,message:'Test alert fired from dashboard',source:'dashboard'})});load();}
async function delRule(id){if(confirm('Delete?')){await fetch(A+'/rules/'+id,{method:'DELETE'});load();}}
async function filterAlerts(s){const r=await fetch(A+'/alerts?status='+s+'&limit=100').then(r=>r.json());alerts=r.alerts||[];render();}
function openRuleForm(){document.getElementById('mdl').innerHTML='<h2>New Alert Rule</h2><div class="fr"><label>Name</label><input id="f-name" placeholder="e.g. High CPU Usage"></div><div class="fr"><label>Severity</label><select id="f-sev"><option value="info">Info</option><option value="warning" selected>Warning</option><option value="critical">Critical</option></select></div><div class="fr"><label>Description</label><input id="f-desc" placeholder="Alert when CPU exceeds 90%"></div><div class="fr"><label>Channel</label><input id="f-ch" placeholder="e.g. #ops-alerts, email"></div><div class="actions"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-primary" onclick="submitRule()">Create</button></div>';document.getElementById('mbg').classList.add('open');}
async function submitRule(){await fetch(A+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-name').value,severity:document.getElementById('f-sev').value,description:document.getElementById('f-desc').value,channel:document.getElementById('f-ch').value,enabled:true})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
function ft(t){if(!t)return'';const d=new Date(t);return d.toLocaleDateString()+' '+d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});}
load();
</script></body></html>`
