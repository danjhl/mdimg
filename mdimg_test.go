package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlDownload(t *testing.T) {
	table := []struct {
		out          string
		contentType  string
		expectedTag  string
		expectedFile string
		expectedErr  string
	}{
		{"", "image/png", "![./img/generated-name.png](./img/generated-name.png)", "./img/generated-name.png", ""},
		{"./img/s/myimg", "image/png", "![./img/s/myimg.png](./img/s/myimg.png)", "./img/s/myimg.png", ""},
		{"", "text/html", "![./img/s/myimg.png](./img/s/myimg.png)", "./img/s/myimg.png", "Cannot accept content-type: text/html"},
		{"./img/myimg.png", "image/png", "![./img/myimg.png](./img/myimg.png)", "./img/myimg.png", ""},
	}

	for i, tst := range table {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", tst.contentType)
			fmt.Fprintf(w, "dummy data")
		}))
		defer server.Close()
		defer os.RemoveAll("./img")

		tag, err := CreateImageTag(server.URL, tst.out, func() string { return "generated-name" })

		if err != nil && err.Error() == tst.expectedErr {
			continue
		}

		if tst.expectedErr == "" && err != nil || err != nil && err.Error() != tst.expectedErr {
			t.Errorf("%d: expectedFile: %s expected err: %s, actual: %v", i, tst.expectedFile, tst.expectedErr, err)
			continue
		}

		if tag != tst.expectedTag {
			t.Errorf("%d: expected tag: %s, actual tag: %s", i, tst.expectedTag, tag)
		}

		_, err = os.Stat(tst.expectedFile)

		if err != nil {
			t.Errorf("%d: expectedFile: %s, %v", i, tst.expectedFile, err)
		}
	}
}

func TestMissingUrlOption(t *testing.T) {
	_, err := CreateImageTag("", "", func() string { return "generated-name" })
	assert.Equal(t, "Requires -u option", err.Error())
}
