package godump

import "testing"

func TestDumper(t *testing.T) {
	testCases := []struct {
		inputVar any
		expected string
	}{
		{int(123), "123"},
		{int8(123), "123"},
		{int16(123), "123"},
		{int32(123), "123"},
		{int64(123), "123"},
		{int(-123), "-123"},
		{int8(-123), "-123"},
		{int16(-123), "-123"},
		{int32(-123), "-123"},
		{int64(-123), "-123"},
		{uint(123), "123"},
		{uint8(123), "123"},
		{uint16(123), "123"},
		{uint32(123), "123"},
		{uint64(123), "123"},
		{float32(12.3), "12.300000190734863"},   // due to float64 casting caused by the `reflect` package
		{float32(-12.3), "-12.300000190734863"}, // due to float64 casting caused by the `reflect` package
		{float64(12.3), "12.3"},
		{float64(-12.3), "-12.3"},
		{complex64(12.3), "(12.300000190734863+0i)"},
		{complex64(-12.3), "(-12.300000190734863+0i)"},
		{complex128(12.3), "(12.3+0i)"},
		{complex128(-12.3), "(-12.3+0i)"},
		{true, "true"},
		{false, "false"},
		{"hello world", `"hello world"`},
		{func(i int) int { return i }, "func(int) int"},
		{func(int) {}, "func(int)"},
		{func() int { return 123 }, "func() int"},
		{func() {}, "func()"},
		{make([]any, 0, 5), "[]interface {}:0:5 {\n}"},
		{make([]any, 3, 5), `[]interface {}:3:5 {
   uninitialized,
   uninitialized,
   uninitialized,
}`},
		{
			[]int{1, 2, -3},
			`[]int:3:3 {
   1,
   2,
   -3,
}`,
		},
		{
			[]int8{1, 2, -3},
			`[]int8:3:3 {
   1,
   2,
   -3,
}`,
		},
		{
			[]int16{1, 2, -3},
			`[]int16:3:3 {
   1,
   2,
   -3,
}`,
		},
		{
			[]int32{1, 2, -3},
			`[]int32:3:3 {
   1,
   2,
   -3,
}`,
		},
		{
			[]int64{1, 2, -3},
			`[]int64:3:3 {
   1,
   2,
   -3,
}`,
		},
	}

	for i, tc := range testCases {
		var d dumper
		d.dump(tc.inputVar)

		if returned := d.buf.String(); returned != tc.expected {
			t.Fatalf(`Case#%d failed, dumper returned unuexpected results : "%s", expected "%s"`, i, returned, tc.expected)
		}
	}

}
