package charmutils

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

// CREDIT: https://gist.github.com/ras0q/9bf5d81544b22302393f61206892e2cd

// OverlayCenter writes the overlay string onto the background string such that the middle of the
// overlay string will be at the middle of the overlay will be at the middle of the background.
//
// See Overlay for more.
func OverlayCenter(bg string, overlay string, ignoreMarginWhitespace bool) (string, error) {
	row := (lipgloss.Height(bg) - lipgloss.Height(overlay)) / 2
	row = max(0, row)
	col := (lipgloss.Width(bg) - lipgloss.Width(overlay)) / 2
	col = max(0, col)
	return Overlay(bg, overlay, row, col, ignoreMarginWhitespace)
}

// Overlay writes the overlay string onto the background string at the specified row and column.
// Background ANSI escape sequences that would be overwritten are truncated on either side of
// the overlay lines.
// In this case, the row and column are zero indexed.
func Overlay(bg, overlay string, row, col int, ignoreMarginWhitespace bool) (string, error) {
	bgLines := strings.Split(bg, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for overlayLineIdx, overlayLine := range overlayLines {
		targetRow := row + overlayLineIdx

		// Ensure the target row exists in the background lines
		for len(bgLines) <= targetRow {
			bgLines = append(bgLines, "")
		}

		bgLine := bgLines[targetRow]
		bgLineWidth := ansi.StringWidth(bgLine)

		// Pad the background line with spaces so that it has the required columns
		if bgLineWidth < col {
			bgLine += strings.Repeat(" ", col-bgLineWidth)
			bgLines[targetRow] = bgLine
		}

		if ignoreMarginWhitespace {
			// Remove leading and trailing whitespace from the overlay line.
			overlayLine = removeMarginWhitespace(bgLine, overlayLine, col)
		}

		// Get the left part of the background line (up to the column)
		bgLeft := ansi.TruncateWc(bgLine, col, "")

		// Get the right part of the background line (from the column + overlay width)
		insertPoint := col + ansi.StringWidth(overlayLine)
		var bgRight string
		if insertPoint < ansi.StringWidth(bgLine) {
			bgRight = ansi.TruncateLeftWc(bgLine, insertPoint, "")
		}

		// Combine the left part, overlay, and right part
		bgLines[targetRow] = bgLeft + overlayLine + bgRight
	}

	return strings.Join(bgLines, "\n"), nil
}

// removeMarginWhitespace preserves the background where the overlay line has leading or trailing whitespace.
// This is done by detecting those empty cells in the overlay string and replacing them with the corresponding background cells.
//
//nolint:gocognit
func removeMarginWhitespace(bgLine, overlayLine string, col int) string {
	var result strings.Builder

	// Variables to track ANSI escape sequences
	inAnsi := false
	ansiSeq := strings.Builder{}

	// Strip ANSI codes to analyze whitespace
	overlayStripped := ansi.Strip(overlayLine)
	overlayRunes := []rune(overlayStripped)

	// Find first and last non-whitespace positions
	firstNonWhitespacePos := -1
	lastNonWhitespacePos := -1
	visualPos := 0
	overlayVisualWidths := make([]int, len(overlayRunes))

	for i, r := range overlayRunes {
		runeWidth := runewidth.RuneWidth(r)
		overlayVisualWidths[i] = runeWidth
		if !unicode.IsSpace(r) {
			if firstNonWhitespacePos == -1 {
				firstNonWhitespacePos = visualPos
			}
			lastNonWhitespacePos = visualPos + runeWidth - 1 // inclusive
		}
		visualPos += runeWidth
	}

	// If all characters are whitespace
	if firstNonWhitespacePos == -1 {
		firstNonWhitespacePos = 0
		lastNonWhitespacePos = -1
	}

	// Now, process the overlayLine, keeping track of visual positions
	visualPos = 0
	runeReader := strings.NewReader(overlayLine)

	for {
		r, _, err := runeReader.ReadRune()
		if err != nil {
			break
		}

		if r == '\x1b' {
			// Start of ANSI escape sequence
			inAnsi = true
			ansiSeq.WriteRune(r)
			continue
		}

		if inAnsi {
			ansiSeq.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				// End of ANSI escape sequence
				inAnsi = false
				result.WriteString(ansiSeq.String())
				ansiSeq.Reset()
			}
			continue
		}

		runeWidth := runewidth.RuneWidth(r)

		// Determine if current position is leading whitespace or trailing whitespace
		var isLeadingWhitespace, isTrailingWhitespace bool

		if visualPos < firstNonWhitespacePos {
			isLeadingWhitespace = true
		} else if visualPos > lastNonWhitespacePos {
			isTrailingWhitespace = true
		}

		if unicode.IsSpace(r) && (isLeadingWhitespace || isTrailingWhitespace) {
			// Preserve background character
			for k := range runeWidth {
				bgChar := getBgCharAt(bgLine, col+visualPos+k)
				result.WriteString(bgChar)
			}
		} else {
			// Include character from overlay (could be a non-whitespace or whitespace character in between)
			result.WriteRune(r)
		}

		visualPos += runeWidth
	}

	return result.String()
}

// getBgCharAt returns the character from the background line at the specified visual index.
func getBgCharAt(bgLine string, visualIndex int) string {
	var result strings.Builder
	displayWidth := 0
	inAnsi := false
	ansiSeq := strings.Builder{}

	runeReader := strings.NewReader(bgLine)
	for {
		r, _, err := runeReader.ReadRune()
		if err != nil {
			break
		}

		if r == '\x1b' {
			// Start of ANSI escape sequence
			inAnsi = true
			ansiSeq.WriteRune(r)
			continue
		}

		if inAnsi {
			ansiSeq.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				// End of ANSI escape sequence
				inAnsi = false
				result.WriteString(ansiSeq.String())
				ansiSeq.Reset()
			}
			continue
		}

		charWidth := runewidth.RuneWidth(r)
		if displayWidth+charWidth > visualIndex {
			// We have reached the desired index
			result.WriteRune(r)
			break
		}

		displayWidth += charWidth
	}

	// If no character found at the position, return a space
	if result.Len() == 0 {
		return " "
	}

	return result.String()
}
