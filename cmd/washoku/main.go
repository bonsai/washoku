package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

type Restaurant struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Area           string   `json:"area"`
	Category       string   `json:"category"`
	Genre          []string `json:"genre"`
	BudgetYen      int      `json:"budget_yen"`
	Meal           []string `json:"meal"`
	Fish           string   `json:"fish"`
	Freshness      string   `json:"freshness"`
	VisualAppeal   string   `json:"visual_appeal"`
	CostPerformance string  `json:"cost_performance"`
	TargetUse      []string `json:"target_use"`
	Tags           []string `json:"tags"`
	MarketDirect   bool     `json:"market_direct"`
}

type Result struct {
	Restaurant
	Score float64
}

func load() ([]Restaurant, error) {
	b, err := os.ReadFile("data/restaurants.json")
	if err != nil { return nil, err }
	var rs []Restaurant
	if err := json.Unmarshal(b, &rs); err != nil { return nil, err }
	return rs, nil
}

func contains(rs []string, q string) bool {
	for _, s := range rs { if strings.Contains(s, q) { return true } }
	return false
}

func score(q string, r Restaurant) float64 {
	q = strings.ToLower(q)
	text := strings.ToLower(strings.Join(append([]string{r.Name, r.Area, r.Category, r.Fish, r.Freshness, r.VisualAppeal}, append(r.Genre, r.Tags...)...), " "))
	var s float64
	for _, term := range tokenize(q) {
		if strings.Contains(text, term) { s += 1 }
	}
	if strings.Contains(q, "川崎") && r.Area == "川崎" { s += 0.25 }
	if (strings.Contains(q, "魚") || strings.Contains(q, "海鮮") || strings.Contains(q, "刺身") || strings.Contains(q, "寿司") || strings.Contains(q, "鮨")) && (r.Fish == "high" || r.Fish == "very_high") { s += 0.15 }
	if strings.Contains(q, "新鮮") || strings.Contains(q, "鮮度") || strings.Contains(q, "市場") || strings.Contains(q, "直送") {
		if r.Freshness == "very_high" { s += 0.2 }
	}
	if strings.Contains(q, "見栄え") || strings.Contains(q, "映え") || strings.Contains(q, "インスタ") || strings.Contains(q, "見た目") {
		if r.VisualAppeal == "high" { s += 0.15 }
	}
	if strings.Contains(q, "一人") && contains(r.TargetUse, "solo") { s += 0.1 }
	return s
}

func tokenize(s string) []string {
	var out []string
	for _, r := range []rune(s) {
		if r <= utf8.RuneSelf && (r >= 'a' && r <= 'z' || r >= '0' && r <= '9') { out = append(out, string(r)) }
	}
	// Japanese matching is primarily handled by direct substring checks above.
	if len(out) == 0 { return []string{s} }
	return out
}

func main() {
	limit := flag.Int("limit", 5, "number of results")
	budget := flag.Int("budget-max", 0, "maximum budget in yen")
	area := flag.String("area", "", "area filter")
	jsonOut := flag.Bool("json", false, "JSON output")
	flag.Parse()
	query := strings.Join(flag.Args(), " ")
	if query == "" { fmt.Fprintln(os.Stderr, "usage: washoku [--limit N] [--budget-max YEN] [--area AREA] QUERY"); os.Exit(2) }
	rs, err := load(); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	results := make([]Result, 0, len(rs))
	for _, r := range rs {
		if *budget > 0 && r.BudgetYen > *budget { continue }
		if *area != "" && r.Area != *area { continue }
		s := score(query, r); if s > 0 { results = append(results, Result{Restaurant:r, Score:s}) }
	}
	sort.Slice(results, func(i, j int) bool { if results[i].Score != results[j].Score { return results[i].Score > results[j].Score }; return results[i].Name < results[j].Name })
	if *limit < len(results) { results = results[:*limit] }
	if *jsonOut { json.NewEncoder(os.Stdout).Encode(results); return }
	if len(results) == 0 { fmt.Println("該当する候補が見つかりませんでした。"); return }
	for i, r := range results { fmt.Printf("%d. %s [%.2f] — 約%d円 / %s / 魚:%s 鮮度:%s\n", i+1, r.Name, r.Score, r.BudgetYen, r.Category, r.Fish, r.Freshness) }
}
