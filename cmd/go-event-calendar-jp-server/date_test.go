package main

import (
	"testing"
)

func TestGetMonthsRange(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month int
		want  []string
	}{
		{
			name:  "middle of year",
			year:  2024,
			month: 6,
			want: []string{
				"202403",
				"202404",
				"202405",
				"202406",
				"202407",
				"202408",
				"202409",
			},
		},
		{
			name:  "start of year",
			year:  2024,
			month: 1,
			want: []string{
				"202310",
				"202311",
				"202312",
				"202401",
				"202402",
				"202403",
				"202404",
			},
		},
		{
			name:  "end of year",
			year:  2024,
			month: 12,
			want: []string{
				"202409",
				"202410",
				"202411",
				"202412",
				"202501",
				"202502",
				"202503",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMonthsRange(tt.year, tt.month)
			if len(got) != len(tt.want) {
				t.Errorf("GetMonthsRange() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetMonthsRange()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
