package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func doesDotGitDirExists(dir string) (bool) {
	gitLocation := fmt.Sprintf("%v/.git", dir)
	fileInfo, err := os.Stat(gitLocation)
	if err != nil {
		return false
	}
	if !fileInfo.IsDir() {
		return false
	}
	return true
}

func main() {
	// Get the current working dir
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Current Working Dir: %v\n", pwd)

	// Find the path where .git is present
	gitFolderPath := pwd
	var pathsWalked []string
	for _ = range 10 {
		// Search for .git here
		pathsWalked = append(pathsWalked, gitFolderPath)
		if doesDotGitDirExists(gitFolderPath) {
			fmt.Printf("Found .git folder in path: %v\n", gitFolderPath)
			break
		}
		parentDir := filepath.Dir(gitFolderPath)
		if parentDir == gitFolderPath {
			fmt.Println( "Searched for .git the following folders: " )
			fmt.Println(pathsWalked)
			fmt.Println("But could not find .git folder")
			os.Exit(1)
		}
		gitFolderPath = parentDir
	}


	headFile, err := os.Open(fmt.Sprintf("%v/.git/HEAD", gitFolderPath))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer headFile.Close()

	ref, err := io.ReadAll(headFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	currRef := string(ref)
	splitCurrRef := strings.Fields(currRef)
	if len(splitCurrRef) != 2 {
		fmt.Printf("Ref in HEAD file has more than 2 word: %v\n", currRef)
		os.Exit(1)
	}

	branchName := strings.TrimPrefix(splitCurrRef[1], "refs/heads/")
	fmt.Printf("current branch: %v", branchName)
}
