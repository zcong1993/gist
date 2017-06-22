package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const config = ".gistrc"

func checkConf(path string) string {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("set key first, use 'gist -s <key>'")
	}
	if string(key) == "" {
		log.Fatal("set key first, use 'gist -s <key>'")
	}
	return string(key)
}

func stringAddress(v string) *string { return &v }

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getFiles(files []string) (map[github.GistFilename]github.GistFile, error) {
	eg := errgroup.Group{}
	results := map[github.GistFilename]github.GistFile{}
	for _, file := range files {
		file := file
		eg.Go(func() error {
			content, err := ioutil.ReadFile(file)
			fmt.Fprintf(os.Stdout, "--> Parsing file: %15s\n", file)
			if err != nil {
				return errors.Wrapf(err,
					"failed to get file content: %s", file)
			}
			results[github.GistFilename(path.Base(file))] = github.GistFile{Filename: stringAddress(string(path.Base(file))), Content: stringAddress(string(content))}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, errors.Wrap(err, "one of the goroutines failed")
	}
	return results, nil
}

func main() {
	home, err := homedir.Dir()
	checkError(err)
	configFile := path.Join(home, config)
	help := flag.Bool("h", false, "show help")
	version := flag.Bool("v", false, "show version")
	setKey := flag.Bool("s", false, "set token for auth")
	delKey := flag.Bool("r", false, "remove token")
	isPublic := flag.Bool("p", false, "create public gist?")
	description := flag.String("d", "published by 'zcong1993/gist' with golang", "add custom description")
	flag.Parse()
	if *version {
		Version()
		os.Exit(0)
	}
	if *help {
		fmt.Println("\nUsage :\n\tgist [flag] [files...]")
		fmt.Println("\nFlags :")
		fmt.Println()
		fmt.Println("\t -s, \t set token for auth")
		fmt.Println("\t -r, \t remove token")
		fmt.Println("\t -p, \t create public gist?")
		fmt.Println("\t -d, \t add custom description, default is `published by 'zcong1993/gist' with golang`")
		fmt.Println("\t -h, \t show help")
		os.Exit(0)
	}
	if *setKey {
		if len(flag.Args()) == 0 {
			log.Fatal("token is required")
		}
		err = ioutil.WriteFile(configFile, []byte(flag.Args()[0]), 0644)
		checkError(err)
		println("token set success")
		os.Exit(0)
	}
	if *delKey {
		err = ioutil.WriteFile(configFile, []byte(""), 0644)
		checkError(err)
		println("token delete success")
		os.Exit(0)
	}
	key := checkConf(configFile)
	files := flag.Args()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: key},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	postFiles, err := getFiles(files)
	checkError(err)
	gist := github.Gist{
		Public:      isPublic,
		Files:       postFiles,
		Description: description,
	}
	g, _, err := client.Gists.Create(context.Background(), &gist)
	checkError(err)
	fmt.Printf("\nDone! The gist url is %s\n", *g.HTMLURL)
}
