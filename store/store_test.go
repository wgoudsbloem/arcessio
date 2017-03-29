package store

import (
	"bytes"
	"testing"
)

type TestSeeker struct {
	bytes.Buffer
}

func (t *TestSeeker) Seek(offset int64, whence int) (int64, error) {
	return int64(t.Len()), nil
}

func (t *TestSeeker) ReadAt(p []byte, off int64) (n int, err error) {
	return t.Read(p)
}

// store a string into a buffer
func TestPut(t *testing.T) {
	in1 := "teststring"
	in2 := "teststring2"
	exp := in1 + "\n" + in2 + "\n"
	var bb TestSeeker
	s := store{&bb, 0}
	_, err := s.Put([]byte(in1))
	if err != nil {
		t.Error(err)
	}
	offset, err := s.Put([]byte(in2))
	if err != nil {
		t.Error(err)
	}
	expLen := int64(len(exp))
	if offset != expLen {
		t.Errorf("want %v, got %v", expLen, offset)
	}
	res, _ := bb.ReadString('\t')
	if exp != res {
		t.Errorf("want %v, got %v", exp, res)
	}
}

func TestGet(t *testing.T) {
	in1 := "teststring"
	in2 := "teststring2"
	var bb TestSeeker
	s := store{&bb, 0}
	_, err := s.Put([]byte(in1))
	if err != nil {
		t.Error(err)
	}
	offset, err := s.Put([]byte(in2))
	if err != nil {
		t.Error(err)
	}
	b, err := s.Get()
	if err != nil {
		t.Error(err)
	}
	if string(b) != in1 {
		t.Errorf("want %v got %v", in1, string(b))
	}
	if string(b) == in2 {
		t.Errorf("want %v got %v", in1, string(b))
	}
	s2 := store{&bb, offset}
	b2, err := s2.Get()
	if err != nil {
		t.Error(err)
	}
	if string(b2) != in2 {
		t.Errorf("want [%v] got [%v]", in2, string(b2))
	}
}

func TestInternalLastEntry(t *testing.T) {
	expVal1 := `{"test1":"val1"}`
	expVal2 := `{"test2":"val2"}`
	testVal1 := expVal1 + "\n"
	testVal2 := expVal2 + "\n"
	testVal3 := testVal1 + testVal2
	res1 := lastEntry([]byte(testVal1))
	if string(res1) != expVal1 {
		t.Errorf("want: '%v' got: '%v'", expVal1, string(res1))
	}
	res2 := lastEntry([]byte(testVal3))
	if string(res2) != expVal2 {
		t.Errorf("want: '%v' got: '%v'", expVal2, string(res2))
	}
}
