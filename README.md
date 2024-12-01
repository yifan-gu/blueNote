# blueNote ![GitHub Actions](https://github.com/yifan-gu/blueNote/actions/workflows/go.yml/badge.svg)
A Notes/Clippings Browser

---

## Build
```
git clone git@github.com:yifan-gu/blueNote.git && go build
./blueNote -h
```


---

## Notes
This tool supports both **HTML highlights** exported via the Kindle App and the `My Clippings.txt` file directly from a Kindle device.

- **My Clippings.txt**:
  - Pros: Contains all highlights and notes.
  - Cons: Lacks chapter information.

- **HTML Highlights**:
  - Pros: Better formatting and includes chapter information.
  - Cons: You can only generate highlights/notes one book at a time.

---

## Usage

### Convert `My Clippings.txt` to JSON and display it in the console
``` 
./blueNote convert -i kindle-my-clippings -o json --json.pretty examples/My\ Clippings.txt
```

### Export notes as html using the Kindle App
![Export Notes From Kindle App](screenshots/export-notes-from-kindle-app.png)


### Convert notes to JSON and display them in the console
```
./blueNote convert -i kindle-html -o json --json.pretty examples/kindle_html_single_book_example.html
```

### Convert notes and store them into MongoDB
```
./blueNote convert -i kindle-html -o mongodb examples/kindle_html_single_book_example.html
```

### Convert notes to org-roam files and save to the current dir
```
./blueNote convert -i kindle-html -o org-roam examples/kindle_html_single_book_example.html ./
```

### Add `-s` if the book is a collection of multiple books
```
./blueNote convert -i kindle-html -o json --json.pretty -s examples/kindle_html_collection_example.html
```


<!--### Browse and edit the notes with tags in Emacs Org
![View and Edit Notes in Emacs Org-roam](screenshots/view-notes-with-emacs-org-roam.png)

### Sync the org-roam database
Remember to run `M-x org-roam-db-sync` to sync the org-roam database.
![Run org-roam-db-sync](screenshots/org-roam-db-sync.png)

### 📖 Happy Notes Searching! 📖
![Search for Notes in Emacs Org-roam](screenshots/search-keywords-with-emacs-org-roam.png)

---

## References

- [Doom Emacs](https://github.com/doomemacs/doomemacs)
- [Org-roam](https://www.orgroam.com/)
- [My Doom Emacs Config](https://github.com/yifan-gu/.doom)
- A [custom Emacs theme](https://github.com/yifan-gu/.doom/blob/master/themes/org-leuven-theme.el) for Org-roam mode, based on [Leuven](https://github.com/fniessen/emacs-leuven-theme). -->

---

## TODO

### Documentation
- [ ] Update README to reflect JSON/MongoDB-based usage.

### Emacs Org-roam
- [ ] <s>[Dropped] Roam module (fix bug).</s>
- [ ] <s>[Dropped] Check roam version.</s>

### Parser/Exporter
- [x] Refactor book module.
- [x] Refactor configs for parser and exporter.
- [x] Fix location output.
- [x] Fix user notes content.
- [x] JSON exporter.
- [x] Optional author/title flag.
- [x] Stacktrace error handling.
- [x] List parsers and exporters.
- [x] MongoDB exporter.
- [x] JSON parser.
- [ ] One-click export from Kindle app.
- [ ] Change parser/exporter type from string to safe type.
- [x] `My Clippings.txt` parser.
- [ ] Diff the previous processed `My Clippings.txt`
- [ ] Support multiple authors.

### Server Backend
- [x] Database storage.
- [ ] <s>Index, unique on digest.</s>
- [x] Update storage interface.
- [ ] <s>Recompute digest.</s>
- [x] Limit on returned marks.
- [x] Add timestamps (created, last modified).
- [ ] Server REST API.
- [x] GraphQL API (READ).
- [x] GraphQL API (CREATE).
- [x] GraphQL API (UPDATE).
- [x] GraphQL API (DELETE).
- [ ] Handle GraphQL null fields.
- [ ] GraphQL API tests with mocked storage.

### Application
- [ ] Create database schema for users.
- [ ] Create database schema for books.
- [ ] Create database schema for mark interactions (user likes, comments, shares).
- [ ] Add search by tags, keywords, book, author.
- [ ] Show random notes/highlights every time.
- [ ] Display connected notes.
- [ ] Add clickable tags, books, authors.
- [ ] Manual tag updates.
- [ ] Ratings system.
- [ ] Support audiobooks.
- [ ] Generate tags automatically.
- [ ] Suggest connected notes.

### Advanced Features
- [ ] User ratings.
- [ ] User comments.
- [ ] User profiles.
- [ ] User-uploaded audiobook readings.
