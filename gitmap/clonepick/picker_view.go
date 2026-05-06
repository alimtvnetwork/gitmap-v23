package clonepick

// picker_view.go: View() implementation for the bubbletea picker.
// Kept in its own file so picker.go can stay under 200 lines and the
// rendering loop has room to breathe.

import (
	"fmt"
	"strings"
)

// View renders the picker as a single string per tea-model contract.
// Layout:
//
//	gitmap clone-pick --ask  (12/487 selected)
//	[x] docs/                       <- cursor row, bracketed
//	[ ] src/cmd/
//	[-] node_modules/    (auto-greyed)
//	...
//	space toggle | a all | n none | s save | q quit
func (m pickerModel) View() string {
	if len(m.paths) == 0 {
		return "clone-pick: repository has no tracked files\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "gitmap clone-pick --ask  (%d/%d selected)\n",
		m.countPicked(), len(m.paths))
	m.renderRows(&b)
	b.WriteString("\nspace toggle | a all | n none | s save | q quit\n")

	return b.String()
}

// renderRows writes one line per visible path. Visible window is the
// full slice in v1 (most repos fit in a normal terminal); a windowed
// scroller is a follow-up if real-world repos blow past the screen.
func (m pickerModel) renderRows(b *strings.Builder) {
	for i, path := range m.paths {
		b.WriteString(formatRow(i == m.cursor, m.picked[i],
			IsAutoExcluded(path), path))
		b.WriteByte('\n')
	}
}

// formatRow returns the single-line representation of one picker
// entry. Cursor row gets a leading ">", everything else gets two
// spaces so columns line up.
func formatRow(isCursor, isPicked, isGreyed bool, path string) string {
	prefix := "  "
	if isCursor {
		prefix = "> "
	}
	mark := pickMark(isPicked, isGreyed)
	suffix := ""
	if isGreyed {
		suffix = "  (auto-greyed)"
	}

	return prefix + mark + " " + path + suffix
}

// pickMark returns the bracketed checkbox glyph for the row state.
// Greyed rows use "[-]" so the user can still see they're toggleable
// (versus the disabled-looking "[ ]").
func pickMark(isPicked, isGreyed bool) string {
	switch {
	case isPicked:
		return "[x]"
	case isGreyed:
		return "[-]"
	default:
		return "[ ]"
	}
}

// countPicked is the header counter. O(rows) once per render, fine
// for the row counts we expect (<10k entries before the windowed
// scroller lands).
func (m pickerModel) countPicked() int {
	n := 0
	for _, on := range m.picked {
		if on {
			n++
		}
	}

	return n
}
