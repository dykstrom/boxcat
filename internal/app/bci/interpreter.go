package bci

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// BoxCat interpreter state
type Interpreter struct {
	boxes   map[string]int
	paw     int
	pc      int
	stack   []int
	labels  map[string]int
	program []string
	rng     *rand.Rand
}

func NewInterpreter(program []string) *Interpreter {
	return &Interpreter{
		boxes:   make(map[string]int),
		paw:     0,
		pc:      0,
		stack:   []int{},
		labels:  make(map[string]int),
		program: program,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Preprocess: find all labels
func (i *Interpreter) findLabels() error {
	for idx, line := range i.program {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, ":") {
			label := strings.TrimSuffix(line, ":")
			label = strings.ToUpper(label)
			if _, found := i.labels[label]; found {
				return withNumber(idx, duplicateLabel(label))
			}
			i.labels[label] = idx
		}
	}
	return nil
}

// Run the program
func (i *Interpreter) Run() error {
	err := i.findLabels()
	if err != nil {
		return err
	}

	for i.pc < len(i.program) {
		line := i.program[i.pc]
		err = i.executeLine(line)
		if err != nil {
			return withNumber(i.pc, err)
		}
		i.pc++
	}
	return nil
}

func (i *Interpreter) executeLine(line string) error {
	// Remove comments
	if idx := strings.Index(line, "#"); idx != -1 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)
	if line == "" || strings.HasSuffix(line, ":") {
		return nil
	}

	// Tokenize
	line = strings.ToUpper(line)
	tokens := strings.Fields(line)

	//fmt.Println("Executing:", tokens, ", Paw:", i.paw, ", Boxes:", i.boxes)
	var err error
	switch tokens[0] {
	case "SIT":
		i.sitIn(getArg(tokens, 2))
	case "JUMP":
		i.jumpOutOf(getArg(tokens, 3))
	case "PEEK":
		i.peekInside(getArg(tokens, 2))
	case "DROP":
		i.dropIn(getArg(tokens, 2))
	case "POUNCE":
		i.pounceOn(getArg(tokens, 2))
	case "PURR":
		i.purrAt(getArg(tokens, 2))
	case "HISS":
		i.hissAt(getArg(tokens, 2))
	case "PLAYFULLY":
		i.playfullyBat(getArg(tokens, 3))
	case "KNOCK":
		err = i.knockOver(getArg(tokens, 2))
	case "LEAVE":
		err = i.leaveGiftIn(getArg(tokens, 5))
	case "MEOW":
		fmt.Printf("%d ", i.paw)
	case "YOWL":
		fmt.Printf("%c", i.paw)
	case "LISTEN":
		err = i.listenForWhisper()
	case "SNIFF":
		err = i.handleSniff(getArg(tokens, 1))
	case "DART":
		err = i.handleDart(tokens)
	case "LEAP":
		i.leapTo(getArg(tokens, 2))
	case "IF":
		err = i.handleIf(line)
	case "SUDDENLY":
		i.suddenlyScratch()
	case "DOZE":
		i.doze()
	case "TAKE":
		os.Exit(0)
	case "GET":
		i.getStuck()
	default:
		err = unknownCommand(line)
	}
	return err
}

// Helper to get argument at position n
func getArg(tokens []string, n int) string {
	if len(tokens) > n {
		return tokens[n]
	}
	return ""
}

// Command implementations

func (i *Interpreter) sitIn(box string) {
	i.boxes[box] = i.paw
	i.paw = 0
}

func (i *Interpreter) jumpOutOf(box string) {
	i.paw = i.boxes[box]
	i.boxes[box] = 0
}

func (i *Interpreter) peekInside(box string) {
	i.paw = i.boxes[box]
}

func (i *Interpreter) dropIn(box string) {
	i.boxes[box] = i.paw
}

func (i *Interpreter) pounceOn(val string) {
	n, _ := strconv.Atoi(val)
	i.paw = n
}

func (i *Interpreter) purrAt(box string) {
	i.paw += i.boxes[box]
}

func (i *Interpreter) hissAt(box string) {
	i.paw -= i.boxes[box]
}

func (i *Interpreter) playfullyBat(box string) {
	i.paw *= i.boxes[box]
}

func (i *Interpreter) knockOver(box string) error {
	if i.boxes[box] == 0 {
		return divisionByZero("Division", box)
	}
	i.paw /= i.boxes[box]
	return nil
}

func (i *Interpreter) leaveGiftIn(box string) error {
	if i.boxes[box] == 0 {
		return divisionByZero("Modulo", box)
	}
	i.paw %= i.boxes[box]
	return nil
}

func (i *Interpreter) listenForWhisper() error {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	trimmed := strings.TrimSpace(text)
	n, err := strconv.Atoi(trimmed)
	if err != nil {
		return notNumber(trimmed)
	}
	i.paw = n
	return nil
}

func (i *Interpreter) handleSniff(target string) error {
	if target == "AROUND" {
		return i.sniffAround()
	} else if target == "CATNIP" {
		i.sniffCatnip()
	}
	return nil
}

func (i *Interpreter) sniffAround() error {
	reader := bufio.NewReader(os.Stdin)
	ch, _, err := reader.ReadRune()
	if err == io.EOF {
		i.paw = -1
	} else if err != nil {
		return fatalError(err)
	} else {
		i.paw = int(ch)
	}
	return nil
}

func (i *Interpreter) sniffCatnip() {
	i.paw = i.rng.Intn(32768)
}

func (i *Interpreter) handleDart(tokens []string) error {
	direction := getArg(tokens, 1)
	if direction == "TO" {
		return i.dartTo(getArg(tokens, 2))
	} else if direction == "BACK" {
		return i.dartBack()
	}
	return nil
}

func (i *Interpreter) dartTo(label string) error {
	i.stack = append(i.stack, i.pc)
	if idx, ok := i.labels[label]; ok {
		i.pc = idx
		return nil
	} else {
		return undefinedLabel(label)
	}
}

func (i *Interpreter) dartBack() error {
	if len(i.stack) == 0 {
		return stackEmpty()
	}
	i.pc = i.stack[len(i.stack)-1]
	i.stack = i.stack[:len(i.stack)-1]
	return nil
}

func (i *Interpreter) leapTo(label string) error {
	if idx, ok := i.labels[label]; ok {
		i.pc = idx
		return nil
	} else {
		return undefinedLabel(label)
	}
}

func (i *Interpreter) handleIf(line string) error {
	// Example: IF CAT CURIOUS, MEOW
	cond, rest, _ := strings.Cut(line[2:], ",")
	cond = strings.TrimSpace(cond)
	fields := strings.Fields(cond)
	cmd := strings.TrimSpace(rest)
	shouldExec := false

	var err error
	switch {
	case strings.HasPrefix(cond, "CAT CURIOUS"):
		shouldExec = i.paw != 0
	case strings.HasPrefix(cond, "CAT BORED"):
		shouldExec = i.paw == 0
	case strings.HasPrefix(cond, "BOX EMPTY"):
		if len(fields) < 3 {
			err = missingBoxPonder(line)
		} else {
			shouldExec = i.boxes[fields[2]] == 0
		}
	case strings.HasPrefix(cond, "BOX NOT EMPTY"):
		if len(fields) < 4 {
			err = missingBoxPonder(line)
		} else {
			shouldExec = i.boxes[fields[3]] != 0
		}
	default:
		err = unknownCommand(line)
	}
	if err != nil {
		return err
	} else if cmd == "" {
		err = missingCommand(line)
	} else if shouldExec {
		err = i.executeLine(cmd)
	}
	return err
}

func (i *Interpreter) suddenlyScratch() {
	i.paw = 0
}

func (i *Interpreter) doze() {
	if i.paw > 0 {
		time.Sleep(time.Duration(i.paw) * time.Millisecond)
	}
}

func (i *Interpreter) getStuck() {
	for {
	}
}

func unknownCommand(text string) error {
	return fmt.Errorf(
		"Hiss! The line '%s' caused the cat to arch its back. Syntax unclear, human servant!", text,
	)
}

func missingCommand(text string) error {
	return fmt.Errorf(
		"Hiss! The line '%s' is missing a command. Cat looks at you expectantly.", text,
	)
}

func missingBoxPonder(text string) error {
	return fmt.Errorf(
		"The command '%s' needs a box name to ponder on! Cat is pacing, waiting for specifics.", text,
	)
}

func duplicateLabel(label string) error {
	return fmt.Errorf(
		"Two '%s's? Cat is torn! Like choosing between the salmon-flavored treats and the chicken-flavored ones. Make up your mind, human!", label,
	)
}

func undefinedLabel(label string) error {
	return fmt.Errorf(
		"Lost kitty! '%s' is nowhere to be found. Cat is now sitting in the middle of the program, looking confused.", label,
	)
}

func stackEmpty() error {
	return fmt.Errorf(
		"Cat is looking for somewhere to dart back, but found nothing. Stack is empty.",
	)
}

func divisionByZero(operation, box string) error {
	return fmt.Errorf(
		"'%s' is as empty as a treat jar after a midnight raid! %s by zero attempted. Cat just blinked slowly.", box, operation,
	)
}

func notNumber(text string) error {
	return fmt.Errorf(
		"The human whispered '%s'. Cat blinked. That's not a number it recognizes.", text,
	)
}

func fatalError(err error) error {
	return fmt.Errorf(
		"Fatal: The catnip may have been too strong. Unfathomable error '%v'. Cat needs a system reboot (and maybe a snack).", err,
	)
}

func withNumber(num int, err error) error {
	return fmt.Errorf("Error on line %d: %v", num+1, err)
}
