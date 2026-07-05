package contracts

import "testing"

func TestStoreContentFromResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    *SendResult
		want bool
	}{
		{"nil result", nil, false},
		{"nil meta", &SendResult{ID: "1"}, false},
		{"false", &SendResult{Meta: map[string]string{MetaKeyStoreContent: metaStoreContentFalse}}, false},
		{"true", &SendResult{Meta: map[string]string{MetaKeyStoreContent: metaStoreContentTrue}}, true},
		{"missing key", &SendResult{Meta: map[string]string{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := StoreContentFromResult(tt.r); got != tt.want {
				t.Fatalf("StoreContentFromResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDispatchModeFromResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    *SendResult
		want string
	}{
		{"nil result", nil, ""},
		{"memory", &SendResult{Meta: map[string]string{MetaKeyDispatchMode: DispatchModeMemory}}, DispatchModeMemory},
		{"provider", &SendResult{Meta: map[string]string{MetaKeyDispatchMode: DispatchModeProvider}}, DispatchModeProvider},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := DispatchModeFromResult(tt.r); got != tt.want {
				t.Fatalf("DispatchModeFromResult() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInboxMessageIDFromResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    *SendResult
		want string
	}{
		{"nil result", nil, ""},
		{"provider without inbox", &SendResult{ID: "prov-1", Meta: map[string]string{MetaKeyDispatchMode: DispatchModeProvider}}, ""},
		{"memory mode retained", &SendResult{
			ID: "inbox-1",
			Meta: map[string]string{
				MetaKeyDispatchMode: DispatchModeMemory,
				MetaKeyStoreContent: metaStoreContentTrue,
			},
		}, "inbox-1"},
		{"memory mode not retained", &SendResult{
			ID: "mem-1",
			Meta: map[string]string{
				MetaKeyDispatchMode: DispatchModeMemory,
				MetaKeyStoreContent: metaStoreContentFalse,
			},
		}, ""},
		{"provider with inbox meta", &SendResult{
			ID: "prov-1",
			Meta: map[string]string{
				MetaKeyDispatchMode:   DispatchModeProvider,
				MetaKeyInboxMessageID: "inbox-2",
			},
		}, "inbox-2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := InboxMessageIDFromResult(tt.r); got != tt.want {
				t.Fatalf("InboxMessageIDFromResult() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetStoreContentMeta(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()
		r := &SendResult{}
		SetStoreContentMeta(r, true)
		if r.Meta[MetaKeyStoreContent] != metaStoreContentTrue {
			t.Fatalf("meta: %q", r.Meta[MetaKeyStoreContent])
		}
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()
		r := &SendResult{Meta: map[string]string{"other": "x"}}
		SetStoreContentMeta(r, false)
		if r.Meta[MetaKeyStoreContent] != metaStoreContentFalse {
			t.Fatalf("meta: %q", r.Meta[MetaKeyStoreContent])
		}
		if r.Meta["other"] != "x" {
			t.Fatalf("unexpected meta overwrite")
		}
	})
}

func TestSetDispatchModeMeta(t *testing.T) {
	t.Parallel()

	r := &SendResult{}
	SetDispatchModeMeta(r, DispatchModeProvider)
	if got := DispatchModeFromResult(r); got != DispatchModeProvider {
		t.Fatalf("dispatch_mode: %q", got)
	}
}

func TestDispatchModeValuesMatchDomain(t *testing.T) {
	t.Parallel()

	// Document the cross-package contract without importing internal/domain.
	if DispatchModeMemory != "memory" || DispatchModeProvider != "provider" {
		t.Fatal("dispatch mode constants must stay aligned with domain.MessageDispatchMode")
	}
}
