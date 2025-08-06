// Package gazelle implements a Gazelle language extension for Ruby.
// This package provides support for generating BUILD files for Ruby projects
// using Bazel's Gazelle tool.
package gazelle

import (
	"flag"
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const rubyName = "ruby"

// RubyLang implements the Gazelle language interface for Ruby.
type RubyLang struct {
	language.BaseLang
}

// NewLanguage returns a new instance of the Ruby language extension.
func NewLanguage() language.Language {
	return &RubyLang{}
}

func (*RubyLang) Name() string {
	return rubyName
}

func (*RubyLang) RegisterFlags(
	fs *flag.FlagSet, cmd string, c *config.Config,
) {
}

func (*RubyLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

func (*RubyLang) KnownDirectives() []string {
	return nil
}

func (*RubyLang) Configure(c *config.Config, rel string, f *rule.File) {}

func (*RubyLang) Imports(
	c *config.Config, r *rule.Rule, f *rule.File,
) []resolve.ImportSpec {
	return nil
}

func (*RubyLang) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

func (*RubyLang) Resolve(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	imports any,
	from label.Label,
) {
}

func (*RubyLang) GenerateRules(
	args language.GenerateArgs,
) language.GenerateResult {
	return language.GenerateResult{}
}

func (*RubyLang) Fix(c *config.Config, f *rule.File) {}

func (*RubyLang) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		"rb_library": {
			MatchAny: false,
			NonEmptyAttrs: map[string]bool{
				"srcs": true,
			},
			MergeableAttrs: map[string]bool{
				"srcs": true,
				"deps": true,
			},
			ResolveAttrs: map[string]bool{
				"deps": true,
			},
		},
		"rb_binary": {
			MatchAny: false,
			NonEmptyAttrs: map[string]bool{
				"srcs": true,
			},
			MergeableAttrs: map[string]bool{
				"srcs": true,
				"deps": true,
			},
			ResolveAttrs: map[string]bool{
				"deps": true,
			},
		},
		"rb_test": {
			MatchAny: false,
			NonEmptyAttrs: map[string]bool{
				"srcs": true,
			},
			MergeableAttrs: map[string]bool{
				"srcs": true,
				"deps": true,
			},
			ResolveAttrs: map[string]bool{
				"deps": true,
			},
		},
	}
}

func (*RubyLang) ApparentLoads(
	moduleToApparentName func(string) string,
) []rule.LoadInfo {
	rulesRuby := moduleToApparentName("rules_ruby")
	if rulesRuby == "" {
		rulesRuby = "rules_ruby"
	}
	return []rule.LoadInfo{
		{
			Name:    fmt.Sprintf("@%s//ruby:defs.bzl", rulesRuby),
			Symbols: []string{"rb_library", "rb_binary", "rb_test"},
		},
	}
}
