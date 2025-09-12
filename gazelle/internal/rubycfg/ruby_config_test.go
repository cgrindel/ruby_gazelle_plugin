package rubycfg_test

import (
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/cgrindel/ruby_gazelle_plugin/gazelle/internal/rubycfg"
)

func TestNewRubyConfig(t *testing.T) {
	rc := rubycfg.NewRubyConfig()

	if rc == nil {
		t.Fatal("NewRubyConfig() returned nil")
	}

	if rc.BundleRepoName != "bundle" {
		t.Errorf("NewRubyConfig().BundleRepoName = %q, want %q",
			rc.BundleRepoName, "bundle")
	}
}

func TestGetRubyConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected *rubycfg.RubyConfig
	}{
		{
			name:     "nil config",
			config:   nil,
			expected: rubycfg.NewRubyConfig(),
		},
		{
			name:     "empty config",
			config:   &config.Config{},
			expected: rubycfg.NewRubyConfig(),
		},
		{
			name: "config with ruby config",
			config: func() *config.Config {
				c := &config.Config{}
				rc := &rubycfg.RubyConfig{
					BundleRepoName: "custom_bundle",
				}
				rubycfg.SetRubyConfig(c, rc)
				return c
			}(),
			expected: &rubycfg.RubyConfig{
				BundleRepoName: "custom_bundle",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rubycfg.GetRubyConfig(tt.config)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetRubyConfig() = %v, want %v",
					got, tt.expected)
			}
		})
	}
}

func TestSetRubyConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		rubyConfig  *rubycfg.RubyConfig
		expectPanic bool
	}{
		{
			name:        "nil config",
			config:      nil,
			rubyConfig:  rubycfg.NewRubyConfig(),
			expectPanic: false,
		},
		{
			name:   "valid config and ruby config",
			config: &config.Config{},
			rubyConfig: &rubycfg.RubyConfig{
				BundleRepoName: "test_bundle",
			},
		},
		{
			name: "config with existing extensions",
			config: &config.Config{
				Exts: map[string]any{"other": "value"},
			},
			rubyConfig: &rubycfg.RubyConfig{
				BundleRepoName: "another_bundle",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rubycfg.SetRubyConfig(tt.config, tt.rubyConfig)

			if tt.config != nil {
				retrieved := rubycfg.GetRubyConfig(tt.config)
				if !reflect.DeepEqual(retrieved, tt.rubyConfig) {
					t.Errorf("After SetRubyConfig, "+
						"GetRubyConfig() = %v, want %v",
						retrieved, tt.rubyConfig)
				}
			}
		})
	}
}

func TestDeepCopyRubyConfig(t *testing.T) {
	tests := []struct {
		name     string
		original *rubycfg.RubyConfig
		expected *rubycfg.RubyConfig
	}{
		{
			name:     "nil config",
			original: nil,
			expected: rubycfg.NewRubyConfig(),
		},
		{
			name:     "default config",
			original: rubycfg.NewRubyConfig(),
			expected: rubycfg.NewRubyConfig(),
		},
		{
			name: "custom config",
			original: &rubycfg.RubyConfig{
				BundleRepoName: "custom",
			},
			expected: &rubycfg.RubyConfig{
				BundleRepoName: "custom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rubycfg.DeepCopyRubyConfig(tt.original)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DeepCopyRubyConfig() = %v, want %v",
					got, tt.expected)
			}

			// Verify it's a different instance (not the same pointer)
			if tt.original != nil && got == tt.original {
				t.Error("DeepCopyRubyConfig() returned the same " +
					"instance, expected a copy")
			}

			// Verify modifying the copy doesn't affect the original
			if tt.original != nil {
				originalValue := tt.original.BundleRepoName
				got.BundleRepoName = "modified"
				if tt.original.BundleRepoName != originalValue {
					t.Error("Modifying copy affected the original")
				}
			}
		})
	}
}
