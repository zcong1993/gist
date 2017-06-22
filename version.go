package main

import (
	"bytes"
	"fmt"
	latest "github.com/tcnksm/go-latest"
	"time"
)

const (
	// AppName is the cli name
	AppName = "gist"
	// AppVersion is current version of cli
	AppVersion = "v2.1.0"

	defaultCheckTimeout = 2 * time.Second
)

// GitCommit is cli current git commit hash
var GitCommit string

// Version show the cli's current version
func Version() {
	version := fmt.Sprintf("\n%s %s", AppName, AppVersion)
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

		res, err := latest.Check(githubTag, fixFunc(AppVersion))
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
