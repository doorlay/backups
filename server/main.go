package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	lockFileName = "ente-sync.lock"
	timeout      = 6 * time.Hour
)

func main() {
	log.SetFlags(log.LstdFlags)

	lockFile := filepath.Join(os.TempDir(), lockFileName)

	exportDir := os.Getenv("EXPORT_DIR")
	secretsPath := os.Getenv("SECRETS_PATH")

	if exportDir == "" {
		log.Fatal("EXPORT_DIR must be set")
	}

	if secretsPath == "" {
		log.Fatal("SECRETS_PATH must be set")
	}

	// Prevent overlapping runs
	lockHandle, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("failed to open lock file: %v", err)
	}
	defer lockHandle.Close()

	if err := syscall.Flock(int(lockHandle.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		log.Println("another sync process is already running, exiting")
		os.Exit(0)
	}
	defer syscall.Flock(int(lockHandle.Fd()), syscall.LOCK_UN)

	if err := os.MkdirAll(filepath.Dir(secretsPath), 0700); err != nil {
		log.Fatalf("failed to create secrets directory: %v", err)
	}

	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Fatalf("failed to create export directory: %v", err)
	}

	os.Setenv("ENTE_CLI_SECRETS_PATH", secretsPath)

	args := []string{"export", exportDir}

	if albums := os.Getenv("ALBUMS"); albums != "" {
		args = append(args, "--albums", albums)
	}

	if os.Getenv("INCLUDE_HIDDEN") == "true" {
		args = append(args, "--hidden=true")
	}

	log.Printf("starting ente export to %s", exportDir)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ente", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("ente export timed out after %s", timeout)
		}
		log.Fatalf("ente export failed: %v", err)
	}

	log.Printf("export completed successfully")
}
