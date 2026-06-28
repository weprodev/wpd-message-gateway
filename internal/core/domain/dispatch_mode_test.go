package domain

import "testing"

func TestSettingValueToDispatchMode(t *testing.T) {
	tests := []struct {
		value string
		want  MessageDispatchMode
		ok    bool
	}{
		{"memory", DispatchMemoryOnly, true},
		{"both", DispatchMemoryAndProvider, true},
		{"memory_database", DispatchMemoryAndProvider, true},
		{"providers", DispatchProviderOnly, true},
		{"provider", DispatchProviderOnly, true},
		{"provider_database", DispatchProviderAndDatabase, true},
		{"memory_only", DispatchMemoryOnly, true},
		{"memory_and_provider", DispatchMemoryAndProvider, true},
		{"provider_only", DispatchProviderOnly, true},
		{"provider_and_database", DispatchProviderAndDatabase, true},
		{"unknown", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, ok := SettingValueToDispatchMode(tt.value)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShouldRetainRequestLog(t *testing.T) {
	tests := []struct {
		mode MessageDispatchMode
		want bool
	}{
		{DispatchMemoryOnly, false},
		{DispatchProviderOnly, false},
		{DispatchMemoryAndProvider, true},
		{DispatchProviderAndDatabase, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			if got := ShouldRetainRequestLog(tt.mode); got != tt.want {
				t.Fatalf("ShouldRetainRequestLog(%q) = %v, want %v", tt.mode, got, tt.want)
			}
		})
	}
}
