package scanner

import "testing"

func TestMethodInError(t *testing.T) {
	tests := []struct {
		name   string
		errStr string
		method string
		want   bool
	}{
		{
			name:   "method present",
			errStr: "ssh: unable to authenticate, attempted methods [none password], no supported methods remain",
			method: "password",
			want:   true,
		},
		{
			name:   "method absent",
			errStr: "ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain",
			method: "password",
			want:   false,
		},
		{
			name:   "no brackets",
			errStr: "connection refused",
			method: "password",
			want:   false,
		},
		{
			name:   "empty string",
			errStr: "",
			method: "password",
			want:   false,
		},
		{
			name:   "partial match does not count",
			errStr: "ssh: attempted methods [none passwordless], no match",
			method: "password",
			want:   false,
		},
		{
			name:   "multiple bracket groups",
			errStr: "attempted methods [none] then [password publickey]",
			method: "password",
			want:   true,
		},
		{
			name:   "trailing comma stripped",
			errStr: "ssh: attempted methods [none, password,], no match",
			method: "password",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := methodInError(tt.errStr, tt.method)
			if got != tt.want {
				t.Errorf("methodInError(%q, %q) = %v, want %v", tt.errStr, tt.method, got, tt.want)
			}
		})
	}
}
