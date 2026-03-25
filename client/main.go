package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
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
	EnvFile    = "client/.env"
	NtfyURL    = "https://ntfy.sh"
)

func main() {
	loadEnv(EnvFile)

	if !isServerReachable(fmt.Sprintf("%s:22", RemoteIP)) {
		log.Printf("Server %s is not reachable on port 22, skipping backup", RemoteIP)
		return
	}

	lockFile, err := acquireLock("/tmp/backup-tool.lock")
	if err != nil {
		log.Printf("Could not acquire lock (another backup may be running): %v", err)
		return
	}
	defer lockFile.Close()

	file, err := os.Open(ConfigFile)
	if err != nil {
		notify(fmt.Sprintf("Failed to open config: %v", err))
		log.Fatalf("Failed to open config: %v", err)
	}
	defer file.Close()
	jobs, err := parseConfig(file)
	if err != nil {
		notify(fmt.Sprintf("Failed to read config: %v", err))
		log.Fatalf("Failed to read config: %v", err)
	}

	day := time.Now().Format("2006-01-02")
	historyDir := fmt.Sprintf(".rsync-history/%s", day)
	remoteTarget := fmt.Sprintf("admin@%s", RemoteIP)

	var failures []string

	for _, job := range jobs {
		fullDest := fmt.Sprintf("%s/%s", DestRoot, job.DestSubdir)
		remoteHist := fmt.Sprintf("%s/%s", fullDest, historyDir)

		log.Printf("Starting backup: %s -> %s", job.Source, fullDest)

		// 1. Ensure remote history directory exists
		mkdirCmd := exec.Command("ssh", remoteTarget, "mkdir", "-p", remoteHist)
		if err := mkdirCmd.Run(); err != nil {
			msg := fmt.Sprintf("Error creating remote dir for %s: %v", job.DestSubdir, err)
			log.Print(msg)
			failures = append(failures, msg)
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
			msg := fmt.Sprintf("Rsync failed for %s: %v", job.Source, err)
			log.Print(msg)
			failures = append(failures, msg)
		}
	}

	if len(failures) > 0 {
		notify(fmt.Sprintf("Backup completed with %d error(s):\n%s", len(failures), strings.Join(failures, "\n")))
	} else {
		log.Printf("All %d backup(s) completed successfully", len(jobs))
		notify(fmt.Sprintf("All %d backup(s) completed successfully", len(jobs)))
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

func loadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(val))
	}
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

func notify(msg string) {
	ntfyTopic := os.Getenv("NTFY_TOPIC")
	if NtfyURL == "" || ntfyTopic == "" {
		return
	}
	url := fmt.Sprintf("%s/%s", NtfyURL, ntfyTopic)
	resp, err := http.Post(url, "text/plain", strings.NewReader(msg))
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return
	}
	resp.Body.Close()
}
