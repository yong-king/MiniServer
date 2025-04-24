package base62

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestInt2String(t *testing.T) {
	tests := []struct {
		name string
		seq  uint64
		want string
	}{
		{name: "case 0", seq: 0, want: "0"},
		{name: "case 61", seq: 61, want: "Z"},
		{name: "case 62", seq: 62, want: "10"},
		{name: "case 6347", seq: 6347, want: "1En"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int2String(tt.seq); got != tt.want {
				t.Errorf("Int2String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString2Int(t *testing.T) {
	c.Convey("case 1En", t, func ()  {
		s :=  "1En"
		s2i := String2Int(s)
		c.So(s2i, c.ShouldEqual, 6347)
	})
	c.Convey("case 0", t, func ()  {
		s :=  "0"
		s2i := String2Int(s)
		c.So(s2i, c.ShouldEqual, 0)
	})
	c.Convey("case Z", t, func ()  {
		s :=  "Z"
		s2i := String2Int(s)
		c.ShouldEqual(s2i, 61)
	})
}
