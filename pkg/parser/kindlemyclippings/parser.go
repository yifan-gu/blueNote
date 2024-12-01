/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package kindlemyclippings

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/model"
)

type KindleMyClippingsParser struct {
	minSimilarity float64 // Threshold for deduplication (0 to 1)
}

func (p *KindleMyClippingsParser) Name() string { return "kindle-my-clippings" }

// LoadConfigs sets up the flags for the parser.
func (p *KindleMyClippingsParser) LoadConfigs(cmd *cobra.Command) {
	cmd.Flags().Float64Var(&p.minSimilarity, "min-similarity", 0.8, "Minimum similarity percentage (0-1) to consider a highlight as duplicate")
}

// Parse processes the MyClippings.txt file and returns a list of Books with deduplicated Marks.
func (p *KindleMyClippingsParser) Parse(inputPath string) ([]*model.Book, error) {
	file, err := os.Open(inputPath)
	defer file.Close()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	var books []*model.Book
	markListMap := make(map[string][]*model.Mark) // Maps book title -> list of marks

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Extract book details
		title, author, err := extractTitleAndAuthor(stripLeadingBOM(scanner.Text()))
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		// Parse metadata
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF, expecting metadata")
		}
		meta := scanner.Text()
		markType, location, createdAt, err := extractMeta(meta)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		scanner.Scan() // Skip the empty line before the actual text.

		var text []string
		for scanner.Scan() {
			line := scanner.Text()
			if line == "==========" {
				break
			}
			text = append(text, line)
		}

		// Parse data or note
		if text == nil { // Empty notes, skip
			return nil, nil
		}

		// Create the mark
		mark := &model.Mark{
			Type:     markType,
			Title:    title,
			Author:   author,
			Location: location,
			Data:     strings.Join(text, "\n"),
			CreatedAt: func() *int64 {
				if createdAt != 0 {
					return &createdAt
				}
				return nil
			}(),
		}

		// Validate the mark
		if err := model.ValidateMark(mark); err != nil {
			continue
		}

		// Add the mark to the book's list of marks
		if _, exists := markListMap[title]; !exists {
			markListMap[title] = []*model.Mark{}
		}
		markListMap[title] = append(markListMap[title], mark)
	}

	// Deduplicate marks for each book
	for title, marks := range markListMap {
		deduplicatedMarks := deduplicateMarks(marks, p.minSimilarity) // Use the configured threshold
		book := &model.Book{
			Title:  title,
			Author: marks[0].Author,
			Marks:  deduplicatedMarks,
		}
		books = append(books, book)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "")
	}

	model.SortBooksByTitle(books)

	return books, nil
}

func stripLeadingBOM(s string) string {
	runes := []rune(s)
	start := 0
	for start < len(runes) && runes[start] == '\uFEFF' {
		start++
	}
	return string(runes[start:])
}

// Deduplicate marks by comparing similarity of content using a minimum common words threshold
func deduplicateMarks(marks []*model.Mark, minSimilarity float64) []*model.Mark {
	var deduplicated []*model.Mark

	for _, mark := range marks {
		duplicateFound := false
		for i, dedupMark := range deduplicated {
			if hasMinCommonWords(dedupMark.Data, mark.Data, minSimilarity) {
				// Keep the more recent mark
				if mark.CreatedAt != nil && (dedupMark.CreatedAt == nil || *mark.CreatedAt > *dedupMark.CreatedAt) {
					deduplicated[i] = mark
				}
				duplicateFound = true
				break
			}
		}
		if !duplicateFound {
			deduplicated = append(deduplicated, mark)
		}
	}

	return deduplicated
}

// Check if two strings share a longest common substring whose percentage of similarity meets the threshold
func hasMinCommonWords(a, b string, minSimilarity float64) bool {
	// Preprocess both strings to handle Chinese text
	a = preprocessString(a)
	b = preprocessString(b)

	// Convert strings to runes for character-level comparison
	aRunes := []rune(a)
	bRunes := []rune(b)

	// Find the length of the longest common substring
	longestCommonLength := findLongestCommonSubstring(aRunes, bRunes)

	// Calculate the percentage of the longest common substring
	minLength := len(aRunes)
	if len(bRunes) < minLength {
		minLength = len(bRunes)
	}
	percentage := float64(longestCommonLength) / float64(minLength)

	// Check if the percentage meets the threshold
	return percentage >= minSimilarity
}

// Find the length of the longest common substring between two rune slices
func findLongestCommonSubstring(a, b []rune) int {
	// Initialize a 2D slice for storing lengths of common substrings
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
	}

	// Track the maximum length of common substring
	maxLength := 0

	// Fill the dp table
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLength {
					maxLength = dp[i][j]

				}
			}
		}
	}

	return maxLength
}

// Preprocess a string by adding spaces between Chinese characters
func preprocessString(s string) string {
	var result strings.Builder
	for _, r := range s {
		if isChinese(r) {
			// Add space before and after Chinese character
			result.WriteRune(' ')
			result.WriteRune(r)
			result.WriteRune(' ')
		} else {
			// Append non-Chinese character as is
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}

// Check if a rune is a Chinese character
func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

// Helper functions to extract title, author, metadata, and timestamps

func extractTitleAndAuthor(line string) (string, string, error) {
	parts := strings.Split(line, "(")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid title and author format: %s", line)
	}
	title := strings.TrimSpace(parts[0])
	author := strings.TrimSuffix(strings.TrimSpace(parts[len(parts)-1]), ")")
	if title == "" || author == "" {
		return "", "", fmt.Errorf("empty title or author: %s", line)
	}
	return title, author, nil
}

func extractMeta(meta string) (string, *model.Location, int64, error) {
	// Check for "Your Highlight" or "Your Note" and trim it
	var markType string
	switch {
	case strings.HasPrefix(meta, "- Your Highlight"):
		markType = model.MarkTypeHighlight
		meta = strings.TrimPrefix(meta, "- Your Highlight on ")
	case strings.HasPrefix(meta, "- Your Note"):
		markType = model.MarkTypeNote
		meta = strings.TrimPrefix(meta, "- Your Note on ")
	case strings.HasPrefix(meta, "- Your Bookmark"):
		markType = model.MarkTypeBookmark
		meta = strings.TrimPrefix(meta, "- Your Bookmark on ")
	default:
		return "", nil, 0, fmt.Errorf("unsupported mark type: %s", meta)
	}

	// Split the remaining metadata into parts using "|"
	parts := strings.Split(meta, "|")
	if len(parts) != 2 && len(parts) != 3 {
		return "", nil, 0, fmt.Errorf("invalid metadata format: %s", meta)
	}

	// Parse location
	location := &model.Location{}
	if len(parts) == 3 {
		// Case: "page" and "Location" both present
		if err := parsePageAndLocation(parts[0], parts[1], location); err != nil {
			return "", nil, 0, err
		}
	} else if len(parts) == 2 {
		// Case: Only "Location" present
		if err := parsePageAndLocation("", parts[0], location); err != nil {
			return "", nil, 0, err
		}
	}

	// Parse timestamp (always in the last part)
	createdAt := parseTimestamp(parts[len(parts)-1])
	if createdAt == 0 {
		return "", nil, 0, fmt.Errorf("invalid timestamp: %s", parts[len(parts)-1])
	}

	return markType, location, createdAt, nil
}

func parsePageAndLocation(pagePart, locationPart string, loc *model.Location) error {
	// Parse "page" field if present
	if strings.Contains(pagePart, "page") {
		pageParts := strings.Split(pagePart, "page")
		if len(pageParts) > 1 {
			page := strings.TrimSpace(pageParts[1])
			loc.Page = parsePageString(page)
			if loc.Page == nil {
				return fmt.Errorf("invalid page: %s", page)
			}
		}
	}

	// Parse "Location" field if present
	if strings.Contains(locationPart, "Location") {
		locationParts := strings.Split(locationPart, "Location")
		if len(locationParts) > 1 {
			loc.Location = parseLocationRangeString(locationParts[1])
			if loc.Location == nil {
				return fmt.Errorf("invalid location: %s", locationPart)
			}
		}
	}

	// Ensure at least one valid field is present
	if loc.Page == nil && loc.Location == nil {
		return fmt.Errorf("no valid page or location")
	}

	return nil
}

func parseLocationRangeString(locationRange string) *int {
	// Split the range (e.g., "541-543" or "541")
	parts := strings.Split(locationRange, "-")
	if len(parts) == 0 {
		return nil
	}

	// Parse the start of the range
	startLocation := 0
	_, err := fmt.Sscanf(parts[0], "%d", &startLocation)
	if err == nil {
		return &startLocation
	}

	return nil
}

func parsePageString(pageStr string) *int {
	parsedPage := 0
	_, err := fmt.Sscanf(pageStr, "%d", &parsedPage)
	if err == nil {
		return &parsedPage
	}
	return nil
}

func parseTimestamp(timestampStr string) int64 {
	// Normalize the string by trimming leading/trailing spaces
	timestampStr = strings.TrimSpace(timestampStr)

	// Strip the "Added on " prefix if present
	const prefix = "Added on "
	if strings.HasPrefix(timestampStr, prefix) {
		timestampStr = strings.TrimPrefix(timestampStr, prefix)
	}

	// Define the layout for parsing
	layout := "Monday, January 2, 2006 3:04:05 PM"
	t, err := time.Parse(layout, strings.TrimSpace(timestampStr))
	if err != nil {
		return 0 // Return 0 for invalid timestamps
	}
	return t.Unix()
}
