package gazelle_test

import (
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/cgrindel/ruby_gazelle_plugin/gazelle"
	"github.com/cgrindel/ruby_gazelle_plugin/gazelle/internal/rubycfg"
)

func TestKnownDirectives(t *testing.T) {
	lang := &gazelle.RubyLang{}

	directives := lang.KnownDirectives()
	expected := []string{"ruby_bundle_repo"}

	if !reflect.DeepEqual(directives, expected) {
		t.Errorf("KnownDirectives() = %v, want %v", directives, expected)
	}
}

func TestApparentLoads(t *testing.T) {
	lang := &gazelle.RubyLang{}

	tests := []struct {
		name                 string
		moduleToApparentName func(string) string
		want                 []rule.LoadInfo
	}{
		{
			name: "custom repo name",
			moduleToApparentName: func(module string) string {
				if module == "rules_ruby" {
					return "my_ruby_rules"
				}
				return ""
			},
			want: []rule.LoadInfo{
				{
					Name:    "@my_ruby_rules//ruby:defs.bzl",
					Symbols: []string{"rb_library", "rb_binary", "rb_test"},
				},
			},
		},
		{
			name: "empty moduleToApparentName returns",
			moduleToApparentName: func(module string) string {
				return ""
			},
			want: []rule.LoadInfo{
				{
					Name:    "@rules_ruby//ruby:defs.bzl",
					Symbols: []string{"rb_library", "rb_binary", "rb_test"},
				},
			},
		},
		{
			name: "nil moduleToApparentName function",
			moduleToApparentName: func(module string) string {
				// Simulate a function that doesn't handle rules_ruby
				if module == "other_module" {
					return "other_apparent_name"
				}
				return ""
			},
			want: []rule.LoadInfo{
				{
					Name:    "@rules_ruby//ruby:defs.bzl",
					Symbols: []string{"rb_library", "rb_binary", "rb_test"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lang.ApparentLoads(tt.moduleToApparentName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApparentLoads() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApparentLoadsStructure(t *testing.T) {
	lang := &gazelle.RubyLang{}

	// Test with a simple identity function
	loads := lang.ApparentLoads(func(module string) string {
		return module
	})

	// Verify we get exactly one LoadInfo
	if len(loads) != 1 {
		t.Errorf("Expected 1 LoadInfo, got %d", len(loads))
	}

	loadInfo := loads[0]

	// Verify the structure
	expectedSymbols := []string{"rb_library", "rb_binary", "rb_test"}
	if !reflect.DeepEqual(loadInfo.Symbols, expectedSymbols) {
		t.Errorf("Expected symbols %v, got %v", expectedSymbols, loadInfo.Symbols)
	}

	// Verify the load name format
	expectedName := "@rules_ruby//ruby:defs.bzl"
	if loadInfo.Name != expectedName {
		t.Errorf("Expected load name %s, got %s", expectedName, loadInfo.Name)
	}
}

func TestConvertRequireRelativeToLabel(t *testing.T) {
	tests := []struct {
		name           string
		currentPackage string
		requiredPath   string
		expected       string
	}{
		{
			name:           "simple file in subdirectory",
			currentPackage: "foo/lib",
			requiredPath:   "hello_world/version",
			expected:       "//foo/lib/hello_world:version",
		},
		{
			name:           "file in current directory",
			currentPackage: "foo/lib",
			requiredPath:   "version",
			expected:       ":version",
		},
		{
			name:           "nested subdirectory",
			currentPackage: "foo/lib",
			requiredPath:   "hello_world/sub/version",
			expected:       "//foo/lib/hello_world/sub:version",
		},
		{
			name:           "root package with subdirectory",
			currentPackage: "",
			requiredPath:   "hello_world/version",
			expected:       "//hello_world:version",
		},
		{
			name:           "root package current directory",
			currentPackage: "",
			requiredPath:   "version",
			expected:       ":version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to access the function, so we'll need to expose it or test through GenerateRules
			// For now, let's create a small wrapper to test the logic
			got := gazelle.ConvertRequireRelativeToLabel(tt.currentPackage, tt.requiredPath)
			if got != tt.expected {
				t.Errorf("convertRequireRelativeToLabel(%q, %q) = %q, want %q",
					tt.currentPackage, tt.requiredPath, got, tt.expected)
			}
		})
	}
}

func TestConfigure(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		rel            string
		directives     []rule.Directive
		expectedBundle string
	}{
		{
			name:           "nil file",
			config:         &config.Config{},
			rel:            "",
			directives:     nil,
			expectedBundle: "bundle", // default value
		},
		{
			name:   "ruby_bundle_repo directive",
			config: &config.Config{},
			rel:    "",
			directives: []rule.Directive{
				{Key: "ruby_bundle_repo", Value: "my_gems"},
			},
			expectedBundle: "my_gems",
		},
		{
			name:   "multiple directives",
			config: &config.Config{},
			rel:    "",
			directives: []rule.Directive{
				{Key: "other_directive", Value: "ignored"},
				{Key: "ruby_bundle_repo", Value: "custom_bundle"},
				{Key: "another_directive", Value: "also_ignored"},
			},
			expectedBundle: "custom_bundle",
		},
		{
			name: "override existing config",
			config: func() *config.Config {
				c := &config.Config{}
				rc := &rubycfg.RubyConfig{BundleRepoName: "old_bundle"}
				rubycfg.SetRubyConfig(c, rc)
				return c
			}(),
			rel: "",
			directives: []rule.Directive{
				{Key: "ruby_bundle_repo", Value: "new_bundle"},
			},
			expectedBundle: "new_bundle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := &gazelle.RubyLang{}

			var f *rule.File
			if tt.directives != nil {
				f = &rule.File{Directives: tt.directives}
			}

			lang.Configure(tt.config, tt.rel, f)

			rubyConfig := rubycfg.GetRubyConfig(tt.config)
			if rubyConfig.BundleRepoName != tt.expectedBundle {
				t.Errorf("Configure() set BundleRepoName = %q, want %q",
					rubyConfig.BundleRepoName, tt.expectedBundle)
			}
		})
	}
}
