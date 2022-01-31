package parser

import (
	"fmt"
	"syscall"
)

// timespecToEmacsTimestamp converts a Timespec{Sec, Nsec} to
// Emacs timestamp{high, low, micro, pico}.
// For more info, see https://www.gnu.org/software/emacs/manual/html_node/elisp/Time-of-Day.html
func timespecToEmacsTimestamp(ts syscall.Timespec) []int64 {
	return []int64{
		ts.Sec / 65536,
		ts.Sec % 65536,
		ts.Nsec / 1000,
		(ts.Nsec % 1000) * 1000,
	}
}

func getEmacsTimestamp(path, typ string) (string, error) {
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return "", err
	}
	var est []int64
	switch typ {
	case "atime":
		est = timespecToEmacsTimestamp(st.Atimespec)
	case "mtime":
		est = timespecToEmacsTimestamp(st.Mtimespec)
	default:
		return "", fmt.Errorf("unexpected timestamp type: %q", typ)
	}
	return fmt.Sprintf("(%d %d %d %d)", est[0], est[1], est[2], est[3]), nil
}

func getAtime(path string) (string, error) {
	return getEmacsTimestamp(path, "atime")
}

func getMtime(path string) (string, error) {
	return getEmacsTimestamp(path, "mtime")
}
