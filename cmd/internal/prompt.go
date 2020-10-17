package internal

import (
	"fmt"

	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/containerd/console"
)

type errMsg error

type Prompt struct {
	p           *tea.Program
	textInput   input.Model
	err         error
	title       string
	answer      string
}

func NewPrompt(title string, placeholder ...string) *Prompt {
	p := &Prompt{
		title: title,
		textInput: input.NewModel(),
	}

	if len(placeholder) > 0 {
		p.textInput.Placeholder = placeholder[0]
	}

	p.p = tea.NewProgram(p)

	return p
}

func (p *Prompt) YesOrNo() (bool, error) {
	answer, err := p.Answer()
	if err != nil {
		return false, err
	}

	return parseBool(answer), nil
}

func parseBool(str string) bool {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "Yes":
		return true
	}
	return false
}

func (p *Prompt) Answer() (result string, err error) {
	if err = checkConsole(); err != nil {
		return
	}

	if err := p.p.Start(); err != nil {
		return "", err
	}
	return p.answer, nil
}

func (p *Prompt) Init() tea.Cmd {
	p.textInput.Focus()

	return input.Blink(p.textInput)
}

func (p *Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			fallthrough
		case tea.KeyEnter:
			p.answer = p.textInput.Value()
			return p, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		p.err = msg
		return p, nil
	}

	p.textInput, cmd = input.Update(msg, p.textInput)
	return p, cmd
}

func (p *Prompt) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n",
		p.title,
		input.View(p.textInput),
		"(esc to quit)",
	)
}

func checkConsole() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	console.Current()

	return
}
