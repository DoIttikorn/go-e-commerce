package logger

import (
	"log"
	"os"

	"github.com/Doittikorn/go-e-commerce/config"
)

type StorageImpl interface {
	InitLocalStorage()
	VerifyEnv() StorageImpl
}

// เอาไว้สร้างไฟล์ที่เก็บ log ไว้
type storage struct {
	env     string
	logPath string
	isMkdir bool
}

func NewStorage(env config.AppConfigImpl) StorageImpl {
	return &storage{
		env:     env.Env(),
		logPath: env.LogPath(),
	}
}

func (s *storage) VerifyEnv() StorageImpl {
	if s.env != "production" && s.env != "staging" {
		s.isMkdir = true
	}
	return s
}

func (s storage) InitLocalStorage() {
	if s.isMkdir {
		if _, err := os.Stat(s.logPath); os.IsNotExist(err) {
			if err := os.MkdirAll(s.logPath, 0755); err != nil {
				log.Fatal(err)
			}
		}
	}
}
