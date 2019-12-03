package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func main() {
	dirA, dirB := parseFlags()
	dirAHierarchy, dirBHierarchy := scanDir(dirA, "."), scanDir(dirB, ".")
	comparisons := compareHierarchies(dirAHierarchy, dirBHierarchy)
	logComparisons(comparisons)
}

type comparison struct {
	hashA string
	hashB string
}

func logComparisons(comparisons map[string]comparison) {
	for relPath, comp := range comparisons {
		if comp.hashA != comp.hashB {
			if comp.hashA != "" && comp.hashB != "" {
				fmt.Printf("%s: %s: %s vs %s\n", colorize("DIFF", color.FgRed, true), relPath, comp.hashA[:8], comp.hashB[:8])
			}
		} else {
			fmt.Printf("%s: %s @ %s\n", colorize("DUPE", color.FgYellow, true), relPath, comp.hashA[:8])
		}
	}
}

func compareHierarchies(a, b map[string]string) map[string]comparison {
	m := make(map[string]comparison)
	for relPath, hashA := range a {
		m[relPath] = comparison{hashA: hashA}
	}
	for relPath, hashB := range b {
		comp, _ := m[relPath]
		comp.hashB = hashB
		m[relPath] = comp
	}
	return m
}

func parseFlags() (dirA, dirB string) {
	flag.Parse()
	if flag.NArg() != 2 {
		fatalf("Usage: duperig <dirA> <dirB>")
	}
	return flag.Arg(0), flag.Arg(1)
}

func scanDir(rootPath string, relPath string) map[string]string {
	m := make(map[string]string)
	absPath := filepath.Join(rootPath, relPath)
	infos, err := ioutil.ReadDir(absPath)
	if err != nil {
		fatalf("Could not read dir \"%s\": %v", scanDir, err)
	}
	for _, info := range infos {
		if info.IsDir() {
			subDirMap := scanDir(rootPath, filepath.Join(relPath, info.Name()))
			insertAll(subDirMap, m)
		} else {
			m[filepath.Join(relPath, info.Name())] = hash(filepath.Join(absPath, info.Name()))
		}
	}
	return m
}

func hash(filePath string) string {
	hash := sha256.New()
	f, err := os.Open(filePath)
	if err != nil {
		fatalf("Could not open file \"%s\" for hashing: %v", filePath, err)
	}
	defer f.Close()

	_, err = io.Copy(hash, f)
	if err != nil {
		fatalf("Could not hash file \"%s\": %v", filePath, err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func insertAll(from, to map[string]string) {
	for k, v := range from {
		to[k] = v
	}
}

func fatalf(formatMessage string, args ...interface{}) {
	fmt.Printf(formatMessage+"\n", args...)
	os.Exit(1)
}
