package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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

// feature/PLSCM-2450-abc-xyz
func getCommitMsgPrepend(branch string) (prepend string, err error) {
	p := strings.Split(branch, "/")
	if len(p) != 2 {
		return prepend, errors.New("Branch has no / in it. Failed to split")
	}

	typeOfBranch, ticketDescription := strings.ToLower(p[0]), strings.ToLower(p[1])

	switch {
	case strings.HasPrefix(typeOfBranch, "feat"):
		prepend = "feat"
	case strings.HasPrefix(typeOfBranch, "bugfix"):
		prepend = "bugfix"
	case strings.HasPrefix(typeOfBranch, "ci"):
		prepend = "ci"
	default:
		return prepend, fmt.Errorf("Branch name: %v does not meet any pattern", branch)
	}

	switch {
	case strings.HasPrefix(ticketDescription, "plscm-"):
		prepend = prepend + " PLSCM-"
		ticketDescription = strings.TrimPrefix(ticketDescription, "plscm-")
	case strings.HasPrefix(strings.ToLower(ticketDescription), "poicm-"):
		ticketDescription = strings.TrimPrefix(ticketDescription, "poicm-")
		prepend = prepend + "POICM-"
	default:
		return prepend, fmt.Errorf("Did not find project name in branch name: %v", branch)
	}

	idx := strings.Index(ticketDescription, "-")
	if idx == -1 {
		return prepend, fmt.Errorf("Did not find ticket number in branch name: %v", branch)
	}
	ticketNumber := ticketDescription[:idx]
	prepend = prepend + ticketNumber
	return

}

func getHeadRef(dotGitPath string) (string, error) {
	headFile, err := os.Open(fmt.Sprintf("%v/HEAD", dotGitPath))
	defer headFile.Close()
	if err != nil {
		return "", err
	}
	refBytes, err := io.ReadAll(headFile)
	if err != nil {
		return "", err
	}
	return string(refBytes), nil
}

func findDotGitFolder(pwd string) string {
	gitParentDir := pwd
	var pathsWalked []string
	for _ = range 10 {
		// Search for .git here
		pathsWalked = append(pathsWalked, gitParentDir)
		if doesDotGitDirExists(gitParentDir) {
			fmt.Printf("Found .git folder in path: %v\n", gitParentDir)
			break
		}
		parentDir := filepath.Dir(gitParentDir)
		if parentDir == gitParentDir {
			fmt.Println( "Searched for .git the following folders: " )
			for _, path := range pathsWalked {
				fmt.Println(path)
			}
			fmt.Println("But could not find .git folder")
			os.Exit(1)
		}
		gitParentDir = parentDir
	}
	return gitParentDir
}

func main() {
	// Get the current working dir
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Current Working Dir: %v\n", pwd)

	// Find the path where .git is present
	gitParentDir := findDotGitFolder(pwd)
	dotGitDirPath := fmt.Sprintf("%v/.git", gitParentDir)

	ref, err := getHeadRef(dotGitDirPath)
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
	fmt.Printf("Current Branch: %v\n", branchName)

	prepend, err := getCommitMsgPrepend(branchName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var commitMsg string
	fmt.Println("Enter your commit after following prompt")
	fmt.Printf(prepend)
	fmt.Scanf("%s", &commitMsg)
	fmt.Println(prepend, commitMsg)
	//cmd := exec.Command("git", "commit", "-m", prepend + commitMsg)
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println("Error: ", err)
	//	fmt.Println(string(output))
	//	return
	//}
}
