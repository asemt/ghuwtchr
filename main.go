// ghuwtchr project main.go
package main

import (
	"bufio"
	"fmt"
	"github.com/asemt/ghuwtchr/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/asemt/ghuwtchr/Godeps/_workspace/src/golang.org/x/oauth2"
	"log"
	"os"
	"strings"
	"sync"
)

// Stores required data to unwatch a repo.
type UwtchData struct {
	owner string
	repo  string
}

// Extract the owner of a repo and the repo name.
// We want to accept two shapes of links here:
//   'https://github.com/<owner>/<repo-name>/subscription'
//   'https://github.com/<owner>/<repo-name>'
func extractUwtchData(ghLnk string) UwtchData {
	ghLnk = strings.TrimSpace(ghLnk)
	repoLnkSpltd := strings.Split(ghLnk, "/")
	var repoName string
	if repoLnkSpltd[len(repoLnkSpltd)-1] == "subscription" {
		repoName = repoLnkSpltd[len(repoLnkSpltd)-2 : len(repoLnkSpltd)-1][0]
	} else {
		repoName = repoLnkSpltd[len(repoLnkSpltd)-1]
	}
	return UwtchData{repo: repoName, owner: repoLnkSpltd[3]}
}

func main() {
	var wg sync.WaitGroup

	ghTkn := os.Getenv("GHACCESSTOKEN")
	if ghTkn == "" {
		log.Fatalf("(main) >>  Error. GHACCESSTOKEN not set in environment.")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: fmt.Sprintf("%s", ghTkn)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	// Create a new, already authenticated client.
	c := github.NewClient(tc)

	// Collect the user input.
	var repoLinks2Unwatch []string
	log.Println("***  Paste GitHub links ending w/o '/subscription' below (e.g. 'https://github.com/<owner>/<repo-name>'). A single '.' on an own line ends the input.  ***")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		l := scanner.Text()
		if l == "." {
			break
		}
		repoLinks2Unwatch = append(repoLinks2Unwatch, l)
	}

	// Extract repo names from GitHub links.
	var repos2Unwatch []UwtchData
	for _, rpltu := range repoLinks2Unwatch {
		rud := extractUwtchData(rpltu)
		repos2Unwatch = append(repos2Unwatch, rud)
	}

	// Unwatch the repos in a concurrent manner.
	for _, rtuw := range repos2Unwatch {
		wg.Add(1)
		/*
			We're using a closure here as a Go routine to make sure the correct value of the 'rtuw' variable
			get's passed to the Go routine:
				By adding uwtchData as a parameter to the closure, val is evaluated at each iteration
				and placed on the stack for the goroutine, so each slice element is available
				to the goroutine when it is eventually executed.

			Adapted quote from:
				https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		*/
		go func(uwtchData UwtchData, wg *sync.WaitGroup) {
			_, err := c.Activity.DeleteRepositorySubscription(uwtchData.owner, uwtchData.repo)
			if err != nil {
				log.Printf("(func) >>  Error: %s", err.Error())
				return
			}
			log.Printf("(func) >>  Successfully unwatched: '%s'\n", uwtchData.repo)
			wg.Done()
		}(rtuw, &wg)
	}
	// Wait for all Go routines to finish.
	wg.Wait()
	log.Println("(main) >>  Done.")

}
