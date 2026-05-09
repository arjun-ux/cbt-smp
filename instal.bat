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

echo [4/4] Membuat Script Peluncur...
(
echo @echo off
echo TITLE Menjalankan CBT SMP
echo for /f "tokens=4" %%%%a in ('route print ^^| findstr 0.0.0.0 ^^| findstr /v "0.0.0.0.0"') do set IP=%%%%a
echo echo ===========================================
echo echo      APLIKASI CBT SEDANG BERJALAN
echo echo ===========================================
echo echo AKSES LOKAL: http://localhost:3000
echo echo AKSES JARINGAN: http://%%IP%%:3000
echo echo ===========================================
echo start http://localhost:3000
echo cbt-app.exe
echo pause
) > jalankan.bat

echo.
echo ===========================================
echo      INSTALASI SELESAI!
echo ===========================================
echo Silakan klik "jalankan.bat" untuk memulai.
echo.
pause
