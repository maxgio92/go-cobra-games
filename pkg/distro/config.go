package distro

import (
	"net/url"
	"strings"
)

func (c *Centos) buildConfig(def Config, user Config) (Config, error) {
	config, err := c.mergeConfig(def, user)
	if err != nil {
		return Config{}, err
	}

	err = c.sanitizeConfig(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// Returns the final configuration by merging the default with the user provided.
//nolint:unparam
func (c *Centos) mergeConfig(def Config, config Config) (Config, error) {
	if len(config.Archs) < 1 {
		config.Archs = def.Archs
	} else {
		for _, arch := range config.Archs {
			if arch == "" {
				config.Archs = def.Archs

				break
			}
		}
	}

	//nolint:nestif
	if len(config.Mirrors) < 1 {
		config.Mirrors = def.Mirrors
	} else {
		for _, mirror := range config.Mirrors {
			if mirror.URL == "" {
				config.Mirrors = def.Mirrors

				break
			}
		}
	}

	if len(config.Repositories) < 1 {
		config.Repositories = c.getDefaultRepositories()
	} else {
		for _, repository := range config.Repositories {
			if repository.URI == "" {
				config.Repositories = c.getDefaultRepositories()

				break
			}
		}
	}

	return config, nil
}

// Returns the final configuration by overriding the default.
//nolint:unparam,unused
func (c *Centos) overrideConfig(def Config, override Config) (Config, error) {
	if len(override.Mirrors) > 0 {
		if override.Mirrors[0].URL != "" {
			return override, nil
		}
	}

	return def, nil
}

func (c *Centos) sanitizeConfig(config *Config) error {
	err := c.sanitizeMirrors(&config.Mirrors)
	if err != nil {
		return err
	}

	return nil
}

func (c *Centos) sanitizeMirrors(mirrors *[]Mirror) error {
	for i, mirror := range *mirrors {
		if !strings.HasSuffix(mirror.URL, "/") {
			(*mirrors)[i].URL = mirror.URL + "/"
		}

		_, err := url.Parse(mirror.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
