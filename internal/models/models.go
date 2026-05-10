package models

import (
	"time"
)

// User merepresentasikan akun login
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);not null" json:"role"` // admin, guru, siswa
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Master Kelas
type MasterKelas struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Tingkat   string `gorm:"type:varchar(2);not null" json:"tingkat"`     // 7, 8, 9
	NamaKelas string `gorm:"type:varchar(20);not null" json:"nama_kelas"` // 7A, 8B
}

// Master Ruang
type MasterRuang struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	KodeRuang string `gorm:"type:varchar(10);uniqueIndex" json:"kode_ruang"`
	NamaRuang string `gorm:"type:varchar(50)" json:"nama_ruang"`
}

// Master Sesi
type MasterSesi struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	NamaSesi     string `gorm:"type:varchar(50)" json:"nama_sesi"`
	WaktuMulai   string `gorm:"type:varchar(10)" json:"waktu_mulai"`
	WaktuSelesai string `gorm:"type:varchar(10)" json:"waktu_selesai"`
}

// Master Siswa
type MasterSiswa struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"index" json:"user_id"`
	User        User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	NISN        string      `gorm:"type:varchar(20);uniqueIndex" json:"nisn"`
	NamaLengkap string      `gorm:"type:varchar(100);not null" json:"nama_lengkap"`
	KelasID     *uint       `gorm:"index" json:"kelas_id"`
	Kelas       MasterKelas `gorm:"foreignKey:KelasID;constraint:OnDelete:SET NULL;" json:"kelas"`
	RuangID     *uint       `gorm:"index" json:"ruang_id"`
	Ruang       MasterRuang `gorm:"foreignKey:RuangID;constraint:OnDelete:SET NULL;" json:"ruang"`
	SesiID      *uint       `gorm:"index" json:"sesi_id"`
	Sesi        MasterSesi  `gorm:"foreignKey:SesiID;constraint:OnDelete:SET NULL;" json:"sesi"`
}

// Master Guru
type MasterGuru struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"index" json:"user_id"`
	User     User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	NIP      string `gorm:"type:varchar(30);uniqueIndex" json:"nip"`
	NamaGuru string `gorm:"type:varchar(100);not null" json:"nama_guru"`
}

// Master Mapel
type MasterMapel struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	KodeMapel string `gorm:"type:varchar(20);uniqueIndex" json:"kode_mapel"`
	NamaMapel string `gorm:"type:varchar(100);not null" json:"nama_mapel"`
}

// CBT Bank Soal
type CBTBankSoal struct {
	ID                uint        `gorm:"primaryKey" json:"id"`
	GuruID            *uint       `gorm:"index" json:"guru_id"`
	Guru              MasterGuru  `gorm:"foreignKey:GuruID" json:"guru"`
	MapelID           uint        `gorm:"index" json:"mapel_id"`
	Mapel             MasterMapel `gorm:"foreignKey:MapelID" json:"mapel"`
	TingkatKelas      string      `gorm:"type:varchar(2)" json:"tingkat_kelas"` // 7, 8, 9
	JudulBankSoal     string      `gorm:"type:varchar(200);not null" json:"judul_bank_soal"`
	DefaultBobotPG    float64     `gorm:"default:1" json:"default_bobot_pg"`
	DefaultBobotEssay float64     `gorm:"default:1" json:"default_bobot_essay"`
	Status            string      `gorm:"type:varchar(20);default:'Draft'" json:"status"` // Draft, Aktif
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// CBT Soal
type CBTSoal struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	BankSoalID   uint        `gorm:"index" json:"bank_soal_id"`
	BankSoal     CBTBankSoal `gorm:"foreignKey:BankSoalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"bank_soal"`
	JenisSoal    string      `gorm:"type:varchar(10)" json:"jenis_soal"` // PG, ESSAY
	Pertanyaan   string      `gorm:"type:text" json:"pertanyaan"`
	OpsiA        string      `gorm:"type:text" json:"opsi_a"`
	OpsiB        string      `gorm:"type:text" json:"opsi_b"`
	OpsiC        string      `gorm:"type:text" json:"opsi_c"`
	OpsiD        string      `gorm:"type:text" json:"opsi_d"`
	KunciJawaban string      `gorm:"type:text;not null" json:"kunci_jawaban"`
	BobotNilai   float64     `gorm:"default:1" json:"bobot_nilai"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CBT Jadwal Ujian
type CBTJadwalUjian struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	BankSoalID   uint        `gorm:"index" json:"bank_soal_id"`
	BankSoal     CBTBankSoal `gorm:"foreignKey:BankSoalID" json:"bank_soal"`
	TanggalUjian string      `gorm:"type:date" json:"tanggal_ujian"`
	WaktuMulai   string      `gorm:"type:time" json:"waktu_mulai"`
	DurasiMenit  int         `json:"durasi_menit"`
	RuangID      *uint       `json:"ruang_id"`
	Ruang        MasterRuang `gorm:"foreignKey:RuangID" json:"ruang"`
	SesiID       *uint       `json:"sesi_id"`
	Sesi         MasterSesi  `gorm:"foreignKey:SesiID" json:"sesi"`
	PengawasID   *uint       `json:"pengawas_id"`
	Pengawas     MasterGuru  `gorm:"foreignKey:PengawasID" json:"pengawas"`
	AcakSoal     bool        `gorm:"default:false" json:"acak_soal"`
	AcakJawaban  bool        `gorm:"default:false" json:"acak_jawaban"`
	TokenUjian   string      `gorm:"type:varchar(10)" json:"token_ujian"`
	Status       string      `gorm:"type:varchar(20);default:'Belum Mulai'" json:"status"`
}

// CBT Peserta Ujian (Siswa yang mengikuti jadwal)
type CBTPesertaUjian struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	JadwalID          uint           `gorm:"index" json:"jadwal_id"`
	Jadwal            CBTJadwalUjian `gorm:"foreignKey:JadwalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"jadwal"`
	SiswaID           uint           `gorm:"index" json:"siswa_id"`
	Siswa             MasterSiswa    `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"siswa"`
	WaktuLogin        *time.Time     `json:"waktu_login"`
	WaktuSelesaiUjian *time.Time     `json:"waktu_selesai_ujian"`
	SisaWaktuDetik    int            `json:"sisa_waktu_detik"`
	StatusUjian       string         `gorm:"type:varchar(20);default:'Sedang Mengerjakan'" json:"status_ujian"`
	NilaiPG           float64        `gorm:"default:0" json:"nilai_pg"`
	NilaiEssay        float64        `gorm:"default:0" json:"nilai_essay"`
	TotalNilai        float64        `gorm:"default:0" json:"total_nilai"`
	IsTerblokir       bool           `gorm:"default:false" json:"is_terblokir"`
}

// CBT Jawaban Siswa
type CBTJawabanSiswa struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	PesertaUjianID uint            `gorm:"index" json:"peserta_ujian_id"`
	PesertaUjian   CBTPesertaUjian `gorm:"foreignKey:PesertaUjianID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"peserta_ujian"`
	SoalID         uint            `gorm:"index" json:"soal_id"`
	Soal           CBTSoal         `gorm:"foreignKey:SoalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"soal"`
	JawabanSiswa   string          `gorm:"type:text" json:"jawaban_siswa"`
	RaguRagu       bool            `gorm:"default:false" json:"ragu_ragu"`
	Skor           float64         `gorm:"default:0" json:"skor"`
}

// Log Sistem
type LogSistem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Aktivitas string    `gorm:"type:varchar(255)" json:"aktivitas"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

// Log Ujian
type LogUjian struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	PesertaUjianID uint            `gorm:"index" json:"peserta_ujian_id"`
	PesertaUjian   CBTPesertaUjian `gorm:"foreignKey:PesertaUjianID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"peserta_ujian"`
	KeteranganLog  string          `gorm:"type:varchar(255)" json:"keterangan_log"` // Login, Pindah Tab, Selesai
	Timestamp      time.Time       `gorm:"autoCreateTime" json:"timestamp"`
}

// CBT Rekap Nilai Permanen (Snapshot)
type CBTRekapNilai struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	JadwalID     uint      `gorm:"index" json:"jadwal_id"`
	NISN         string    `gorm:"type:varchar(20)" json:"nisn"`
	NamaSiswa    string    `gorm:"type:varchar(100)" json:"nama_siswa"`
	NamaKelas    string    `gorm:"type:varchar(20)" json:"nama_kelas"`
	NamaMapel    string    `gorm:"type:varchar(100)" json:"nama_mapel"`
	JudulUjian   string    `gorm:"type:varchar(200)" json:"judul_ujian"`
	TanggalUjian string    `json:"tanggal_ujian"`
	NilaiPG      float64   `json:"nilai_pg"`
	NilaiEssay   float64   `json:"nilai_essay"`
	TotalNilai   float64   `json:"total_nilai"`
	CreatedAt    time.Time `json:"created_at"`
}
