# sitemap

This repo contains standalone, multi-threaded CLI application `smgen`, which based on html-parser written on Go to build site maps as suggested by [sitemaps.org](https://www.sitemaps.org/).

## `smgen` features

* сross-platform application as well as Go
* extracting links only from href-attributes of a-elements
* building maps and indexes in XML format only
* auto-splitting results into chunks
* auto-compressing results into gzip if needed
* support only for `lastmod` tag for maps and indexes

## Install `smgen` from source

At first, you should install [Go language distribution](https://golang.org/dl/) and make sure Go is installed correctly:

```cli
$ go.exe version
go version go1.12.5 windows/amd64
```

> Extension `.exe` is only required if you use WSL console like me, when Go was installed for Windows and outside WSL.

At second, you must provide GOPATH environment variable and prepare underlying directory layout as specified at [GitHub Wiki](https://github.com/golang/go/wiki/GOPATH), but in two words:

* make custom base directory and set it path as GOPATH
* make three sub-dirs under GOPATH: `bin`, `src` and `pkg`

Then run:

```cli
$ go.exe install github.com/wtask/sitemap/cmd/smgen
```

> You may need to install [git](https://git-scm.com/) also.

CLI tool now should be installed under GOPATH `bin` subdirectory.

## Binary distribution

Not provided and not planned.

## Docker image

You can use ready to run docker image from https://hub.docker.com/r/wtask/smgen to test `smgen` tool.

Or you may build docker image by yourself and run CLI-tool container as follows:

1. Clone or download this repository or release.
2. Jump into repository or release folder and build image:

```cli
docker build --no-cache -f ./cmd/smgen/.Dockerfile -t smgen .
```

3. Run container as:

```cli
docker run -it --rm -v "{your-local-dir}":/data smgen -depth={desired-depth} -num-workers={desired-num} {desired-URI}
```

>`{your-local-dir}` is absolute path to the directory where to save XML-sitemap and index;
>
>`{desired-depth}`, `{desired-num}`, `{desired-URI}` are CLI-tool argument values.

## Usage

Run `smgen` with `-help` option to get a quick reference:

```cli
.../bin$ smgen.exe -help
Generate site map suggested by https://www.sitemaps.org/protocol.html, starting from given URI:

        smgen [options] URI

Options:

  -depth uint
        Maximum depth of link-junctions from start URL to render site map. (default 1)
  -h
  -help
        Print usage help.
  -index-limit int
        Limit number of entries per index file. (default 50000)
  -index-name string
        Base name for site map INDEX. (default "sitemap_index")
  -map-limit int
        Limit number of entries per site map file. (default 50000)
  -map-name string
        Base name for site map FILE. (default "sitemap")
  -num-workers uint
        Number of allowed concurrent workers to build site map. (default 1)
  -output-dir string
        Output directory where site map and index will be generated. (default "C:\\Go\\bin")
  -size-limit int
        Maximum size of any generated file in bytes. If file size is greater than limitation, file is compressed into gzip. (default 52428800000)
```

Map generation example:

```cli
.../bin$ smgen.exe -depth=2 -num-workers=4 https://www.sitemaps.org/zh_CN/
smgen 2019/05/29 01:48:53 Started for "https://www.sitemaps.org/zh_CN/", depth: 2, workers: 4, output format: "xml", output dir: C:\Go\bin
smgen 2019/05/29 01:48:53 Parser has launched...
smgen 2019/05/29 01:48:55 Completed, num of links found: 5
smgen 2019/05/29 01:48:55 Started saving site map...
smgen 2019/05/29 01:48:55 MAP OK C:\Go\bin\sitemap.xml
smgen 2019/05/29 01:48:55 All done
```

## Feature plans

Fixing bugs as they are detected and minor improvements.
