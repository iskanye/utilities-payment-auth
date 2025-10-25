package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/iskanye/utilities-payment-auth/internal/storage"
	"github.com/iskanye/utilities-payment/pkg/config"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Admins []Admin `yaml:"admins"`
}

type Admin struct {
	Email    string `yaml:"email" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
}

func main() {
	var storagePath, adminsPath string

	flag.StringVar(&storagePath, "storage", "", "path to storage")
	flag.StringVar(&adminsPath, "admins", "", "path to admins credentials")
	flag.Parse()

	cfg := config.MustLoadPath[Config](adminsPath, config.NoModyfing)

	s, err := storage.New(storagePath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	for _, i := range cfg.Admins {
		passHash, err := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		userId, err := s.SaveUser(ctx, i.Email, passHash, true)
		if err != nil {
			fmt.Println("Admin with that credentials already in database:", i.Email)
			continue
		}

		fmt.Println("New admin id ", userId)
	}
}
