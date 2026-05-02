package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMaxRecvMsgSize(t *testing.T) {
	t.Run("empty value should return default", func(t *testing.T) {
		t.Setenv(gaugeMaxMessageSize, "")
		v := GetMaxMessageSize()
		if v != 1024 {
			t.Errorf("Expected 1024, got %d", v)
		}
	})

	t.Run("non-numeric should return default", func(t *testing.T) {
		t.Setenv(gaugeMaxMessageSize, "abcd")
		v := GetMaxMessageSize()
		if v != 1024 {
			t.Errorf("Expected 1024, got %d", v)
		}
	})

	t.Run("numeric should return set value", func(t *testing.T) {
		t.Setenv(gaugeMaxMessageSize, "2048")
		v := GetMaxMessageSize()
		if v != 2048 {
			t.Errorf("Expected 2048, got %d", v)
		}
	})
}

// isEnvSet is the shared truthiness check behind every Should* helper. The
// canonical "true" value is the lowercase string "true"; everything else is
// false. The case-insensitive comparison is the surprising bit, so guard it
// here.
func TestIsEnvSet(t *testing.T) {
	const key = "TEST_ENV_FLAG"
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{"true-lower", "true", true},
		{"true-upper", "TRUE", true},
		{"true-mixed", "True", true},
		{"false", "false", false},
		{"empty", "", false},
		{"garbage", "yes", false},
		{"numeric", "1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(key, tt.val)
			if got := isEnvSet(key); got != tt.want {
				t.Errorf("isEnvSet(%q) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

func TestShouldOverwriteReports(t *testing.T) {
	t.Setenv(OverwriteReportsEnvProperty, "true")
	if !ShouldOverwriteReports() {
		t.Error("ShouldOverwriteReports() = false, want true when env=true")
	}
	t.Setenv(OverwriteReportsEnvProperty, "false")
	if ShouldOverwriteReports() {
		t.Error("ShouldOverwriteReports() = true, want false when env=false")
	}
}

func TestShouldUseNestedSpecs(t *testing.T) {
	t.Setenv(UseNestedSpecs, "true")
	if !ShouldUseNestedSpecs() {
		t.Error("ShouldUseNestedSpecs() = false, want true when env=true")
	}
	t.Setenv(UseNestedSpecs, "")
	if ShouldUseNestedSpecs() {
		t.Error("ShouldUseNestedSpecs() = true, want false when env unset")
	}
}

// CreateDirectory is a thin os.MkdirAll wrapper, but the wrapper is what
// production calls — verify it actually creates the path and is idempotent
// on a path that already exists.
func TestCreateDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "c")
	CreateDirectory(dir)
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected directory created at %s, got err: %v", dir, err)
	}
	// Idempotence: a second call on an existing path should not fail.
	CreateDirectory(dir)
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("idempotent CreateDirectory broke %s: %v", dir, err)
	}
}

// getDefaultPropertiesFile composes <projectRoot>/env/default/default.properties.
// Stub GetProjectRoot so the path is deterministic without requiring a real
// project layout.
func TestGetDefaultPropertiesFile(t *testing.T) {
	saved := GetProjectRoot
	GetProjectRoot = func() string { return "/synthetic/proj" }
	t.Cleanup(func() { GetProjectRoot = saved })

	got := getDefaultPropertiesFile()
	want := filepath.Join("/synthetic/proj", "env", "default", "default.properties")
	if got != want {
		t.Errorf("getDefaultPropertiesFile() = %q, want %q", got, want)
	}
}

// PluginKillTimeout converts a millisecond env value to seconds. The
// fallthroughs (unset, non-numeric) both return 0.
func TestPluginKillTimeout(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv("plugin_kill_timeout", "")
		if got := PluginKillTimeout(); got != 0 {
			t.Errorf("PluginKillTimeout() = %d, want 0 when unset", got)
		}
	})
	t.Run("non-numeric", func(t *testing.T) {
		t.Setenv("plugin_kill_timeout", "abc")
		if got := PluginKillTimeout(); got != 0 {
			t.Errorf("PluginKillTimeout() = %d, want 0 when non-numeric", got)
		}
	})
	t.Run("ms-to-seconds", func(t *testing.T) {
		t.Setenv("plugin_kill_timeout", "5000")
		if got := PluginKillTimeout(); got != 5 {
			t.Errorf("PluginKillTimeout() = %d, want 5 (5000ms / 1000)", got)
		}
	})
}
