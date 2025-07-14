package log

import (
	"fmt"
	"strings"
)

// TestOutput はテスト用の出力キャプチャ
type TestOutput struct {
	Messages       []string
	ErrorMessages  []string
	ProgressMsgs   []string
	SuccessMsgs    []string
	VerboseMsgs    []string
	ConfirmPrompts []string
	ConfirmReturn  bool // Confirmメソッドの戻り値を制御
	Verbose        bool // 詳細出力モードの制御
}

// NewTestOutput は新しいTestOutputを作成する
func NewTestOutput(verbose bool) *TestOutput {
	return &TestOutput{
		Messages:       []string{},
		ErrorMessages:  []string{},
		ProgressMsgs:   []string{},
		SuccessMsgs:    []string{},
		VerboseMsgs:    []string{},
		ConfirmPrompts: []string{},
		ConfirmReturn:  false,
		Verbose:        verbose,
	}
}

// Print outputs a standard message
func (t *TestOutput) Print(msg string) {
	t.Messages = append(t.Messages, msg)
}

// Printf outputs a formatted standard message
func (t *TestOutput) Printf(format string, args ...interface{}) {
	t.Messages = append(t.Messages, fmt.Sprintf(format, args...))
}

// PrintProgress outputs a progress message
func (t *TestOutput) PrintProgress(msg string) {
	t.ProgressMsgs = append(t.ProgressMsgs, msg)
}

// PrintSuccess outputs a success message
func (t *TestOutput) PrintSuccess(msg string) {
	t.SuccessMsgs = append(t.SuccessMsgs, msg)
}

// PrintError outputs an error message
func (t *TestOutput) PrintError(err error) {
	if err != nil {
		t.ErrorMessages = append(t.ErrorMessages, err.Error())
	} else {
		t.ErrorMessages = append(t.ErrorMessages, "nil error")
	}
}

// PrintVerbose outputs a message only in verbose mode
func (t *TestOutput) PrintVerbose(msg string) {
	if t.Verbose {
		t.VerboseMsgs = append(t.VerboseMsgs, msg)
	}
}

// Confirm simulates user confirmation
func (t *TestOutput) Confirm(prompt string) bool {
	t.ConfirmPrompts = append(t.ConfirmPrompts, prompt)
	return t.ConfirmReturn
}

// SetConfirmReturn sets the return value for Confirm method
func (t *TestOutput) SetConfirmReturn(val bool) {
	t.ConfirmReturn = val
}

// テストヘルパーメソッド

// ContainsMessage checks if any standard message contains the given text
func (t *TestOutput) ContainsMessage(partialMsg string) bool {
	for _, msg := range t.Messages {
		if strings.Contains(msg, partialMsg) {
			return true
		}
	}
	return false
}

// ContainsError checks if any error message contains the given text
func (t *TestOutput) ContainsError(partialMsg string) bool {
	for _, msg := range t.ErrorMessages {
		if strings.Contains(msg, partialMsg) {
			return true
		}
	}
	return false
}

// ContainsProgress checks if any progress message contains the given text
func (t *TestOutput) ContainsProgress(partialMsg string) bool {
	for _, msg := range t.ProgressMsgs {
		if strings.Contains(msg, partialMsg) {
			return true
		}
	}
	return false
}

// ContainsSuccess checks if any success message contains the given text
func (t *TestOutput) ContainsSuccess(partialMsg string) bool {
	for _, msg := range t.SuccessMsgs {
		if strings.Contains(msg, partialMsg) {
			return true
		}
	}
	return false
}

// ContainsVerbose checks if any verbose message contains the given text
func (t *TestOutput) ContainsVerbose(partialMsg string) bool {
	for _, msg := range t.VerboseMsgs {
		if strings.Contains(msg, partialMsg) {
			return true
		}
	}
	return false
}

// ContainsConfirmPrompt checks if any confirmation prompt contains the given text
func (t *TestOutput) ContainsConfirmPrompt(partialMsg string) bool {
	for _, msg := range t.ConfirmPrompts {
		if strings.Contains(msg, partialMsg) {
			return true
		}
	}
	return false
}

// Clear resets all captured messages
func (t *TestOutput) Clear() {
	t.Messages = []string{}
	t.ErrorMessages = []string{}
	t.ProgressMsgs = []string{}
	t.SuccessMsgs = []string{}
	t.VerboseMsgs = []string{}
	t.ConfirmPrompts = []string{}
}
