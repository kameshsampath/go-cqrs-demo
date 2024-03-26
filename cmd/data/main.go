package main

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/go-resty/resty/v2"
	"github.com/kameshsampath/go-cqrs-demo/config"
)

func main() {
	log := config.Log
	url := fmt.Sprintf("http://localhost:%s", os.Getenv("APP_PORT"))
	cwd, _ := os.Getwd()
	f, err := os.OpenFile(path.Join(cwd, "samples", "todos.jsonl"), os.O_RDONLY, os.ModePerm)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		todoJSON := scanner.Text()
		client := resty.New()
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(todoJSON).
			EnableTrace().
			Post(url)
		if err != nil || res.StatusCode() != 201 {
			log.Errorf("error creating %s,%v", todoJSON, err)
		}
		log.Infof("Created record:%s", string(res.Body()))
	}

}
