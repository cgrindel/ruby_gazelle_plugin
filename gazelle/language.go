// Package gazelle implements a Gazelle language extension for Ruby.
// This package provides support for generating BUILD files for Ruby projects
// using Bazel's Gazelle tool.
package gazelle

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const rubyName = "ruby"

// Regular expression to match require_relative statements in Ruby files
var requireRelativeRegex = regexp.MustCompile(`require_relative\s+["']([^"']+)["']`)

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
	var result language.GenerateResult
	var rubyFiles []string
	var dependencies []string

	// Find all Ruby files in the directory
	for _, f := range args.RegularFiles {
		if strings.HasSuffix(f, ".rb") {
			rubyFiles = append(rubyFiles, f)
		}
	}

	// If no Ruby files found, return empty result
	if len(rubyFiles) == 0 {
		return result
	}

	// Parse each Ruby file for require_relative statements
	for _, rubyFile := range rubyFiles {
		filePath := filepath.Join(args.Dir, rubyFile)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}

		// Find all require_relative statements
		matches := requireRelativeRegex.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			if len(match) >= 2 {
				requiredPath := match[1]
				// Convert require_relative path to Bazel label
				depLabel := ConvertRequireRelativeToLabel(args.Rel, requiredPath)
				if depLabel != "" {
					dependencies = append(dependencies, depLabel)
				}
			}
		}
	}

	// Remove duplicates from dependencies
	dependencies = removeDuplicates(dependencies)

	// Generate rb_library rule if we have Ruby files
	if len(rubyFiles) > 0 {
		r := rule.NewRule("rb_library", "")
		r.SetAttr("srcs", rubyFiles)
		if len(dependencies) > 0 {
			r.SetAttr("deps", dependencies)
		}
		result.Gen = append(result.Gen, r)
	}

	return result
}

// ConvertRequireRelativeToLabel converts a require_relative path to a Bazel label
// For example: require_relative "hello_world/version" in package //foo/lib
// becomes //foo/lib/hello_world:version
func ConvertRequireRelativeToLabel(currentPackage, requiredPath string) string {
	// Split the required path into directory and file parts
	dir := filepath.Dir(requiredPath)
	base := filepath.Base(requiredPath)

	// Handle current directory case
	if dir == "." {
		// require_relative "version" -> :version
		return ":" + base
	}

	// Build the full package path
	var packagePath string
	if currentPackage == "" {
		packagePath = "//" + dir
	} else {
		packagePath = "//" + currentPackage + "/" + dir
	}

	// Return the full label
	return packagePath + ":" + base
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
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
