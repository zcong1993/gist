package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	config = ".gistrc"
	url    = "https://api.github.com/gists"
)

// Data is post data struct
type Data struct {
	Public      bool            `json:"public"`
	Files       map[string]File `json:"files"`
	Description string          `json:"description"`
}

// File is single file struct
type File struct {
	Content string `json:"content"`
}

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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
	if *version {
		fmt.Println(Version())
		return
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
	var data Data
	data.Public = *isPublic
	data.Files = map[string]File{}
	data.Description = *description
	files := flag.Args()
	if len(files) == 0 {
		log.Fatal("should add some files")
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		checkError(err)
		file = path.Base(file)
		data.Files[file] = File{string(content)}
	}
	js, err := json.Marshal(&data)
	checkError(err)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(js)))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+key)
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	json, err := simplejson.NewJson([]byte(response))
	link, err := json.Get("html_url").String()
	_, user := json.CheckGet("owner")
	checkError(err)
	fmt.Printf("\nsuccess: link is %s\n", link)
	if !user {
		fmt.Println("\nwarning: gist owner is null, maybe your token is not correct!")
	}
}
