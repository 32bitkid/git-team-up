package main

import "fmt"
import "flag"
import "strings"
import "os"
import "log"
import "github.com/32bitkid/gitcmd"
import "encoding/json"

var repoFlag string
var repo *GitCommands

type Config map[string][]string

type GitCommands struct {
	Update         gitcmd.Action     `gitcmd:"remote update"`
	ShortenHash    gitcmd.SingleLine `gitcmd:"rev-parse --short"`
	RemoteBranches gitcmd.MultiLine  `gitcmd:"for-each-ref --format=%(refname) refs/remotes"`
	MergeBase      gitcmd.SingleLine `gitcmd:"merge-base --octopus"`
	CurrentBranch  gitcmd.SingleLine `gitcmd:"symbolic-ref --short -q HEAD"`
	Checkout       gitcmd.Action     `gitcmd:"checkout"`
	Merge          gitcmd.Action     `gitcmd:"merge -s octopus"`
	Reset          gitcmd.Action     `gitcmd:"reset --hard"`
	Branch         gitcmd.Action     `gitcmd:"branch"`
}

func init() {
	log.SetFlags(0)
	flag.StringVar(&repoFlag, "repo", "", "")
}

func filterBranches(root string) []string {

	return []string{
		"refs/remotes/origin/refresh_historical_velocity_S-43455",
		"refs/remotes/origin/modernize_move_to_epic_S-41366",
		"refs/remotes/origin/scriptsite_planningroom_S-42691",
	}

	matches := make([]string, 0)

	if branches, err := repo.RemoteBranches(); err == nil {
		for _, name := range branches {
			if strings.HasPrefix(name, root) {
				matches = append(matches, name)
			}
		}
	} else {
		os.Exit(-1)
	}
	return matches
}

func main() {
	flag.Parse()

	file, _ := os.Open("team-up.json")
	decoder := json.NewDecoder(file)
	var configuration Config
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatalf("Error parsing `team-up.json` -- %s", err)
	}

	os.Exit(0)

	repo = &GitCommands{}
	gitcmd.InitRepo(repoFlag, repo)

	repo.Update()

	targetBranches := filterBranches("refs/remotes/origin/team")
	fmt.Printf("%#v\n", targetBranches)

	base, mergeBaseErr := repo.MergeBase(targetBranches...)
	if mergeBaseErr != nil {
		log.Fatalf("\tCannot find a merge-base: %s\n\n", mergeBaseErr)
	}

	current, currentBranchErr := repo.CurrentBranch()

	if currentBranchErr == nil {
		fmt.Printf("Switching from \"%s\" to \"%s\".\n", current, base)
	} else {
		fmt.Printf("Checking out \"%s\".\n", base)
	}

	if err := repo.Checkout(base); err != nil {
		log.Fatalf("Could not checkout \"%s\": %s\n\n", base, err)
	}

	if err := repo.Merge(targetBranches...); err != nil {
		log.Fatalf("Merge failed....\n\n")
	}

	if err := repo.Branch("-f", "-q", "team/imua", "HEAD"); err != nil {
		log.Fatalf("Error updating branch....\n\n")
		os.Exit(-1)
	}

	if currentBranchErr == nil {
		fmt.Printf("Switching back to \"%s\".\n", current)
		repo.Checkout(current)
	}
}
