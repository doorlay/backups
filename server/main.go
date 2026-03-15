package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	LockFile = "/run/lock/ente-sync.lock"
	LogFile  = "ente-sync.log"
)

func main() {
	// 1. Load and Expand Environment Variables
	exportDir := os.ExpandEnv(os.Getenv("EXPORT_DIR"))
	secretsPath := os.ExpandEnv(os.Getenv("SECRETS_PATH"))
	if exportDir == "" || secretsPath == "" {
		log.Fatal("EXPORT_DIR and SECRETS_PATH must be set in environment")
	}

	// 2. Prevent Overlapping Runs
	file, err := os.OpenFile(LockFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not create lock file: %v", err)
	}
	defer file.Close()

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		fmt.Println("Another sync process is already running. Exiting.")
		os.Exit(0)
	}

	// 3. Setup Environment for Headless Execution
	// This tells Ente to use a flat file for keys instead of a GUI keyring
	os.Setenv("ENTE_CLI_SECRETS_PATH", secretsPath)
	os.MkdirAll(filepath.Dir(secretsPath), 0700)

	// 4. Ensure Export Directory Exists
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Fatalf("Failed to create export dir: %v", err)
	}

	// 5. Build Ente Arguments
	args := []string{"export"}
	if albums := os.Getenv("ALBUMS"); albums != "" {
		args = append(args, "--albums", albums)
	}
	// Note: Ente CLI flags use =true/false or are boolean flags
	if os.Getenv("INCLUDE_HIDDEN") == "true" {
		args = append(args, "--hidden=true")
	}

	// 6. Execute Export
	fmt.Printf("[%s] Starting Ente export...\n", time.Now().Format(time.RFC3339))

	cmd := exec.Command("ente", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Ente export failed: %v", err)
	}

	fmt.Printf("[%s] Export completed successfully.\n", time.Now().Format(time.RFC3339))
}
