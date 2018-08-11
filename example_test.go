package task_test

import (
	"bytes"
	"fmt"
	"github.com/metakeule/task"
	"strings"
	"testing"
)

func newExample(input string) *example {
	return &example{inputs: strings.Split(input, "\n")}
}

type example struct {
	target string
	inputs []string
	inPtr  int
	result bytes.Buffer
}

func (s *example) String() string {
	return s.result.String()
}

func (s *example) getInput() string {
	if s.inPtr > len(s.inputs)-1 {
		return ""
	}

	defer func() {
		s.inPtr++
	}()

	return s.inputs[s.inPtr]
}

func (s *example) Run() (task.Task, error) {
	return s.dispatchByName()
}

func (s *example) dispatchByName() (task.Task, error) {
	name := s.getInput()
	if name == "" {
		return nil, nil
	}

	switch name {
	case "Bob":
		s.result.WriteString("Name: Bob\n")
		return s.handleMale, nil
	case "Anne":
		s.result.WriteString("Name: Anne\n")
		return s.handleFemale, nil
	default:
		return nil, fmt.Errorf("unknown name: %#v", name)
	}

}

func (s *example) handleMale() (task.Task, error) {
	car := s.getInput()
	if car == "" {
		return nil, fmt.Errorf("missing car")
	}

	s.result.WriteString("Car: " + car + "\n")
	return s.cheers, nil
}

func (s *example) cheers() (task.Task, error) {
	s.result.WriteString("cheers\n")
	return nil, nil
}

func (s *example) handleFemale() (task.Task, error) {
	pet := s.getInput()
	if pet == "" {
		return nil, fmt.Errorf("missing pet")
	}

	s.result.WriteString("Pet: " + pet + "\n")
	return s.byebye, nil
}

func (s *example) byebye() (task.Task, error) {
	s.result.WriteString("byebye\n")
	return nil, nil
}

func TestEx(t *testing.T) {

	tests := []struct {
		input    string
		expected string
		err      error
	}{
		{"Bob\nMini", "Name: Bob\nCar: Mini\ncheers\n", nil},
		{"Anne\ncat", "Name: Anne\nPet: cat\nbyebye\n", nil},
		{"Bob\n", "Name: Bob\n", fmt.Errorf("missing car")},
		{"Anne\n", "Name: Anne\n", fmt.Errorf("missing cat")},
		{"Paul\n", "", fmt.Errorf("unknown name: \nPaul\n")},
	}

	for _, test := range tests {

		ex := newExample(test.input)
		err := task.Run(ex.Run)

		if err == nil && test.err != nil {
			t.Errorf("ex := newExample(%v); err := task.Run(ex.Run); err = nil; want %#v", test.input, test.err.Error())
		}

		if err != nil && test.err == nil {
			t.Errorf("ex := newExample(%v); err := task.Run(ex.Run); err := task.Run(); err = %#v; want nil", test.input, err)
		}

		if got, want := ex.String(), test.expected; got != want {
			t.Errorf("ex := newExample(%v); err := task.Run(ex.Run); ex.String() = %#v; want %#v", test.input, got, want)
		}
	}

}

func Example() {
	fmt.Println("")

	var bf bytes.Buffer

	// For sake of clarity, we use closures here.
	// In real code you might want to use methods of
	// a struct that contains your shared data.

	second := func() (task.Task, error) {
		bf.WriteString("Car: Mini\n")
		return nil, nil
	}

	first := func() (task.Task, error) {
		bf.WriteString("Name: Bob\n")
		return second, nil
	}

	err := task.Run(first)

	if err != nil {
		return
	}

	fmt.Println(bf.String())

	// Output:
	// Name: Bob
	// Car: Mini
}
