package items

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/transaction"
)

type Action struct {
	Name string
	Args map[string]interface{}
}

func RunScript(script string, state *game.States, input ItemInput) ([]ItemEffect, error) {
	if script == "" {
		return nil, nil
	}

	params := map[string]any{
		"x":         input.X,
		"y":         input.Y,
		"direction": input.Direction,
		"item_id":   input.ItemID,
	}

	var scriptObj struct {
		Input   string        `json:"input"`
		Actions []interface{} `json:"actions"`
	}
	if err := json.Unmarshal([]byte(script), &scriptObj); err != nil {
		return nil, err
	}

	tx := transaction.NewTransaction()
	var lastRand int
	var effects []ItemEffect
	var openCmds []*game.OpenCellCommand
	var healCmds []*game.HealShipCommand
	openCmdCount := 0

	resolveIntWithCtx := func(val interface{}, params map[string]any) (int, bool) {
		return resolveIntWithRand(val, params, &lastRand)
	}

	addEffect := func(effectType string, coord game.Coord) {
		for i := range effects {
			if effects[i].Type == effectType {
				effects[i].Coords = append(effects[i].Coords, coord)
				return
			}
		}
		effects = append(effects, ItemEffect{Type: effectType, Coords: []game.Coord{coord}})
	}

	processAction := func(actionName string, args map[string]interface{}) {
		x, _ := resolveIntWithCtx(args["x"], params)
		y, _ := resolveIntWithCtx(args["y"], params)
		coord := game.Coord{X: x, Y: y}

		switch actionName {
		case "OPEN_CELL":
			if !state.EnemyState.IsInside(x, y) {
				return
			}
			cmd := game.NewOpenCellCommand(coord)
			tx.Add(cmd)
			openCmds = append(openCmds, cmd)
			addEffect("open", coord)
			openCmdCount++

		case "SET_CELL_STATUS":
			status, _ := args["status"].(string)
			cmd := &setCellStatusCommand{X: x, Y: y, Status: status}
			tx.Add(cmd)

		case "REMOVE_SHIP":
			cmd := game.NewRemoveShipCommand(coord)
			tx.Add(cmd)

		case "PLACE_SHIP":
			cmd := game.NewPlaceShipCommand(1, coord, false)
			tx.Add(cmd)

		case "HEAL_SHIP":
			cmd := game.NewHealShipCommand(coord)
			tx.Add(cmd)
			healCmds = append(healCmds, cmd)
			addEffect("heal", coord)

		case "SHOOT":
			cmd := game.NewShootCommand(coord)
			tx.Add(cmd)
			addEffect("shoot", coord)
		}
	}

	// обработка скрипта
	for _, actRaw := range scriptObj.Actions {
		actMap := map[string]interface{}{}
		if m, ok := actRaw.(map[string]interface{}); ok {
			actMap = m
		} else {
			b, _ := json.Marshal(actRaw)
			_ = json.Unmarshal(b, &actMap)
		}

		if name, ok := actMap["Name"].(string); ok {
			if args, ok := actMap["Args"].(map[string]interface{}); ok {
				processAction(name, args)
			}
		} else {
			for k, v := range actMap {
				if args, ok := v.(map[string]interface{}); ok {
					if k == "SWITCH_CASE" || k == "SWICH_CASE" {
						caseKey := "1"
						if dir, ok := params["direction"]; ok {
							caseKey = fmt.Sprintf("%v", dir)
						}
						if caseVal, ok := args[caseKey]; ok {
							if arr, ok := caseVal.([]interface{}); ok {
								for _, sub := range arr {
									subMap, _ := sub.(map[string]interface{})
									subName, _ := subMap["Name"].(string)
									subArgs, _ := subMap["Args"].(map[string]interface{})
									processAction(subName, subArgs)
								}
							}
						}
					} else {
						processAction(k, args)
					}
				}
			}
		}
	}

	// если были только OPEN_CELL и все оказались вне поля — ошибка, то есть если чел вообще не попал никакой координатой, то эррор. А если попал хоть 1 то айтем юзнется
	if openCmdCount == 0 && len(openCmds) == 0 && len(effects) == 0 {
		return nil, fmt.Errorf("all coordinates are out of bounds")
	}

	// выполняем все команды
	if err := tx.Execute(state); err != nil {
		return nil, err
	}

	// добавим is_ship для open
	for i := range effects {
		if effects[i].Type != "open" {
			continue
		}
		for j := range effects[i].Coords {
			coord := &effects[i].Coords[j]
			for _, cmd := range openCmds {
				if cmd.Coords.X == coord.X && cmd.Coords.Y == coord.Y {
					coord.IsShip = &cmd.ShipFound
					break
				}
			}
		}
	}

	// добавим is_ship для heal (всегда true, раз прошло Apply)
	for i := range effects {
		if effects[i].Type != "heal" {
			continue
		}
		for j := range effects[i].Coords {
			coord := &effects[i].Coords[j]
			for _, cmd := range healCmds {
				if cmd.Coords.X == coord.X && cmd.Coords.Y == coord.Y {
					trueVal := true
					coord.IsShip = &trueVal
					break
				}
			}
		}
	}

	return effects, nil
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
		}
		if op == "-" {
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
		}
		if op == "-" {
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
