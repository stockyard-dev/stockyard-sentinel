package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Sentinel</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#4a7ec9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.tabs{display:flex;border-bottom:1px solid var(--bg3);padding:0 1.5rem;font-size:.75rem}
.tab{padding:.7rem 1.2rem;cursor:pointer;color:var(--cm);border-bottom:2px solid transparent}.tab:hover{color:var(--cream)}.tab.active{color:var(--rust);border-color:var(--rust)}
.main{padding:1.5rem;max-width:900px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.6rem;margin-bottom:1.2rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.4rem}.st-l{font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.1rem}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.5rem;padding:.8rem 1rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.sev{font-size:.55rem;padding:.1rem .4rem;text-transform:uppercase;letter-spacing:1px}
.sev-critical{background:#c9444433;color:#ff6b6b;border:1px solid #c9444455}
.sev-warning{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.sev-info{background:#4a7ec922;color:var(--blue);border:1px solid #4a7ec944}
.status-firing{color:var(--red)}.status-acked{color:var(--orange)}.status-resolved{color:var(--green)}
.btn{font-size:.6rem;padding:.2rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.toggle{position:relative;width:36px;height:18px;cursor:pointer;display:inline-block;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}
.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.meta{font-size:.6rem;color:var(--cm);margin-top:.3rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic}
</style></head><body>
<div class="hdr"><h1>SENTINEL</h1><div id="hdr-stats" style="font-size:.7rem;color:var(--cm)"></div></div>
<div class="tabs"><div class="tab active" onclick="show('alerts')">Alerts</div><div class="tab" onclick="show('rules')">Rules</div></div>
<div class="main" id="main"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let tab='alerts',rules=[],alerts=[],stats={};
async function load(){
  const[r,a,s]=await Promise.all([fetch(A+'/rules').then(r=>r.json()),fetch(A+'/alerts?limit=100').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
  rules=r.rules||[];alerts=a.alerts||[];stats=s;
  document.getElementById('hdr-stats').innerHTML='<span style="color:var(--red)">'+s.firing+' firing</span> · '+s.acked+' acked · '+s.resolved+' resolved';
  render();
}
function show(t){tab=t;document.querySelectorAll('.tab').forEach((e,i)=>e.classList.toggle('active',['alerts','rules'][i]===t));render();}
function render(){const m=document.getElementById('main');if(tab==='rules')renderRules(m);else renderAlerts(m);}
function renderAlerts(m){
  let h='<div class="stats"><div class="st"><div class="st-v" style="color:var(--red)">'+stats.firing+'</div><div class="st-l">Firing</div></div><div class="st"><div class="st-v" style="color:var(--orange)">'+stats.acked+'</div><div class="st-l">Acked</div></div><div class="st"><div class="st-v" style="color:var(--green)">'+stats.resolved+'</div><div class="st-l">Resolved</div></div><div class="st"><div class="st-v">'+stats.total+'</div><div class="st-l">Total</div></div></div>';
  h+='<div style="display:flex;gap:.4rem;margin-bottom:1rem"><button class="btn" onclick="fireTest()">Test Fire</button></div>';
  const firing=alerts.filter(a=>a.status==='firing'),acked=alerts.filter(a=>a.status==='acked'),resolved=alerts.filter(a=>a.status==='resolved');
  if(firing.length){h+='<div style="font-size:.6rem;color:var(--red);margin-bottom:.3rem;text-transform:uppercase;letter-spacing:1px">FIRING ('+firing.length+')</div>';firing.forEach(a=>{h+=alertCard(a)});}
  if(acked.length){h+='<div style="font-size:.6rem;color:var(--orange);margin:1rem 0 .3rem;text-transform:uppercase;letter-spacing:1px">ACKNOWLEDGED ('+acked.length+')</div>';acked.forEach(a=>{h+=alertCard(a)});}
  if(resolved.length){h+='<div style="font-size:.6rem;color:var(--green);margin:1rem 0 .3rem;text-transform:uppercase;letter-spacing:1px">RESOLVED ('+resolved.length+')</div>';resolved.slice(0,20).forEach(a=>{h+=alertCard(a)});}
  if(!alerts.length)h+='<div class="empty">No alerts. All quiet.</div>';
  m.innerHTML=h;
}
function alertCard(a){
  let h='<div class="card"><div class="card-top"><div><span class="sev sev-'+a.severity+'">'+a.severity+'</span> <span style="font-size:.8rem;margin-left:.3rem">'+(a.rule_name||'Manual')+'</span></div><div style="display:flex;gap:.3rem">';
  if(a.status==='firing')h+='<button class="btn" onclick="ack(\''+a.id+'\')">Ack</button><button class="btn" onclick="resolve(\''+a.id+'\')">Resolve</button>';
  if(a.status==='acked')h+='<button class="btn" onclick="resolve(\''+a.id+'\')">Resolve</button>';
  h+='</div></div>';
  if(a.message)h+='<div style="font-size:.75rem;color:var(--cd);margin-top:.3rem">'+esc(a.message)+'</div>';
  h+='<div class="meta">Fired '+ft(a.fired_at);
  if(a.source)h+=' · Source: '+esc(a.source);
  if(a.acked_at)h+=' · Acked '+ft(a.acked_at);
  if(a.resolved_at)h+=' · Resolved '+ft(a.resolved_at);
  h+='</div></div>';return h;
}
function renderRules(m){
  let h='<div style="display:flex;justify-content:space-between;margin-bottom:1rem"><span style="font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px">Alert Rules ('+rules.length+')</span><button class="btn btn-p" onclick="openRuleForm()">+ New Rule</button></div>';
  if(!rules.length)h+='<div class="empty">No alert rules configured.</div>';
  rules.forEach(r=>{
    h+='<div class="card"><div class="card-top"><div><span class="sev sev-'+r.severity+'">'+r.severity+'</span> <span style="font-size:.8rem;margin-left:.3rem">'+esc(r.name)+'</span></div><div style="display:flex;gap:.4rem;align-items:center"><label class="toggle"><input type="checkbox" '+(r.enabled?'checked':'')+' onchange="toggleRule(\''+r.id+'\')"><span class="sl"></span></label><button class="btn" onclick="fireRule(\''+r.id+'\')">Fire</button><button class="btn" onclick="delRule(\''+r.id+'\')" style="color:var(--red)">✕</button></div></div>';
    if(r.description)h+='<div style="font-size:.72rem;color:var(--cm);margin-top:.2rem">'+esc(r.description)+'</div>';
    h+='<div class="meta">Fired '+r.fire_count+'x';if(r.last_fired)h+=' · Last: '+ft(r.last_fired);if(r.channel)h+=' · Channel: '+esc(r.channel);h+='</div></div>';
  });
  m.innerHTML=h;
}
async function ack(id){await fetch(A+'/alerts/'+id+'/ack',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({by:'admin'})});load();}
async function resolve(id){await fetch(A+'/alerts/'+id+'/resolve',{method:'POST'});load();}
async function toggleRule(id){await fetch(A+'/rules/'+id+'/toggle',{method:'POST'});load();}
async function delRule(id){if(confirm('Delete?')){await fetch(A+'/rules/'+id,{method:'DELETE'});load();}}
async function fireRule(id){await fetch(A+'/alerts/fire',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({rule_id:id,message:'Manual test fire',source:'dashboard'})});load();}
async function fireTest(){await fetch(A+'/alerts/fire',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({message:'Test alert from dashboard',source:'sentinel-ui',severity:'info'})});load();}
function openRuleForm(){
  document.getElementById('mdl').innerHTML='<h2>New Alert Rule</h2><div class="fr"><label>Name</label><input id="f-n" placeholder="e.g. High Error Rate"></div><div class="fr"><label>Severity</label><select id="f-s"><option value="info">Info</option><option value="warning" selected>Warning</option><option value="critical">Critical</option></select></div><div class="fr"><label>Description</label><input id="f-d" placeholder="When this fires..."></div><div class="fr"><label>Channel</label><input id="f-c" placeholder="e.g. slack, email, pagerduty"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="submitRule()">Create</button></div>';
  document.getElementById('mbg').classList.add('open');
}
async function submitRule(){await fetch(A+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-n').value,severity:document.getElementById('f-s').value,description:document.getElementById('f-d').value,channel:document.getElementById('f-c').value,enabled:true})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
function ft(t){if(!t)return'';const d=new Date(t);return d.toLocaleDateString()+' '+d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});}
load();
</script></body></html>`
