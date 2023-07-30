package loader

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Default set
var Default = Set{50 * time.Millisecond, []string{"+", "\\", "|", "!", "/", "-", "x"}}

// Moon set
var Moon = Set{100 * time.Millisecond, []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}}

// Arrows set
var Arrows = Set{100 * time.Millisecond, []string{"⇢", "⇨", "⇒", "⇉", "⇶"}}

// Dots set
var Dots = Set{100 * time.Millisecond, []string{".", "·", "•", "¤", "°", "¤", "•", "·"}}

// Loader is the spinner struct
type Loader struct {
	stch chan bool
	stop bool
	set  Set
	wg   *sync.WaitGroup
}

// Set defines animation chars and delay
type Set struct {
	Delay time.Duration
	Chars []string
}

// New gets a new spinner
func New(sets ...Set) *Loader {
	sets = append(sets, Default)
	return &Loader{stch: make(chan bool), set: sets[0]}
}

// WithWait attaches a wait group
func (s *Loader) WithWait(wg *sync.WaitGroup) *Loader {
	wg.Add(1)
	s.wg = wg
	return s
}

// Start starts the spinner
func (s *Loader) Start() {
	if err := tput("civis"); err != nil {
		fmt.Print("\033[?25l")
	}
	s.doSpin()
}

// Stop stops the spinner
func (s *Loader) Stop() {
	if s.wg != nil {
		defer s.wg.Done()
	}
	s.stop = true
	if err := tput("cvvis"); err != nil {
		fmt.Print("\033[?25h")
	}
}

func (s *Loader) doSpin() {
	for {
	outer:
		select {
		case _, ok := <-s.stch:
			if ok {
				fmt.Print("\010")
				break outer
			}
		default:
			for _, c := range s.set.Chars {
				if s.stop {
					s.stch <- true
				} else if len(c) > 0 {
					fmt.Print(c, "\010")
					time.Sleep(s.set.Delay)
				}
			}
		}
	}
}

func tput(arg string) error {
	cmd := exec.Command("tput", arg)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
