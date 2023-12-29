package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/f1bonacc1/glippy"
	"github.com/stretchr/testify/assert"
)

func TestInvalidCalls(t *testing.T) {
	_, err := CreateImageTag("url", "", true, false, func() string { return "" })
	assert.Equal(t, "Cannot use -u with -i", err.Error())

	_, err = CreateImageTag("abc", "", false, true, func() string { return "" })
	assert.Equal(t, "Cannot use -u with -c", err.Error())

	_, err = CreateImageTag("", "", false, false, func() string { return "" })
	assert.Equal(t, "Requires -u, -c or -i option", err.Error())

	_, err = CreateImageTag(" ", "", false, false, func() string { return "" })
	assert.Equal(t, "Could not get url: ' '", err.Error())

	_, err = CreateImageTag("abc", "", false, false, func() string { return "" })
	assert.Equal(t, "Could not get url: 'abc'", err.Error())
}

func TestUrlDownload(t *testing.T) {
	table := []struct {
		o            string
		contentType  string
		expectedFile string
	}{
		{"", "image/png", "./img/generated.png"},
		{"", "image/jpg", "./img/generated.jpg"},
		{"./img/my", "image/png", "./img/my.png"},
		{"./img/s/my", "image/png", "./img/s/my.png"},
	}

	for i, tst := range table {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", tst.contentType)
			fmt.Fprintf(w, "dummy data")
		}))
		defer server.Close()
		defer os.RemoveAll("./img")

		tag, err := CreateImageTag(server.URL, tst.o, false, false, func() string { return "generated" })
		assert.Nil(t, err, "%d: err: %v", i, err)
		expectedTag := fmt.Sprintf("![%s](%s)", tst.expectedFile, tst.expectedFile)
		assert.Equal(t, expectedTag, tag, "%d: tag: ", i, tag)
		assert.FileExists(t, tst.expectedFile)
	}
}

func TestDownloadFailure(t *testing.T) {
	table := []struct {
		contentType    string
		responseStatus int
		expectedError  string
	}{
		{"text/html", 200, "Cannot accept content-type: text/html"},
		{"image", 200, "Cannot accept content-type: image"},
		{"image/png", 402, "Url request responded with: 402"},
		{"image/png", 500, "Url request responded with: 500"},
	}

	for i, tst := range table {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", tst.contentType)
			w.WriteHeader(tst.responseStatus)
			fmt.Fprintf(w, "dummy data")
		}))

		_, err := CreateImageTag(server.URL, "", false, false, func() string { return "" })
		assert.NotNil(t, err, "%d", i)
		assert.Equal(t, tst.expectedError, err.Error(), "%d: err: %v", i, err)
	}
}

func TestClipboardOption(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "image/png")
		fmt.Fprintf(w, "dummy data")
	}))
	defer server.Close()
	defer os.RemoveAll("./img")

	glippy.Set(server.URL)
	tag, err := CreateImageTag("", "", false, true, func() string { return "generated" })

	assert.Nil(t, err)
	assert.Equal(t, "![./img/generated.png](./img/generated.png)", tag)
	assert.FileExists(t, "./img/generated.png")
}

func TestClipboardOptionWithEmptyClipboard(t *testing.T) {
	glippy.Set("")
	tag, err := CreateImageTag("", "", false, true, func() string { return "generated" })

	assert.Equal(t, "Clipboard is empty", err.Error())
	assert.Equal(t, "", tag)
}

func TestRawClipboardImageOption(t *testing.T) {
	command := exec.Command("xclip", "-i", "-selection", "clipboard", "./test_resources/test.png")
	err := command.Run()
	assert.Nil(t, err)

	defer os.RemoveAll("./img")

	tag, err := CreateImageTag("", "", true, false, func() string { return "generated" })
	assert.Nil(t, err)
	assert.Equal(t, "![./img/generated.png](./img/generated.png)", tag)
	assert.FileExists(t, "./img/generated.png")

	file, err := os.Open("./img/generated.png")
	assert.Nil(t, err)
	info, err := file.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(119), info.Size())
}
