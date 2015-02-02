package sudoku

import "testing"

func TestSolve(t *testing.T) {
	grid := NewGrid()
	if ok, err := grid.Set("_______6_28______4__7__58__5__34__2_4__5_1__8_1__76__3__51__2__3______81_9_______"); !ok {
		t.Error(err)
	}
	if !grid.Solve() {
		t.Error("grid should be solvable")
	}
	if grid.String() != "153498762289763514647215839578349126436521978912876453865134297324957681791682345" {
		t.Error("bad solution")
	}
}
