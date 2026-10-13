package updater

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func serve(t *testing.T, status int, body string) {
	t.Helper()
	var userAgent string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	old := latestReleaseURL
	latestReleaseURL = srv.URL
	t.Cleanup(func() {
		srv.Close()
		latestReleaseURL = old
		if userAgent != "" && userAgent != "SimpleAC-Updater" {
			t.Errorf("User-Agent = %q", userAgent)
		}
	})
}

func TestCheckForUpdates(t *testing.T) {
	tests := []struct {
		name    string
		current string
		body    string
		want    string // expected tag, empty = no update
	}{
		{"newer release", "1.0.0", `{"tag_name":"v1.1.0","html_url":"https://example.com/r"}`, "v1.1.0"},
		{"same version", "1.1.0", `{"tag_name":"v1.1.0"}`, ""},
		{"same version, v prefix on current", "v1.1.0", `{"tag_name":"v1.1.0"}`, ""},
		{"same version, no prefix on tag", "1.1.0", `{"tag_name":"1.1.0"}`, ""},
		{"dev build is offered the release", "dev", `{"tag_name":"v1.0.0"}`, "v1.0.0"},
		{"empty tag", "1.0.0", `{"tag_name":""}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serve(t, http.StatusOK, tt.body)
			got, err := CheckForUpdates(tt.current)
			if err != nil {
				t.Fatal(err)
			}
			switch {
			case tt.want == "" && got != nil:
				t.Errorf("unexpected release %+v", got)
			case tt.want != "" && (got == nil || got.TagName != tt.want):
				t.Errorf("got %+v, want tag %s", got, tt.want)
			}
		})
	}
}

func TestCheckForUpdatesReturnsReleaseFields(t *testing.T) {
	serve(t, http.StatusOK, `{"tag_name":"v2.0.0","name":"Two","body":"notes","html_url":"https://example.com/r","published_at":"2026-01-01T00:00:00Z"}`)
	got, err := CheckForUpdates("1.0.0")
	if err != nil || got == nil {
		t.Fatalf("got %v, %v", got, err)
	}
	if got.Name != "Two" || got.Body != "notes" || got.HTMLURL != "https://example.com/r" || got.PublishedAt != "2026-01-01T00:00:00Z" {
		t.Errorf("fields not decoded: %+v", got)
	}
}

func TestCheckForUpdatesErrors(t *testing.T) {
	t.Run("http error status", func(t *testing.T) {
		serve(t, http.StatusNotFound, `{"message":"Not Found"}`)
		if rel, err := CheckForUpdates("1.0.0"); err == nil || rel != nil {
			t.Errorf("got %v, %v; want an error", rel, err)
		}
	})
	t.Run("rate limited", func(t *testing.T) {
		serve(t, http.StatusForbidden, `{"message":"rate limit"}`)
		if _, err := CheckForUpdates("1.0.0"); err == nil {
			t.Error("want an error")
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		serve(t, http.StatusOK, `not json`)
		if _, err := CheckForUpdates("1.0.0"); err == nil {
			t.Error("want an error")
		}
	})
	t.Run("unreachable server", func(t *testing.T) {
		old := latestReleaseURL
		latestReleaseURL = "http://127.0.0.1:1"
		t.Cleanup(func() { latestReleaseURL = old })
		if _, err := CheckForUpdates("1.0.0"); err == nil {
			t.Error("want an error")
		}
	})
}
