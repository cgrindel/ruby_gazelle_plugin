package gazelle_test

import (
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/cgrindel/ruby_gazelle_plugin/gazelle"
)

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
