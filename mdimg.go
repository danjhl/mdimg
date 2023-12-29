package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/f1bonacc1/glippy"
	"github.com/google/uuid"
	"github.com/skanehira/clipboard-image/v2"
)

func main() {
	u := flag.String("u", "", "url of image to insert")
	c := flag.Bool("c", false, "retrieve url from clipboard")
	o := flag.String("o", "", "output file")
	i := flag.Bool("i", false, "save image from clipboard")

	flag.CommandLine.Parse(os.Args[1:])

	tag, err := CreateImageTag(*u, *o, *i, *c, uuid.NewString)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(tag)
}

func CreateImageTag(u string, o string, i bool, c bool, uuid func() string) (string, error) {
	if u != "" && i {
		return "", errors.New("Cannot use -u with -i")
	} else if u != "" && c {
		return "", errors.New("Cannot use -u with -c")
	} else if i {
		return CreateImageTagFromRaw(o, uuid)
	} else if u == "" && !i && !c {
		return "", errors.New("Requires -u, -c or -i option")
	} else {
		url, err := GetUrl(u, c)
		if err == nil {
			return CreateImageTagFromUrl(url, o, uuid)
		}
		return "", err
	}
}

func CreateImageTagFromRaw(out string, uuid func() string) (string, error) {
	reader, err := clipboard.Read()
	if err != nil {
		return "", errors.New("Could not get image copy from clipboard, make sure you have copied an image")
	}

	fileName := out
	if fileName == "" {
		fileName = "./img/" + uuid()
	}
	if !strings.HasSuffix(fileName, ".png") {
		fileName = fileName + ".png"
	}

	err = createDirectories(fileName)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("![%s](%s)", fileName, fileName), nil
}

func GetUrl(url string, clipboard bool) (string, error) {
	if url != "" {
		return url, nil
	} else if url == "" && !clipboard {
		return "", errors.New("Requires -u or -c option")
	} else {
		clipUrl, err := glippy.Get()
		if err != nil {
			return "", errors.New("Could not read url from clipboard")
		}
		if clipUrl == "" {
			return "", errors.New("Clipboard is empty")
		}
		return clipUrl, nil
	}
}

func CreateImageTagFromUrl(url string, out string, uuid func() string) (string, error) {
	err := createDirectories(out)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Could not get url: '" + url + "'")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Url request responded with: " + strconv.Itoa(resp.StatusCode))
	}

	ext, err := getFileExt(resp.Header)
	if err != nil {
		return "", err
	}

	var fileName string
	if out != "" {
		fileName = out
		if !strings.HasSuffix(out, "."+ext) {
			fileName = fileName + "." + ext
		}
	} else {
		fileName = "./img/" + uuid() + "." + ext
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	file.Write(image)

	return fmt.Sprintf("![%s](%s)", fileName, fileName), nil
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
	if strings.HasPrefix(ct, "image") {
		es := strings.Split(ct, "/")
		if len(es) > 1 {
			return es[1], nil
		}
	}
	return "", errors.New("Cannot accept content-type: " + ct)
}
