#!/usr/bin/env bash

set -eu

ROOT_DIR="$(git rev-parse --show-toplevel)"

echo "Test converting a single book to json"
output="$(go run ./... convert -i kindle-html --json.pretty -ojson \
    "${ROOT_DIR}/examples/kindle_html_single_book_example.html")"
echo "${output}" | diff "${ROOT_DIR}/tests/single_book_output.json" -

echo "Test converting a single book to json with title and author override "
output="$(go run ./... convert -i kindle-html \
    --json.pretty --kindle-html.author="Ernest Hemingway" \
    --kindle-html.title="The Sun Also Rises" \
    -ojson \
    "${ROOT_DIR}/examples/kindle_html_single_book_example.html")"
echo "${output}" | diff "${ROOT_DIR}/tests/single_book_output_author_title.json" -

echo "Test converting a book collection to json "
output="$(go run ./...  convert -i kindle-html --json.pretty -s -ojson \
    "${ROOT_DIR}/examples/kindle_html_collection_example.html")"
echo "${output}" | diff "${ROOT_DIR}/tests/book_collection_split.json" -

echo "Passed!"
