package charmutils_test

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Broderick-Westrope/charmutils"
)

func TestOverlayCenter(t *testing.T) {
	tt := map[string]struct {
		bg                     string
		overlay                string
		ignoreMarginWhitespace bool
		want                   string
	}{
		"simple": {
			bg: heredoc.Doc(`
Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic.
Amet quidem reprehenderit ex. Error illum sit est expedita sapiente neque. Laborum vero necessitatibus similique suscipit nam.
Tempore occaecati eligendi accusamus eos similique harum impedit. Quas nam molestiae architecto quam.
Accusamus pariatur facilis ea nostrum exercitationem quam. Sit ipsam aperiam aspernatur hic fugit officia inventore.
Reiciendis doloribus ut eius id. Repellendus eum enim. Reprehenderit veritatis nulla molestiae nulla veniam.
Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.`),
			overlay:                "*********\n*****",
			ignoreMarginWhitespace: false,
			want: heredoc.Doc(`
Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic.
Amet quidem reprehenderit ex. Error illum sit est expedita sapiente neque. Laborum vero necessitatibus similique suscipit nam.
Tempore occaecati eligendi accusamus eos similique harum i*********uas nam molestiae architecto quam.
Accusamus pariatur facilis ea nostrum exercitationem quam.*****ipsam aperiam aspernatur hic fugit officia inventore.
Reiciendis doloribus ut eius id. Repellendus eum enim. Reprehenderit veritatis nulla molestiae nulla veniam.
Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.`),
		},
		"padded; enforce margins": {
			bg: heredoc.Doc(`
Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic.
Amet quidem reprehenderit ex. Error illum sit est expedita sapiente neque. Laborum vero necessitatibus similique suscipit nam.
Tempore occaecati eligendi accusamus eos similique harum impedit. Quas nam molestiae architecto quam.
Accusamus pariatur facilis ea nostrum exercitationem quam. Sit ipsam aperiam aspernatur hic fugit officia inventore.
Reiciendis doloribus ut eius id. Repellendus eum enim. Reprehenderit veritatis nulla molestiae nulla veniam.
Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.`),
			overlay:                lipgloss.NewStyle().Padding(1, 3).Render("*********\n*****"),
			ignoreMarginWhitespace: false,
			want: heredoc.Doc(`
Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic.
Amet quidem reprehenderit ex. Error illum sit est exped               que. Laborum vero necessitatibus similique suscipit nam.
Tempore occaecati eligendi accusamus eos similique haru   *********    nam molestiae architecto quam.
Accusamus pariatur facilis ea nostrum exercitationem qu   *****       periam aspernatur hic fugit officia inventore.
Reiciendis doloribus ut eius id. Repellendus eum enim.                eritatis nulla molestiae nulla veniam.
Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.`),
		},
		"padded; ignore margins": {
			bg: heredoc.Doc(`
Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic.
Amet quidem reprehenderit ex. Error illum sit est expedita sapiente neque. Laborum vero necessitatibus similique suscipit nam.
Tempore occaecati eligendi accusamus eos similique harum impedit. Quas nam molestiae architecto quam.
Accusamus pariatur facilis ea nostrum exercitationem quam. Sit ipsam aperiam aspernatur hic fugit officia inventore.
Reiciendis doloribus ut eius id. Repellendus eum enim. Reprehenderit veritatis nulla molestiae nulla veniam.
Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.`),
			overlay:                lipgloss.NewStyle().Padding(1, 3).Render("*********\n*****"),
			ignoreMarginWhitespace: true,
			want: heredoc.Doc(`
Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic.
Amet quidem reprehenderit ex. Error illum sit est expedita sapiente neque. Laborum vero necessitatibus similique suscipit nam.
Tempore occaecati eligendi accusamus eos similique harum i*********uas nam molestiae architecto quam.
Accusamus pariatur facilis ea nostrum exercitationem quam.*****ipsam aperiam aspernatur hic fugit officia inventore.
Reiciendis doloribus ut eius id. Repellendus eum enim. Reprehenderit veritatis nulla molestiae nulla veniam.
Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.`),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result, err := charmutils.OverlayCenter(tc.bg, tc.overlay, tc.ignoreMarginWhitespace)
			require.NoError(t, err)
			assert.Equal(t, tc.want, result)
		})
	}
}

func TestOverlay(t *testing.T) {
	tt := map[string]struct {
		bg                     string
		overlay                string
		row                    int
		col                    int
		ignoreMarginWhitespace bool
		want                   string
	}{
		"single line; start": {
			bg:                     "Nostrum libero modi velit neque dolores.",
			overlay:                "*********",
			row:                    0,
			col:                    0,
			ignoreMarginWhitespace: false,
			want:                   "*********ibero modi velit neque dolores.",
		},
		"single line; middle": {
			bg:                     "Nostrum libero modi velit neque dolores.",
			overlay:                "*********",
			row:                    0,
			col:                    10,
			ignoreMarginWhitespace: false,
			want:                   "Nostrum li********* velit neque dolores.",
		},
		"single line; beyond final column": {
			bg:                     "Nostrum libero modi velit neque dolores.",
			overlay:                "*********",
			row:                    0,
			col:                    35,
			ignoreMarginWhitespace: false,
			want:                   "Nostrum libero modi velit neque dol*********",
		},
		"single line; beyond final row": {
			bg:                     "Nostrum libero modi velit neque dolores.",
			overlay:                "*********",
			row:                    3,
			col:                    0,
			ignoreMarginWhitespace: false,
			want:                   "Nostrum libero modi velit neque dolores.\n\n\n*********",
		},
		"single line; lipgloss styled": {
			bg:                     "Nostrum libero modi velit neque dolores.",
			overlay:                lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("1")).Render("*****"),
			row:                    0,
			col:                    5,
			ignoreMarginWhitespace: false,
			want:                   "Nostr*****bero modi velit neque dolores.",
		},
		"single line; manual escape code": {
			bg:                     "Nostrum libero modi velit neque dolores.",
			overlay:                "\x1b[31m*****\x1b[0m",
			row:                    0,
			col:                    5,
			ignoreMarginWhitespace: false,
			want:                   "Nostr\u001B[31m*****\u001B[0mbero modi velit neque dolores.",
		},
		"multi-line background; overlay middle line": {
			bg:                     "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
			overlay:                "*****",
			row:                    2,
			col:                    0,
			ignoreMarginWhitespace: false,
			want:                   "Line 1\nLine 2\n*****3\nLine 4\nLine 5",
		},
		"multi-line overlay; beyond background": {
			bg:                     "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
			overlay:                "*******\n*******",
			row:                    1,
			col:                    5,
			ignoreMarginWhitespace: false,
			want:                   "Line 1\nLine *******\nLine *******\nLine 4\nLine 5",
		},
		"multi-line overlay; enforce margins": {
			bg:                     "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
			overlay:                lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1).Render("***\n***"),
			row:                    0,
			col:                    0,
			ignoreMarginWhitespace: false,
			want:                   "     1\n  ***2\n  ***3\nLine 4\nLine 5",
		},
		"multi-line overlay; ignore margins": {
			bg:                     "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
			overlay:                lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1).Render("***\n***"),
			row:                    0,
			col:                    0,
			ignoreMarginWhitespace: true,
			want:                   "Line 1\nLi***2\nLi***3\nLine 4\nLine 5",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result, err := charmutils.Overlay(tc.bg, tc.overlay, tc.row, tc.col, tc.ignoreMarginWhitespace)
			require.NoError(t, err)
			assert.Equal(t, tc.want, result)
		})
	}
}
