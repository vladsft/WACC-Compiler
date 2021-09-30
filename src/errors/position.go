package errors

import "fmt"

//Position is a pair of integers representing the line and column of some code
type Position struct {
	startLine, startCol int
	endLine, endCol     int
}

//NewPosition creates a new position
func NewPosition(startLine, startCol, endLine, endCol int) Position {
	return Position{
		startLine: startLine,
		startCol:  startCol,
		endLine:   endLine,
		endCol:    endCol,
	}
}

func (p Position) String() string {
	return fmt.Sprintf("Line [%d:%d-%d:%d]", p.startLine, p.startCol, p.endLine, p.endCol)
}
