package oviewer

import (
	"fmt"
	"log"
	"strings"
)

// Go to the top line.
func (root *Root) moveTop() {
	root.moveLine(0)
}

// Go to the bottom line.
func (root *Root) moveBottom() {
	tx, tn := root.bottomLineNum(root.Doc.BufEndNum())
	root.Doc.topLN = tn
	root.Doc.topLX = tx
}

// Move to the specified line.
func (root *Root) moveLine(lN int) {
	root.resetSelect()
	root.Doc.topLN = lN
	root.Doc.topLX = 0
}

// Move up one screen.
func (root *Root) movePgUp() {
	root.resetSelect()
	root.moveNumUp(root.statusPos - root.headerLen())
}

// Moves down one screen.
func (root *Root) movePgDn() {
	root.resetSelect()

	y := root.bottomLN - root.Doc.Header
	x := root.bottomLX
	root.limitMoveDown(x, y)
}

func (root *Root) limitMoveDown(x int, y int) {
	m := root.Doc
	endNum := root.Doc.BufEndNum()
	if y+root.vHight >= endNum {
		tx, tn := root.bottomLineNum(endNum)
		if y > tn || (y == tn && x > tx) {
			if m.topLN < tn || (m.topLN == tn && m.topLX < tx) {
				m.topLN = tn
				m.topLX = tx
			}
			return
		}
	}
	m.topLN = y
	m.topLX = x
}

// Moves up half a screen.
func (root *Root) moveHfUp() {
	root.resetSelect()
	root.moveNumUp((root.statusPos - root.headerLen()) / 2)
}

// Moves down half a screen.
func (root *Root) moveHfDn() {
	root.resetSelect()
	root.moveNumDown((root.statusPos - root.headerLen()) / 2)
}

// numOfSlice returns what number x is in slice.
func numOfSlice(listX []int, x int) int {
	for n, v := range listX {
		if v >= x {
			return n
		}
	}
	return len(listX) - 1
}

// numOfReverseSlice returns what number x is from the back of slice.
func numOfReverseSlice(listX []int, x int) int {
	for n := len(listX) - 1; n >= 0; n-- {
		if listX[n] <= x {
			return n
		}
	}
	return 0
}

// Moves up by the specified number of y.
func (root *Root) moveNumUp(moveY int) {
	if !root.Doc.WrapMode {
		root.Doc.topLN -= moveY
		return
	}

	// WrapMode
	num := root.Doc.topLN + root.Doc.Header
	root.Doc.topLX, num = root.findNumUp(root.Doc.topLX, num, moveY)
	root.Doc.topLN = num - root.Doc.Header
}

// Moves down by the specified number of y.
func (root *Root) moveNumDown(moveY int) {
	m := root.Doc
	num := m.topLN + m.Header
	if !m.WrapMode {
		root.limitMoveDown(0, num+moveY)
		return
	}

	// WrapMode
	x := m.topLX
	listX, err := root.leftMostX(num)
	if err != nil {
		log.Println(err, num)
		return
	}
	n := numOfReverseSlice(listX, x)

	for y := 0; y < moveY; y++ {
		if n >= len(listX) {
			num++
			if num > m.BufEndNum() {
				break
			}
			listX, err = root.leftMostX(num)
			if err != nil {
				log.Println(err, num)
				return
			}
			n = 0
		}
		x = 0
		if len(listX) > 0 && n < len(listX) {
			x = listX[n]
		}
		n++
	}

	root.limitMoveDown(x, num-m.Header)
}

// Move up one line.
func (root *Root) moveUp() {
	root.resetSelect()

	m := root.Doc
	if m.topLN == 0 && m.topLX == 0 {
		return
	}

	if !m.WrapMode {
		m.topLN--
		m.topLX = 0
		return
	}

	// WrapMode.
	// Same line.
	if m.topLX > 0 {
		listX, err := root.leftMostX(m.topLN + m.Header)
		if err != nil {
			log.Println(err)
			return
		}
		for n, x := range listX {
			if x >= m.topLX {
				m.topLX = listX[n-1]
				return
			}
		}
	}

	// Previous line.
	m.topLN--
	if m.topLN < 0 {
		m.topLN = 0
		m.topLX = 0
		return
	}
	listX, err := root.leftMostX(m.topLN + m.Header)
	if err != nil {
		log.Println(err)
		return
	}
	if len(listX) > 0 {
		m.topLX = listX[len(listX)-1]
		return
	}
	m.topLX = 0
}

// Move down one line.
func (root *Root) moveDown() {
	root.resetSelect()

	m := root.Doc
	num := m.topLN

	if !m.WrapMode {
		num++
		root.limitMoveDown(0, num)
		return
	}

	// WrapMode
	listX, err := root.leftMostX(m.topLN + m.Header)
	if err != nil {
		log.Println(err)
		return
	}

	for _, x := range listX {
		if x > m.topLX {
			root.limitMoveDown(x, num)
			return
		}
	}

	// Next line.
	num++
	m.topLX = 0
	root.limitMoveDown(m.topLX, num)
}

// Move to the left.
func (root *Root) moveLeft() {
	root.resetSelect()

	m := root.Doc
	if m.ColumnMode {
		if m.columnNum > 0 {
			m.columnNum--
			m.x = root.columnModeX()
		}
		return
	}
	if m.WrapMode {
		return
	}
	m.x--
	if m.x < root.minStartX {
		m.x = root.minStartX
	}
}

// Move to the right.
func (root *Root) moveRight() {
	root.resetSelect()

	m := root.Doc
	if m.ColumnMode {
		m.columnNum++
		m.x = root.columnModeX()
		return
	}
	if m.WrapMode {
		return
	}
	m.x++
}

// columnModeX returns the actual x from m.columnNum.
func (root *Root) columnModeX() int {
	m := root.Doc
	// m.Header+10 = Maximum columnMode target.
	for i := 0; i < m.Header+10; i++ {
		lc, err := m.lineToContents(m.topLN+m.Header+i, m.TabWidth)
		if err != nil {
			continue
		}
		lineStr, byteMap := contentsToStr(lc)
		// Skip lines that do not contain a delimiter.
		if !strings.Contains(lineStr, m.ColumnDelimiter) {
			continue
		}

		start, end := rangePosition(lineStr, m.ColumnDelimiter, m.columnNum)
		if start < 0 || end < 0 {
			m.columnNum--
			start, end = rangePosition(lineStr, m.ColumnDelimiter, m.columnNum)
		}
		sx := byteMap[start]
		ex := byteMap[end] + 10
		if root.vWidth > ex {
			return 0
		}
		if ex-root.vWidth > 0 {
			return ex - root.vWidth
		}
		return sx
	}
	m.columnNum = 0
	return 0
}

// Move to the left by half a screen.
func (root *Root) moveHfLeft() {
	m := root.Doc

	if m.WrapMode {
		return
	}
	root.resetSelect()
	moveSize := (root.vWidth / 2)
	if m.x > 0 && (m.x-moveSize) < 0 {
		m.x = 0
	} else {
		m.x -= moveSize
		if m.x < root.minStartX {
			m.x = root.minStartX
		}
	}
}

// Move to the right by half a screen.
func (root *Root) moveHfRight() {
	m := root.Doc

	if m.WrapMode {
		return
	}
	root.resetSelect()
	if m.x < 0 {
		m.x = 0
	} else {
		m.x += (root.vWidth / 2)
	}
}

// bottomLineNum returns the display start line
// when the last line number as an argument.
func (root *Root) bottomLineNum(lN int) (int, int) {
	hight := (root.vHight - root.headerLen()) - 2
	if lN < root.headerLen() {
		return 0, 0
	}

	if !root.Doc.WrapMode {
		return 0, lN - (hight + root.headerLen())
	}

	// WrapMode
	lX, lN := root.findNumUp(0, lN, hight)
	return lX, lN - root.Doc.Header
}

// findNumUp finds lX, lN when the number of lines is moved up from lX, lN.
func (root *Root) findNumUp(lX int, lN int, upY int) (int, int) {
	listX, err := root.leftMostX(lN)
	n := 0
	if err != nil {
		// lN has no lines.
		root.debugMessage(fmt.Sprintf("%s:%d", err.Error(), lN))
	} else {
		n = numOfSlice(listX, lX)
	}

	for y := upY; y > 0; y-- {
		if n <= 0 {
			lN--
			if lN < root.Doc.Header {
				lN = 0
				lX = 0
				break
			}
			listX, err = root.leftMostX(lN)
			if err != nil {
				log.Println(err, "findNumUp", lN)
				return 0, 0
			}
			n = len(listX)
		}
		if n > 0 {
			lX = listX[n-1]
		} else {
			lX = 0
		}
		n--
	}
	return lX, lN
}
