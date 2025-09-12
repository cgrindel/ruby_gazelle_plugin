package gazelle_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/cgrindel/ruby_gazelle_plugin/gazelle"
	"github.com/cgrindel/ruby_gazelle_plugin/gazelle/internal/rubycfg"
)

func TestGenerateRules(t *testing.T) {
	lang := &gazelle.RubyLang{}

	tests := []struct {
		name         string
		files        map[string]string
		regularFiles []string
		rel          string
		want         int
		wantSrcs     []string
		wantDeps     []string
	}{
		{
			name: "no Ruby files",
			files: map[string]string{
				"README.md": "# Test project",
			},
			regularFiles: []string{"README.md"},
			rel:          "",
			want:         0,
		},
		{
			name: "single Ruby file without dependencies",
			files: map[string]string{
				"main.rb": "puts 'Hello, World!'",
			},
			regularFiles: []string{"main.rb"},
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     nil,
		},
		{
			name: "multiple Ruby files without dependencies",
			files: map[string]string{
				"main.rb":   "puts 'Hello, World!'",
				"helper.rb": "def help; end",
			},
			regularFiles: []string{"main.rb", "helper.rb"},
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb", "helper.rb"},
			wantDeps:     nil,
		},
		{
			name: "Ruby file with require_relative",
			files: map[string]string{
				"main.rb": "require_relative 'helper'\nputs 'Hello, World!'",
			},
			regularFiles: []string{"main.rb"},
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{":helper"},
		},
		{
			name: "Ruby file with subdirectory require_relative",
			files: map[string]string{
				"main.rb": "require_relative 'lib/version'\nputs 'Hello, World!'",
			},
			regularFiles: []string{"main.rb"},
			rel:          "foo",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{"//foo/lib:version"},
		},
		{
			name: "multiple Ruby files with dependencies",
			files: map[string]string{
				"main.rb":   "require_relative 'helper'\nrequire_relative 'lib/version'",
				"helper.rb": "require_relative 'utils'",
			},
			regularFiles: []string{"main.rb", "helper.rb"},
			rel:          "foo",
			want:         1,
			wantSrcs:     []string{"main.rb", "helper.rb"},
			wantDeps:     []string{":helper", "//foo/lib:version", ":utils"},
		},
		{
			name: "duplicate dependencies removed across files",
			files: map[string]string{
				"main.rb":  "require_relative 'helper'\nrequire_relative 'utils'",
				"other.rb": "require_relative 'helper'\nrequire_relative 'config'",
			},
			regularFiles: []string{"main.rb", "other.rb"},
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb", "other.rb"},
			wantDeps:     []string{":helper", ":utils", ":config"},
		},
		{
			name: "Ruby files with double quotes",
			files: map[string]string{
				"main.rb": "require_relative \"helper\"\nrequire_relative \"lib/version\"",
			},
			regularFiles: []string{"main.rb"},
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{":helper", "//lib:version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "gazelle_test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir) // nolint:errcheck

			for fileName, content := range tt.files {
				filePath := filepath.Join(tmpDir, fileName)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", fileName, err)
				}
			}

			args := language.GenerateArgs{
				Dir:          tmpDir,
				Rel:          tt.rel,
				RegularFiles: tt.regularFiles,
			}

			result := lang.GenerateRules(args)

			if len(result.Gen) != tt.want {
				t.Errorf("GenerateRules() generated %d rules, want %d", len(result.Gen), tt.want)
				return
			}

			if tt.want == 0 {
				return
			}

			rule := result.Gen[0]
			if rule.Kind() != "rb_library" {
				t.Errorf("Generated rule kind = %s, want rb_library", rule.Kind())
			}

			srcs := rule.AttrStrings("srcs")
			if !reflect.DeepEqual(srcs, tt.wantSrcs) {
				t.Errorf("Generated rule srcs = %v, want %v", srcs, tt.wantSrcs)
			}

			deps := rule.AttrStrings("deps")
			if !reflect.DeepEqual(deps, tt.wantDeps) {
				t.Errorf("Generated rule deps = %v, want %v", deps, tt.wantDeps)
			}
		})
	}
}

func TestGenerateRulesWithMixedFiles(t *testing.T) {
	lang := &gazelle.RubyLang{}

	tmpDir, err := os.MkdirTemp("", "gazelle_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir) // nolint:errcheck

	files := map[string]string{
		"main.rb":    "require_relative 'helper'",
		"helper.rb":  "def help; end",
		"README.md":  "# Documentation",
		"config.yml": "key: value",
		"test.py":    "print('not ruby')",
	}

	for fileName, content := range files {
		filePath := filepath.Join(tmpDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fileName, err)
		}
	}

	args := language.GenerateArgs{
		Dir:          tmpDir,
		Rel:          "",
		RegularFiles: []string{"main.rb", "helper.rb", "README.md", "config.yml", "test.py"},
	}

	result := lang.GenerateRules(args)

	if len(result.Gen) != 1 {
		t.Errorf("GenerateRules() generated %d rules, want 1", len(result.Gen))
		return
	}

	rule := result.Gen[0]
	srcs := rule.AttrStrings("srcs")
	expectedSrcs := []string{"main.rb", "helper.rb"}

	if !reflect.DeepEqual(srcs, expectedSrcs) {
		t.Errorf("Generated rule srcs = %v, want %v", srcs, expectedSrcs)
	}
}

func TestGenerateRulesErrorHandling(t *testing.T) {
	lang := &gazelle.RubyLang{}

	tmpDir, err := os.MkdirTemp("", "gazelle_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir) // nolint:errcheck

	rubyFile := filepath.Join(tmpDir, "main.rb")
	if err := os.WriteFile(rubyFile, []byte("require_relative 'helper'"), 0644); err != nil {
		t.Fatalf("Failed to write Ruby file: %v", err)
	}

	// Test with unreadable file (simulate by removing after creating)
	if err := os.Remove(rubyFile); err != nil {
		t.Fatalf("Failed to remove Ruby file: %v", err)
	}

	args := language.GenerateArgs{
		Dir:          tmpDir,
		Rel:          "",
		RegularFiles: []string{"main.rb"},
	}

	result := lang.GenerateRules(args)

	// Should still generate a rule even if file reading fails
	if len(result.Gen) != 1 {
		t.Errorf("GenerateRules() generated %d rules, want 1", len(result.Gen))
		return
	}

	rule := result.Gen[0]
	srcs := rule.AttrStrings("srcs")
	if !reflect.DeepEqual(srcs, []string{"main.rb"}) {
		t.Errorf("Generated rule srcs = %v, want [main.rb]", srcs)
	}

	// Should have no dependencies due to read error
	deps := rule.AttrStrings("deps")
	if len(deps) != 0 {
		t.Errorf("Generated rule deps = %v, want []", deps)
	}
}

func TestGenerateRulesWithRequireStatements(t *testing.T) {
	tests := []struct {
		name         string
		files        map[string]string
		regularFiles []string
		bundleRepo   string
		rel          string
		want         int
		wantSrcs     []string
		wantDeps     []string
	}{
		{
			name: "single require statement",
			files: map[string]string{
				"main.rb": "require 'json'\nputs 'Hello, World!'",
			},
			regularFiles: []string{"main.rb"},
			bundleRepo:   "my_gems",
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{"@my_gems"},
		},
		{
			name: "multiple require statements",
			files: map[string]string{
				"main.rb": "require 'json'\nrequire 'yaml'\nputs 'Hello!'",
			},
			regularFiles: []string{"main.rb"},
			bundleRepo:   "bundle",
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{"@bundle"},
		},
		{
			name: "require and require_relative mixed",
			files: map[string]string{
				"main.rb": "require 'json'\nrequire_relative 'helper'",
			},
			regularFiles: []string{"main.rb"},
			bundleRepo:   "my_bundle",
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{":helper", "@my_bundle"},
		},
		{
			name: "no bundle repo configured",
			files: map[string]string{
				"main.rb": "require 'json'\nputs 'Hello!'",
			},
			regularFiles: []string{"main.rb"},
			bundleRepo:   "", // no bundle repo explicitly set
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{"@bundle"}, // uses default "bundle"
		},
		{
			name: "duplicate require statements",
			files: map[string]string{
				"main.rb":  "require 'json'\nrequire 'yaml'",
				"other.rb": "require 'json'\nrequire 'csv'",
			},
			regularFiles: []string{"main.rb", "other.rb"},
			bundleRepo:   "gems",
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb", "other.rb"},
			wantDeps:     []string{"@gems"},
		},
		{
			name: "require with double quotes",
			files: map[string]string{
				"main.rb": "require \"json\"\nrequire \"yaml\"",
			},
			regularFiles: []string{"main.rb"},
			bundleRepo:   "bundle",
			rel:          "",
			want:         1,
			wantSrcs:     []string{"main.rb"},
			wantDeps:     []string{"@bundle"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "gazelle_test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir) // nolint:errcheck

			for fileName, content := range tt.files {
				filePath := filepath.Join(tmpDir, fileName)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", fileName, err)
				}
			}

			// Create config with bundle repo
			c := &config.Config{}
			if tt.bundleRepo != "" {
				rubyConfig := &rubycfg.RubyConfig{
					BundleRepoName: tt.bundleRepo,
				}
				rubycfg.SetRubyConfig(c, rubyConfig)
			}

			args := language.GenerateArgs{
				Config:       c,
				Dir:          tmpDir,
				Rel:          tt.rel,
				RegularFiles: tt.regularFiles,
			}

			lang := &gazelle.RubyLang{}
			result := lang.GenerateRules(args)

			if len(result.Gen) != tt.want {
				t.Errorf("GenerateRules() generated %d rules, want %d",
					len(result.Gen), tt.want)
				return
			}

			if tt.want == 0 {
				return
			}

			rule := result.Gen[0]
			if rule.Kind() != "rb_library" {
				t.Errorf("Generated rule kind = %s, want rb_library",
					rule.Kind())
			}

			srcs := rule.AttrStrings("srcs")
			if !reflect.DeepEqual(srcs, tt.wantSrcs) {
				t.Errorf("Generated rule srcs = %v, want %v",
					srcs, tt.wantSrcs)
			}

			deps := rule.AttrStrings("deps")
			if !reflect.DeepEqual(deps, tt.wantDeps) {
				t.Errorf("Generated rule deps = %v, want %v",
					deps, tt.wantDeps)
			}
		})
	}
}
