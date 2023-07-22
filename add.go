package main

import (
	"log"
)

type AddCmd struct {
	Key   string `arg:"positional"`
	Value string `arg:"positional"`
}

func (cmd *AddCmd) Run(s Store) {
	err := s.Put(cmd.Key, cmd.Value)
	if err != nil {
		log.Fatalln(err)
	}
}
