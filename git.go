package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// returns most recent commit hash where SHA-256 of file rootPath/relPath matches sha
func commitForSha(repoRoot, rootPath, relPath, sha string) string {
	shas := shasOfFile(repoRoot, rootPath, relPath)
	for contentSha, commit := range shas {
		if sha == contentSha {
			return commit
		}
	}
	return ""
}

func gitRoot(atPath string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	workDir, err := filepath.Abs(atPath)
	if err != nil {
		panic(err)
	}
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return abs(strings.Trim(strings.TrimSpace(string(out)), "\n\t\r"))
}

// map-key is SHA-256 of file content; map-value is commit SHA
func shasOfFile(repoRoot, rootPath, relPath string) map[string]string {
	fileCommits := getCommitsForFile(repoRoot, rel(repoRoot, filepath.Join(rootPath, relPath)))
	m := make(map[string]string)
	for _, commit := range fileCommits {
		cmd := exec.Command("git", "show", commit+":"+rel(repoRoot, filepath.Join(rootPath, relPath)))
		var err error
		cmd.Dir, err = filepath.Abs(rootPath)
		if err != nil {
			panic(err)
		}
		hash := sha256.New()
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(hash, pipe)
		if err != nil {
			panic(err)
		}
		err = cmd.Wait()
		if err == nil {
			// err will be != nil for example if the file was removed, i.e. became untracked with this commit
			m[hex.EncodeToString(hash.Sum(nil))] = commit
		}
	}
	return m
}

func getCommitsForFile(repoRoot, repoRelPath string) []string {
	cmd := exec.Command("git", "log", "--pretty=tformat:\"%H\"", repoRelPath)
	var err error
	cmd.Dir = repoRoot
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	commits := strings.FieldsFunc(strings.ReplaceAll(string(output), `"`, ``), func(r rune) bool {
		return r == '\n'
	})
	return commits
}
