package models

// CBTSetting adalah tabel key-value untuk menyimpan konfigurasi aplikasi dan identitas sekolah
type CBTSetting struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}
