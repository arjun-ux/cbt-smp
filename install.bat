@echo off
TITLE Instalasi CBT SMP
echo ===========================================
echo      PROSES INSTALASI APLIKASI CBT
echo ===========================================
echo.

echo [1/4] Mengunduh dependencies Frontend...
cd frontend
call npm install
echo.

echo [2/4] Membangun Frontend (Production Build)...
call npm run build
cd ..
echo.

echo [3/4] Mengompilasi Backend (Go)...
go build -o cbt-app.exe ./cmd/server/main.go
echo.

echo [4/4] Finalisasi Instalasi...
echo.

echo ===========================================
echo      INSTALASI SELESAI!
echo ===========================================
echo Silakan jalankan "cbt-app.exe" untuk memulai server.
echo Setelah server menyala, instruksi akses akan muncul di terminal.
echo.
pause
