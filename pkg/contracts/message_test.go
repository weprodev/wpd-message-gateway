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
