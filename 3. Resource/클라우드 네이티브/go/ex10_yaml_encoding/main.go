package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

func yamlMarshal() {
	c := Config{
		Host: "localhost",
		Port: 8080,
		Tags: map[string]string{"env": "dev"},
	}

	bytes, _ := yaml.Marshal(c)
	fmt.Println("\n# yamlMarshal")
	fmt.Println(string(bytes))
}

func yamlUnmarshal() {
	// 들여쓰기는 탭이 아니라 스페이스 사용
	bytes := []byte(`
host: localhost
port: 1234
tags:
    foo: bar
`)
	c := Config{}
	yaml.Unmarshal(bytes, &c)
	fmt.Println("\n# yamlUnmarshal")
	fmt.Println(c)
}

func main() {
	yamlMarshal()
	yamlUnmarshal()
}
