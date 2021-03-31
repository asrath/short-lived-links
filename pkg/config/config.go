package config

import (
	"sync"

	"github.com/asrath/short-lived-links/pkg/storage/paste"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var lock = &sync.Mutex{}

type Config struct {
	App struct {
		Title            string
		LogoText         string
		PasteStoragePath string
		Expirations      map[string]string
	}
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		lock.Lock()
		defer lock.Unlock()
		if cfg == nil {
			cfg = &Config{}
			config.WithOptions(config.ParseEnv)

			// add driver for support yaml content
			config.AddDriver(yaml.Driver)
			err := config.LoadFiles("app.yaml")
			if err != nil {
				err := config.LoadFiles("configs/app.yaml")
				if err != nil {
					panic(err)
				}
			}
			cfg.App.Title = config.String("app.title")
			cfg.App.LogoText = config.String("app.logoText")
			cfg.App.PasteStoragePath = config.String("app.pasteStoragePath", paste.DefaultStoragePath)
			cfg.App.Expirations = config.StringMap("app.expirations")
		}
	}

	return cfg
}
