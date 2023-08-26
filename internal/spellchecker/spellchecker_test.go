package spellchecker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "correct text",
			text: "correct text",
			want: "correct text",
		},
		{
			name: "correct text with punctuation",
			text: "correct, text.",
			want: "correct, text.",
		},
		{
			name: "correct text (single) with undercase",
			text: "correct_text",
			want: "correct_text",
		},
		{
			name: "incorrect text with punctuation",
			text: "corrct,..! правилно!!!!",
			want: "correct,..! правильно!!!!",
		},
		{
			name: "incorrect text with missing symbol and undercase",
			text: "правльный_текст",
			want: "правильный_текст",
		},
		{
			name: "incorrect text with missing symbol",
			text: "правльный",
			want: "правильный",
		},
		{
			name: "incorrect text with two wrong words and undercase",
			text: "правльный_ткст",
			want: "правильный_текст",
		},
		{
			name: "empty",
			text: "",
			want: "",
		},
		{
			name: "correct text with digits (without space)",
			text: "correct9",
			want: "correct 9",
		},
		{
			name: "mixed incorrect text",
			text: "corrct правилно",
			want: "correct правильно",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixedText, err := Handle(tt.text)
			require.NoError(t, err)

			assert.Equal(t, tt.want, fixedText)
		})
	}
}
