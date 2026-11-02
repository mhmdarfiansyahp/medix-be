package dto

import (
	"encoding/json"
	"testing"
)

func TestUserStatusUnmarshal(t *testing.T) {
	cases := []struct {
		input    string
		expected UserStatus
	}{
		{`{"status":1}`, "aktif"},
		{`{"status":"aktif"}`, "aktif"},
		{`{"status":0}`, "nonaktif"},
		{`{"status":"nonaktif"}`, "nonaktif"},
	}

	for _, tc := range cases {
		var req CreateUserRequest
		if err := json.Unmarshal([]byte(tc.input), &req); err != nil {
			t.Fatalf("failed to unmarshal %s: %v", tc.input, err)
		}
		if req.Status != tc.expected {
			t.Fatalf("expected %q, got %q for %s", tc.expected, req.Status, tc.input)
		}
	}
}
