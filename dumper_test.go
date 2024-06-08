package godump

import (
	"bytes"
	"math/cmplx"
	"os"
	"testing"
)

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
		{float32(12.3), "12.3"},
		{float32(-12.3), "-12.3"},
		{float64(12.3), "12.3"},
		{float64(-12.3), "-12.3"},
		{complex64(12.3), "(12.3+0i)"},
		{complex64(-12.3), "(-12.3+0i)"},
		{complex128(12.3), "(12.3+0i)"},
		{complex128(-12.3), "(-12.3+0i)"},
		{true, "true"},
		{false, "false"},
		{"hello world", `"hello world"`},
		{func(i int) int { return i }, "func(int) int"},
		{func(int) {}, "func(int)"},
		{func() int { return 123 }, "func() int"},
		{func() {}, "func()"},
		{make([]any, 0, 5), "[]interface {}:0:5 {}"},
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
   1.2,
   3.4,
   5.6,
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
   (2.3+0i),
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
		{make(map[any]any), "map[interface {}]interface {}:0 {}"},
		{map[string]int{"x": 123}, `map[string]int:1 {
   "x": 123,
}`},
	}

	type User struct {
		Name   string
		Friend *User
	}

	person1 := User{"test", nil}
	person2 := User{"test 2", &person1}
	person3 := User{"test 3", &person2}
	person1.Friend = &person3

	testCases = append(testCases, struct {
		inputVar any
		expected string
	}{
		person3,
		`godump.User {
   Name: "test 3",
   Friend: #1&godump.User {
      Name: "test 2",
      Friend: &godump.User {
         Name: "test",
         Friend: &godump.User {
            Name: "test 3",
            Friend: &@1,
         },
      },
   },
}`,
	})

	for i, tc := range testCases {
		var d dumper
		d.dump(tc.inputVar)

		if returned := string(d.buf); returned != tc.expected {
			t.Fatalf(`Case#%d failed, dumper returned unuexpected results : "%s" (%d), expected "%s" (%d)`, i, returned, len(returned), tc.expected,
				len(tc.expected))
		}
	}

}

// Define the complex nested structure
type Node struct {
	ID              int
	Int8Field       int8
	Int16Field      int16
	Int32Field      int32
	Int64Field      int64
	UintField       uint
	Uint8Field      uint8
	Uint16Field     uint16
	Uint32Field     uint32
	Uint64Field     uint64
	UintptrField    uintptr
	Float32Field    float32
	Float64Field    float64
	Complex64Field  complex64
	Complex128Field complex128
	StringField     string
	BoolField       bool
	Children        []*Node
	Attributes      map[string]interface{}
	IntPointer      *int
	StringPointer   *string
	MapPointer      *map[string]int
	ChannelInt      chan int
	ChannelStr      chan string
	ChannelStruct   chan Detail
	Array           [3]int
	SubNode         *SubNode
	privateSubNode  *SubNode
	CyclicReference **Node
	NamedTypes      NamedTypes
}

type SubNode struct {
	Code       int
	Parent     *Node
	Data       []*Node
	SubDetails []Detail
	ChannelSub chan *SubNode
}

type Detail struct {
	Info        string
	Count       int
	Next        *Detail
	DetailMap   map[string]bool
	DetailSlice []complex64
}

type NamedTypes struct {
	NamedInt   int
	NamedFloat float64
	NamedStr   string
	NamedMap   map[string]string
}

func TestDumperWithComplexDataStructure(t *testing.T) {
	var intValue = 42
	var strValue = "example"
	var mapValue = map[string]int{"key1": 1}

	channelInt := make(chan int)
	channelStr := make(chan string)
	channelStruct := make(chan Detail)
	channelSub := make(chan *SubNode)

	var cyclicNode *Node

	root := &Node{
		ID:              1,
		Int8Field:       int8(8),
		Int16Field:      int16(16),
		Int32Field:      int32(32),
		Int64Field:      int64(64),
		UintField:       uint(100),
		Uint8Field:      uint8(200),
		Uint16Field:     uint16(300),
		Uint32Field:     uint32(400),
		Uint64Field:     uint64(500),
		UintptrField:    uintptr(600),
		Float32Field:    float32(123.456),
		Float64Field:    789.012,
		Complex64Field:  complex64(1 + 2i),
		Complex128Field: complex128(cmplx.Exp(1 + 2i)),
		StringField:     "RootNode",
		BoolField:       true,
		Array:           [3]int{1, 2, 3},
		Children: []*Node{
			{
				ID:              2,
				Int8Field:       int8(18),
				Int16Field:      int16(116),
				Int32Field:      int32(132),
				Int64Field:      int64(164),
				UintField:       uint(1100),
				Uint8Field:      uint8(100),
				Uint16Field:     uint16(1300),
				Uint32Field:     uint32(1400),
				Uint64Field:     uint64(1500),
				UintptrField:    uintptr(1600),
				Float32Field:    float32(223.456),
				Float64Field:    1789.012,
				Complex64Field:  complex64(2 + 3i),
				Complex128Field: complex128(cmplx.Exp(2 + 3i)),
				StringField:     "ChildNode1",
				BoolField:       false,
				Array:           [3]int{4, 5, 6},
				IntPointer:      &intValue,
				StringPointer:   &strValue,
				MapPointer:      &mapValue,
				ChannelInt:      channelInt,
				ChannelStr:      channelStr,
				ChannelStruct:   channelStruct,
				Attributes: map[string]interface{}{
					"attr": []float64{1.1, 2.2, 3.3},
				},
				SubNode: &SubNode{
					Code: 100,
					SubDetails: []Detail{
						{
							Info:  "Detail1",
							Count: 1,
							Next: &Detail{
								Info:        "Detail2",
								Count:       2,
								Next:        nil,
								DetailMap:   map[string]bool{"key1": true},
								DetailSlice: []complex64{1 + 1i, 2 + 2i},
							},
							DetailMap:   map[string]bool{"key2": false},
							DetailSlice: []complex64{3 + 3i, 4 + 4i},
						},
					},
					ChannelSub: channelSub,
				},
				privateSubNode: &SubNode{
					Code: 100,
					SubDetails: []Detail{
						{
							Info:  "Detail1",
							Count: 1,
							Next: &Detail{
								Info:        "Detail2",
								Count:       2,
								Next:        nil,
								DetailMap:   map[string]bool{"key1": true},
								DetailSlice: []complex64{1 + 1i, 2 + 2i},
							},
							DetailMap:   map[string]bool{"key2": false},
							DetailSlice: []complex64{3 + 3i, 4 + 4i},
						},
					},
					ChannelSub: channelSub,
				},
			},
			{
				ID:              3,
				Int8Field:       int8(28),
				Int16Field:      int16(216),
				Int32Field:      int32(232),
				Int64Field:      int64(264),
				UintField:       uint(2100),
				Uint8Field:      uint8(220),
				Uint16Field:     uint16(2300),
				Uint32Field:     uint32(2400),
				Uint64Field:     uint64(2500),
				UintptrField:    uintptr(2600),
				Float32Field:    float32(323.456),
				Float64Field:    2789.012,
				Complex64Field:  complex64(3 + 4i),
				Complex128Field: complex128(cmplx.Exp(3 + 4i)),
				StringField:     "ChildNode2",
				BoolField:       true,
				Array:           [3]int{7, 8, 9},
				IntPointer:      &intValue,
				StringPointer:   &strValue,
				MapPointer:      &mapValue,
				ChannelInt:      channelInt,
				ChannelStr:      channelStr,
				ChannelStruct:   channelStruct,
				Attributes: map[string]interface{}{
					"attr": []string{"a", "b", "c"},
				},
				SubNode: &SubNode{
					Code: 200,
					SubDetails: []Detail{
						{
							Info:  "DetailA",
							Count: 10,
							Next: &Detail{
								Info:        "DetailB",
								Count:       20,
								Next:        nil,
								DetailMap:   map[string]bool{"key3": true},
								DetailSlice: []complex64{5 + 5i, 6 + 6i},
							},
							DetailMap:   map[string]bool{"key4": false},
							DetailSlice: []complex64{7 + 7i, 8 + 8i},
						},
					},
					ChannelSub: channelSub,
				},
			},
		},
		Attributes: map[string]interface{}{
			"globalAttr": []int{100, 200, 300},
		},
		IntPointer:    &intValue,
		StringPointer: &strValue,
		MapPointer:    &mapValue,
		ChannelInt:    channelInt,
		ChannelStr:    channelStr,
		ChannelStruct: channelStruct,
		SubNode: &SubNode{
			Code: 999,
			Data: []*Node{
				{
					ID:              4,
					Int8Field:       int8(38),
					Int16Field:      int16(316),
					Int32Field:      int32(332),
					Int64Field:      int64(364),
					UintField:       uint(3100),
					Uint8Field:      uint8(200),
					Uint16Field:     uint16(3300),
					Uint32Field:     uint32(3400),
					Uint64Field:     uint64(3500),
					UintptrField:    uintptr(3600),
					Float32Field:    float32(423.456),
					Float64Field:    3789.012,
					Complex64Field:  complex64(4 + 5i),
					Complex128Field: complex128(cmplx.Exp(4 + 5i)),
					StringField:     "GrandChildNode1",
					BoolField:       false,
					Array:           [3]int{10, 11, 12},
					Children: []*Node{
						{
							ID:              5,
							Int8Field:       int8(48),
							Int16Field:      int16(416),
							Int32Field:      int32(432),
							Int64Field:      int64(464),
							UintField:       uint(4100),
							Uint8Field:      uint8(200),
							Uint16Field:     uint16(4300),
							Uint32Field:     uint32(4400),
							Uint64Field:     uint64(4500),
							UintptrField:    uintptr(4600),
							Float32Field:    float32(523.456),
							Float64Field:    4789.012,
							Complex64Field:  complex64(5 + 6i),
							Complex128Field: complex128(cmplx.Exp(5 + 6i)),
							StringField:     "GreatGrandChildNode1",
							BoolField:       true,
							Array:           [3]int{13, 14, 15},
							Attributes: map[string]interface{}{
								"greatAttr": "greatValue",
							},
							IntPointer:    &intValue,
							StringPointer: &strValue,
							MapPointer:    &mapValue,
							ChannelInt:    channelInt,
							ChannelStr:    channelStr,
							ChannelStruct: channelStruct,
						},
					},
					Attributes: map[string]interface{}{
						"grandAttr": "grandValue1",
					},
					IntPointer:    &intValue,
					StringPointer: &strValue,
					MapPointer:    &mapValue,
					ChannelInt:    channelInt,
					ChannelStr:    channelStr,
					ChannelStruct: channelStruct,
				},
				{
					ID:              6,
					Int8Field:       int8(58),
					Int16Field:      int16(516),
					Int32Field:      int32(532),
					Int64Field:      int64(564),
					UintField:       uint(5100),
					Uint8Field:      uint8(200),
					Uint16Field:     uint16(5300),
					Uint32Field:     uint32(5400),
					Uint64Field:     uint64(5500),
					UintptrField:    uintptr(5600),
					Float32Field:    float32(623.456),
					Float64Field:    5789.012,
					Complex64Field:  complex64(6 + 7i),
					Complex128Field: complex128(cmplx.Exp(6 + 7i)),
					StringField:     "GrandChildNode2",
					BoolField:       true,
					Array:           [3]int{16, 17, 18},
					Children: []*Node{
						{
							ID:              7,
							Int8Field:       int8(68),
							Int16Field:      int16(616),
							Int32Field:      int32(632),
							Int64Field:      int64(664),
							UintField:       uint(6100),
							Uint8Field:      uint8(200),
							Uint16Field:     uint16(6300),
							Uint32Field:     uint32(6400),
							Uint64Field:     uint64(6500),
							UintptrField:    uintptr(6600),
							Float32Field:    float32(723.456),
							Float64Field:    6789.012,
							Complex64Field:  complex64(7 + 8i),
							Complex128Field: complex128(cmplx.Exp(7 + 8i)),
							StringField:     "GreatGrandChildNode2",
							BoolField:       false,
							Array:           [3]int{19, 20, 21},
							Attributes: map[string]interface{}{
								"greatAttr": "greatValue2",
							},
							IntPointer:    &intValue,
							StringPointer: &strValue,
							MapPointer:    &mapValue,
							ChannelInt:    channelInt,
							ChannelStr:    channelStr,
							ChannelStruct: channelStruct,
						},
					},
					Attributes: map[string]interface{}{
						"grandAttr": "grandValue2",
					},
					IntPointer:    &intValue,
					StringPointer: &strValue,
					MapPointer:    &mapValue,
					ChannelInt:    channelInt,
					ChannelStr:    channelStr,
					ChannelStruct: channelStruct,
				},
			},
			ChannelSub: channelSub,
		},
		NamedTypes: NamedTypes{
			NamedInt:   99,
			NamedFloat: 99.99,
			NamedStr:   "NamedValue",
			NamedMap:   map[string]string{"keyA": "valueA"},
		},
		CyclicReference: &cyclicNode,
	}

	cyclicNode = root

	expectedOutput, err := os.ReadFile("./testdata/output.txt")
	if err != nil {
		t.Fatal(err)
	}

	var d dumper
	d.dump(root)
	returned := d.buf

	r_lines := bytes.Split(returned, []byte("\n"))
	e_lines := bytes.Split(expectedOutput, []byte("\n"))

	if len(r_lines) != len(e_lines) {
		t.Fatalf("expected %d lines, got %d", len(e_lines), len(r_lines))
	}

	for i, line := range e_lines {
		if len(line) != len(r_lines[i]) {
			t.Fatalf(`mismatche at line %d:
--- "%s"
+++ "%s"`, i+1, line, r_lines[i])
		}

		for j, ch := range line {
			if ch != r_lines[i][j] {
				t.Fatalf(`expected "%c", got "%c" at line %d:%d"`, ch, r_lines[i][j], i+1, j)
			}
		}
	}
}
