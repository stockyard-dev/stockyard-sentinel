package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Sentinel</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.stats{display:flex;gap:1.2rem;font-family:var(--mono);font-size:.7rem}
.stat-fire{color:var(--red)}.stat-ack{color:var(--orange)}.stat-ok{color:var(--green)}
.tabs{display:flex;border-bottom:1px solid var(--bg3);padding:0 1.5rem;font-family:var(--mono);font-size:.75rem}
.tab{padding:.7rem 1.2rem;cursor:pointer;color:var(--cm);border-bottom:2px solid transparent}.tab:hover{color:var(--cream)}.tab.active{color:var(--rust);border-color:var(--rust)}
.ct{padding:1.5rem;max-width:900px;margin:0 auto}
.card{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.5rem;padding:.8rem 1rem}
.card-top{display:flex;justify-content:space-between;align-items:center}
.badge{font-family:var(--mono);font-size:.55rem;padding:.12rem .4rem;text-transform:uppercase;letter-spacing:1px}
.b-critical{background:#c9444433;color:#ff6b6b;border:1px solid #c9444455}
.b-warning{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.b-info{background:#4a7ec922;color:#6ba3e8;border:1px solid #4a7ec944}
.b-firing{background:#c9444422;color:var(--red);border:1px solid #c9444444}
.b-acked{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.b-resolved{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}
.meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.2rem}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-p:hover{opacity:.85}
.toggle{position:relative;width:36px;height:18px;cursor:pointer;display:inline-block;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:2rem;color:var(--cm);font-style:italic}
.pulse{display:inline-block;width:6px;height:6px;border-radius:50%;background:var(--red);margin-right:.3rem;animation:p 1.5s infinite}@keyframes p{0%,100%{opacity:1}50%{opacity:.3}}
</style></head><body>
<div class="hdr"><h1>SENTINEL</h1><div class="stats" id="st"></div></div>
<div class="tabs"><div class="tab active" onclick="show('alerts')">Alerts</div><div class="tab" onclick="show('rules')">Rules</div></div>
<div class="ct" id="main"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let tab='alerts',rules=[],alerts=[];
async function ld(){const[r,a,s]=await Promise.all([fetch(A+'/rules').then(r=>r.json()),fetch(A+'/alerts').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);rules=r.rules||[];alerts=a.alerts||[];
document.getElementById('st').innerHTML=(s.firing?'<span class="stat-fire"><span class="pulse"></span>'+s.firing+' firing</span>':'')+'<span class="stat-ack">'+s.acked+' acked</span><span class="stat-ok">'+s.resolved+' resolved</span><span style="color:var(--cm)">'+s.rules+' rules</span>';rn();}
function show(t){tab=t;document.querySelectorAll('.tab').forEach((e,i)=>e.classList.toggle('active',['alerts','rules'][i]===t));rn();}
function rn(){const m=document.getElementById('main');if(tab==='rules'){rRules(m)}else{rAlerts(m)}}
function rAlerts(m){let h='<div style="display:flex;justify-content:space-between;margin-bottom:1rem"><div style="display:flex;gap:.3rem"><button class="btn'+(filt===''?' btn-p':'')+'" onclick="setF(\'\')">All</button><button class="btn'+(filt==='firing'?' btn-p':'')+'" onclick="setF(\'firing\')">Firing</button><button class="btn'+(filt==='acked'?' btn-p':'')+'" onclick="setF(\'acked\')">Acked</button><button class="btn'+(filt==='resolved'?' btn-p':'')+'" onclick="setF(\'resolved\')">Resolved</button></div><button class="btn btn-p" onclick="oFire()">Fire Alert</button></div>';
if(!alerts.length){h+='<div class="empty">No alerts. That\'s good.</div>';}
alerts.forEach(a=>{h+='<div class="card"><div class="card-top"><div><span class="badge b-'+a.severity+'">'+a.severity+'</span> <span style="font-family:var(--mono);font-size:.8rem">'+esc(a.rule_name||'Manual')+'</span></div><span class="badge b-'+a.status+'">'+a.status+'</span></div>';
if(a.message)h+='<div style="font-size:.78rem;color:var(--cd);margin-top:.3rem">'+esc(a.message)+'</div>';
h+='<div class="meta">'+esc(a.source||'')+' &middot; fired '+ft(a.fired_at);if(a.acked_at)h+=' &middot; acked '+ft(a.acked_at);if(a.resolved_at)h+=' &middot; resolved '+ft(a.resolved_at);h+='</div>';
h+='<div style="display:flex;gap:.3rem;margin-top:.4rem">';if(a.status==='firing')h+='<button class="btn" onclick="ack(\''+a.id+'\')">Acknowledge</button>';if(a.status!=='resolved')h+='<button class="btn" onclick="res(\''+a.id+'\')">Resolve</button>';h+='</div></div>';});m.innerHTML=h;}
function rRules(m){let h='<div style="display:flex;justify-content:space-between;margin-bottom:1rem"><span style="font-family:var(--mono);font-size:.75rem;color:var(--leather)">ALERT RULES</span><button class="btn btn-p" onclick="oRule()">+ New Rule</button></div>';
if(!rules.length){h+='<div class="empty">No rules configured. Create your first alert rule.</div>';}
rules.forEach(r=>{h+='<div class="card"><div class="card-top"><div><span class="badge b-'+r.severity+'">'+r.severity+'</span> <span style="font-family:var(--mono);font-size:.8rem">'+esc(r.name)+'</span> <label class="toggle"><input type="checkbox" '+(r.enabled?'checked':'')+' onchange="tog(\''+r.id+'\')"><span class="sl"></span></label></div><button class="btn" onclick="del(\''+r.id+'\')" style="font-size:.55rem;color:var(--red)">Delete</button></div>';
if(r.description)h+='<div style="font-size:.78rem;color:var(--cd);margin-top:.2rem">'+esc(r.description)+'</div>';
h+='<div class="meta">channel: '+esc(r.channel||'none')+' &middot; fired '+r.fire_count+'x';if(r.last_fired)h+=' &middot; last: '+ft(r.last_fired);h+='</div></div>';});m.innerHTML=h;}
let filt='';async function setF(f){filt=f;const r=await fetch(A+'/alerts?status='+f).then(r=>r.json());alerts=r.alerts||[];rn();}
async function ack(id){await fetch(A+'/alerts/'+id+'/ack',{method:'PATCH',headers:{'Content-Type':'application/json'},body:JSON.stringify({by:'admin'})});ld();}
async function res(id){await fetch(A+'/alerts/'+id+'/resolve',{method:'PATCH'});ld();}
async function tog(id){await fetch(A+'/rules/'+id+'/toggle',{method:'PATCH'});ld();}
async function del(id){if(confirm('Delete rule?')){await fetch(A+'/rules/'+id,{method:'DELETE'});ld();}}
function oRule(){document.getElementById('mdl').innerHTML='<h2>New Alert Rule</h2><div class="fr"><label>Name</label><input id="rn" placeholder="e.g. High Error Rate"></div><div class="fr"><label>Severity</label><select id="rs"><option value="info">Info</option><option value="warning" selected>Warning</option><option value="critical">Critical</option></select></div><div class="fr"><label>Channel</label><input id="rc" placeholder="e.g. slack, email, webhook"></div><div class="fr"><label>Description</label><input id="rd" placeholder="When to fire this rule"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sRule()">Create</button></div>';document.getElementById('mbg').classList.add('open');}
async function sRule(){await fetch(A+'/rules',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('rn').value,severity:document.getElementById('rs').value,channel:document.getElementById('rc').value,description:document.getElementById('rd').value})});cm();ld();}
function oFire(){let opts=rules.map(r=>'<option value="'+r.id+'">'+esc(r.name)+'</option>').join('');document.getElementById('mdl').innerHTML='<h2>Fire Alert</h2><div class="fr"><label>Rule</label><select id="fr"><option value="">Manual (no rule)</option>'+opts+'</select></div><div class="fr"><label>Message</label><input id="fm" placeholder="What happened"></div><div class="fr"><label>Source</label><input id="fs" placeholder="e.g. api-server, monitor"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sFire()">Fire</button></div>';document.getElementById('mbg').classList.add('open');}
async function sFire(){await fetch(A+'/alerts/fire',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({rule_id:document.getElementById('fr').value,message:document.getElementById('fm').value,source:document.getElementById('fs').value})});cm();ld();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
function ft(t){if(!t)return'';const d=new Date(t);return d.toLocaleDateString()+' '+d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});}
ld();
</script></body></html>`
