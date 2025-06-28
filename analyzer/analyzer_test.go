package analyzer

import (
	"fadingrose/tx-trace-lsp/parser"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"go.lsp.dev/protocol"
)

// byRange 是为了给 []protocol.FoldingRange 排序，以确保测试的稳定性
type byRange []protocol.FoldingRange

func (r byRange) Len() int      { return len(r) }
func (r byRange) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byRange) Less(i, j int) bool {
	if r[i].StartLine != r[j].StartLine {
		return r[i].StartLine < r[j].StartLine
	}
	// 如果起始行相同，结束行更晚的范围更大（在外层）
	return r[i].EndLine > r[j].EndLine
}

func TestCalculateFoldingRanges(t *testing.T) {
	p := parser.NewParser()
	if p == nil {
		t.Fatal("Failed to create parser instance")
	}

	testFilePath := filepath.Join("testdata", "sample.tx")
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", testFilePath, err)
	}

	tree := p.MustParse(content)
	if tree == nil {
		t.Fatal("Failed to parse the test file")
	}

	if tree.RootNode() == nil {
		t.Fatal("Parsed tree has no root node")
	}

	actualRanges := CalculateFoldingRanges(tree.RootNode(), content)

	// - call C: start 2, end 3
	// - call B: start 1, end 4
	// - call D: start 5, end 6
	// - call A: start 0, end 7 (文件末尾)
	expectedRanges := []protocol.FoldingRange{
		{StartLine: 2, EndLine: 3, Kind: protocol.RegionFoldingRange},
		{StartLine: 1, EndLine: 4, Kind: protocol.RegionFoldingRange},
		{StartLine: 5, EndLine: 6, Kind: protocol.RegionFoldingRange},
		{StartLine: 0, EndLine: 7, Kind: protocol.RegionFoldingRange},
	}

	// 6. 验证结果
	// 为了避免因实现细节导致顺序不一致，我们在比较前先对两个切片进行排序
	sort.Sort(byRange(actualRanges))
	sort.Sort(byRange(expectedRanges))

	if !reflect.DeepEqual(actualRanges, expectedRanges) {
		log.Printf("Actual Ranges: %+v\n", actualRanges)
		log.Printf("Expected Ranges: %+v\n", expectedRanges)
		t.Error("Calculated folding ranges do not match expected ranges.")
	}
}
