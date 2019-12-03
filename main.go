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
	logComparisons(dirA, dirB, comparisons)
}

type comparison struct {
	hashA string
	hashB string
}

func logComparisons(dirA, dirB string, comparisons map[string]comparison) {
	dupeCount, diffCount, riggedCount := 0, 0, 0
	for relPath, comp := range comparisons {
		if comp.hashA != comp.hashB {
			if comp.hashA != "" && comp.hashB != "" {
				labeText, labelColor := "Diff", color.FgYellow
				commitA, commitB := commitForSha(dirA, relPath, comp.hashA), commitForSha(dirA, relPath, comp.hashB)
				if commitA == "" || commitB == "" {
					riggedCount++
					labeText, labelColor = "Diff", color.FgRed
				}
				commitAStr, commitBStr := richCommit(commitA, color.FgHiBlack), richCommit(commitB, color.FgYellow)
				fmt.Printf("%s: %s: %s %s vs %s %s\n", colorize(labeText, labelColor, true), relPath,
					colorize(comp.hashA[:10], color.FgHiBlack, true), commitAStr, colorize(comp.hashB[:10], color.FgHiBlack, true), commitBStr)
				diffCount++
			}
		} else {
			fmt.Printf("%s: %s @ %s\n", colorize("Dupe", color.FgHiWhite, true), relPath, colorize(comp.hashA[:10], color.FgHiBlack, true))
			dupeCount++
		}
	}
	fmt.Printf("There are %d coinciding paths. Out of these, %d have matching files and %d have differing files.\n", dupeCount+diffCount, dupeCount, diffCount)
	if riggedCount > 0 {
		fmt.Printf("Out of the %d different files, "+colorize("%d files have modifications unknown to the repository at %s", color.FgRed, true)+".\n",
			diffCount, riggedCount, gitRoot(dirA))
	}
}

func richCommit(commit string, col color.Attribute) string {
	if commit == "" {
		return colorize("(NO MATCHING COMMIT)", color.FgRed, true)
	}
	return fmt.Sprintf("(Commit: %s)", colorize(commit[:10], col, true))
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
		fatalf("Could not read dir \"%s\": %v", absPath, err)
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
