const axios = require('axios');

const API_URL = 'http://localhost:3000/api';
const NUM_STUDENTS = 100;      
const SYNC_CYCLES = 10;        
const TEST_INTERVAL = 5000;    
const EXAM_TOKEN = 'TEST100';  
const JADWAL_ID = 1;           

async function runStressTest() {
    console.log(`\n🚀 === BATTLEFIELD EXTREME (Images & Essays) START ===`);
    console.log(`👥 Target       : ${NUM_STUDENTS} Siswa`);
    console.log(`🔄 Siklus Sync  : ${SYNC_CYCLES} kali`);
    console.log(`📅 Jadwal ID    : ${JADWAL_ID}`);
    
    const students = [];

    // 1. LOGIN & PREPARATION MASSAL
    console.log(`\n[1/4] 🔐 Tahap Login & Analisis Soal (Gambar/Essay)...`);
    const loginStartTime = Date.now();
    
    for (let i = 1; i <= NUM_STUDENTS; i++) {
        const username = `siswa${i}`;
        try {
            const res = await axios.post(`${API_URL}/auth/login`, {
                username: username,
                password: 'siswa123'
            });
            
            // Format API Anda dibungkus dalam properti "data" oleh helper SendSuccess
            const apiResponse = res.data.data;
            if (!apiResponse || !apiResponse.token) {
                throw new Error(`Respon login tidak mengandung token: ${JSON.stringify(res.data)}`);
            }

            const token = apiResponse.token;
            const user = apiResponse.user;

            const valRes = await axios.post(`${API_URL}/siswa/validate?jadwalId=${JADWAL_ID}`, {
                token: EXAM_TOKEN,
                nisn: user.username
            }, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            const pesertaId = valRes.data.data.peserta_id;

            // Ambil Daftar Soal & Deteksi Gambar
            const soalRes = await axios.get(`${API_URL}/siswa/soal/${JADWAL_ID}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            
            const rawSoals = soalRes.data.data.items;
            const soalDetails = rawSoals.map(s => {
                // Cari URL Gambar dalam pertanyaan atau opsi
                const imgRegex = /src="([^"]+)"/g;
                const images = [];
                let match;
                const fullText = s.pertanyaan + s.opsi_a + s.opsi_b + s.opsi_c + s.opsi_d;
                while ((match = imgRegex.exec(fullText)) !== null) {
                    images.push(match[1]);
                }
                return { id: s.id, type: s.jenis_soal, images: images };
            });

            students.push({
                id: i,
                token: token,
                peserta_id: pesertaId,
                soals: soalDetails
            });

            if (i % 20 === 0) console.log(`   > ${i} siswa siap...`);
        } catch (err) {
            console.error(`   X Siswa ${i} gagal: ${err.message}`);
        }
    }

    // 2. SIMULASI PENGERJAAN & DOWNLOAD GAMBAR
    console.log(`\n[2/4] ⚡ Tahap Simulasi (Sync Jawaban & Download Gambar)...`);
    
    for (let cycle = 1; cycle <= SYNC_CYCLES; cycle++) {
        const cycleStartTime = Date.now();
        let syncSuccess = 0;
        let imgSuccess = 0;
        let failed = 0;

        const requests = students.map(async (s) => {
            try {
                // A. Simulasi Download Gambar (Hanya di siklus pertama/acak)
                if (cycle === 1 || Math.random() > 0.7) {
                    const allImages = s.soals.flatMap(soal => soal.images);
                    if (allImages.length > 0) {
                        const targetImg = allImages[Math.floor(Math.random() * allImages.length)];
                        const baseUrl = API_URL.replace('/api', '');
                        const imgUrl = targetImg.startsWith('http') ? targetImg : `${baseUrl}${targetImg}`;
                        
                        await axios.get(imgUrl, { responseType: 'arraybuffer' });
                        imgSuccess++;
                    }
                }

                // B. Simulasi Sinkronisasi Jawaban (PG & Essay)
                const items = [];
                for(let k=0; k<5; k++) {
                    const targetSoal = s.soals[Math.floor(Math.random() * s.soals.length)];
                    let jawaban = "";
                    
                    if (targetSoal.type === 'PG') {
                        jawaban = ['A', 'B', 'C', 'D'][Math.floor(Math.random() * 4)];
                    } else {
                        // Simulasi Jawaban Essay (Teks Panjang)
                        jawaban = "Ini adalah jawaban essay simulasi yang cukup panjang untuk mengetes performa penyimpanan teks di database. " + 
                                  "Siswa sedang menjelaskan jawaban dengan detail agar mendapatkan nilai maksimal dari guru pengampu.";
                    }

                    items.push({
                        soal_id: targetSoal.id,
                        jawaban_teks: jawaban,
                        ragu_ragu: Math.random() > 0.9
                    });
                }

                await axios.post(`${API_URL}/siswa/sync`, {
                    peserta_ujian_id: s.peserta_id,
                    items: items,
                    sisa_waktu: 3600 - (cycle * 30)
                }, {
                    headers: { 'Authorization': `Bearer ${s.token}` }
                });
                
                syncSuccess++;
            } catch (err) {
                failed++;
            }
        });

        await Promise.all(requests);
        const duration = Date.now() - cycleStartTime;
        console.log(`   🔄 Siklus #${cycle}: ${syncSuccess} Sync, ${imgSuccess} Gambar OK | Durasi: ${duration}ms`);
        
        if (cycle < SYNC_CYCLES) await new Promise(r => setTimeout(r, TEST_INTERVAL));
    }

    // 3. FINAL SUBMIT
    console.log(`\n[3/4] 🏁 Tahap Submit Ujian Serentak...`);
    const submitStartTime = Date.now();
    let subSuccess = 0;
    
    await Promise.all(students.map(async (s) => {
        try {
            await axios.post(`${API_URL}/siswa/submit/${s.peserta_id}`, {}, {
                headers: { 'Authorization': `Bearer ${s.token}` }
            });
            subSuccess++;
        } catch {}
    }));

    const submitDuration = Date.now() - submitStartTime;
    console.log(`   └─ Hasil Submit: ${subSuccess} Berhasil | Total Waktu: ${submitDuration}ms`);

    // 4. SUMMARY
    console.log(`\n📊 === KESIMPULAN PENGUJIAN ===`);
    console.log(`- Beban Statis (Gambar) : Berhasil disimulasikan`);
    console.log(`- Beban Data (Essay)    : Berhasil disimulasikan`);
    console.log(`- Avg Submit Time       : ${(submitDuration/students.length).toFixed(2)}ms`);
    console.log(`================================\n`);
}

runStressTest();
