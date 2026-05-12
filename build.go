package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("      CBT SMP - MULTI-PLATFORM BUILD       ")
	fmt.Println("===========================================")

	// 1. Build Frontend
	fmt.Println("\n[1/3] Membangun Frontend...")
	if err := runCommand("./frontend", "npm", "run", "build"); err != nil {
		fmt.Printf("Gagal membangun frontend: %v\n", err)
		return
	}

	// 2. Build Backend
	fmt.Println("\n[2/3] Membangun Backend...")
	
	// Windows
	fmt.Println("- Windows (cbt-app-win.exe)...")
	if err := buildBackend("windows", "amd64", "cbt-app-win.exe"); err != nil {
		fmt.Printf("Gagal build Windows: %v\n", err)
	}

	// Linux
	fmt.Println("- Linux (cbt-app-linux)...")
	if err := buildBackend("linux", "amd64", "cbt-app-linux"); err != nil {
		fmt.Printf("Gagal build Linux: %v\n", err)
	}

	fmt.Println("\n[3/3] Selesai!")
	fmt.Println("Hasil build tersedia di direktori root.")
	fmt.Println("===========================================")
}

func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if name == "npm" && runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", name, args[0], args[1])
	}
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildBackend(osName, arch, output string) error {
	cmd := exec.Command("go", "build", "-o", output, "./cmd/server/main.go")
	cmd.Env = append(os.Environ(), "GOOS="+osName, "GOARCH="+arch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
