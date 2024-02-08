package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	BaseAPIURL = "https://www.useragents.me/api"
)

type apiAnswer struct {
	Data []userAgent `json:"data"`
}
type userAgent struct {
	Ua  string  `json:"ua"`
	Pct float64 `json:"-"`
}

func main() {
	out, _ := os.Create(os.Args[1])
	defer out.Close()

	fmt.Fprintln(out, "// DO NOT EDIT")
	fmt.Fprintln(out, "//")
	fmt.Fprintln(out, "//nolint:lll")
	fmt.Fprintln(out, "package scraper")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "var (")
	fmt.Fprintln(out, "\tuserAgents = []string{")

	defer func() {
		fmt.Fprintln(out, "\t}")
		fmt.Fprintln(out, `)`)
	}()

	resp, err := http.Get(BaseAPIURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ua := apiAnswer{}
	_ = json.Unmarshal(data, &ua)

	for i := range ua.Data {
		fmt.Fprintf(out, "\t\t\"%s\",\n", ua.Data[i].Ua)
	}
}
