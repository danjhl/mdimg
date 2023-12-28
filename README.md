# mdimg - Quickly add an image tag to your markdown for a copied url

# Usage
```
mdimg -u <url> [-o <output>]
mdimg -c [-o <output>]
```

`-u` specifies the url of the image that should be added.
`-o` specifies the output file. By default a file with random name is
saved to `./img`.
`-c` will use an url from the clipboard if present.

`mdimg` will then output the markdown image tag for the downloaded file.
