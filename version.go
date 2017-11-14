package main

import (
	"bytes"
	"fmt"
	"time"

	latest "github.com/tcnksm/go-latest"
)

const (
	// AppName is the cli name
	AppName = "gist"

	defaultCheckTimeout = 2 * time.Second
)

// GitCommit is cli current git commit hash
var GitCommit string

// Config.
var version = "master"

// Version show the cli's current version
func Version() {
	version := fmt.Sprintf("\n%s %s", AppName, version)
	if len(GitCommit) != 0 {
		version += fmt.Sprintf(" (%s)", GitCommit)
	}
	version += "\nCopyright (c) 2017, zcong1993."
	fmt.Println(version)
	var buf bytes.Buffer
	verCheckCh := make(chan *latest.CheckResponse)
	go func() {
		fixFunc := latest.DeleteFrontV()
		githubTag := &latest.GithubTag{
			Owner:             "zcong1993",
			Repository:        "gist",
			FixVersionStrFunc: fixFunc,
		}

		res, err := latest.Check(githubTag, version)
		if err != nil {
			// Don't return error
			return
		}
		verCheckCh <- res
	}()

	select {
	case <-time.After(defaultCheckTimeout):
	case res := <-verCheckCh:
		if res.Outdated {
			fmt.Fprintf(&buf,
				"Latest version of gist is v%s, please upgrade!\n",
				res.Current)
		}
	}
	fmt.Print(buf.String())
}
