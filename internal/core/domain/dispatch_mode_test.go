package domain

import "testing"

func TestParseMessageDispatchMode(t *testing.T) {
	tests := []struct {
		input string
		want  MessageDispatchMode
		ok    bool
	}{
		{string(DispatchMemoryOnly), DispatchMemoryOnly, true},
		{string(DispatchProviderOnly), DispatchProviderOnly, true},
		{string(DispatchProviderAndDatabase), DispatchProviderAndDatabase, true},
		{string(DispatchMemoryAndProvider), DispatchMemoryAndProvider, true},
		{"invalid", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, ok := ParseMessageDispatchMode(tt.input)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("mode = %q, want %q", got, tt.want)
			}
		})
	}
}
