package items

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/transaction"
)

func RunScript(script string, state *game.States, params map[string]any) (string, error) {
	if script == "" {
		return "", nil
	}

	var scriptObj struct {
		Input   string        `json:"input"`
		Actions []interface{} `json:"actions"`
	}
	if err := json.Unmarshal([]byte(script), &scriptObj); err != nil {
		return "", err
	}

	tx := transaction.NewTransaction()
	var lastRand int

	resolveIntWithCtx := func(val interface{}, params map[string]any) (int, bool) {
		return resolveIntWithRand(val, params, &lastRand)
	}

	for _, actRaw := range scriptObj.Actions {
		var actMap map[string]interface{}
		if m, ok := actRaw.(map[string]interface{}); ok {
			actMap = m
		} else {
			b, _ := json.Marshal(actRaw)
			json.Unmarshal(b, &actMap)
		}

		for k, v := range actMap {
			var actionName string
			var args map[string]interface{}
			if k == "Name" {
				actionName = v.(string)
				if a, ok := actMap["Args"]; ok {
					args, _ = a.(map[string]interface{})
				}
			} else {
				actionName = k
				if a, ok := v.(map[string]interface{}); ok {
					args = a
				}
			}

			switch actionName {
			case "OPEN_CELL":
				x, _ := resolveIntWithCtx(args["x"], params)
				y, _ := resolveIntWithCtx(args["y"], params)
				if x < 0 || x >= 10 || y < 0 || y >= 10 {
					fmt.Printf("[OPEN_CELL] SKIP: out of bounds x=%d y=%d\n", x, y)
					continue
				}
				fmt.Printf("[OPEN_CELL] x=%d y=%d\n", x, y)
				cmd := game.NewOpenCellCommand(game.Coord{X: x, Y: y})
				tx.Add(cmd)
			case "SET_CELL_STATUS":
				x, _ := resolveIntWithCtx(args["x"], params)
				y, _ := resolveIntWithCtx(args["y"], params)
				status, _ := args["status"].(string)
				cmd := &setCellStatusCommand{X: x, Y: y, Status: status}
				tx.Add(cmd)
			case "END_PLAYER_ACTION":
				// no-op or handle as needed
			case "REMOVE_SHIP":
				x, _ := resolveIntWithCtx(args["x"], params)
				y, _ := resolveIntWithCtx(args["y"], params)
				cmd := game.NewRemoveShipCommand(game.Coord{X: x, Y: y})
				tx.Add(cmd)
			case "PLACE_SHIP":
				x, _ := resolveIntWithCtx(args["x"], params)
				y, _ := resolveIntWithCtx(args["y"], params)
				cmd := game.NewPlaceShipCommand(1, game.Coord{X: x, Y: y}, false) // TODO: get len/bearings from args
				tx.Add(cmd)
			case "HEAL_SHIP":
				x, _ := resolveIntWithCtx(args["x"], params)
				y, _ := resolveIntWithCtx(args["y"], params)
				cmd := game.NewHealShipCommand(game.Coord{X: x, Y: y})
				tx.Add(cmd)
			case "SHOOT":
				x, _ := resolveIntWithCtx(args["x"], params)
				y, _ := resolveIntWithCtx(args["y"], params)
				cmd := game.NewShootCommand(game.Coord{X: x, Y: y})
				tx.Add(cmd)
			case "SWITCH_CASE", "SWICH_CASE":
				caseKey := "1"
				if dir, ok := params["direction"]; ok {
					caseKey = fmt.Sprintf("%v", dir)
				}
				fmt.Printf("[SWITCH_CASE] direction=%v, caseKey=%s\n", params["direction"], caseKey)
				if caseVal, ok := args[caseKey]; ok {
					fmt.Printf("[SWITCH_CASE] caseVal type=%T, value=%#v\n", caseVal, caseVal)
					if arr, ok := caseVal.([]interface{}); ok {
						fmt.Printf("[SWITCH_CASE] arr len=%d\n", len(arr))
						for idx, subAct := range arr {
							fmt.Printf("[SWITCH_CASE] subAct[%d] type=%T, value=%#v\n", idx, subAct, subAct)
							subMap, ok := subAct.(map[string]interface{})
							if !ok {
								fmt.Printf("[SWITCH_CASE] subAct[%d] is not map[string]interface{}\n", idx)
								continue
							}
							subActionName, ok := subMap["Name"].(string)
							if !ok {
								fmt.Printf("[SWITCH_CASE] subMap has no string Name, value=%#v\n", subMap["Name"])
								continue
							}
							subArgs, ok := subMap["Args"].(map[string]interface{})
							if !ok {
								fmt.Printf("[SWITCH_CASE] subMap has no map Args, value=%#v\n", subMap["Args"])
								continue
							}
							fmt.Printf("[SWITCH_CASE] subActionName=%s, subArgs=%#v\n", subActionName, subArgs)
							switch subActionName {
							case "OPEN_CELL":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								if x < 0 || x >= 10 || y < 0 || y >= 10 {
									fmt.Printf("[OPEN_CELL] SKIP: out of bounds x=%d y=%d\n", x, y)
									continue
								}
								fmt.Printf("[OPEN_CELL] x=%d y=%d\n", x, y)
								cmd := game.NewOpenCellCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							case "SET_CELL_STATUS":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								status, _ := subArgs["status"].(string)
								cmd := &setCellStatusCommand{X: x, Y: y, Status: status}
								tx.Add(cmd)
							case "END_PLAYER_ACTION":
								// no-op or handle as needed
							case "REMOVE_SHIP":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewRemoveShipCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							case "PLACE_SHIP":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewPlaceShipCommand(1, game.Coord{X: x, Y: y}, false) // TODO: get len/bearings from args
								tx.Add(cmd)
							case "HEAL_SHIP":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewHealShipCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							case "SHOOT":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewShootCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							}
						}
					} else if subMap, ok := caseVal.(map[string]interface{}); ok {
						for subK, subV := range subMap {
							subActionName := subK
							subArgs, _ := subV.(map[string]interface{})
							switch subActionName {
							case "OPEN_CELL":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								if x < 0 || x >= 10 || y < 0 || y >= 10 {
									fmt.Printf("[OPEN_CELL] SKIP: out of bounds x=%d y=%d\n", x, y)
									continue
								}
								fmt.Printf("[OPEN_CELL] x=%d y=%d\n", x, y)
								cmd := game.NewOpenCellCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							case "SET_CELL_STATUS":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								status, _ := subArgs["status"].(string)
								cmd := &setCellStatusCommand{X: x, Y: y, Status: status}
								tx.Add(cmd)
							case "END_PLAYER_ACTION":
								// no-op or handle as needed
							case "REMOVE_SHIP":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewRemoveShipCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							case "PLACE_SHIP":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewPlaceShipCommand(1, game.Coord{X: x, Y: y}, false) // TODO: get len/bearings from args
								tx.Add(cmd)
							case "HEAL_SHIP":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewHealShipCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							case "SHOOT":
								x, _ := resolveIntWithCtx(subArgs["x"], params)
								y, _ := resolveIntWithCtx(subArgs["y"], params)
								cmd := game.NewShootCommand(game.Coord{X: x, Y: y})
								tx.Add(cmd)
							}
						}
					}
				}
			}
		}
	}

	err := tx.Execute(state)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

var simpleExprRe = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*([+-])\s*(\d+)$`)

func evalSimpleExpr(expr string, params map[string]any) (int, bool) {
	m := simpleExprRe.FindStringSubmatch(strings.ReplaceAll(expr, " ", ""))
	if len(m) == 4 {
		key := m[1]
		op := m[2]
		delta, _ := strconv.Atoi(m[3])
		base, ok := resolveInt(key, params)
		if !ok {
			return 0, false
		}
		if op == "+" {
			return base + delta, true
		}
		return base - delta, true
	}
	return 0, false
}

func resolveInt(val interface{}, params map[string]any) (int, bool) {
	switch v := val.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case string:
		if strings.HasPrefix(v, "$") {
			return resolveInt(v[1:], params)
		}
		if p, ok := params[v]; ok {
			return resolveInt(p, params)
		}
		if n, ok := evalSimpleExpr(v, params); ok {
			return n, true
		}
		// Новый универсальный парсер выражений с несколькими слагаемыми
		if n, ok := evalComplexExpr(v, params); ok {
			return n, true
		}
	}
	if m, ok := val.(map[string]interface{}); ok {
		if name, ok := m["Name"]; ok && name == "RAND" {
			rand.Seed(time.Now().UnixNano())
			return rand.Intn(10), true // 0..9
		}
	}
	if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
		return resolveInt(arr[0], params)
	}
	return 0, false
}

// Универсальный парсер выражений с + и - (например, "{'Name': 'RAND', ...} - $FIELD_SIZE + $x")
func evalComplexExpr(expr string, params map[string]any) (int, bool) {
	tokens := splitExprTokens(expr)
	if len(tokens) == 0 {
		return 0, false
	}
	res, ok := resolveIntToken(tokens[0], params)
	if !ok {
		return 0, false
	}
	for i := 1; i < len(tokens)-1; i += 2 {
		op := tokens[i]
		val, ok := resolveIntToken(tokens[i+1], params)
		if !ok {
			return 0, false
		}
		if op == "+" {
			res += val
		} else if op == "-" {
			res -= val
		}
	}
	return res, true
}

// Разделяет выражение на токены: операнды и операторы
func splitExprTokens(expr string) []string {
	expr = strings.ReplaceAll(expr, " ", "")
	tokens := []string{}
	buf := ""
	inObj := 0
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if c == '{' {
			inObj++
		}
		if c == '}' {
			inObj--
		}
		if (c == '+' || c == '-') && inObj == 0 {
			if buf != "" {
				tokens = append(tokens, buf)
				buf = ""
			}
			tokens = append(tokens, string(c))
			continue
		}
		buf += string(c)
	}
	if buf != "" {
		tokens = append(tokens, buf)
	}
	return tokens
}

// Преобразует токен в значение int (учитывает объекты, переменные, числа)
func resolveIntToken(token string, params map[string]any) (int, bool) {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(token, "$") {
		return resolveInt(token[1:], params)
	}
	if n, ok := strconv.Atoi(token); ok == nil {
		return n, true
	}
	if p, ok := params[token]; ok {
		return resolveInt(p, params)
	}
	// Попытка распарсить как JSON-объект (например, {'Name': 'RAND', ...})
	if strings.HasPrefix(token, "{") && strings.HasSuffix(token, "}") {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(strings.ReplaceAll(strings.ReplaceAll(token, "'", "\""), " ", "")), &m)
		if err == nil {
			return resolveInt(m, params)
		}
	}
	return 0, false
}

type setCellStatusCommand struct {
	X, Y   int
	Status string
}

func (c *setCellStatusCommand) Apply(states *game.States) error {
	gs := states.PlayerState
	if !gs.IsInside(c.X, c.Y) {
		return nil
	}
	var s int
	switch c.Status {
	case "open":
		s = game.Open
	case "close":
		s = game.Close
	default:
		s = game.Open
	}
	gs.Field[c.X][c.Y].State = s
	return nil
}

func (c *setCellStatusCommand) Undo(states *game.States) {}

// resolveInt с поддержкой lastRand и FIELD_SIZE
func resolveIntWithRand(val interface{}, params map[string]any, lastRand *int) (int, bool) {
	switch v := val.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case string:
		if v == "FIELD_SIZE" {
			return 9, true
		}
		if v == "PREV_RAND" {
			return *lastRand, true
		}
		if strings.HasPrefix(v, "$") {
			return resolveIntWithRand(v[1:], params, lastRand)
		}
		if p, ok := params[v]; ok {
			return resolveIntWithRand(p, params, lastRand)
		}
		if n, ok := evalSimpleExpr(v, params); ok {
			return n, true
		}
		if n, ok := evalComplexExprWithRand(v, params, lastRand); ok {
			return n, true
		}
		// Новый блок: если строка выглядит как объект, пробуем парсить
		if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			jsonStr := strings.ReplaceAll(strings.ReplaceAll(v, "'", "\""), " ", "")
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &m); err == nil {
				return resolveIntWithRand(m, params, lastRand)
			}
		}
	}
	if m, ok := val.(map[string]interface{}); ok {
		if name, ok := m["Name"]; ok && name == "RAND" {
			rand.Seed(time.Now().UnixNano())
			*lastRand = rand.Intn(10)
			return *lastRand, true
		}
	}
	if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
		return resolveIntWithRand(arr[0], params, lastRand)
	}
	return 0, false
}

// evalComplexExprWithRand аналогична evalComplexExpr, но с поддержкой lastRand
func evalComplexExprWithRand(expr string, params map[string]any, lastRand *int) (int, bool) {
	tokens := splitExprTokens(expr)
	if len(tokens) == 0 {
		return 0, false
	}
	res, ok := resolveIntTokenWithRand(tokens[0], params, lastRand)
	if !ok {
		return 0, false
	}
	for i := 1; i < len(tokens)-1; i += 2 {
		op := tokens[i]
		val, ok := resolveIntTokenWithRand(tokens[i+1], params, lastRand)
		if !ok {
			return 0, false
		}
		if op == "+" {
			res += val
		} else if op == "-" {
			res -= val
		}
	}
	return res, true
}

// resolveIntTokenWithRand аналогична resolveIntToken, но с поддержкой lastRand
func resolveIntTokenWithRand(token string, params map[string]any, lastRand *int) (int, bool) {
	token = strings.TrimSpace(token)
	if token == "FIELD_SIZE" {
		return 9, true
	}
	if token == "PREV_RAND" {
		return *lastRand, true
	}
	if strings.HasPrefix(token, "$") {
		return resolveIntWithRand(token[1:], params, lastRand)
	}
	if n, ok := strconv.Atoi(token); ok == nil {
		return n, true
	}
	if p, ok := params[token]; ok {
		return resolveIntWithRand(p, params, lastRand)
	}
	if strings.HasPrefix(token, "{") && strings.HasSuffix(token, "}") {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(strings.ReplaceAll(strings.ReplaceAll(token, "'", "\""), " ", "")), &m)
		if err == nil {
			return resolveIntWithRand(m, params, lastRand)
		}
	}
	return 0, false
}
