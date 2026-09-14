package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	in, err := os.Open("rag/metadata.jsonl")
	if err != nil { panic(err) }
	defer in.Close()

	var restaurants []json.RawMessage
	s := bufio.NewScanner(in)
	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 { continue }
		var v json.RawMessage
		if err := json.Unmarshal(line, &v); err != nil { panic(err) }
		restaurants = append(restaurants, v)
	}
	if err := s.Err(); err != nil { panic(err) }

	out, err := os.Create("data/restaurants.json")
	if err != nil { panic(err) }
	defer out.Close()
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(restaurants); err != nil { panic(err) }
	fmt.Printf("updated %d RAG records\n", len(restaurants))
}
