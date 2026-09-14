package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Restaurant struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Area            string   `json:"area"`
	Category        string   `json:"category"`
	Genre           []string `json:"genre"`
	BudgetYen       int      `json:"budget_yen"`
	Meal            []string `json:"meal"`
	Fish            string   `json:"fish"`
	Freshness       string   `json:"freshness"`
	VisualAppeal    string   `json:"visual_appeal"`
	CostPerformance string   `json:"cost_performance"`
	TargetUse       []string `json:"target_use"`
	Tags            []string `json:"tags"`
	MarketDirect    bool     `json:"market_direct"`
}

type Intent struct {
	BudgetMax int
	Area      string
	Instagram bool
	Visual    bool
	Fish      bool
	Fresh     bool
	Solo      bool
	Dinner    bool
	Date      bool
}

type Result struct {
	Restaurant
	Score  float64
	Reason []string
}

func load() ([]Restaurant, error) {
	b, err := os.ReadFile("data/restaurants.json")
	if err != nil {
		return nil, err
	}
	var rs []Restaurant
	if err := json.Unmarshal(b, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func has(rs []string, q string) bool {
	for _, s := range rs {
		if strings.Contains(s, q) {
			return true
		}
	}
	return false
}

// intent converts a short natural-language request into deterministic
// retrieval constraints. Context words such as "デート" expand into a
// practical restaurant persona rather than matching only the literal word.
func intent(q string) Intent {
	q = strings.ToLower(q)
	i := Intent{}

	for _, n := range []int{1000, 1500, 2000, 2500, 3000, 4000, 5000} {
		if strings.Contains(q, strconv.Itoa(n)) {
			i.BudgetMax = n
			break
		}
	}
	if strings.Contains(q, "川崎") {
		i.Area = "川崎"
	}

	i.Instagram = strings.Contains(q, "インスタ") || strings.Contains(q, "instagram") || strings.Contains(q, "ストーリー") || strings.Contains(q, "映え")
	i.Visual = i.Instagram || strings.Contains(q, "見栄え") || strings.Contains(q, "見た目") || strings.Contains(q, "おしゃれ")
	i.Fish = strings.Contains(q, "魚") || strings.Contains(q, "海鮮") || strings.Contains(q, "刺身") || strings.Contains(q, "寿司") || strings.Contains(q, "鮨")
	i.Fresh = strings.Contains(q, "新鮮") || strings.Contains(q, "鮮度") || strings.Contains(q, "市場") || strings.Contains(q, "直送")
	i.Solo = strings.Contains(q, "一人") || strings.Contains(q, "ひとり")
	i.Dinner = strings.Contains(q, "夜") || strings.Contains(q, "夕食") || strings.Contains(q, "ディナー")
	i.Date = strings.Contains(q, "デート") || strings.Contains(q, "date")

	// "デート" is a composite intent. In the absence of explicit
	// constraints, assume the current washoku persona: affordable dinner,
	// fish-oriented, and visually/socially appealing. Explicit query terms
	// still win; e.g. "デート 5000円" keeps the user's 5000-yen budget.
	if i.Date {
		if i.BudgetMax == 0 {
			i.BudgetMax = 2500
		}
		i.Dinner = true
		i.Fish = true
		i.Visual = true
		i.Instagram = true
	}

	return i
}

func score(q string, r Restaurant) (float64, []string) {
	i := intent(q)
	var s float64
	var reasons []string

	if i.Area != "" && r.Area == i.Area {
		s += 0.5
		reasons = append(reasons, "エリア一致")
	}
	if i.BudgetMax > 0 && r.BudgetYen <= i.BudgetMax {
		s += 0.3
		reasons = append(reasons, fmt.Sprintf("%d円以内", i.BudgetMax))
	}
	if i.Fish {
		if r.Fish == "very_high" {
			s += 1.0
			reasons = append(reasons, "魚評価が非常に高い")
		} else if r.Fish == "high" {
			s += 0.6
			reasons = append(reasons, "魚評価が高い")
		}
	}
	if i.Fresh {
		if r.Freshness == "very_high" {
			s += 1.0
			reasons = append(reasons, "鮮度評価が非常に高い")
		} else if r.Freshness == "high" {
			s += 0.5
			reasons = append(reasons, "鮮度評価が高い")
		}
	}
	if i.Visual {
		if r.VisualAppeal == "high" {
			s += 1.2
			reasons = append(reasons, "見栄えが高い")
		} else if r.VisualAppeal == "medium" {
			s += 0.4
			reasons = append(reasons, "見栄えが中程度")
		}
	}
	if i.Instagram && r.VisualAppeal == "high" {
		s += 0.8
		reasons = append(reasons, "Instagram向き")
	}
	if i.Date {
		if has(r.TargetUse, "date") || has(r.TargetUse, "casual") {
			s += 0.8
			reasons = append(reasons, "デート利用向き")
		}
		if r.VisualAppeal == "high" {
			s += 0.6
			reasons = append(reasons, "デートで映えやすい")
		}
		if has(r.TargetUse, "solo") {
			s -= 0.2
		}
	}
	if i.Solo && has(r.TargetUse, "solo") {
		s += 0.5
		reasons = append(reasons, "一人利用向き")
	}
	if i.Dinner && has(r.Meal, "dinner") {
		s += 0.3
		reasons = append(reasons, "夜利用可能")
	}
	if i.Instagram && r.CostPerformance == "high" {
		s += 0.3
		reasons = append(reasons, "コスパ良")
	}
	if i.Instagram && has(r.Tags, "2000円") {
		s += 0.1
	}
	return s, reasons
}

func normalize(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, s)
}

func main() {
	limit := flag.Int("limit", 5, "number of results")
	budget := flag.Int("budget-max", 0, "maximum budget in yen")
	area := flag.String("area", "", "area filter")
	jsonOut := flag.Bool("json", false, "JSON output")
	flag.Parse()

	query := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if query == "" {
		fmt.Fprintln(os.Stderr, "usage: washoku [--limit N] [--budget-max YEN] [--area AREA] QUERY")
		os.Exit(2)
	}

	q := normalize(query)
	rs, err := load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	results := make([]Result, 0, len(rs))
	for _, r := range rs {
		if *budget > 0 && r.BudgetYen > *budget {
			continue
		}
		if *area != "" && r.Area != *area {
			continue
		}
		s, reasons := score(q, r)
		if s > 0 {
			results = append(results, Result{Restaurant: r, Score: s, Reason: reasons})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Name < results[j].Name
	})
	if *limit < len(results) {
		results = results[:*limit]
	}

	if *jsonOut {
		json.NewEncoder(os.Stdout).Encode(results)
		return
	}
	if len(results) == 0 {
		fmt.Println("該当する候補が見つかりませんでした。")
		return
	}

	i := intent(q)
	fmt.Printf("意図: %s\n", query)
	if i.Date {
		fmt.Println("複合解釈: デート → 予算2500円以内・夜・魚・映え")
	}
	fmt.Println()
	for n, r := range results {
		fmt.Printf("%d. %s [%.2f] — 約%d円 / %s\n", n+1, r.Name, r.Score, r.BudgetYen, r.Category)
		if len(r.Reason) > 0 {
			fmt.Printf("   → %s\n", strings.Join(r.Reason, "・"))
		}
	}
}
