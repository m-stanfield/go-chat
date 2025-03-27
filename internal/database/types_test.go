package database_test

import (
	"testing"

	"go-chat-react/internal/database"
)

func TestParseIntToID(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      int
		want    database.Id
		wantErr bool
	}{
		{
			name:    "Positive ID",
			id:      1,
			want:    1,
			wantErr: false,
		},
		{
			name:    "Zero ID",
			id:      0,
			want:    0,
			wantErr: true,
		},
		{
			name:    "Negative ID",
			id:      -1,
			want:    0,
			wantErr: true,
		},
		{
			name:    "Positive ID",
			id:      100,
			want:    100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.ParseIntToID(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseIntToID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseIntToID() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("ParseIntToID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStringToID(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      string
		want    database.Id
		wantErr bool
	}{
		{
			name:    "Positive ID",
			id:      "1",
			want:    1,
			wantErr: false,
		},
		{
			name: "Zero ID",
			id:   "0",
			want: 0,

			wantErr: true,
		},
		{
			name:    "Negative ID",
			id:      "-1",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.ParseStringToID(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseStringToID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseStringToID() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("ParseStringToID() = %v, want %v", got, tt.want)
			}
		})
	}
}
