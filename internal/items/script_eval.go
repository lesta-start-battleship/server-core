package items

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func replaceSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "\"")
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	}
	return 0, false
}

func evalExpr(expr string, params map[string]interface{}, prevRand float64) (interface{}, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, nil
	}

	expr = strings.TrimPrefix(expr, "$")
	if n, err := strconv.ParseFloat(expr, 64); err == nil {
		return n, nil
	}

	if val, ok := params[expr]; ok {
		return val, nil
	}

	if expr == "FIELD_SIZE" {
		if val, ok := params["FIELD_SIZE"]; ok {
			return val, nil
		}
		return 10, nil
	}

	if strings.Contains(expr, "RAND") {
		randRe := regexp.MustCompile(`\{\s*\"RAND\"\s*:\s*\"None\"\s*\}`)
		if randRe.MatchString(expr) {
			val := float64(rand.Intn(10))
			return val, nil
		}
	}

	if strings.Contains(expr, "PREV_RAND") {
		prevRandRe := regexp.MustCompile(`\{\s*\"PREV_RAND\"\s*:\s*\"None\"\s*\}`)
		if prevRandRe.MatchString(expr) {
			return prevRand, nil
		}
	}

	if strings.Contains(expr, "+") || strings.Contains(expr, "-") {
		repl := func(s string) string {
			s = strings.TrimSpace(s)
			s = strings.TrimPrefix(s, "$")
			if n, err := strconv.ParseFloat(s, 64); err == nil {
				return fmt.Sprintf("%f", n)
			}

			if val, ok := params[s]; ok {
				return fmt.Sprintf("%v", val)
			}

			if s == "FIELD_SIZE" {
				if val, ok := params["FIELD_SIZE"]; ok {
					return fmt.Sprintf("%v", val)
				}
				return "10"
			}
			return s
		}

		re := regexp.MustCompile(`([+-])`)
		parts := re.Split(expr, -1)
		ops := re.FindAllString(expr, -1)
		var nums []float64

		for _, p := range parts {
			p = strings.TrimSpace(p)
			p = strings.TrimPrefix(p, "$")

			if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
				v, _ := evalExpr(p, params, prevRand)
				if f, ok := toFloat(v); ok {
					nums = append(nums, f)
					continue
				}
			}

			n, err := strconv.ParseFloat(repl(p), 64)
			if err != nil {
				return nil, fmt.Errorf("cannot parse %s", p)
			}

			nums = append(nums, n)
		}
		if len(nums) == 0 {
			return nil, fmt.Errorf("cannot parse expr: %s", expr)
		}

		res := nums[0]
		for i, op := range ops {
			if op == "+" {
				res += nums[i+1]
			}

			if op == "-" {
				res -= nums[i+1]
			}
		}

		return res, nil
	}
	if strings.HasPrefix(expr, "{") && strings.HasSuffix(expr, "}") {
		if strings.Contains(expr, "RAND") {
			val := float64(rand.Intn(10))

			return val, nil
		}
		if strings.Contains(expr, "PREV_RAND") {
			return prevRand, nil
		}
	}
	return expr, nil
}
