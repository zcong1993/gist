# gist
[![wercker status](https://app.wercker.com/status/a69dff6e96eca8984e750f81af414cc4/s/master "wercker status")](https://app.wercker.com/project/byKey/a69dff6e96eca8984e750f81af414cc4)
[![Build Status](https://travis-ci.org/zcong1993/gist.svg?branch=master)](https://travis-ci.org/zcong1993/gist)

> gist cli in go

## Usage

[download](https://github.com/zcong1993/gist/releases) the package and put in any `$PATH` folder.

```bash
# in your work folder
$ gist [flags] [files...]
# example
# set token
$ gist -s <your token> 
# create gist with file1 and file2
$ gist file1 file2 
# will return gist web link 
# create public gist
$ gist -p file 
# add your custom description
$ gist -d="your custom description" 
```

You can get your gist `token` here [https://github.com/settings/tokens/new](https://github.com/settings/tokens/new), remember to select `gist` scope.

## Build

```bash
$ git clone https://github.com/zcong1993/gist.git
$ cd gist
$ go build gist.go
# then move the output to your `$PATH` folder.
```

## License

MIT &copy; zcong1993
