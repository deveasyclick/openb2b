package parseuint

import (
	"fmt"
	"testing"
)

func TestParseUint(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		field     string
		want      uint
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid number",
			input:     "123",
			field:     "age",
			want:      123,
			expectErr: false,
		},
		{
			name:      "zero value",
			input:     "0",
			field:     "count",
			want:      0,
			expectErr: false,
		},
		{
			name:      "large valid number",
			input:     fmt.Sprintf("%d", uint64(^uint(0))), // max uint for the platform
			field:     "limit",
			want:      ^uint(0),
			expectErr: false,
		},
		{
			name:      "negative number",
			input:     "-5",
			field:     "balance",
			want:      0,
			expectErr: true,
			errMsg:    `error parsing balance: "-5": strconv.ParseUint: parsing "-5": invalid syntax`,
		},
		{
			name:      "non-numeric",
			input:     "abc",
			field:     "id",
			want:      0,
			expectErr: true,
			errMsg:    `error parsing id: "abc": strconv.ParseUint: parsing "abc": invalid syntax`,
		},
		{
			name:      "value too large for uint",
			input:     "18446744073709551616", // > max uint64
			field:     "count",
			want:      0,
			expectErr: true,
			errMsg:    `error parsing count: "18446744073709551616": strconv.ParseUint: parsing "18446744073709551616": value out of range`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint(tt.input, tt.field)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("unexpected error message:\n got:  %q\n want: %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("got %d, want %d", got, tt.want)
				}
			}
		})
	}
}
