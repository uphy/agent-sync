package log

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// OutputWriter はユーザー向け出力を扱うインターフェース
type OutputWriter interface {
	// 標準メッセージを出力
	Print(msg string)
	Printf(format string, args ...interface{})

	// 進捗状況を出力
	PrintProgress(msg string)

	// 成功メッセージを出力
	PrintSuccess(msg string)

	// エラーメッセージを出力
	PrintError(err error)

	// 詳細モード時のみ出力
	PrintVerbose(msg string)

	// ユーザーの確認を得る
	Confirm(prompt string) bool
}

// ConsoleOutput は標準出力へのOutputWriter実装
type ConsoleOutput struct {
	Verbose bool // 詳細出力モード
	Color   bool // カラー出力の有効/無効
}

// Print outputs a standard message
func (c *ConsoleOutput) Print(msg string) {
	fmt.Println(msg)
}

// Printf outputs a formatted standard message
func (c *ConsoleOutput) Printf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// PrintProgress outputs a progress message
func (c *ConsoleOutput) PrintProgress(msg string) {
	if c.Color {
		_, err := color.New(color.FgBlue).Println("➜ " + msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print progress message: %v\n", err)
		}
	} else {
		fmt.Println("➜ " + msg)
	}
}

// PrintSuccess outputs a success message
func (c *ConsoleOutput) PrintSuccess(msg string) {
	if c.Color {
		_, err := color.New(color.FgGreen).Println("✓ " + msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print success message: %v\n", err)
		}
	} else {
		fmt.Println("✓ " + msg)
	}
}

// PrintError outputs an error message
func (c *ConsoleOutput) PrintError(err error) {
	if c.Color {
		_, printErr := color.New(color.FgRed).Fprintf(os.Stderr, "✗ Error: %v\n", err)
		if printErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to print error message: %v\n", printErr)
		}
	} else {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
	}
}

// PrintVerbose outputs a message only in verbose mode
func (c *ConsoleOutput) PrintVerbose(msg string) {
	if c.Verbose {
		if c.Color {
			_, err := color.New(color.FgCyan).Println("ℹ " + msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to print verbose message: %v\n", err)
			}
		} else {
			fmt.Println("ℹ " + msg)
		}
	}
}

// Confirm asks for user confirmation
func (c *ConsoleOutput) Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		if c.Color {
			_, err := color.New(color.FgYellow).Printf("%s [y/N]: ", prompt)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to print confirmation prompt: %v\n", err)
				return false
			}
		} else {
			fmt.Printf("%s [y/N]: ", prompt)
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "", "n", "no":
			return false
		}
	}
}
