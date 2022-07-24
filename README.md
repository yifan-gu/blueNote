BlueNote
============
Organize reading notes and clippings


## Usage

```
git clone git@github.com:yifan-gu/BlueNote.git && go build -o bluenote
./bluenote -h
```

## Example

### Generate org-roam files to the local folder from a kindle html clipping file
```
./bluenote -i kindlehtml -o org-roam -s examples/kindle_html_example.html ./
```

## TODO

Emacs org roam related:
- [ ] Roam module (fix bug)
- [ ] Check roam version

Parser:
- [x] Refactor book module
- [ ] Refactor configs for parser and exporter
- [x] Fix location output
- [ ] Fix user notes content
- [ ] Marshal and unmarshal json
- [ ] Stacktrace error
- [ ] Optional author, title flag

Web/App:
- [ ] Database storage
- [ ] Search by tags, keywords, book, author
- [ ] Show random notes/highlights every time
- [ ] Show connected notes
- [ ] Clickable with tags, book, author
- [ ] Manual update tags
- [ ] Ratings?
- [ ] Audio book
- [ ] Generate tags automatically
- [ ] Connected notes suggestion

Social:
- [ ] User ratings
- [ ] User comments
- [ ] User profiles
- [ ] User upload audiobook readings?
