package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/tjhorner/fs/shorty"
)

type config struct {
	ProjectID    string         `json:"projectId"`
	BucketName   string         `json:"bucketName"`
	Host         string         `json:"host"`
	Shorten      bool           `json:"shorten"`
	ShortyConfig *shorty.Config `json:"shorty"`
}

func newConfig() *config {
	return &config{
		ProjectID:  "",
		BucketName: "",
		Host:       "",
		Shorten:    false,
		ShortyConfig: &shorty.Config{
			BaseURL: "",
		},
	}
}

func defaultConfigDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	return path.Join(home, ".config", "fs")
}

func defaultConfigPath() string {
	return path.Join(defaultConfigDirectory(), "config.json")
}

func defaultServiceAccountPath() string {
	return path.Join(defaultConfigDirectory(), "service_account.json")
}

func loadConfig(path string) *config {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		saveConfig(newConfig(), path)
		return newConfig()
	}

	var conf config
	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return newConfig()
	}

	return &conf
}

func saveConfig(conf *config, fp string) error {
	enc, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	os.MkdirAll(path.Dir(fp), os.ModePerm)
	err = ioutil.WriteFile(fp, enc, 0644)
	if err != nil {
		return err
	}

	return nil
}
