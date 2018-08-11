package task_test

import (
	"bytes"
	"fmt"
	"github.com/metakeule/task"
	"strings"
	"testing"
)

func newSystem(input string) *system {
	return &system{inputs: strings.Split(input, "\n")}
}

type system struct {
	target string
	inputs []string
	inPtr  int
	result bytes.Buffer
}

func (s *system) String() string {
	return s.result.String()
}

func (s *system) getInput() string {
	if s.inPtr > len(s.inputs)-1 {
		return ""
	}

	defer func() {
		s.inPtr++
	}()

	return s.inputs[s.inPtr]
}

func (s *system) Run() (task.Command, error) {
	return s.dispatchByName()
}

func (s *system) dispatchByName() (task.Command, error) {
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

func (s *system) handleMale() (task.Command, error) {
	car := s.getInput()
	if car == "" {
		return nil, fmt.Errorf("missing car")
	}

	s.result.WriteString("Car: " + car + "\n")
	return s.cheers, nil
}

func (s *system) cheers() (task.Command, error) {
	s.result.WriteString("cheers\n")
	return nil, nil
}

func (s *system) handleFemale() (task.Command, error) {
	pet := s.getInput()
	if pet == "" {
		return nil, fmt.Errorf("missing pet")
	}

	s.result.WriteString("Pet: " + pet + "\n")
	return s.byebye, nil
}

func (s *system) byebye() (task.Command, error) {
	s.result.WriteString("byebye\n")
	return nil, nil
}

func TestSystem(t *testing.T) {

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

		sys := newSystem(test.input)
		err := task.New(sys.Run).Run()

		if err == nil && test.err != nil {
			t.Errorf("sys := newSystem(%v); err := task.New(sys.Run).Run(); err = nil; want %#v", test.input, test.err.Error())
		}

		if err != nil && test.err == nil {
			t.Errorf("sys := newSystem(%v); err := task.New(sys.Run).Run(); err := task.Run(); err = %#v; want nil", test.input, err)
		}

		if got, want := sys.String(), test.expected; got != want {
			t.Errorf("sys := newSystem(%v); err := task.New(sys.Run).Run(); sys.String() = %#v; want %#v", test.input, got, want)
		}
	}

}

func Example() {
	fmt.Println("")

	sys := newSystem("Bob\nMini")

	err := task.New(sys.Run).Run()

	if err != nil {
		return
	}

	fmt.Println(sys.String())

	// Output:
	// Name: Bob
	// Car: Mini
	// cheers
}
