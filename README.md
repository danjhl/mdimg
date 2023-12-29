# mdimg - copy images into markdown

# Usage
```
mdimg -i [-o <output>]
mdimg -c [-o <output>]
mdimg -u <url> [-o <output>]
```

`-i` specifies that an image copied to the clipboard should be used.

`-u` specifies the url of the image that should be added.

`-o` specifies the output file. By default a file with random name is
saved to `./img`.

`-c` will use an url from the clipboard if present.

`mdimg` will then output the markdown image tag for the downloaded file.

# Dependencies
`xclip` on linux.
