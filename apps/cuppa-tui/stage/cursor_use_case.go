package stage

// SetCursor puts the keyboard cursor on a canvas cell and says whether to draw
// it. The pane scrolls so the cell is in view.
func (m *Model) SetCursor(x, y int, show bool) {
	moved := m.cursor != [2]int{x, y} || !m.showCursor
	m.cursor, m.showCursor = [2]int{x, y}, show
	// Scrolling with the wheel must not be undone every frame.
	if !show || !moved || m.w == 0 || m.h == 0 {
		return
	}
	switch {
	case x < m.offX:
		m.offX = x
	case x >= m.offX+m.w:
		m.offX = x - m.w + 1
	}
	switch {
	case y < m.offY:
		m.offY = y
	case y >= m.offY+m.h:
		m.offY = y - m.h + 1
	}
	m.clampOffset()
}
