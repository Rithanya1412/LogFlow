package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type LogEvent struct {
	Service string
	Level   string
	Message string
}

func main() {
	f, err := os.Open("app.log")
	if err != nil {
		fmt.Println("error opening file", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		evt := LogEvent{}

		for i, p := range parts {
			if strings.HasPrefix(p, "message=") {
				msgPart := strings.Join(parts[i:], " ")
				val := strings.TrimPrefix(msgPart, "message=")
				if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
					val = val[1 : len(val)-1]
				}
				evt.Message = val
				break
			}

			kv := strings.SplitN(p, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key, val := kv[0], kv[1]

			switch key {
			case "service":
				evt.Service = val
			case "level":
				evt.Level = val
			}
		}
		data, err := json.Marshal(evt)
		if err != nil {
			fmt.Println("error marshaling:", err)
			continue
		}

		resp, err := http.Post("http://localhost:8080/ingest", "application/json", bytes.NewReader(data))
		if err != nil {
			fmt.Println("error posting:", err)
			continue
		}
		resp.Body.Close()

		fmt.Println("SENT:", string(data))

		// TODO: send evt as JSON to http://localhost:8080/ingest
	}
}
