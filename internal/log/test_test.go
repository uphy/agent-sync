package log

import (
	"errors"
	"testing"
)

func TestTestOutput(t *testing.T) {
	// 通常モード（Verbose = false）でのテスト
	output := NewTestOutput(false)

	// 基本的なメッセージ出力テスト
	output.Print("test message")
	if len(output.Messages) != 1 || output.Messages[0] != "test message" {
		t.Errorf("Print did not store message correctly, got: %v", output.Messages)
	}

	// フォーマット付きメッセージ出力テスト
	output.Printf("formatted %s", "message")
	if len(output.Messages) != 2 || output.Messages[1] != "formatted message" {
		t.Errorf("Printf did not store formatted message correctly, got: %v", output.Messages)
	}

	// 進捗メッセージテスト
	output.PrintProgress("in progress")
	if len(output.ProgressMsgs) != 1 || output.ProgressMsgs[0] != "in progress" {
		t.Errorf("PrintProgress did not store message correctly, got: %v", output.ProgressMsgs)
	}

	// 成功メッセージテスト
	output.PrintSuccess("success")
	if len(output.SuccessMsgs) != 1 || output.SuccessMsgs[0] != "success" {
		t.Errorf("PrintSuccess did not store message correctly, got: %v", output.SuccessMsgs)
	}

	// エラーメッセージテスト
	testErr := errors.New("test error")
	output.PrintError(testErr)
	if len(output.ErrorMessages) != 1 || output.ErrorMessages[0] != "test error" {
		t.Errorf("PrintError did not store error correctly, got: %v", output.ErrorMessages)
	}

	// nil エラーのテスト
	output.PrintError(nil)
	if len(output.ErrorMessages) != 2 || output.ErrorMessages[1] != "nil error" {
		t.Errorf("PrintError did not handle nil error correctly, got: %v", output.ErrorMessages)
	}

	// Verbose モードが無効の場合、メッセージは保存されないこと
	output.PrintVerbose("verbose message")
	if len(output.VerboseMsgs) != 0 {
		t.Errorf("PrintVerbose stored message when verbose mode is off, got: %v", output.VerboseMsgs)
	}

	// 確認プロンプトテスト
	output.SetConfirmReturn(true)
	result := output.Confirm("confirm?")
	if !result {
		t.Error("Confirm did not return the value set by SetConfirmReturn")
	}
	if len(output.ConfirmPrompts) != 1 || output.ConfirmPrompts[0] != "confirm?" {
		t.Errorf("Confirm did not store prompt correctly, got: %v", output.ConfirmPrompts)
	}

	// ContainsMessage メソッドのテスト
	if !output.ContainsMessage("test") {
		t.Error("ContainsMessage failed to find existing message")
	}
	if output.ContainsMessage("nonexistent") {
		t.Error("ContainsMessage found non-existent message")
	}

	// ContainsError メソッドのテスト
	if !output.ContainsError("test error") {
		t.Error("ContainsError failed to find existing error")
	}

	// ContainsProgress メソッドのテスト
	if !output.ContainsProgress("progress") {
		t.Error("ContainsProgress failed to find existing progress message")
	}

	// ContainsSuccess メソッドのテスト
	if !output.ContainsSuccess("success") {
		t.Error("ContainsSuccess failed to find existing success message")
	}

	// ContainsConfirmPrompt メソッドのテスト
	if !output.ContainsConfirmPrompt("confirm") {
		t.Error("ContainsConfirmPrompt failed to find existing prompt")
	}

	// Clear メソッドのテスト
	output.Clear()
	if len(output.Messages) != 0 || len(output.ErrorMessages) != 0 || len(output.ProgressMsgs) != 0 ||
		len(output.SuccessMsgs) != 0 || len(output.VerboseMsgs) != 0 || len(output.ConfirmPrompts) != 0 {
		t.Error("Clear did not reset all message slices")
	}
}

func TestVerboseMode(t *testing.T) {
	// Verbose モード有効でのテスト
	output := NewTestOutput(true)

	// Verbose モードが有効の場合、メッセージが保存されること
	output.PrintVerbose("verbose message")
	if len(output.VerboseMsgs) != 1 || output.VerboseMsgs[0] != "verbose message" {
		t.Errorf("PrintVerbose did not store message when verbose mode is on, got: %v", output.VerboseMsgs)
	}

	// ContainsVerbose メソッドのテスト
	if !output.ContainsVerbose("verbose") {
		t.Error("ContainsVerbose failed to find existing verbose message")
	}
}
