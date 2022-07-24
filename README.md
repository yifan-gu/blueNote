blueNote
============
Organize reading notes and clippings


## Build

```
git clone git@github.com:yifan-gu/blueNote.git && go build
./blueNote -h
```

## Usage

#### Export notes to html with Kindle App
![Export Notes From Kindle App](screenshots/export-notes-from-kindle-app.png)

#### Convert notes to org-roam files
```
./blueNote -i kindlehtml -o org-roam examples/kindle_html_single_book_example.html ./
```

#### Add `-s` if the book is a collection of multiple books
```
./blueNote -i kindlehtml -o org-roam -s examples/kindle_html_collection_example.html ./
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

User Interaction:
- [ ] User ratings
- [ ] User comments?
- [ ] User profiles
- [ ] User upload audiobook readings?
