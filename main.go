package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/google/go-github/v60/github"
)

var (
	// Authentication
	token string = os.Getenv("GH_TOKEN")

	// Setup
	repoOwner   string = os.Getenv("GITHUB_REPOSITORY_OWNER")
	repoName    string = strings.TrimPrefix(os.Getenv("GITHUB_REPOSITORY"), fmt.Sprintf("%s/", repoOwner))
	sha         string = os.Getenv("GITHUB_SHA")
	outputPath  string = os.Getenv("GITHUB_OUTPUT")
	summaryPath string = os.Getenv("GITHUB_STEP_SUMMARY")

	// Configuration
	patchLabels  string = getEnvDefault("PATCH_LABELS", "patch")
	minorLabels  string = getEnvDefault("MINOR_LABELS", "minor")
	majorLabels  string = getEnvDefault("MAJOR_LABELS", "major")
	noBumpLabels string = getEnvDefault("NO_BUMP_LABELS", "no bump")
	defaultBump  string = getEnvDefault("DEFAULT_BUMP", "none")
)

type Summary struct {
	FinalBump    string
	AllLabels    []string
	DidFindPR    bool
	DidFindLabel bool
	PRLink       string
	PRNumber     int
}

//go:embed summary.md
var summaryTmplRaw string
var summaryTmpl = template.Must(template.New("summary").Parse(summaryTmplRaw))

func getEnvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getFoundLabel(labels []*github.Label) (string, bool) {
	patchLs := strings.Split(patchLabels, ",")
	minorLs := strings.Split(minorLabels, ",")
	majorLs := strings.Split(majorLabels, ",")
	noBumpLs := strings.Split(noBumpLabels, ",")

	if contains(labels, noBumpLs) {
		return "none", true
	}
	if contains(labels, majorLs) {
		return "major", true
	}
	if contains(labels, minorLs) {
		return "minor", true
	}
	if contains(labels, patchLs) {
		return "patch", true
	}

	return "", false
}

func contains(s []*github.Label, e []string) bool {
	for _, a := range s {
		for _, b := range e {
			if a.GetName() == b {
				return true
			}
		}
	}
	return false
}

func writeOutput(variable, value string) {
	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s=%s\n", variable, value))
	if err != nil {
		panic(err)
	}
}

func writeSummary(summary *Summary) {
	f, err := os.OpenFile(summaryPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = summaryTmpl.Execute(f, summary)
	if err != nil {
		panic(err)
	}
}

func main() {

	if token == "" {
		fmt.Println("GH_TOKEN not set")
		os.Exit(1)
	}

	client := github.NewClient(nil).WithAuthToken(token)
	summary := &Summary{
		FinalBump:    defaultBump,
		DidFindLabel: false,
		DidFindPR:    false,
	}

	prs, _, err := client.PullRequests.ListPullRequestsWithCommit(context.Background(), repoOwner, repoName, sha, nil)
	if err != nil {
		panic(err)
	}

	if len(prs) > 0 {
		summary.DidFindPR = true
		pr := prs[0]
		summary.PRLink = pr.GetHTMLURL()
		summary.PRNumber = pr.GetNumber()
		for _, label := range pr.Labels {
			summary.AllLabels = append(summary.AllLabels, label.GetName())
		}

		if label, found := getFoundLabel(pr.Labels); found {
			summary.FinalBump = label
			summary.DidFindLabel = true
		}
	}

	writeOutput("bump", summary.FinalBump)
	writeSummary(summary)

}
