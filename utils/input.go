package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/unicornfairy864/unf-harness/constants"
)

var stdin = bufio.NewReader(os.Stdin)

// InputUserPrompt 读用户输入，支持 "\" 续行；EOF 返回 constants.CliExitCommand。
func InputUserPrompt() string {

	var sb strings.Builder

	for {
		if sb.Len() == 0 {
			fmt.Print(constants.CliUserPrompt)
		} else {
			fmt.Print(constants.CliContPrompt)
		}

		line, err := stdin.ReadString('\n')

		// 不用管是哪个平台：右侧的 \r 和 \n 一起删干净（cutset，不是子串）
		line = strings.TrimRight(line, "\r\n")

		// 续行判定：删掉行末空白后以 \ 结尾，且 \ 前面是空白（或整行就一个 \）
		body := strings.TrimRight(line, " \t")
		head := strings.TrimSuffix(body, `\`)
		isCont := strings.HasSuffix(body, `\`) &&
			(head == "" || strings.HasSuffix(head, " ") || strings.HasSuffix(head, "\t"))

		if isCont {
			sb.WriteString(strings.TrimRight(head, " \t")) // 续行符连同它前面的空白一起收走
			sb.WriteString("\n")
			if err != nil {
				break // 续行符后直接 EOF：攒下的内容照常返回
			}
			continue
		}

		// 非续行：纯空输入才提示重来（有累积内容时，空行算有效的一行）
		if sb.Len() == 0 && strings.TrimSpace(line) == "" {
			if err != nil {
				return constants.CliExitCommand
			}
			fmt.Println(constants.CliInputAgainPrompt)
			continue
		}

		sb.WriteString(line)
		break
	}

	return sb.String()
}

// InputLine 读一行原始输入，去掉行尾 \r\n；不处理续行、不拦空行。
// ok 为 false 表示 EOF。给「y/n/d」这类单行确认用。
func InputLine() (string, bool) {
	line, err := stdin.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	if err != nil && line == "" {
		return "", false
	}
	return line, true
}
