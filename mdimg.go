package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	url := flag.String("u", "", "url of image to insert")
	out := flag.String("o", "", "output file")

	flag.CommandLine.Parse(os.Args[1:])
	tag, err := CreateImageTag(*url, *out, uuid.NewString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(tag)
}

func CreateImageTag(url string, out string, uuid func() string) (string, error) {
	if url == "" {
		return "", errors.New("Requires -u option")
	}

	err := createDirectories(out)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ext, err := getFileExt(resp.Header)
	if err != nil {
		return "", err
	}

	var download string
	if out != "" {
		download = out
		if !strings.HasSuffix(out, "."+ext) {
			download = download + "." + ext
		}
	} else {
		download = "./img/" + uuid() + "." + ext
	}

	file, err := os.Create(download)
	if err != nil {
		return "", err
	}
	defer file.Close()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	file.Write(image)

	return fmt.Sprintf("![%s](%s)", download, download), nil
}

func createDirectories(out string) error {
	dir := "./img"

	if out != "" {
		s := strings.Split(out, "/")
		if len(s) == 1 {
			return nil
		}
		dir = strings.Join(s[:len(s)-1], "/")
	}

	return os.MkdirAll(dir, os.ModePerm)
}

func getFileExt(header http.Header) (string, error) {
	ct := header.Get("content-type")
	if ct != "" {
		es := strings.Split(ct, "/")
		if len(es) > 1 && es[0] == "image" {
			return es[1], nil
		}
	}
	return "", errors.New("Cannot accept content-type: " + ct)
}
