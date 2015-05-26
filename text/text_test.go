package text

import (
	"fmt"
	"testing"
)

func TestOverall(t *testing.T) {
	txt := New(nil)
	txt.checkPiecesCnt(t, 2)
	txt.checkContent(t, "")

	txt.insertString(0, "")
	txt.checkPiecesCnt(t, 2)
	txt.checkContent(t, "")

	txt.insertString(0, "All work makes John a dull boy")
	txt.checkPiecesCnt(t, 3)
	txt.checkContent(t, "All work makes John a dull boy")

	txt.insertString(9, "and no playing ")
	txt.checkPiecesCnt(t, 6)
	txt.checkContent(t, "All work and no playing makes John a dull boy")

	txt.Delete(20, 14)
	txt.checkContent(t, "All work and no play a dull boy")

	txt.insertString(20, " makes Jack")
	txt.checkContent(t, "All work and no play makes Jack a dull boy")

	txt.Undo()
	txt.checkContent(t, "All work and no play a dull boy")
	txt.Undo()
	txt.checkContent(t, "All work and no playing makes John a dull boy")
	txt.Undo()
	txt.checkContent(t, "All work makes John a dull boy")

	txt.Redo()
	txt.checkContent(t, "All work and no playing makes John a dull boy")
	txt.Redo()
	txt.checkContent(t, "All work and no play a dull boy")
	txt.Redo()
	txt.checkContent(t, "All work and no play makes Jack a dull boy")
	txt.Redo()
	txt.checkContent(t, "All work and no play makes Jack a dull boy")
}

func TestCacheInsertAndDelete(t *testing.T) {
	txt := New([]byte("testing insertation"))
	txt.checkPiecesCnt(t, 3)
	txt.checkContent(t, "testing insertation")

	txt.cacheInsertString(8, "caching")
	txt.checkPiecesCnt(t, 6)
	txt.checkContent(t, "testing cachinginsertation")

	txt.cacheInsertString(15, " ")
	txt.checkPiecesCnt(t, 6)
	txt.checkContent(t, "testing caching insertation")

	txt.cacheDelete(12, 3)
	txt.checkPiecesCnt(t, 6)
	txt.checkContent(t, "testing cach insertation")

	txt.cacheInsertString(12, "ed")
	txt.checkPiecesCnt(t, 6)
	txt.checkContent(t, "testing cached insertation")
}

func TestDelete(t *testing.T) {
	txt := New([]byte("and what is a dream?"))
	txt.insertString(9, "exactly ")
	txt.checkContent(t, "and what exactly is a dream?")

	txt.delete(22, 2000)
	txt.checkContent(t, "and what exactly is a ")
	txt.insertString(22, "joke?")
	txt.checkContent(t, "and what exactly is a joke?")

	cases := []struct {
		pos, len int
		expected string
	}{
		{9, 8, "and what is a joke?"},
		{9, 13, "and what joke?"},
		{5, 6, "and wactly is a joke?"},
		{9, 14, "and what oke?"},
		{11, 3, "and what exly is a joke?"},
	}
	for _, c := range cases {
		txt.delete(c.pos, c.len)
		txt.checkContent(t, c.expected)
		txt.Undo()
		txt.checkContent(t, "and what exactly is a joke?")
	}
}

func TestGroupChanges(t *testing.T) {
	txt := New([]byte("group 1, group 2, group 3"))
	txt.checkPiecesCnt(t, 3)
	txt.GroupChanges()

	txt.cacheDelete(0, 6)
	txt.checkContent(t, "1, group 2, group 3")

	txt.cacheDelete(3, 6)
	txt.checkContent(t, "1, 2, group 3")

	txt.cacheDelete(6, 6)
	txt.checkContent(t, "1, 2, 3")

	txt.Undo()
	txt.checkContent(t, "group 1, group 2, group 3")
	txt.Undo()
	txt.checkContent(t, "group 1, group 2, group 3")

	txt.Redo()
	txt.checkContent(t, "1, 2, 3")
	txt.Redo()
	txt.checkContent(t, "1, 2, 3")
}

func TextSaving(t *testing.T) {
	txt := New(nil)

	txt.checkModified(t, false)
	txt.insertString(0, "stars can frighten")
	txt.checkModified(t, true)

	txt.Save()
	txt.checkModified(t, false)

	txt.insertString(0, "Neptun, Titan, ")
	txt.checkModified(t, true)
	txt.Undo()
	txt.checkModified(t, false)

	txt.Redo()
	txt.checkModified(t, true)

	txt.Save()
	txt.checkModified(t, false)

	txt = New([]byte("my book is closed"))
	txt.checkModified(t, false)

	txt.insertString(17, ", I read no more")
	txt.checkModified(t, true)
	txt.Undo()
	txt.checkModified(t, false)

	txt.Save()
	txt.checkModified(t, false)
}

func (txt *Text) checkPiecesCnt(t *testing.T, expected int) {
	if txt.piecesCnt != expected {
		t.Errorf("got %d pieces, want %d", txt.piecesCnt, expected)
	}
}

func (txt *Text) checkContent(t *testing.T, expected string) {
	c := txt.AllContent()
	if c != expected {
		t.Errorf("got '%s', want '%s'", c, expected)
	}
}

func (t *Text) insertString(pos int, data string) {
	t.CommitChanges()
	t.cacheInsertString(pos, data)
}

func (t *Text) cacheInsertString(pos int, data string) {
	err := t.Insert(pos, []byte(data))
	if err != nil {
		panic(err)
	}
}

func (t *Text) delete(pos, length int) {
	t.CommitChanges()
	t.cacheDelete(pos, length)
}

func (t *Text) cacheDelete(pos, length int) {
	t.Delete(pos, length)
}

func (t *Text) printPieces() {
	for p := t.begin; p != nil; p = p.next {
		prev, next := 0, 0
		if p.prev != nil {
			prev = p.prev.id
		}
		if p.next != nil {
			next = p.next.id
		}
		fmt.Printf("%d, p:%d, n:%d = %s\n", p.id, prev, next, string(p.data))
	}
	fmt.Println()
}

func (txt *Text) checkModified(t *testing.T, expected bool) {
	if txt.Modified() != expected {
		if expected {
			t.Errorf("text should be modified")
		} else {
			t.Errorf("text should not be modified")
		}
	}
}
