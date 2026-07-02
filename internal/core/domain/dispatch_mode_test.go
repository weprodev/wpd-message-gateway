package domain

import "testing"

func TestSettingValueToDispatchConfig(t *testing.T) {
	tests := []struct {
		value string
		want  MessageDispatchConfig
		ok    bool
	}{
		{"memory_only", MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: false}, true},
		{"memory_and_database", MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: true}, true},
		{"provider_only", MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: false}, true},
		{"provider_and_database", MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: true}, true},
		{"memory", MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: false}, true},
		{"both", MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: true}, true},
		{"memory_database", MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: true}, true},
		{"memory_and_provider", MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: true}, true},
		{"provider", MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: false}, true},
		{"provider_database", MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: true}, true},
		{"unknown", MessageDispatchConfig{}, false},
		{"", MessageDispatchConfig{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, ok := SettingValueToDispatchConfig(tt.value)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestMessageDispatchConfig_APIValue(t *testing.T) {
	tests := []struct {
		cfg  MessageDispatchConfig
		want MessageDispatchAPIValue
	}{
		{MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: false}, APIMemoryOnly},
		{MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: true}, APIMemoryAndDatabase},
		{MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: false}, APIProviderOnly},
		{MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: true}, APIProviderAndDatabase},
	}

	for _, tt := range tests {
		t.Run(string(tt.want), func(t *testing.T) {
			if got := tt.cfg.APIValue(); got != tt.want {
				t.Fatalf("APIValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseMessageDispatchConfig_roundTrip(t *testing.T) {
	apiValues := []MessageDispatchAPIValue{
		APIMemoryOnly,
		APIMemoryAndDatabase,
		APIProviderOnly,
		APIProviderAndDatabase,
	}

	for _, apiValue := range apiValues {
		t.Run(string(apiValue), func(t *testing.T) {
			cfg, ok := ParseMessageDispatchConfig(string(apiValue))
			if !ok {
				t.Fatal("expected ok")
			}
			if cfg.APIValue() != apiValue {
				t.Fatalf("round-trip got %q, want %q", cfg.APIValue(), apiValue)
			}
		})
	}
}
