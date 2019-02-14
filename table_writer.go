package jsonmap

type tableWriter struct {
	deep    int
	buffers []*rowBuffer
}

// newTableWriter to create a new writer
func newTableWriter() *tableWriter {
	return &tableWriter{}
}

// openRow to create a new row
func (writer *tableWriter) openRow() {
	if writer.deep == 0 {
		writer.buffers = append(writer.buffers, newRowBuffer())
	}
	writer.deep++
}

// closeRow to close the current row and write values
func (writer *tableWriter) closeRow() {
	// PANIC IF NO OPENED ?
	buffer := writer.buffers[len(writer.buffers)-1]
	buffer.release(writer.deep)
	writer.deep--
}

// cell to create a new cell in row
func (writer *tableWriter) cell(val *nodeValue) {
	if len(writer.buffers) == 0 {
		writer.openRow()
	}
	buffer := writer.buffers[len(writer.buffers)-1]
	buffer.bufferize(writer.deep, val)
}

// write to into a dictionary our tabulars datas.
func (writer *tableWriter) write() []map[string]interface{} {
	var datas []map[string]interface{}

	if writer.deep > 0 {
		writer.closeRow()
	}

	for _, r := range writer.buffers {
		datas = append(datas, r.write()...)
	}
	return datas
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

func (trb *rowBuffer) extractValues(deep int, to map[string]interface{}) {
	list := trb.values[deep]
	for _, values := range list {
		// Copy values
		for _, v := range values {
			to[v.key] = v.value
		}
	}
}

func (trb *rowBuffer) write() []map[string]interface{} {
	var datas []map[string]interface{}
	max := trb.getMaxDeep()

	list := trb.values[max]
	for _, values := range list {
		row := make(map[string]interface{})

		// Copy current values.
		for _, v := range values {
			row[v.key] = v.value
		}

		for i := max - 1; i > 0; i-- {
			trb.extractValues(i, row)
		}

		datas = append(datas, row)
	}

	return datas
}
