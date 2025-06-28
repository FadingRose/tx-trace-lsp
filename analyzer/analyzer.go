package analyzer

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"go.lsp.dev/protocol"
)

// callInfo 用于存储调用堆栈上的信息
type callInfo struct {
	startLine int // 该调用的起始行号 (0-indexed)
	indent    int // 该调用的缩进级别
}

// CalculateFoldingRanges 接收语法树的根节点，返回所有语义折叠范围
func CalculateFoldingRanges(rootNode *sitter.Node, code []byte) []protocol.FoldingRange {
	var ranges []protocol.FoldingRange
	var callStack []callInfo

	// 遍历语法树中的每一行 (trace_line 节点)
	for i := 0; i < int(rootNode.ChildCount()); i++ {
		traceNode := rootNode.Child(i)

		// 我们只关心 trace_line 节点
		if traceNode.Type() != "trace_line" {
			continue
		}

		// 1. 获取当前行的缩进级别
		prefixNode := traceNode.ChildByFieldName("prefix")
		currentIndent := getIndentLevel(prefixNode, code)
		currentLine := int(traceNode.StartPoint().Row)

		// 2. 检查缩进，判断是否有函数调用结束
		// 如果当前行的缩进小于或等于栈顶调用的缩进，说明栈顶的调用已经结束
		for len(callStack) > 0 && currentIndent <= callStack[len(callStack)-1].indent {
			// 从堆栈中弹出已结束的调用
			poppedCall := callStack[len(callStack)-1]
			callStack = callStack[:len(callStack)-1]

			// 结束行是当前行的前一行
			endLine := currentLine - 1

			// 只有当范围有效时才添加 (至少跨一行)
			if endLine > poppedCall.startLine {
				ranges = append(ranges, protocol.FoldingRange{
					StartLine: uint32(poppedCall.startLine),
					EndLine:   uint32(endLine),
					Kind:      protocol.RegionFoldingRange, // 'region' 是一种通用的折叠类型
				})
			}
		}

		// 3. 检查当前行是否是一个新的调用
		callNode := traceNode.ChildByFieldName("call")
		if callNode != nil {
			// 将新调用压入堆栈
			callStack = append(callStack, callInfo{
				startLine: currentLine,
				indent:    currentIndent,
			})
		}
	}

	// 4. 处理文件末尾仍未闭合的调用
	endOfFile_Line := int(rootNode.EndPoint().Row)
	for len(callStack) > 0 {
		poppedCall := callStack[len(callStack)-1]
		callStack = callStack[:len(callStack)-1]
		if endOfFile_Line > poppedCall.startLine {
			ranges = append(ranges, protocol.FoldingRange{
				StartLine: uint32(poppedCall.startLine),
				EndLine:   uint32(endOfFile_Line),
				Kind:      protocol.RegionFoldingRange,
			})
		}
	}

	return ranges
}

// getIndentLevel 是一个辅助函数，用于计算行的缩进级别
// 这里的实现很简单，直接计算前缀字符串的长度。
// 您可以根据需要实现更复杂的逻辑。
func getIndentLevel(prefixNode *sitter.Node, code []byte) int {
	if prefixNode == nil {
		return 0
	}
	// 简单地使用前缀的字符长度作为缩进级别
	// 注意：为了与 ASCII 和多字节字符（如 │）兼容，最好计算 Rune 的数量
	return len([]rune(strings.TrimRight(prefixNode.Content(code), " ")))
}
