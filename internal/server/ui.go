package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Sentinel</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#4a7ec9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.header{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.header h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.tabs{display:flex;border-bottom:1px solid var(--bg3);padding:0 1.5rem;font-family:var(--mono);font-size:.75rem}
.tab{padding:.7rem 1.2rem;cursor:pointer;color:var(--cm);border-bottom:2px solid transparent}.tab:hover{color:var(--cream)}.tab.active{color:var(--rust);border-color:var(--rust)}
.content{padding:1.5rem;max-width:900px;margin:0 auto}
.stats-row{display:grid;grid-template-columns:repeat(auto-fit,minmax(120px,1fr));gap:.8rem;margin-bottom:1.5rem}
.stat{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem;text-align:center}
.stat-val{font-family:var(--mono);font-size:1.4rem}.stat-label{font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.5rem;padding:.8rem 1rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.card-name{font-family:var(--mono);font-size:.8rem}
.card-meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.2rem}
.badge{font-family:var(--mono);font-size:.55rem;padding:.1rem .4rem;text-transform:uppercase;letter-spacing:1px}
.badge-critical{background:#c9444433;color:#ff6b6b;border:1px solid #c9444455}
.badge-warning{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.badge-info{background:#4a7ec922;color:var(--blue);border:1px solid #4a7ec944}
.badge-firing{background:#c9444433;color:var(--red);border:1px solid #c9444444}
.badge-acked{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.badge-resolved{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-primary:hover{opacity:.85}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.toggle{position:relative;display:inline-block;width:36px;height:18px;cursor:pointer;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}
.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.form-row{margin-bottom:.7rem}
.form-row label{display:block;font-family:var(--mono);font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.form-row input,.form-row select,.form-row textarea{width:100%;padding:.45rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem}
.actions{display:flex;gap:.5rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:2rem;color:var(--cm);font-style:italic}
.fire-pulse{animation:fpulse 1.5s infinite}
@keyframes fpulse{0%,100%{opacity:1}50%{opacity:.5}}
</style></head><body>
<div class="header"><h1>SENTINEL</h1><div id="firing" style="font-family:var(--mono);font-size:.75rem"></div></div>
<div class="tabs"><div class="tab active" onclick="showTab('alerts')">Alerts</div><div class="tab" onclick="showTab('rules')">Rules</div></div>
<div class="content" id="main"></div>
<div class="modal-bg" id="modalBg" onclick="if(event.target===this)closeModal()"><div class="modal" id="modal"></div></div>

<script>
const API='/api';let tab='alerts',rules=[],alerts=[],stats={};

async function load(){
  const[r,a,s]=await Promise.all([fetch(API+'/rules').then(r=>r.json()),fetch(API+'/alerts').then(r=>r.json()),fetch(API+'/stats').then(r=>r.json())]);
  rules=r.rules||[];alerts=a.alerts||[];stats=s;
  const f=document.getElementById('firing');
  f.innerHTML=stats.firing>0?'<span class="fire-pulse" style="color:var(--red)">&#9679;</span> '+stats.firing+' firing':'<span style="color:var(--green)">&#9679;</span> All clear';
  render();
}

function showTab(t){tab=t;document.querySelectorAll('.tab').forEach((el,i)=>el.classList.toggle('active',['alerts','rules'][i]===t));render();}

function render(){
  const m=document.getElementById('main');
  if(tab==='rules')renderRules(m);else renderAlerts(m);
}

function renderAlerts(m){
  let h='<div class="stats-row"><div class="stat"><div class="stat-val" style="color:var(--red)">'+stats.firing+'</div><div class="stat-label">Firing</div></div><div class="stat"><div class="stat-val" style="color:var(--orange)">'+stats.acked+'</div><div class="stat-label">Acknowledged</div></div><div class="stat"><div class="stat-val" style="color:var(--green)">'+stats.resolved+'</div><div class="stat-label">Resolved</div></div><div class="stat"><div class="stat-val">'+stats.total+'</div><div class="stat-label">Total</div></div></div>';
  h+='<div style="display:flex;gap:.3rem;margin-bottom:1rem;font-family:var(--mono);font-size:.65rem">';
  ['all','firing','acked','resolved'].forEach(s=>{h+='<button class="btn'+(s==='all'?' btn-primary':'')+'" onclick="filterAlerts(\''+s+'\')">'+s+'</button>';});
  h+='</div>';
  if(!alerts||!alerts.length){h+='<div class="empty">No alerts. Create rules and fire them via POST /api/alerts/fire.</div>';}
  else{(alerts||[]).forEach(a=>{
    h+='<div class="card"><div class="card-top"><div><span class="badge badge-'+a.severity+'">'+a.severity+'</span> <span class="card-name">'+esc(a.rule_name||a.rule_id||'manual')+'</span></div><span class="badge badge-'+a.status+'">'+a.status+'</span></div>';
    if(a.message)h+='<div style="font-size:.78rem;color:var(--cd);margin-top:.3rem">'+esc(a.message)+'</div>';
    h+='<div class="card-meta">'+fmtTime(a.fired_at);
    if(a.source)h+=' &middot; '+esc(a.source);
    if(a.acked_by)h+=' &middot; acked by '+esc(a.acked_by);
    h+='</div><div style="display:flex;gap:.3rem;margin-top:.4rem">';
    if(a.status==='firing')h+='<button class="btn btn-sm" onclick="ack(\''+a.id+'\')">Acknowledge</button>';
    if(a.status!=='resolved')h+='<button class="btn btn-sm" onclick="resolve(\''+a.id+'\')">Resolve</button>';
    h+='</div></div>';
  });}
  m.innerHTML=h;
}

function renderRules(m){
  let h='<div style="display:flex;justify-content:space-between;margin-bottom:1rem"><h2 style="font-family:var(--mono);font-size:.75rem;color:var(--leather)">ALERT RULES</h2><button class="btn btn-primary" onclick="openRuleForm()">+ New Rule</button></div>';
  if(!rules||!rules.length){h+='<div class="empty">No rules yet. Define alert rules to route and manage incidents.</div>';}
  else{(rules||[]).forEach(r=>{
    h+='<div class="card"><div class="card-top"><div><span class="badge badge-'+r.severity+'">'+r.severity+'</span> <span class="card-name">'+esc(r.name)+'</span> <label class="toggle"><input type="checkbox" '+(r.enabled?'checked':'')+' onchange="toggleRule(\''+r.id+'\')"><span class="sl"></span></label></div><button class="btn btn-sm" onclick="fireRule(\''+r.id+'\')">Fire Test</button></div>';
    if(r.description)h+='<div style="font-size:.75rem;color:var(--cm);margin-top:.2rem">'+esc(r.description)+'</div>';
    h+='<div class="card-meta">channel: '+(r.channel||'default')+' &middot; fired '+r.fire_count+'× &middot; last: '+(r.last_fired?fmtTime(r.last_fired):'never')+'</div>';
    h+='<div style="margin-top:.3rem"><button class="btn btn-sm" onclick="delRule(\''+r.id+'\')" style="color:var(--red)">Delete</button></div></div>';
  });}
  m.innerHTML=h;
}

async function filterAlerts(s){const r=await fetch(API+'/alerts?status='+s).then(r=>r.json());alerts=r.alerts||[];render();}
async function ack(id){await fetch(API+'/alerts/'+id+'/ack',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({by:'admin'})});load();}
async function resolve(id){await fetch(API+'/alerts/'+id+'/resolve',{method:'POST'});load();}
async function toggleRule(id){await fetch(API+'/rules/'+id+'/toggle',{method:'POST'});load();}
async function delRule(id){if(confirm('Delete?')){await fetch(API+'/rules/'+id,{method:'DELETE'});load();}}
async function fireRule(id){await fetch(API+'/alerts/fire',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({rule_id:id,message:'Test alert fired from dashboard',source:'dashboard'})});load();}

function openRuleForm(){
  document.getElementById('modal').innerHTML='<h2>New Alert Rule</h2><div class="form-row"><label>Name</label><input id="f-name" placeholder="e.g. High Error Rate"></div><div class="form-row"><label>Severity</label><select id="f-sev"><option value="info">Info</option><option value="warning" selected>Warning</option><option value="critical">Critical</option></select></div><div class="form-row"><label>Description</label><input id="f-desc" placeholder="What triggers this alert"></div><div class="form-row"><label>Channel</label><input id="f-chan" placeholder="e.g. slack, email, pagerduty"></div><div class="actions"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-primary" onclick="submitRule()">Create</button></div>';
  document.getElementById('modalBg').classList.add('open');
}
async function submitRule(){await fetch(API+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-name').value,severity:document.getElementById('f-sev').value,description:document.getElementById('f-desc').value,channel:document.getElementById('f-chan').value,enabled:true})});closeModal();load();}

function closeModal(){document.getElementById('modalBg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
function fmtTime(t){if(!t)return'';const d=new Date(t);return d.toLocaleDateString()+' '+d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});}
load();setInterval(load,30000);
</script></body></html>`
