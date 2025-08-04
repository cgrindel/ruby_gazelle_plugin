package ruby

import (
	"flag"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const rubyName = "ruby"

type rubyLang struct{}

// NewLanguage returns a new instance of the Ruby language extension.
func NewLanguage() language.Language {
	return &rubyLang{}
}

func (*rubyLang) Name() string {
	return rubyName
}

func (*rubyLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {}

func (*rubyLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

func (*rubyLang) KnownDirectives() []string {
	return nil
}

func (*rubyLang) Configure(c *config.Config, rel string, f *rule.File) {}

func (*rubyLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	return nil
}

func (*rubyLang) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

func (*rubyLang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
}

func (*rubyLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	return language.GenerateResult{}
}

func (*rubyLang) Fix(c *config.Config, f *rule.File) {}

func (*rubyLang) Kinds() map[string]rule.KindInfo {
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

func (*rubyLang) Loads() []rule.LoadInfo {
	return []rule.LoadInfo{
		{
			Name:    "@rules_ruby//ruby:defs.bzl",
			Symbols: []string{"rb_library", "rb_binary", "rb_test"},
		},
	}
}
