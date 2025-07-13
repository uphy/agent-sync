package agent

import (
	"fmt"
	"strings"
)

func formatMCP(agent, command string, args ...string) string {
	a := ""
	if len(args) == 0 {
		a = ""
	} else {
		a = fmt.Sprintf("(%s)", strings.Join(args, ", "))
	}
	return fmt.Sprintf("MCP tool `%s.%s%s`", agent, command, a)
}
