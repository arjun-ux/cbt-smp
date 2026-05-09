import{_ as $,u as C,a as N,x as T,c as o,b as t,g as M,t as d,f as G,F as y,y as k,r as h,o as r,n as U,e as z,D as J,l as V}from"./index-DMmg1YdV.js";const I={class:"space-y-6 pb-12"},L={class:"bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden"},E={class:"p-4 border-b border-slate-50 bg-slate-50/30 flex justify-between items-center"},R={class:"text-xs font-bold text-blue-500 bg-blue-50 px-3 py-1 rounded-lg"},F={class:"overflow-x-auto"},H={class:"min-w-full divide-y divide-slate-100"},O={class:"bg-white divide-y divide-slate-50"},W={key:0},K={key:1},j={class:"px-6 py-5 text-sm font-bold text-slate-400"},q={class:"px-6 py-5"},Q={class:"flex flex-col"},X={class:"font-bold text-slate-800"},Y={class:"text-[10px] font-black text-blue-500 uppercase tracking-tighter"},Z={class:"px-6 py-5"},tt={class:"flex flex-col"},et={class:"text-sm font-bold text-slate-700 leading-tight"},at={class:"flex items-center gap-2 mt-1"},st={class:"px-6 py-5"},lt={class:"px-2.5 py-1 bg-slate-100 text-slate-600 rounded-md text-[10px] font-bold border border-slate-200"},nt={class:"px-6 py-5"},ot={class:"relative"},rt=["onUpdate:modelValue","onChange","disabled"],dt=["value"],it={class:"absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none"},pt={key:0,class:"w-4 h-4 text-slate-400",fill:"none",stroke:"currentColor",viewBox:"0 0 24 24"},ut={key:1,class:"w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"},ct={__name:"PlottingPengawas",setup(gt){const c=C(),x=N(),g=h([]),_=h([]),m=h(!1),b=h(null),P=V(()=>{var s;return((s=c.user)==null?void 0:s.role)==="admin"?"/api/admin":"/api/guru"}),A=async()=>{var a;m.value=!0;const s={Authorization:`Bearer ${c.token}`},e=((a=c.user)==null?void 0:a.role)==="guru";try{const[l,i]=await Promise.all([fetch(e?"/api/guru/pengawas-jadwal":"/api/admin/jadwal",{headers:s}),fetch(e?"/api/guru/guru":"/api/admin/guru",{headers:s})]);if(l.ok&&i.ok){const u=await l.json(),p=await i.json();g.value=u.data||[],_.value=p.data||[]}}catch{x.showAlert("Gagal memuat data","error")}finally{m.value=!1}},D=async(s,e)=>{var a;if(((a=c.user)==null?void 0:a.role)==="admin"){b.value=s;try{(await fetch(`${P.value}/jadwal/${s}/pengawas`,{method:"POST",headers:{Authorization:`Bearer ${c.token}`,"Content-Type":"application/json"},body:JSON.stringify({pengawas_id:e?parseInt(e):null})})).ok?x.showAlert("Pengawas berhasil diplot","success"):x.showAlert("Gagal mengupdate pengawas","error")}catch{x.showAlert("Kesalahan sistem","error")}finally{b.value=null}}},f=s=>s?new Date(s).toLocaleDateString("id-ID",{day:"2-digit",month:"2-digit",year:"numeric"}):"-",S=()=>{const s=window.open("","_blank");let e="";g.value.forEach((l,i)=>{var p,n,w,v;const u=((p=_.value.find(B=>B.id===l.pengawas_id))==null?void 0:p.nama_guru)||"-";e+=`
      <tr>
        <td style="border: 1px solid #000; padding: 8px; text-align: center;">${i+1}</td>
        <td style="border: 1px solid #000; padding: 8px;">${f(l.tanggal_ujian)}</td>
        <td style="border: 1px solid #000; padding: 8px; text-align: center;">${((n=l.ruang)==null?void 0:n.nama_ruang)||"-"}</td>
        <td style="border: 1px solid #000; padding: 8px; text-align: center;">${((w=l.sesi)==null?void 0:w.nama_sesi)||l.sesi_id}</td>
        <td style="border: 1px solid #000; padding: 8px;">${(v=l.bank_soal)==null?void 0:v.judul_bank_soal}</td>
        <td style="border: 1px solid #000; padding: 8px; font-weight: bold;">${u}</td>
      </tr>
    `});const a=`
    <html>
      <head>
        <title>Daftar Pengawas Ujian</title>
        <style>
          body { font-family: sans-serif; padding: 20px; }
          table { width: 100%; border-collapse: collapse; margin-top: 20px; }
          th { border: 1px solid #000; padding: 10px; background: #eee; }
          h2 { text-align: center; margin-bottom: 5px; }
          p { text-align: center; margin-top: 0; font-size: 14px; }
        </style>
      </head>
      <body>
        <h2>DAFTAR PLOTTING PENGAWAS UJIAN</h2>
        <p>CBT SMP APP - Tanggal Cetak: ${new Date().toLocaleString("id-ID")}</p>
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Tanggal</th>
              <th>Ruang</th>
              <th>Sesi</th>
              <th>Mata Pelajaran</th>
              <th>Nama Pengawas</th>
            </tr>
          </thead>
          <tbody>${e}</tbody>
        </table>
        <div style="margin-top: 40px; text-align: right; padding-right: 50px;">
          <p>Panitia Ujian,</p>
          <br><br><br>
          <p>( ____________________ )</p>
        </div>
        <script>window.print(); window.close();<\/script>
      </body>
    </html>
  `;s.document.write(a),s.document.close()};return T(A),(s,e)=>(r(),o("div",I,[t("div",{class:"flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4"},[e[1]||(e[1]=t("div",{class:"flex items-center gap-4"},[t("div",null,[t("h3",{class:"text-xl font-bold text-slate-800"},"Plotting Pengawas"),t("p",{class:"text-sm text-slate-500 mt-0.5"},"Tugaskan guru untuk mengawas ujian di setiap ruang dan sesi")])],-1)),t("button",{onClick:S,class:"flex items-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-emerald-100"},[...e[0]||(e[0]=[t("svg",{class:"w-4 h-4",fill:"none",stroke:"currentColor",viewBox:"0 0 24 24"},[t("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:"M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"})],-1),M(" Cetak Daftar ",-1)])])]),t("div",L,[t("div",E,[e[2]||(e[2]=t("div",{class:"text-xs font-bold text-slate-400 uppercase tracking-widest"},"Daftar Jadwal Aktif",-1)),t("div",R,d(g.value.length)+" Jadwal",1)]),t("div",F,[t("table",H,[e[7]||(e[7]=t("thead",{class:"bg-slate-50"},[t("tr",null,[t("th",{class:"px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider w-16"},"No"),t("th",{class:"px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider"},"Jadwal & Ruang"),t("th",{class:"px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider"},"Mata Pelajaran"),t("th",{class:"px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider"},"Sesi"),t("th",{class:"px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider w-72"},"Pengawas")])],-1)),t("tbody",O,[m.value?(r(),o("tr",W,[...e[3]||(e[3]=[t("td",{colspan:"5",class:"px-6 py-12 text-center text-slate-400"},[t("div",{class:"flex flex-col items-center gap-2 animate-pulse"},[t("div",{class:"w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"}),t("span",{class:"text-sm font-medium"},"Memuat data...")])],-1)])])):g.value.length===0?(r(),o("tr",K,[...e[4]||(e[4]=[t("td",{colspan:"5",class:"px-6 py-12 text-center text-slate-400 italic font-medium"},"Belum ada jadwal ujian yang dibuat.",-1)])])):G("",!0),(r(!0),o(y,null,k(g.value,(a,l)=>{var i,u,p;return r(),o("tr",{key:a.id,class:"hover:bg-slate-50/50 transition-colors group"},[t("td",j,d(l+1),1),t("td",q,[t("div",Q,[t("span",X,d(f(a.tanggal_ujian)),1),t("span",Y,"RUANG: "+d(((i=a.ruang)==null?void 0:i.nama_ruang)||"-"),1)])]),t("td",Z,[t("div",tt,[t("span",et,d((u=a.bank_soal)==null?void 0:u.judul_bank_soal),1),t("div",at,[t("span",{class:U([{"bg-amber-100 text-amber-700 border-amber-200":a.status==="Belum Mulai","bg-blue-100 text-blue-700 border-blue-200":a.status==="Berlangsung","bg-emerald-100 text-emerald-700 border-emerald-200":a.status==="Selesai"||a.status==="Diarsipkan"},"px-2 py-0.5 border rounded text-[8px] font-black uppercase"])},d(a.status),3)])])]),t("td",st,[t("span",lt,d(((p=a.sesi)==null?void 0:p.nama_sesi)||"SESI "+a.sesi_id),1)]),t("td",nt,[t("div",ot,[z(t("select",{"onUpdate:modelValue":n=>a.pengawas_id=n,onChange:n=>D(a.id,n.target.value),disabled:b.value===a.id||a.status==="Berlangsung"||a.status==="Diarsipkan",class:"w-full pl-4 pr-10 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold focus:ring-2 focus:ring-blue-500 transition-all outline-none appearance-none disabled:opacity-50"},[e[5]||(e[5]=t("option",{value:null},"-- Pilih Pengawas --",-1)),(r(!0),o(y,null,k(_.value,n=>(r(),o("option",{key:n.id,value:n.id},d(n.nama_guru),9,dt))),128))],40,rt),[[J,a.pengawas_id]]),t("div",it,[b.value!==a.id?(r(),o("svg",pt,[...e[6]||(e[6]=[t("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:"M19 9l-7 7-7-7"},null,-1)])])):(r(),o("div",ut))])])])])}),128))])])])])]))}},bt=$(ct,[["__scopeId","data-v-1a834aaf"]]);export{bt as default};
