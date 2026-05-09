const axios = require('axios');

const BASE_URL = 'http://localhost:3000';
const API_URL = `${BASE_URL}/api`;
const NUM_STUDENTS = 100;
const TEST_INTERVAL = 5000; // SANGAT CEPAT: 5 detik (Aslinya 15 detik)
const EXAM_TOKEN = 'TEST100'; 
const JADWAL_ID = 1;        

// Daftar gambar asli untuk simulasi beban bandwidth
const SAMPLE_IMAGES = [
    '/uploads/soal/f9edfd39-8fb8-4679-b72f-7409246f4c79.jpeg',
    '/uploads/soal/e003e882-5070-4140-98bf-37adf2ce49f2.jpeg',
    '/uploads/soal/b2eaa3b5-59a3-4ee7-9486-8d7912ef5167.jpeg',
    '/uploads/soal/8e655470-7833-4bd3-b56d-137b058a0d04.jpeg'
];

const LONG_TEXT = "Ini adalah simulasi jawaban essay yang sangat panjang. Siswa sedang mengetik jawaban dengan sangat detail untuk menjelaskan teori yang diminta oleh guru. Jawaban ini mengandung banyak karakter untuk mengetes apakah database sanggup menampung teks besar dari 100 siswa sekaligus tanpa mengalami perlambatan atau hang. ".repeat(5);

async function runExtremeTest() {
    console.log(`\x1b[31m=== CBT MAXIMUM EXTREME STRESS TEST ===\x1b[00m`);
    console.log(`Target   : ${NUM_STUDENTS} Siswa serentak`);
    console.log(`Interval : ${TEST_INTERVAL/1000} detik (Mode Siksa)`);
    console.log(`Beban    : Essay Panjang + Download Gambar\n`);
    
    const students = [];

    // 1. LOGIN MASSAL
    console.log(`[1/3] Memulai Login Massal...`);
    for (let i = 1; i <= NUM_STUDENTS; i++) {
        const username = `siswa${i}`;
        try {
            const res = await axios.post(`${API_URL}/auth/login`, { username, password: 'siswa123' });
            const token = res.data.token;
            const user = res.data.user;

            const valRes = await axios.post(`${API_URL}/siswa/validate?jadwalId=${JADWAL_ID}`, {
                token: EXAM_TOKEN, nisn: user.username
            }, { headers: { 'Authorization': `Bearer ${token}` } });

            const soalRes = await axios.get(`${API_URL}/siswa/soal/${JADWAL_ID}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const soalIds = soalRes.data.data.map(s => s.id);

            students.push({ id: i, token, peserta_id: valRes.data.data.peserta_id, soal_ids: soalIds });
            if (i % 20 === 0) console.log(`> ${i} siswa siap menyerbu...`);
        } catch (err) {
            const errorMsg = err.response?.data?.error || err.message;
            console.error(`X Siswa ${i} gagal: ${errorMsg}`);
        }
    }

    console.log(`\n[2/3] Login Selesai. Total Siswa Aktif: ${students.length}`);
    console.log(`[3/3] Memulai Siksaan (Looping 5 detik)...`);

    let cycle = 1;
    setInterval(async () => {
        console.log(`\n--- SIKSAAN SIKLUS #${cycle} ---`);
        const startTime = Date.now();
        let successSync = 0;
        let successImg = 0;
        let failed = 0;

        const requests = students.map(async (s) => {
            try {
                // A. SIMULASI DOWNLOAD 2 GAMBAR ACAK
                const img1 = SAMPLE_IMAGES[Math.floor(Math.random() * SAMPLE_IMAGES.length)];
                const img2 = SAMPLE_IMAGES[Math.floor(Math.random() * SAMPLE_IMAGES.length)];
                await axios.get(`${BASE_URL}${img1}`);
                await axios.get(`${BASE_URL}${img2}`);
                successImg += 2;

                // B. SIMULASI KIRIM 2 JAWABAN ESSAY PANJANG
                const sId1 = s.soal_ids[0] || 1;
                const sId2 = s.soal_ids[1] || 2;
                
                await axios.post(`${API_URL}/siswa/sync`, {
                    peserta_ujian_id: s.peserta_id,
                    items: [
                        { soal_id: sId1, jawaban_teks: LONG_TEXT, ragu_ragu: false },
                        { soal_id: sId2, jawaban_teks: LONG_TEXT, ragu_ragu: false }
                    ],
                    sisa_waktu: 7200
                }, { headers: { 'Authorization': `Bearer ${s.token}` } });
                
                successSync++;
            } catch (err) {
                failed++;
                console.error(`! Gagal di siswa ${s.id}: ${err.message}`);
            }
        });

        await Promise.all(requests);
        const duration = Date.now() - startTime;
        
        console.log(`\x1b[33mLaporan Siklus #${cycle}:\x1b[00m`);
        console.log(`- Berhasil Simpan Essay : ${successSync}`);
        console.log(`- Berhasil Tarik Gambar : ${successImg}`);
        console.log(`- Gagal Total           : ${failed}`);
        console.log(`- Kecepatan Proses      : \x1b[1m${duration}ms\x1b[00m`);
        
        if (duration > TEST_INTERVAL) {
            console.warn(`\x1b[31m! BAHAYA: Server mulai kewalahan, durasi melebihi interval!\x1b[00m`);
        }
        cycle++;
    }, TEST_INTERVAL);
}

runExtremeTest();
