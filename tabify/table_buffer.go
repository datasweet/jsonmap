package tabify

type tableBuffer struct {
	deep    int
	buffers map[int][]*rowBuffer
}

// newTableBuffer to create a new table in memory
func newTableBuffer() *tableBuffer {
	buffers := make(map[int][]*rowBuffer)
	return &tableBuffer{
		buffers: buffers,
	}
}

// openRow to create a new row
func (tb *tableBuffer) openRow() {
	rb := &rowBuffer{
		parent: len(tb.buffers[tb.deep]) - 1,
	}

	tb.deep++
	tb.buffers[tb.deep] = append(tb.buffers[tb.deep], rb)
}

// closeRow to close the current row and write values
func (tb *tableBuffer) closeRow() {
	if tb.deep > 0 {
		tb.deep--
	}
}

// cell to create a new cell in row
func (tb *tableBuffer) cell(val *nodeValue) {
	if len(tb.buffers) == 0 {
		tb.openRow()
	}
	buffers := tb.buffers[tb.deep]
	curr := buffers[len(buffers)-1]
	curr.values = append(curr.values, val)
}

func (tb *tableBuffer) getMaxDeep() int {
	var max int
	for d := range tb.buffers {
		if d > max {
			max = d
		}
	}
	return max
}

// write our buffer into a tablewriter
func (tb *tableBuffer) write(tw TableWriter) {
	if tb.deep > 0 {
		tb.closeRow()
	}

	max := tb.getMaxDeep()

	for _, row := range tb.buffers[max] {
		tw.OpenRow()

		for _, cell := range row.values {
			tw.Cell(cell.key, cell.value, cell.deep)
		}

		curr := row
		for deep := max - 1; deep > 0; deep-- {
			prow := tb.buffers[deep][curr.parent]
			for _, pcell := range prow.values {
				tw.Cell(pcell.key, pcell.value, pcell.deep)
			}
			curr = prow
		}
		tw.CloseRow()
	}
}

type rowBuffer struct {
	parent int
	values []*nodeValue
}
