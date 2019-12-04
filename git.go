package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// returns most recent commit hash where SHA-256 of file rootPath/relPath matches sha
func commitForSha(repoRoot, rootPath, relPath, sha string) string {
	// fmt.Println("commitForSha", rootPath, relPath, sha)
	shas := shasOfFile(repoRoot, rootPath, relPath)
	// fmt.Println("SHAs of file", filepath.Join(rootPath, relPath), "=", shas)
	for contentSha, commit := range shas {
		if sha == contentSha {
			return commit
		}
	}
	return ""
}

func gitRoot(atPath string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	abs, err := filepath.Abs(atPath)
	if err != nil {
		panic(err)
	}
	cmd.Dir = abs
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return strings.Trim(strings.TrimSpace(string(out)), "\n\t\r")
}

// map-key is SHA-256 of file content; map-value is commit SHA
func shasOfFile(repoRoot, rootPath, relPath string) map[string]string {
	fileCommits := getCommitsForFile(repoRoot, rel(repoRoot, filepath.Join(rootPath, relPath)))
	// fmt.Println("Commits of file", filepath.Join(rootPath, relPath), "=", fileCommits)
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
		if err != nil {
			fmt.Println("git show", commit+":"+relPath, "in", cmd.Dir)
			panic(err)
		}
		m[hex.EncodeToString(hash.Sum(nil))] = commit
	}
	return m
}

func rel(basePath, targPath string) string {
	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		panic(err)
	}
	absTargPath, err := filepath.Abs(targPath)
	if err != nil {
		panic(err)
	}
	relPath, err := filepath.Rel(absBasePath, absTargPath)
	if err != nil {
		panic(err)
	}
	return relPath
}

func getCommitsForFile(repoRoot, repoRelPath string) []string {
	cmd := exec.Command("git", "log", "--pretty=tformat:\"%H\"", repoRelPath)
	var err error
	cmd.Dir = repoRoot
	// fmt.Println("Run", cmd.Path, cmd.Args, "in", cmd.Dir)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	// fmt.Println("output=", string(output))
	commits := strings.FieldsFunc(strings.ReplaceAll(string(output), `"`, ``), func(r rune) bool {
		return r == '\n'
	})
	// fmt.Println("commits=", commits)
	return commits
}
