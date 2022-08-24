blueNote ![GitHub Actions](https://github.com/yifan-gu/blueNote/actions/workflows/go.yml/badge.svg)
============
A Notes/Clippings Browser


## Build

```
git clone git@github.com:yifan-gu/blueNote.git && go build
./blueNote -h
```

## Usage

#### Export notes as html using the Kindle App
![Export Notes From Kindle App](screenshots/export-notes-from-kindle-app.png)

#### Convert notes to org-roam files
```
./blueNote convert -i kindlehtml -o org-roam examples/kindle_html_single_book_example.html ./
```

#### Add `-s` if the book is a collection of multiple books
```
./blueNote convert -i kindlehtml -o org-roam -s examples/kindle_html_collection_example.html ./
```

#### Browse and edit the notes with tags
![View and Edit Notes in Emacs Org-roam](screenshots/view-notes-with-emacs-org-roam.png)

#### Remember to run `M-x org-roam-db-sync` to sync the org-roam database
![Run org-roam-db-sync](screenshots/org-roam-db-sync.png)

#### 📖 Happy Notes Searching! 📖
![Search for Notes in Emacs Org-roam](screenshots/search-keywords-with-emacs-org-roam.png)


### References

- Doom Emacs: https://github.com/doomemacs/doomemacs
- Org-roam: https://www.orgroam.com/
- My Doom Emacs config: https://github.com/yifan-gu/.doom
- A [custom Emacs theme](https://github.com/yifan-gu/.doom/blob/master/themes/org-leuven-theme.el) I made for Org-roam mode based on [leuven](https://github.com/fniessen/emacs-leuven-theme)


## TODO

Emacs org roam related:
- [ ] <s>[Dropped] Roam module (fix bug)</s>
- [ ] <s>[Dropped] Check roam version</s>

Parser/Exporter:
- [x] Refactor book module
- [x] Refactor configs for parser and exporter
- [x] Fix location output
- [x] Fix user notes content
- [x] JSON exporter
- [x] Optional author, title flag
- [x] Stacktrace error
- [x] List parsers and exporters
- [x] MongoDB exporter
- [x] JSON parser

Server Backend:
- [x] Database storage
- [ ] Index, unique on digest
- [ ] Limit on return marks
- [ ] Timestamp(created, updated)
- [ ] Server REST API
- [ ] Graphql API?

App:
- [ ] Search by tags, keywords, book, author
- [ ] Show random notes/highlights every time
- [ ] Show connected notes
- [ ] Clickable with tags, book, author
- [ ] Manual update tags
- [ ] Ratings?
- [ ] Audio book
- [ ] Generate tags automatically
- [ ] Connected notes suggestion

Advanced functions:
- [ ] User ratings
- [ ] User comments?
- [ ] User profiles
- [ ] User upload audiobook readings?
