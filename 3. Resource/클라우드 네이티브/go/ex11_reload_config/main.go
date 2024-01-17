package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

var config Config

// 설정 파일 변경 시 Config 구조체 업데이트
func loadConfig(filepath string) (Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func startListening(updates <-chan string, errors <-chan error) {
	for {
		select {
		case filepath := <-updates:
			c, err := loadConfig(filepath)
			if err != nil {
				log.Println("error loading config:", err)
			} else {
				config = c
			}
		case err := <-errors:
			log.Println("error watching config:", err)
		}
	}
}

func init() {
	updates, errors, err := watchConfig("config.yaml") // or watchConfigNotify("config.yaml")

	if err != nil {
		panic(err)
	}

	go startListening(updates, errors)
}

// 1. 폴링으로 주기적으로 변경 감시
func calculateFileHash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return sum, nil
}

func watchConfig(filepath string) (<-chan string, <-chan error, error) {
	changes := make(chan string)
	errs := make(chan error)
	hash := ""

	go func() {
		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			newHash, err := calculateFileHash(filepath)
			if err != nil {
				errs <- err
			} else if hash != newHash {
				hash = newHash
				changes <- filepath
			}
		}
	}()

	return changes, errs, nil
}

// 2. os 파일시스템 노티피케이션 감시하기
func watchConfigNotify(filepath string) (<-chan string, <-chan error, error) {
	changes := make(chan string)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	err = watcher.Add(filepath)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		changes <- filepath

		for event := range watcher.Events {
			if event.Op&fsnotify.Write == fsnotify.Write {
				changes <- event.Name
			}
		}
	}()

	return changes, watcher.Errors, nil
}

func main() {

}
