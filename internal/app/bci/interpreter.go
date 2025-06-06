package bci

import (
	"bufio"
	"fmt"
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
func (i *Interpreter) findLabels() {
	for idx, line := range i.program {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, ":") {
			label := strings.TrimSuffix(line, ":")
			i.labels[strings.ToUpper(label)] = idx
		}
	}
}

// Run the program
func (i *Interpreter) Run() {
	i.findLabels()
	for i.pc < len(i.program) {
		line := i.program[i.pc]
		i.executeLine(line)
		i.pc++
	}
}

func (i *Interpreter) executeLine(line string) {
	// Remove comments
	if idx := strings.Index(line, "#"); idx != -1 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)
	if line == "" || strings.HasSuffix(line, ":") {
		return
	}

	// Tokenize
	upper := strings.ToUpper(line)
	tokens := strings.Fields(upper)
	if len(tokens) == 0 {
		return
	}

	//fmt.Println("Executing:", tokens, ", Paw:", i.paw, ", Boxes:", i.boxes)
	switch tokens[0] {
	case "SIT", "SITIN":
		i.sitIn(getArg(tokens, 2))
	case "JUMP", "JUMPOUTOF":
		i.jumpOutOf(getArg(tokens, 3))
	case "PEEK", "PEEKINSIDE":
		i.peekInside(getArg(tokens, 2))
	case "DROP", "DROPIN":
		i.dropIn(getArg(tokens, 2))
	case "POUNCE", "POUNCEON":
		i.pounceOn(getArg(tokens, 2))
	case "PURR", "PURRAT":
		i.purrAt(getArg(tokens, 2))
	case "HISS", "HISSAT":
		i.hissAt(getArg(tokens, 2))
	case "PLAYFULLY", "PLAYFULLYBAT":
		i.playfullyBat(getArg(tokens, 3))
	case "KNOCK", "KNOCKOVER":
		i.knockOver(getArg(tokens, 2))
	case "LEAVE", "LEAVEA":
		i.leaveGiftIn(getArg(tokens, 5))
	case "MEOW":
		fmt.Println(i.paw)
	case "YOWL":
		fmt.Printf("%c", i.paw)
	case "LISTEN", "LISTENFOR":
		i.listenForWhisper()
	case "SNIFF", "SNIFFAROUND":
		if tokens[1] == "AROUND" {
			i.sniffAround()
		} else if tokens[1] == "CATNIP" {
			i.sniffCatnip()
		}
	case "DART", "DARTTO":
		i.dartTo(getArg(tokens, 2))
	case "DARTBACK":
		i.dartBack()
	case "LEAP", "LEAPTO":
		i.leapTo(getArg(tokens, 2))
	case "IF":
		i.handleIf(line)
	case "SUDDENLY", "SUDDENLYSCRATCH":
		i.suddenlyScratch()
	case "DOZE":
		i.doze()
	case "TAKE", "TAKEANAP":
		os.Exit(0)
	case "GET", "GETSTUCK":
		i.getStuck()
	}
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

func (i *Interpreter) knockOver(box string) {
	if i.boxes[box] == 0 {
		fmt.Fprintln(os.Stderr, "Angry Cat: Division by zero!")
		os.Exit(1)
	}
	i.paw /= i.boxes[box]
}

func (i *Interpreter) leaveGiftIn(box string) {
	if i.boxes[box] == 0 {
		fmt.Fprintln(os.Stderr, "Angry Cat: Modulo by zero!")
		os.Exit(1)
	}
	i.paw %= i.boxes[box]
}

func (i *Interpreter) listenForWhisper() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	n, _ := strconv.Atoi(strings.TrimSpace(text))
	i.paw = n
}

func (i *Interpreter) sniffAround() {
	reader := bufio.NewReader(os.Stdin)
	ch, _, _ := reader.ReadRune()
	i.paw = int(ch)
}

func (i *Interpreter) sniffCatnip() {
	i.paw = i.rng.Intn(32768)
}

func (i *Interpreter) dartTo(label string) {
	i.stack = append(i.stack, i.pc)
	if idx, ok := i.labels[label]; ok {
		i.pc = idx
	} else {
		fmt.Fprintf(os.Stderr, "Unknown label: %s\n", label)
		os.Exit(1)
	}
}

func (i *Interpreter) dartBack() {
	if len(i.stack) == 0 {
		fmt.Fprintln(os.Stderr, "DART BACK with empty stack!")
		os.Exit(1)
	}
	i.pc = i.stack[len(i.stack)-1]
	i.stack = i.stack[:len(i.stack)-1]
}

func (i *Interpreter) leapTo(label string) {
	if idx, ok := i.labels[label]; ok {
		i.pc = idx
	} else {
		fmt.Fprintf(os.Stderr, "Unknown label: %s\n", label)
		os.Exit(1)
	}
}

func (i *Interpreter) handleIf(line string) {
	// Example: IF CAT CURIOUS, MEOW
	cond, rest, _ := strings.Cut(line[2:], ",")
	cond = strings.TrimSpace(cond)
	cmd := strings.TrimSpace(rest)
	shouldExec := false

	switch {
	case strings.HasPrefix(cond, "CAT CURIOUS"):
		shouldExec = i.paw != 0
	case strings.HasPrefix(cond, "CAT BORED"):
		shouldExec = i.paw == 0
	case strings.HasPrefix(cond, "BOX EMPTY"):
		box := strings.Fields(cond)[2]
		shouldExec = i.boxes[box] == 0
	case strings.HasPrefix(cond, "BOX NOT EMPTY"):
		box := strings.Fields(cond)[3]
		shouldExec = i.boxes[box] != 0
	}
	if shouldExec && cmd != "" {
		i.executeLine(cmd)
	}
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
