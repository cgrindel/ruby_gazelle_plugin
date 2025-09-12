package rubycfg

import (
	"github.com/bazelbuild/bazel-gazelle/config"
)

const rubyConfigKey = "_ruby_config"

// RubyConfig represents Ruby-specific configuration for Gazelle.
type RubyConfig struct {
	// BundleRepoName specifies the name of the bundle repository to use
	// for resolving require statements to gem dependencies.
	BundleRepoName string
}

// NewRubyConfig creates a new RubyConfig with default values.
func NewRubyConfig() *RubyConfig {
	return &RubyConfig{
		BundleRepoName: "bundle",
	}
}

// GetRubyConfig retrieves the RubyConfig from the provided config.Config.
// If no RubyConfig exists, it returns a new one with default values.
func GetRubyConfig(c *config.Config) *RubyConfig {
	if c == nil {
		return NewRubyConfig()
	}

	if rubyConfig, ok := c.Exts[rubyConfigKey]; ok {
		if rc, ok := rubyConfig.(*RubyConfig); ok {
			return rc
		}
	}

	return NewRubyConfig()
}

// SetRubyConfig stores the RubyConfig in the provided config.Config.
func SetRubyConfig(c *config.Config, rc *RubyConfig) {
	if c == nil {
		return
	}

	if c.Exts == nil {
		c.Exts = make(map[string]any)
	}

	c.Exts[rubyConfigKey] = rc
}

// DeepCopyRubyConfig creates a deep copy of the provided RubyConfig.
func DeepCopyRubyConfig(rc *RubyConfig) *RubyConfig {
	if rc == nil {
		return NewRubyConfig()
	}

	return &RubyConfig{
		BundleRepoName: rc.BundleRepoName,
	}
}
