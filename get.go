package main

import (
	"log"
	"time"

	"github.com/pquerna/otp/totp"

	. "github.com/kepkin/async-script"
)

type GetCmd struct {
	Key string `arg:"positional"`
}

func (cmd *GetCmd) Run(s Store) {
	secret, err := s.Get(cmd.Key)
	if err != nil {
		log.Fatalln(err)
	}

	out, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	MustRun(
		Execf("echo %v", out),
		Execf("pbcopy"),
	)
}
