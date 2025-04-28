package validation

import "testing"

func TestValidateVideoURL(t *testing.T) {
	cases := []struct {
		input   string
		wantErr error
	}{
		{"", ErrEmptyURL},
		{"   ", ErrEmptyURL},
		{"notaurl", ErrInvalidURL},
		{"https://google.com", ErrNotYouTube},
		{"https://youtu.be/abc123", nil},
		{"https://www.youtube.com/watch?v=abc123", nil},
	}
	for _, tc := range cases {
		err := ValidateVideoURL(tc.input)
		if err != tc.wantErr {
			t.Errorf("ValidateVideoURL(%q) = %v, want %v", tc.input, err, tc.wantErr)
		}
	}
}
