# mdimg - copy images into markdown
![build](https://github.com/danjhl/mdimg/actions/workflows/build.yml/badge.svg)

# Usage
```
mdimg -i [-o <output>]              create tag from image copy in clipboard
mdimg -c [-o <output>]              create tag from url in clipboard
mdimg -u <url> [-o <output>]        create tag from url
```

`-o` specifies the output file. By default a file with random name is
saved to `./img`.

`mdimg` will then output the markdown image tag for the saved file.
