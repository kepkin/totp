package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v3"
)

type Store interface {
	Get(key string) (string, error)
	Put(key, value string) error
}

type yamlStore struct {
	path    string
	Default string            `yaml:"default"`
	Secrets map[string]string `yaml:"secrets"`
}

func NewStoreFromYaml(path string) (*yamlStore, error) {
	res := yamlStore{
		path:    path,
		Secrets: map[string]string{},
	}

	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// skip
		return &res, nil
	} else if err != nil {
		return &res, err
	}

	err = yaml.Unmarshal(content, &res)
	return &res, err
}

func (ys *yamlStore) Get(key string) (string, error) {
	if len(key) == 0 {
		key = ys.Default
	}
	v, ok := ys.Secrets[key]
	if !ok {
		return v, fmt.Errorf("key doesn't exists")
	}

	return v, nil
}

func (ys *yamlStore) Dump() error {
	f, err := os.Create(ys.path)
	if err != nil {
		return fmt.Errorf("can not dump store %w", err)
	}
	defer f.Close()

	content, err := yaml.Marshal(ys)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	return err
}

func (ys *yamlStore) Put(k, v string) error {
	ys.Secrets[k] = v
	return ys.Dump()
}

func main() {
	var args struct {
		Add    *AddCmd    `arg:"subcommand:add"`
		Get    *GetCmd    `arg:"subcommand:get"`
		Decode *DecodeCmd `arg:"subcommand:decode"`
	}

	arg.MustParse(&args)

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	store, err := NewStoreFromYaml(cfgDir + "/totp.yaml")
	if err != nil {
		log.Fatalf("Failed to init store: %v", err)
	}

	if args.Add != nil {
		args.Add.Run(store)
	} else if args.Get != nil {
		args.Get.Run(store)
	} else if args.Decode != nil {
		args.Decode.Run()
	} else {
		args.Get.Run(store)
	}
}
