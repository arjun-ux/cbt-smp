const axios = require('axios');

const API_URL = 'http://localhost:3000/api';
const NUM_STUDENTS = 100;
const TEST_INTERVAL = 15000; // 15 detik
const EXAM_TOKEN = 'TEST100'; // Sesuaikan dengan token jadwal yang Bapak buat
const JADWAL_ID = 1;        // GANTI INI sesuai ID Jadwal yang Bapak buat

async function runStressTest() {
    console.log(`=== CBT STRESS TEST SIMULATOR ===`);
    console.log(`Target: ${NUM_STUDENTS} Siswa serentak`);
    console.log(`Interval: ${TEST_INTERVAL/1000} detik`);
    console.log(`Jadwal ID: ${JADWAL_ID}`);
    
    const students = [];

    // 1. LOGIN MASSAL
    console.log(`\n[1/3] Memulai Login Massal...`);
    for (let i = 1; i <= NUM_STUDENTS; i++) {
        const username = `siswa${i}`;
        try {
            const res = await axios.post(`${API_URL}/auth/login`, {
                username: username,
                password: 'siswa123'
            });
            
            // Login langsung kirim token di root (tidak dibungkus data)
            const token = res.data.token;
            const user = res.data.user;

            // 2. VALIDASI TOKEN UJIAN & AMBIL SOAL
            const valRes = await axios.post(`${API_URL}/siswa/validate?jadwalId=${JADWAL_ID}`, {
                token: EXAM_TOKEN,
                nisn: user.username
            }, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            const pesertaId = valRes.data.data.peserta_id;

            // Ambil daftar ID Soal yang asli dari server
            const soalRes = await axios.get(`${API_URL}/siswa/soal/${JADWAL_ID}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const soalIds = soalRes.data.data.map(s => s.id);

            students.push({
                id: i,
                username: username,
                token: token,
                peserta_id: pesertaId,
                soal_ids: soalIds
            });

            if (i % 20 === 0) console.log(`> ${i} siswa berhasil login & ambil soal...`);
        } catch (err) {
            console.error(`X Siswa ${i} gagal: ${err.response?.data?.error || err.message}`);
        }
    }

    console.log(`\n[2/3] Login Selesai. Total Siswa Aktif: ${students.length}`);
    console.log(`\n[3/3] Memulai Simulasi Sinkronisasi Jawaban (Looping)...`);

    let cycle = 1;
    setInterval(async () => {
        console.log(`\n--- Siklus Sinkronisasi #${cycle} ---`);
        const startTime = Date.now();
        let success = 0;
        let failed = 0;

        const requests = students.map(async (s) => {
            try {
                // Kirim 2 jawaban acak per siklus
                const randomSoalId = s.soal_ids[Math.floor(Math.random() * s.soal_ids.length)];
                const payload = {
                    peserta_ujian_id: s.peserta_id,
                    items: [
                        {
                            soal_id: randomSoalId,
                            jawaban_teks: ['A', 'B', 'C', 'D'][Math.floor(Math.random() * 4)],
                            ragu_ragu: Math.random() > 0.8
                        }
                    ],
                    sisa_waktu: 7200 - (cycle * 15)
                };

                await axios.post(`${API_URL}/siswa/sync`, payload, {
                    headers: { 'Authorization': `Bearer ${s.token}` }
                });
                success++;
            } catch (err) {
                failed++;
            }
        });

        await Promise.all(requests);
        const duration = Date.now() - startTime;
        
        console.log(`Hasil Siklus #${cycle}:`);
        console.log(`- Berhasil: ${success}`);
        console.log(`- Gagal   : ${failed}`);
        console.log(`- Durasi  : ${duration}ms`);
        
        if (failed > 0) {
            console.warn(`! PERINGATAN: Ada ${failed} request yang gagal. Cek beban CPU server.`);
        }

        cycle++;
    }, TEST_INTERVAL);
}

runStressTest();
