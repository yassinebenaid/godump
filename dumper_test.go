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
   nil,
   nil,
   nil,
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
		{
			[]uint{1, 2, 3},
			`[]uint:3:3 {
   1,
   2,
   3,
}`,
		},
		{
			[]uint8{1, 2, 3},
			`[]uint8:3:3 {
   1,
   2,
   3,
}`,
		},
		{
			[]uint16{1, 2, 3},
			`[]uint16:3:3 {
   1,
   2,
   3,
}`,
		},
		{
			[]uint32{1, 2, 3},
			`[]uint32:3:3 {
   1,
   2,
   3,
}`,
		},
		{
			[]uint64{1, 2, 3},
			`[]uint64:3:3 {
   1,
   2,
   3,
}`,
		},
		{
			[]float32{1.2, 3.4, 5.6},
			`[]float32:3:3 {
   1.2000000476837158,
   3.4000000953674316,
   5.599999904632568,
}`,
		},
		{
			[]float64{1.2, 3.4, 5.6},
			`[]float64:3:3 {
   1.2,
   3.4,
   5.6,
}`,
		},
		{
			[]complex64{1, 2.3, -4},
			`[]complex64:3:3 {
   (1+0i),
   (2.299999952316284+0i),
   (-4+0i),
}`,
		},
		{
			[]complex128{1, 2.3, -4},
			`[]complex128:3:3 {
   (1+0i),
   (2.3+0i),
   (-4+0i),
}`,
		},
		{
			[]bool{true, false},
			`[]bool:2:2 {
   true,
   false,
}`,
		},
		{
			[]any{
				func(i int) int { return i },
				func(int) {},
				func() int { return 123 },
			},
			`[]interface {}:3:3 {
   func(int) int,
   func(int),
   func() int,
}`,
		},
		{make(map[any]any), "map[interface {}]interface {}:0 {\n}"},
		{map[string]int{"x": 123, "y": 456}, `map[string]int:2 {
   "x": 123,
   "y": 456,
}`},
	}

	for i, tc := range testCases {
		var d dumper
		d.dump(tc.inputVar)

		if returned := d.buf.String(); returned != tc.expected {
			t.Fatalf(`Case#%d failed, dumper returned unuexpected results : "%s" (%d), expected "%s" (%d)`, i, returned, len(returned), tc.expected,
				len(tc.expected))
		}
	}

}
