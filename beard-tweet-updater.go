package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const defaultRootDir = "/appdata"

func makeRequest(url string) ([]byte, error) {
	fmt.Printf("Making request to: [%s]\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func writeFile(rootDir string) {
	ioutil.WriteFile(rootDir+"/last_run_date", []byte(time.Now().Format("2006-01-02-150405")+"\n"), 0644)
}

func main() {

	var rootDir = defaultRootDir
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	writeFile(rootDir)
}
