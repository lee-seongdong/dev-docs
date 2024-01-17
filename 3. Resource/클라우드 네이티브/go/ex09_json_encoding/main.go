package main

import (
	"encoding/json"
	"fmt"
)

// 마샬링/언마샬링 하려면 패키지 외부로 보낼 수 있도록 public 타입으로 작성해야 한다.
// 타입뒤에 백틱으로 메타데이터를 추가할 수 있다.
type Config struct {
	Host string `json:"host,omitempty"`
	Port uint16
	Tags map[string]string
}

func jsonMarshal() {
	c := Config{
		Host: "localhost",
		Port: 8080,
		Tags: map[string]string{"env": "dev"},
	}

	bytes, _ := json.Marshal(c)
	fmt.Println("\n# jsonMarshal")
	fmt.Println(string(bytes))
}

func jsonPrettyMarshal() {
	c := Config{
		Host: "localhost",
		Port: 8080,
		Tags: map[string]string{"env": "dev"},
	}

	bytes, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println("\n# jsonPrettyMarshal")
	fmt.Println(string(bytes))
}

func jsonUnmarshal() {
	c := Config{}
	bytes := []byte(`{"Host":"localhost","Port":8080,"Tags":{"env":"dev"}}`)
	json.Unmarshal(bytes, &c)
	fmt.Println("\n# jsonUnmarshal")
	fmt.Println(c)
}

func jsonUnmarshalCherryPick() {
	c := Config{}
	bytes := []byte(`{"Host":"localhost","Food": "Pizza"}`)
	json.Unmarshal(bytes, &c)
	fmt.Println("\n# jsonUnmarshalCherryPick")
	fmt.Println(c)
}

func jsonUnmarshalAny() {
	var c any
	bytes := []byte(`{"Host":"localhost","Food": "Pizza"}`)
	json.Unmarshal(bytes, &c)
	fmt.Println("\n# jsonUnmarshalAny")
	fmt.Println(c)
}

func jsonUnmarshalMapAny() {
	var c map[string]any
	bytes := []byte(`{"Host":"localhost","Food": "Pizza"}`)
	json.Unmarshal(bytes, &c)
	fmt.Println("\n# jsonUnmarshalMapAny")
	fmt.Println(c)
}

func main() {
	jsonMarshal()
	jsonPrettyMarshal()
	jsonUnmarshal()
	jsonUnmarshalCherryPick()
	jsonUnmarshalAny()
	jsonUnmarshalMapAny()
}
