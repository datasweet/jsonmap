package tabify

type tableBuffer struct {
	deep    int
	buffers []*rowBuffer
}

// newTableBuffer to create a new table in memory
func newTableBuffer() *tableBuffer {
	return &tableBuffer{}
}

// openRow to create a new row
func (writer *tableBuffer) openRow() {
	if writer.deep == 0 {
		writer.buffers = append(writer.buffers, newRowBuffer())
	}
	writer.deep++
}

// closeRow to close the current row and write values
func (writer *tableBuffer) closeRow() {
	// PANIC IF NO OPENED ?
	buffer := writer.buffers[len(writer.buffers)-1]
	buffer.release(writer.deep)
	writer.deep--
}

// cell to create a new cell in row
func (writer *tableBuffer) cell(val *nodeValue) {
	if len(writer.buffers) == 0 {
		writer.openRow()
	}
	buffer := writer.buffers[len(writer.buffers)-1]
	buffer.bufferize(writer.deep, val)
}

// write our buffer into a tablewriter
func (writer *tableBuffer) write(tw TableWriter) {
	if writer.deep > 0 {
		writer.closeRow()
	}

	for _, r := range writer.buffers {
		r.write(tw)
	}
}

type rowBuffer struct {
	buffer map[int][]*nodeValue
	values map[int][][]*nodeValue
}

func newRowBuffer() *rowBuffer {
	return &rowBuffer{
		buffer: make(map[int][]*nodeValue),
		values: make(map[int][][]*nodeValue),
	}
}

func (trb *rowBuffer) clear(deep int) {
	delete(trb.buffer, deep)
}

func (trb *rowBuffer) bufferize(deep int, value *nodeValue) {
	trb.buffer[deep] = append(trb.buffer[deep], value)
}

func (trb *rowBuffer) release(deep int) {
	trb.values[deep] = append(trb.values[deep], trb.buffer[deep])
	trb.clear(deep)
}

func (trb *rowBuffer) getMaxDeep() int {
	var max int
	for d := range trb.values {
		if d > max {
			max = d
		}
	}
	return max
}

func (trb *rowBuffer) extractValues(deep int, tw TableWriter) {
	list := trb.values[deep]
	for _, values := range list {
		// Copy values
		for _, v := range values {
			tw.Cell(v.key, v.value)
		}
	}
}

func (trb *rowBuffer) write(tw TableWriter) {
	max := trb.getMaxDeep()
	list := trb.values[max]
	for _, values := range list {
		tw.OpenRow()

		// Copy current values.
		for _, v := range values {
			tw.Cell(v.key, v.value)
		}

		for i := max - 1; i > 0; i-- {
			trb.extractValues(i, tw)
		}

		tw.CloseRow()
	}
}
