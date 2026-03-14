---
title: "Viper Configuration Integration"
version: "0.1.0"
status: "Final"
date: "2026-03-15"
last_updated: "2026-03-15"
authors:
  - "raiworks"
supersedes: ""
---

# Viper Configuration Integration

## Abstract

RapidGo ships with a lightweight configuration system: `godotenv` for
`.env` loading and `LoadConfig[T]()` for typed struct binding. This
guide shows how to use `spf13/viper` alongside the built-in system
when you need YAML/TOML config files, remote config, or live reloading.

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [When to Use Viper](#2-when-to-use-viper)
3. [Install Viper](#3-install-viper)
4. [Create a Config Provider](#4-create-a-config-provider)
5. [Loading Config Files](#5-loading-config-files)
6. [Using with LoadConfig[T]](#6-using-with-loadconfigt)
7. [Live Reloading](#7-live-reloading)
8. [References](#8-references)

## 1. Prerequisites

- RapidGo v2.6.0+ (for `LoadConfig[T]`)
- Go 1.21+

## 2. When to Use Viper

| Need | Built-in | Viper |
|------|----------|-------|
| `.env` file loading | `config.Load()` | Overkill |
| Typed struct binding with validation | `LoadConfig[T]()` | Overkill |
| Single env var read | `config.Env()` | Overkill |
| YAML/TOML/JSON config files | Not supported | **Use Viper** |
| Remote config (etcd, Consul) | Not supported | **Use Viper** |
| Live config reloading | Not supported | **Use Viper** |
| Multiple config sources merged | Not supported | **Use Viper** |

For most projects, RapidGo's built-in config is sufficient. Consider
Viper when you need file-based config or remote config stores.

## 3. Install Viper

```bash
go get github.com/spf13/viper
```

## 4. Create a Config Provider

Create `app/providers/viper_provider.go`:

```go
package providers

import (
	"fmt"

	"github.com/raiworks/rapidgo/v2/core/container"
	"github.com/spf13/viper"
)

type ViperProvider struct{}

func (p *ViperProvider) Register(c *container.Container) {
	c.Singleton("viper", func(_ *container.Container) interface{} {
		v := viper.New()

		// Search paths
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.SetConfigName("config")
		v.SetConfigType("yaml")

		// Environment variable override
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				panic(fmt.Sprintf("viper config error: %v", err))
			}
			// Config file not found — rely on env vars and defaults
		}

		return v
	})
}

func (p *ViperProvider) Boot(c *container.Container) {}
```

## 5. Loading Config Files

Create `config/config.yaml` in your project root:

```yaml
app:
  name: "MyApp"
  debug: false

cache:
  driver: "redis"
  ttl: "5m"

features:
  enable_websockets: true
  max_upload_mb: 25
```

Resolve and use in your services:

```go
v := container.MustMake[*viper.Viper](c, "viper")

appName := v.GetString("app.name")          // "MyApp"
debug := v.GetBool("app.debug")             // false
ttl := v.GetDuration("cache.ttl")           // 5m0s
maxUpload := v.GetInt("features.max_upload_mb") // 25
```

Environment variables override file values. Viper maps `APP_NAME` to
`app.name` automatically when `AutomaticEnv()` is enabled.

## 6. Using with LoadConfig[T]

You can use both systems together. RapidGo's `LoadConfig[T]()` handles
env-var-driven config (database, auth, etc.), while Viper handles
file-based config (feature flags, business rules):

```go
// Environment-driven (RapidGo built-in)
dbCfg, err := config.LoadConfig[database.DBConfig]()

// File-driven (Viper)
v := container.MustMake[*viper.Viper](c, "viper")
enableWS := v.GetBool("features.enable_websockets")
```

## 7. Live Reloading

Viper supports watching config files for changes:

```go
func (p *ViperProvider) Boot(c *container.Container) {
	v := container.MustMake[*viper.Viper](c, "viper")
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log := container.MustMake[logger.Logger](c, "logger")
		log.Info("config reloaded", "file", e.Name)
	})
}
```

Install the `fsnotify` dependency:

```bash
go get github.com/fsnotify/fsnotify
```

## 8. References

- [Viper documentation](https://pkg.go.dev/github.com/spf13/viper)
- [RapidGo Configuration](../references/configuration.md)
- [RapidGo Extending the Framework](extending-framework.md)
