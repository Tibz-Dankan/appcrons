package services

import "testing"

func TestRemoveSpaces(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no spaces", "https://myapp.onrender.com/active", "https://myapp.onrender.com/active"},
		{"leading and trailing spaces", "  https://myapp.onrender.com/active ", "https://myapp.onrender.com/active"},
		{"interior space", "https://myapp.onrender.com/ active", "https://myapp.onrender.com/active"},
		{"tab and newline", "https://myapp.onrender.com\t/active\n", "https://myapp.onrender.com/active"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveSpaces(tt.input); got != tt.want {
				t.Errorf("RemoveSpaces(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
