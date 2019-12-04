package main

import "path/filepath"

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

func abs(filePath string) string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	return absPath
}
