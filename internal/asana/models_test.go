package asana

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCustomTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		checkFn func(*CustomTime) bool
	}{
		{
			name:    "null",
			json:    `null`,
			wantErr: false,
			checkFn: func(ct *CustomTime) bool { return ct.IsZero() },
		},
		{
			name:    "date only",
			json:    `"2026-02-19"`,
			wantErr: false,
			checkFn: func(ct *CustomTime) bool {
				return ct.Year() == 2026 && ct.Month() == 2 && ct.Day() == 19
			},
		},
		{
			name:    "RFC3339",
			json:    `"2026-02-19T10:30:00Z"`,
			wantErr: false,
			checkFn: func(ct *CustomTime) bool {
				return ct.Year() == 2026
			},
		},
		{
			name:    "invalid format",
			json:    `"not-a-date"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime
			err := json.Unmarshal([]byte(tt.json), &ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.checkFn != nil && !tt.checkFn(&ct) {
				t.Error("CustomTime value doesn't match expected")
			}
		})
	}
}

func TestCustomTimeIsZero(t *testing.T) {
	var ct CustomTime
	if !ct.IsZero() {
		t.Error("empty CustomTime should be zero")
	}

	ct.Time = time.Now()
	if ct.IsZero() {
		t.Error("CustomTime with time should not be zero")
	}
}

func TestTaskStructure(t *testing.T) {
	// Just verify the struct can be created
	task := &Task{
		GID:       "test-gid",
		Name:      "Test Task",
		Completed: false,
	}

	if task.GID != "test-gid" {
		t.Error("Task GID not set correctly")
	}

	if task.Name != "Test Task" {
		t.Error("Task Name not set correctly")
	}
}