package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type BasicInterpreter struct {
	program        map[int]string
	variables      map[string]interface{}
	programCounter int
	lineNumbers    []int
	forStack       []ForLoop
	output         []string
}

type ForLoop struct {
	variable string
	end      float64
	step     float64
	line     int
}

func NewBasicInterpreter() *BasicInterpreter {
	return &BasicInterpreter{
		program:   make(map[int]string),
		variables: make(map[string]interface{}),
		forStack:  make([]ForLoop, 0),
		output:    make([]string, 0),
	}
}

func (bi *BasicInterpreter) LoadProgram(programText string) error {
	bi.program = make(map[int]string)
	bi.variables = make(map[string]interface{})
	bi.forStack = make([]ForLoop, 0)
	bi.output = make([]string, 0)

	lines := strings.Split(strings.TrimSpace(programText), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		lineNum, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		bi.program[lineNum] = parts[1]
	}

	bi.lineNumbers = make([]int, 0, len(bi.program))
	for lineNum := range bi.program {
		bi.lineNumbers = append(bi.lineNumbers, lineNum)
	}
	sort.Ints(bi.lineNumbers)

	return nil
}

func (bi *BasicInterpreter) Run(programText string) error {
	if err := bi.LoadProgram(programText); err != nil {
		return err
	}
	return bi.Execute()
}

func (bi *BasicInterpreter) Execute() error {
	if len(bi.lineNumbers) == 0 {
		return nil
	}

	bi.programCounter = 0

	for bi.programCounter < len(bi.lineNumbers) {
		lineNum := bi.lineNumbers[bi.programCounter]
		statement := bi.program[lineNum]

		shouldContinue, err := bi.executeStatement(statement)
		if err != nil {
			return fmt.Errorf("error at line %d: %v", lineNum, err)
		}

		if !shouldContinue {
			break
		}

		bi.programCounter++
	}

	return nil
}

func (bi *BasicInterpreter) executeStatement(statement string) (bool, error) {
	statement = strings.TrimSpace(statement)

	if strings.HasPrefix(statement, "PRINT") {
		return true, bi.executePrint(statement)
	} else if strings.HasPrefix(statement, "LET") {
		return true, bi.executeLet(statement)
	} else if strings.HasPrefix(statement, "GOTO") {
		return true, bi.executeGoto(statement)
	} else if strings.HasPrefix(statement, "IF") {
		return true, bi.executeIf(statement)
	} else if strings.HasPrefix(statement, "FOR") {
		return true, bi.executeFor(statement)
	} else if strings.HasPrefix(statement, "NEXT") {
		return true, bi.executeNext(statement)
	} else if strings.HasPrefix(statement, "INPUT") {
		return true, bi.executeInput(statement)
	} else if strings.HasPrefix(statement, "REM") {
		return true, nil // Comment
	} else if strings.HasPrefix(statement, "END") {
		return false, nil
	} else {
		return false, fmt.Errorf("syntax error: unknown command '%s'", statement)
	}
}

func (bi *BasicInterpreter) executePrint(statement string) error {
	expr := strings.TrimSpace(statement[5:])

	if expr == "" {
		bi.output = append(bi.output, "")
		fmt.Println()
		return nil
	}

	parts := bi.parsePrintParts(expr)
	outputParts := make([]string, 0)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == ";" {
			continue
		}

		if strings.HasPrefix(part, "\"") && strings.HasSuffix(part, "\"") {
			outputParts = append(outputParts, part[1:len(part)-1])
		} else {
			result, err := bi.evaluateExpression(part)
			if err != nil {
				return fmt.Errorf("error evaluating expression '%s': %v", part, err)
			}
			outputParts = append(outputParts, bi.formatValue(result))
		}
	}

	output := strings.Join(outputParts, " ")
	bi.output = append(bi.output, output)
	fmt.Println(output)
	return nil
}

func (bi *BasicInterpreter) executeLet(statement string) error {
	expr := strings.TrimSpace(statement[3:])
	parts := strings.SplitN(expr, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid LET syntax")
	}

	varName := strings.TrimSpace(parts[0])
	valueExpr := strings.TrimSpace(parts[1])

	value, err := bi.evaluateExpression(valueExpr)
	if err != nil {
		return err
	}

	bi.variables[varName] = value
	return nil
}

func (bi *BasicInterpreter) executeGoto(statement string) error {
	lineNumStr := strings.TrimSpace(statement[4:])
	targetLine, err := strconv.Atoi(lineNumStr)
	if err != nil {
		return fmt.Errorf("invalid GOTO syntax")
	}

	for i, lineNum := range bi.lineNumbers {
		if lineNum == targetLine {
			bi.programCounter = i - 1
			return nil
		}
	}

	return fmt.Errorf("undefined line number %d in GOTO statement", targetLine)
}

func (bi *BasicInterpreter) executeIf(statement string) error {
	expr := strings.TrimSpace(statement[2:])
	parts := strings.Split(expr, " THEN ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid IF syntax")
	}

	condition := strings.TrimSpace(parts[0])
	thenPart := strings.TrimSpace(parts[1])

	conditionResult, err := bi.evaluateCondition(condition)
	if err != nil {
		return err
	}

	if conditionResult {
		_, err := bi.executeStatement(thenPart)
		return err
	}

	return nil
}

func (bi *BasicInterpreter) executeFor(statement string) error {
	expr := strings.TrimSpace(statement[3:])
	parts := strings.Fields(expr)
	if len(parts) < 5 || parts[1] != "=" || parts[3] != "TO" {
		return fmt.Errorf("invalid FOR syntax")
	}

	varName := parts[0]
	startValue, err := bi.evaluateExpression(parts[2])
	if err != nil {
		return err
	}
	endValue, err := bi.evaluateExpression(parts[4])
	if err != nil {
		return err
	}

	stepValue := 1.0
	if len(parts) >= 7 && parts[5] == "STEP" {
		step, err := bi.evaluateExpression(parts[6])
		if err != nil {
			return err
		}
		stepValue = bi.toFloat(step)
	}

	bi.variables[varName] = startValue
	currentLine := bi.lineNumbers[bi.programCounter]
	bi.forStack = append(bi.forStack, ForLoop{
		variable: varName,
		end:      bi.toFloat(endValue),
		step:     stepValue,
		line:     currentLine,
	})

	return nil
}

func (bi *BasicInterpreter) executeNext(statement string) error {
	if len(bi.forStack) == 0 {
		return fmt.Errorf("NEXT without FOR")
	}

	var varName string
	if len(statement) > 4 {
		varName = strings.TrimSpace(statement[4:])
	}

	loopInfo := bi.forStack[len(bi.forStack)-1]

	if varName != "" && varName != loopInfo.variable {
		return fmt.Errorf("NEXT %s doesn't match FOR %s", varName, loopInfo.variable)
	}

	currentValue := bi.toFloat(bi.variables[loopInfo.variable])
	newValue := currentValue + loopInfo.step
	bi.variables[loopInfo.variable] = newValue

	if (loopInfo.step > 0 && newValue <= loopInfo.end) ||
		(loopInfo.step < 0 && newValue >= loopInfo.end) {
		for i, lineNum := range bi.lineNumbers {
			if lineNum == loopInfo.line {
				bi.programCounter = i
				break
			}
		}
	} else {
		bi.forStack = bi.forStack[:len(bi.forStack)-1]
	}

	return nil
}

func (bi *BasicInterpreter) executeInput(statement string) error {
	expr := strings.TrimSpace(statement[5:])

	var prompt string
	var varName string

	if strings.Contains(expr, ";") {
		parts := strings.SplitN(expr, ";", 2)
		prompt = strings.TrimSpace(parts[0])
		varName = strings.TrimSpace(parts[1])

		if strings.HasPrefix(prompt, "\"") && strings.HasSuffix(prompt, "\"") {
			prompt = prompt[1 : len(prompt)-1]
			fmt.Print(prompt)
		}
	} else {
		varName = expr
		fmt.Print("? ")
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)

	if value, err := strconv.ParseFloat(input, 64); err == nil {
		if value == float64(int(value)) {
			bi.variables[varName] = int(value)
		} else {
			bi.variables[varName] = value
		}
	} else {
		bi.variables[varName] = input
	}

	return nil
}

func (bi *BasicInterpreter) evaluateExpression(expr string) (interface{}, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") {
		return expr[1 : len(expr)-1], nil
	}

	if value, exists := bi.variables[expr]; exists {
		return value, nil
	}

	if value, err := strconv.ParseFloat(expr, 64); err == nil {
		if value == float64(int(value)) {
			return int(value), nil
		}
		return value, nil
	}

	return bi.evaluateArithmetic(expr)
}

func (bi *BasicInterpreter) evaluateArithmetic(expr string) (interface{}, error) {
	expr = strings.TrimSpace(expr)

	// Handle addition and subtraction
	for i := len(expr) - 1; i >= 0; i-- {
		if expr[i] == '+' || expr[i] == '-' {
			if i > 0 && !strings.ContainsAny(string(expr[i-1]), "*/+-(<>=") {
				left, err := bi.evaluateExpression(expr[:i])
				if err != nil {
					return nil, err
				}
				right, err := bi.evaluateExpression(expr[i+1:])
				if err != nil {
					return nil, err
				}

				leftFloat := bi.toFloat(left)
				rightFloat := bi.toFloat(right)

				if expr[i] == '+' {
					result := leftFloat + rightFloat
					if result == float64(int(result)) {
						return int(result), nil
					}
					return result, nil
				} else {
					result := leftFloat - rightFloat
					if result == float64(int(result)) {
						return int(result), nil
					}
					return result, nil
				}
			}
		}
	}

	// Handle multiplication and division
	for i := len(expr) - 1; i >= 0; i-- {
		if expr[i] == '*' || expr[i] == '/' {
			left, err := bi.evaluateExpression(expr[:i])
			if err != nil {
				return nil, err
			}
			right, err := bi.evaluateExpression(expr[i+1:])
			if err != nil {
				return nil, err
			}

			leftFloat := bi.toFloat(left)
			rightFloat := bi.toFloat(right)

			if expr[i] == '*' {
				result := leftFloat * rightFloat
				if result == float64(int(result)) {
					return int(result), nil
				}
				return result, nil
			} else {
				if rightFloat == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				result := leftFloat / rightFloat
				if result == float64(int(result)) {
					return int(result), nil
				}
				return result, nil
			}
		}
	}

	if value, exists := bi.variables[expr]; exists {
		return value, nil
	}

	if value, err := strconv.ParseFloat(expr, 64); err == nil {
		if value == float64(int(value)) {
			return int(value), nil
		}
		return value, nil
	}

	return nil, fmt.Errorf("cannot evaluate expression: %s", expr)
}

func (bi *BasicInterpreter) evaluateCondition(condition string) (bool, error) {
	condition = strings.TrimSpace(condition)

	operators := []string{">", "<", "="}
	for _, op := range operators {
		if strings.Contains(condition, op) {
			parts := strings.SplitN(condition, op, 2)
			if len(parts) == 2 {
				left, err := bi.evaluateExpression(strings.TrimSpace(parts[0]))
				if err != nil {
					return false, err
				}
				right, err := bi.evaluateExpression(strings.TrimSpace(parts[1]))
				if err != nil {
					return false, err
				}

				leftFloat := bi.toFloat(left)
				rightFloat := bi.toFloat(right)

				switch op {
				case ">":
					return leftFloat > rightFloat, nil
				case "<":
					return leftFloat < rightFloat, nil
				case "=":
					return leftFloat == rightFloat, nil
				}
			}
		}
	}

	return false, nil
}

func (bi *BasicInterpreter) parsePrintParts(expr string) []string {
	parts := make([]string, 0)
	currentPart := ""
	inQuotes := false

	for _, char := range expr {
		if char == '"' {
			inQuotes = !inQuotes
			currentPart += string(char)
		} else if char == ';' && !inQuotes {
			if strings.TrimSpace(currentPart) != "" {
				parts = append(parts, strings.TrimSpace(currentPart))
			}
			currentPart = ""
		} else {
			currentPart += string(char)
		}
	}

	if strings.TrimSpace(currentPart) != "" {
		parts = append(parts, strings.TrimSpace(currentPart))
	}

	return parts
}

func (bi *BasicInterpreter) toFloat(value interface{}) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		return 0
	default:
		return 0
	}
}

func (bi *BasicInterpreter) formatValue(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		if v == float64(int(v)) {
			return strconv.Itoa(int(v))
		}
		return strconv.FormatFloat(v, 'g', -1, 64)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (bi *BasicInterpreter) GetOutput() []string {
	return bi.output
}

func main() {
	interpreter := NewBasicInterpreter()
	
	program := `10 PRINT "Hello, World!"
20 LET A = 42
30 PRINT A`
	
	if err := interpreter.Run(program); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}