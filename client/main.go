package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type BackupJob struct {
	Source     string
	DestSubdir string
}

const (
	RemoteIP   = "192.168.1.216"
	DestRoot   = "/data/backups"
	ConfigFile = "client/backups.conf"
)

func main() {
	if !isServerReachable(fmt.Sprintf("%s:22", RemoteIP)) {
		return
	}

	lockFile, err := acquireLock("/tmp/backup-tool.lock")
	if err != nil {
		return
	}
	defer lockFile.Close()

	file, err := os.Open(ConfigFile)
	if err != nil {
		log.Fatalf("Failed to open config: %v", err)
	}
	defer file.Close()
	jobs, err := parseConfig(file)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	day := time.Now().Format("2006-01-02")
	historyDir := fmt.Sprintf(".rsync-history/%s", day)
	remoteTarget := fmt.Sprintf("admin@%s", RemoteIP)

	for _, job := range jobs {
		fullDest := fmt.Sprintf("%s/%s", DestRoot, job.DestSubdir)
		remoteHist := fmt.Sprintf("%s/%s", fullDest, historyDir)

		fmt.Printf(">> Starting backup: %s -> %s\n", job.Source, fullDest)

		// 1. Ensure remote history directory exists
		mkdirCmd := exec.Command("ssh", remoteTarget, "mkdir", "-p", remoteHist)
		if err := mkdirCmd.Run(); err != nil {
			log.Printf("Error creating remote dir for %s: %v", job.DestSubdir, err)
			continue
		}

		// 2. Execute rsync
		// Note: Using trailing slash on source to sync contents, not the directory itself
		rsyncArgs := []string{
			"-a", "--delete", "--backup", "--inplace",
			"--backup-dir=" + historyDir,
			"--exclude=/.rsync-history/",
			"-e", "ssh",
			job.Source + "/",
			fmt.Sprintf("%s:%s/", remoteTarget, fullDest),
		}

		rsyncCmd := exec.Command("rsync", rsyncArgs...)
		rsyncCmd.Stdout = os.Stdout
		rsyncCmd.Stderr = os.Stderr

		if err := rsyncCmd.Run(); err != nil {
			log.Printf("Rsync failed for %s: %v", job.Source, err)
		}
	}
}

func parseConfig(file *os.File) ([]BackupJob, error) {
	var jobs []BackupJob
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		jobs = append(jobs, BackupJob{
			Source:     os.ExpandEnv(strings.TrimSpace(parts[0])),
			DestSubdir: os.ExpandEnv(strings.TrimSpace(parts[1])),
		})
	}
	return jobs, scanner.Err()
}

func acquireLock(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	// LOCK_EX = Exclusive lock
	// LOCK_NB = Non-blocking (returns error immediately if held)
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return nil, err
	}
	return file, nil
}

func isServerReachable(address string) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
