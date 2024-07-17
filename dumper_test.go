package godump_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/yassinebenaid/godump"
)

func TestCanDumpPrimitives(t *testing.T) {
	type IntType int
	type Int8Type int8
	type Int16Type int16
	type Int32Type int32
	type Int64Type int64
	type UintType uint
	type Uint8Type uint8
	type Uint16Type uint16
	type Uint32Type uint32
	type Uint64Type uint64
	type Float32Type float32
	type Float64Type float64
	type Complex64Type complex64
	type Complex128Type complex128
	type Bool1Type bool
	type Bool2Type bool
	type StringType string
	type UintptrType uintptr

	type IntPtrType *int
	type Int8PtrType *int8
	type Int16PtrType *int16
	type Int32PtrType *int32
	type Int64PtrType *int64
	type UintPtrType *uint
	type Uint8PtrType *uint8
	type Uint16PtrType *uint16
	type Uint32PtrType *uint32
	type Uint64PtrType *uint64
	type Float32PtrType *float32
	type Float64PtrType *float64
	type Complex64PtrType *complex64
	type Complex128PtrType *complex128
	type Bool1PtrType *bool
	type Bool2PtrType *bool
	type StringPtrType *string
	type UintptrPtrType *uintptr

	type FuncType func()
	type Func2Type func(int) float64
	type Func3Type func(...*any) any
	type Func4Type func(byte, ...[]*complex128) bool

	type ChanType chan struct{}
	type Chan1Type <-chan struct{}
	type Chan2Type chan<- struct{}

	type Node struct {
		Int        int
		Int8       int8
		Int16      int16
		Int32      int32
		Int64      int64
		Uint       uint
		Uint8      uint8
		Uint16     uint16
		Uint32     uint32
		Uint64     uint64
		Float32    float32
		Float64    float64
		Complex64  complex64
		Complex128 complex128
		Bool1      bool
		Bool2      bool
		String     string

		Uintptr uintptr

		IntPtr        *int
		Int8Ptr       *int8
		Int16Ptr      *int16
		Int32Ptr      *int32
		Int64Ptr      *int64
		UintPtr       *uint
		Uint8Ptr      *uint8
		Uint16Ptr     *uint16
		Uint32Ptr     *uint32
		Uint64Ptr     *uint64
		Float32Ptr    *float32
		Float64Ptr    *float64
		Complex64Ptr  *complex64
		Complex128Ptr *complex128
		Bool1Ptr      *bool
		Bool2Ptr      *bool
		StringPtr     *string

		UintptrPtr *uintptr

		TypedInt        IntType
		TypedInt8       Int8Type
		TypedInt16      Int16Type
		TypedInt32      Int32Type
		TypedInt64      Int64Type
		TypedUint       UintType
		TypedUint8      Uint8Type
		TypedUint16     Uint16Type
		TypedUint32     Uint32Type
		TypedUint64     Uint64Type
		TypedFloat32    Float32Type
		TypedFloat64    Float64Type
		TypedComplex64  Complex64Type
		TypedComplex128 Complex128Type
		TypedBool1      Bool1Type
		TypedBool2      Bool2Type
		TypedString     StringType

		TypedUintptr UintptrType

		TypedIntPtr        IntPtrType
		TypedInt8Ptr       Int8PtrType
		TypedInt16Ptr      Int16PtrType
		TypedInt32Ptr      Int32PtrType
		TypedInt64Ptr      Int64PtrType
		TypedUintPtr       UintPtrType
		TypedUint8Ptr      Uint8PtrType
		TypedUint16Ptr     Uint16PtrType
		TypedUint32Ptr     Uint32PtrType
		TypedUint64Ptr     Uint64PtrType
		TypedFloat32Ptr    Float32PtrType
		TypedFloat64Ptr    Float64PtrType
		TypedComplex64Ptr  Complex64PtrType
		TypedComplex128Ptr Complex128PtrType
		TypedBool1Ptr      Bool1PtrType
		TypedBool2Ptr      Bool2PtrType
		TypedStringPtr     StringPtrType

		TypedUintptrPtr UintptrPtrType

		PtrTypedInt        *IntType
		PtrTypedInt8       *Int8Type
		PtrTypedInt16      *Int16Type
		PtrTypedInt32      *Int32Type
		PtrTypedInt64      *Int64Type
		PtrTypedUint       *UintType
		PtrTypedUint8      *Uint8Type
		PtrTypedUint16     *Uint16Type
		PtrTypedUint32     *Uint32Type
		PtrTypedUint64     *Uint64Type
		PtrTypedFloat32    *Float32Type
		PtrTypedFloat64    *Float64Type
		PtrTypedComplex64  *Complex64Type
		PtrTypedComplex128 *Complex128Type
		PtrTypedBool1      *Bool1Type
		PtrTypedBool2      *Bool2Type
		PtrTypedString     *StringType

		PtrTypedUintptr *UintptrType

		Nil *any

		Func  func()
		Func2 func(int) float64
		Func3 func(...*any) any
		Func4 func(byte, ...[]*complex128) bool

		FuncPtr  *func()
		Func2Ptr *func(int) float64
		Func3Ptr *func(...*any) any
		Func4Ptr *func(byte, ...[]*complex128) bool

		TypedFunc  FuncType
		TypedFunc2 Func2Type
		TypedFunc3 Func3Type
		TypedFunc4 Func4Type

		PtrTypedFunc  *FuncType
		PtrTypedFunc2 *Func2Type
		PtrTypedFunc3 *Func3Type
		PtrTypedFunc4 *Func4Type

		Chan  chan struct{}
		Chan1 <-chan struct{}
		Chan2 chan<- struct{}

		ChanPtr  *chan struct{}
		Chan1Ptr *<-chan struct{}
		Chan2Ptr *chan<- struct{}

		TypedChan  ChanType
		TypedChan1 Chan1Type
		TypedChan2 Chan2Type

		PtrTypedChan  *ChanType
		PtrTypedChan1 *Chan1Type
		PtrTypedChan2 *Chan2Type

		UnsafePointer1 unsafe.Pointer
		UnsafePointer2 *unsafe.Pointer
	}

	node := Node{
		Int:        123,
		Int8:       -45,
		Int16:      6789,
		Int32:      -987,
		Int64:      3849876543247876432,
		Uint:       837,
		Uint8:      38,
		Uint16:     3847,
		Uint32:     9843,
		Uint64:     2834,
		Float32:    123.475,
		Float64:    -12345.09876,
		Complex64:  12.987i,
		Complex128: -473i,
		Bool1:      true,
		Bool2:      false,
		String:     "foo bar",

		Uintptr: 1234567890,

		TypedInt:        IntType(123),
		TypedInt8:       Int8Type(-45),
		TypedInt16:      Int16Type(6789),
		TypedInt32:      Int32Type(-987),
		TypedInt64:      Int64Type(3849876543247876432),
		TypedUint:       UintType(837),
		TypedUint8:      Uint8Type(38),
		TypedUint16:     Uint16Type(3847),
		TypedUint32:     Uint32Type(9843),
		TypedUint64:     Uint64Type(2834),
		TypedFloat32:    Float32Type(123.475),
		TypedFloat64:    Float64Type(-12345.09876),
		TypedComplex64:  Complex64Type(12.987i),
		TypedComplex128: Complex128Type(-473i),
		TypedBool1:      Bool1Type(true),
		TypedBool2:      Bool2Type(false),
		TypedString:     StringType("foo bar"),

		TypedUintptr: UintptrType(1234567890),

		Nil: nil,

		UnsafePointer1: nil,
	}

	node.IntPtr = &node.Int
	node.Int8Ptr = &node.Int8
	node.Int16Ptr = &node.Int16
	node.Int32Ptr = &node.Int32
	node.Int64Ptr = &node.Int64
	node.UintPtr = &node.Uint
	node.Uint8Ptr = &node.Uint8
	node.Uint16Ptr = &node.Uint16
	node.Uint32Ptr = &node.Uint32
	node.Uint64Ptr = &node.Uint64
	node.Float32Ptr = &node.Float32
	node.Float64Ptr = &node.Float64
	node.Complex64Ptr = &node.Complex64
	node.Complex128Ptr = &node.Complex128
	node.Bool1Ptr = &node.Bool1
	node.Bool2Ptr = &node.Bool2
	node.StringPtr = &node.String

	node.UintptrPtr = &node.Uintptr

	node.TypedIntPtr = node.IntPtr
	node.TypedInt8Ptr = node.Int8Ptr
	node.TypedInt16Ptr = node.Int16Ptr
	node.TypedInt32Ptr = node.Int32Ptr
	node.TypedInt64Ptr = node.Int64Ptr
	node.TypedUintPtr = node.UintPtr
	node.TypedUint8Ptr = node.Uint8Ptr
	node.TypedUint16Ptr = node.Uint16Ptr
	node.TypedUint32Ptr = node.Uint32Ptr
	node.TypedUint64Ptr = node.Uint64Ptr
	node.TypedFloat32Ptr = node.Float32Ptr
	node.TypedFloat64Ptr = node.Float64Ptr
	node.TypedComplex64Ptr = node.Complex64Ptr
	node.TypedComplex128Ptr = node.Complex128Ptr
	node.TypedBool1Ptr = node.Bool1Ptr
	node.TypedBool2Ptr = node.Bool2Ptr
	node.TypedStringPtr = node.StringPtr

	node.TypedUintptrPtr = node.UintptrPtr

	node.PtrTypedInt = &node.TypedInt
	node.PtrTypedInt8 = &node.TypedInt8
	node.PtrTypedInt16 = &node.TypedInt16
	node.PtrTypedInt32 = &node.TypedInt32
	node.PtrTypedInt64 = &node.TypedInt64
	node.PtrTypedUint = &node.TypedUint
	node.PtrTypedUint8 = &node.TypedUint8
	node.PtrTypedUint16 = &node.TypedUint16
	node.PtrTypedUint32 = &node.TypedUint32
	node.PtrTypedUint64 = &node.TypedUint64
	node.PtrTypedFloat32 = &node.TypedFloat32
	node.PtrTypedFloat64 = &node.TypedFloat64
	node.PtrTypedComplex64 = &node.TypedComplex64
	node.PtrTypedComplex128 = &node.TypedComplex128
	node.PtrTypedBool1 = &node.TypedBool1
	node.PtrTypedBool2 = &node.TypedBool2
	node.PtrTypedString = &node.TypedString

	node.PtrTypedUintptr = &node.TypedUintptr

	node.FuncPtr = &node.Func
	node.Func2Ptr = &node.Func2
	node.Func3Ptr = &node.Func3
	node.Func4Ptr = &node.Func4
	node.PtrTypedFunc = &node.TypedFunc
	node.PtrTypedFunc2 = &node.TypedFunc2
	node.PtrTypedFunc3 = &node.TypedFunc3
	node.PtrTypedFunc4 = &node.TypedFunc4

	ch := make(chan struct{})
	var ch2 <-chan struct{} = ch
	var ch3 chan<- struct{} = ch

	tch := ChanType(ch)
	tch1 := Chan1Type(ch2)
	tch2 := Chan2Type(ch3)

	node.ChanPtr = &ch
	node.Chan1Ptr = &ch2
	node.Chan2Ptr = &ch3
	node.TypedChan = ch
	node.TypedChan1 = ch2
	node.TypedChan2 = ch3
	node.PtrTypedChan = &tch
	node.PtrTypedChan1 = &tch1
	node.PtrTypedChan2 = &tch2

	node.UnsafePointer2 = (*unsafe.Pointer)(unsafe.Pointer(&node))

	var d godump.Dumper
	result := d.Sprint(node)

	checkFromFeed(t, []byte(result), "./testdata/primitives.txt")
}

func TestCanDumpStructs(t *testing.T) {
	type Number int

	type Child1 struct {
		X int
		Y float64
		Z Number
	}

	type Child struct {
		Field1 Child1

		Field2 *Child
	}

	type Node struct {
		Inline struct {
			Field1 struct {
				X int
				Y float64
				Z Number
			}

			Field2 Child
		}

		Typed Child

		Ptr   **int
		Empty struct{}

		Ref *Node
	}

	num := 123
	numaddr := &num
	node := Node{
		Inline: struct {
			Field1 struct {
				X int
				Y float64
				Z Number
			}
			Field2 Child
		}{
			Field1: struct {
				X int
				Y float64
				Z Number
			}{
				X: 123,
				Y: 123.456,
				Z: Number(987),
			},

			Field2: Child{
				Field1: Child1{
					X: 12344,
					Y: 578,
					Z: Number(9876543),
				},
				Field2: &Child{
					Field1: Child1{
						X: 12344,
						Y: 578,
						Z: Number(9876543),
					},
				},
			},
		},
		Ptr: &numaddr,
	}

	node.Inline.Field2.Field2.Field2 = node.Inline.Field2.Field2

	node.Typed.Field2 = &node.Inline.Field2
	node.Ref = &node

	var d godump.Dumper
	result := d.Sprint(node)

	checkFromFeed(t, []byte(result), "./testdata/structs.txt")
}

func TestCannotDumpPrivateStructsWhenHidingOptionIsEnabled(t *testing.T) {
	type number int

	type child1 struct {
		x int
		y float64
		z number
	}

	type child struct {
		field1 child1

		field2 *child
	}

	type node struct {
		inline struct {
			field1 struct {
				x int
				y float64
				z number
			}

			field2 child
		}

		typed child

		empty struct{}

		ref *node
	}

	n := node{
		inline: struct {
			field1 struct {
				x int
				y float64
				z number
			}
			field2 child
		}{
			field1: struct {
				x int
				y float64
				z number
			}{
				x: 123,
				y: 123.456,
				z: number(987),
			},

			field2: child{
				field1: child1{
					x: 12344,
					y: 578,
					z: number(9876543),
				},
				field2: &child{
					field1: child1{
						x: 12344,
						y: 578,
						z: number(9876543),
					},
				},
			},
		},
		empty: struct{}{},
	}

	n.inline.field2.field2.field2 = n.inline.field2.field2

	n.typed.field2 = &n.inline.field2
	n.ref = &n

	var d godump.Dumper
	d.HidePrivateFields = true

	result := d.Sprint(n)

	if result != "godump_test.node {}" {
		t.Fatalf("unexpected result when trying to dump a private struct with hide private fields option enabled, expected `godump.node {}`, got `%v`", result)
	}
}

func TestCanDumpPrivateStructs(t *testing.T) {
	type number int

	type child1 struct {
		x int
		y float64
		z number
	}

	type child struct {
		field1 child1

		field2 *child
	}

	type node struct {
		inline struct {
			field1 struct {
				x int
				y float64
				z number
			}

			field2 child
		}

		typed child

		empty struct{}

		ref *node
	}

	n := node{
		inline: struct {
			field1 struct {
				x int
				y float64
				z number
			}
			field2 child
		}{
			field1: struct {
				x int
				y float64
				z number
			}{
				x: 123,
				y: 123.456,
				z: number(987),
			},

			field2: child{
				field1: child1{
					x: 12344,
					y: 578,
					z: number(9876543),
				},
				field2: &child{
					field1: child1{
						x: 12344,
						y: 578,
						z: number(9876543),
					},
				},
			},
		},
		empty: struct{}{},
	}

	n.inline.field2.field2.field2 = n.inline.field2.field2

	n.typed.field2 = &n.inline.field2
	n.ref = &n

	var d godump.Dumper
	result := d.Sprint(n)

	checkFromFeed(t, []byte(result), "./testdata/private-structs.txt")
}

func TestCanDumpSlices(t *testing.T) {
	type Slice []any

	foo := "foo"
	bar := "bar"
	baz := "baz"

	s := Slice{
		1,
		2.3,
		true,
		false,
		nil,
		[]*string{
			&foo,
			&bar,
			&baz,
		},
		[]any{},
		&[]bool{
			true,
			false,
		},
		make([]any, 3, 8),
	}
	s = append(s, &s)

	var d godump.Dumper
	result := d.Sprint(s)

	checkFromFeed(t, []byte(result), "./testdata/slices.txt")
}

func TestCanDumpMaps(t *testing.T) {
	type SomeMap map[*SomeMap]*SomeMap
	sm := &SomeMap{}

	m := map[any]any{12: 34}
	maps := []any{
		make(map[string]string),
		map[any]int{
			&m: 123,
		},
		map[string]any{
			"cyclic": &m,
		},
		SomeMap{
			&SomeMap{}: &SomeMap{sm: sm},
		},
	}

	var d godump.Dumper
	result := d.Sprint(maps)

	checkFromFeed(t, []byte(result), "./testdata/maps.txt")
}

func TestCanCustomizeIndentation(t *testing.T) {
	type User struct {
		Name       string
		Age        int
		hobbies    []string
		bestFriend *User
	}

	me := User{
		Name: "yassinebenaid",
		Age:  22,
		hobbies: []string{
			"Dev",
			"Go",
			"Web",
			"DevOps",
		},
	}
	me.bestFriend = &me

	d := godump.Dumper{
		Indentation: "            ",
	}
	result := d.Sprint(me)

	checkFromFeed(t, []byte(result), "./testdata/indentation.txt")
}

type CSSColor struct {
	R, G, B int
}

func (c CSSColor) Apply(s string) string {
	return fmt.Sprintf(`<div style="color: rgb(%d, %d, %d); display: inline-block">%s</div>`, c.R, c.G, c.B, s)
}

func TestCanCustomizeTheme(t *testing.T) {
	type User struct {
		Name       string
		Age        int
		hobbies    []string
		bestFriend *User
	}

	me := User{
		Name: "yassinebenaid",
		Age:  22,
		hobbies: []string{
			"Dev",
			"Go",
			"Web",
			"DevOps",
		},
	}
	me.bestFriend = &me

	var d godump.Dumper

	d.Theme = godump.Theme{
		String:        CSSColor{138, 201, 38},
		Quotes:        CSSColor{112, 214, 255},
		Bool:          CSSColor{249, 87, 56},
		Number:        CSSColor{10, 178, 242},
		Types:         CSSColor{0, 150, 199},
		Address:       CSSColor{205, 93, 0},
		PointerTag:    CSSColor{110, 110, 110},
		Nil:           CSSColor{219, 57, 26},
		Func:          CSSColor{160, 90, 220},
		Fields:        CSSColor{189, 176, 194},
		Chan:          CSSColor{195, 154, 76},
		UnsafePointer: CSSColor{89, 193, 180},
		Braces:        CSSColor{185, 86, 86},
	}

	result := d.Sprint(me)

	checkFromFeed(t, []byte(result), "./testdata/theme.txt")
}

func TestDumperPrint_Sprint_And_Fprint(t *testing.T) {
	type User struct {
		Name    string
		Age     int
		hobbies []string
	}

	me := User{
		Name: "yassinebenaid",
		Age:  22,
		hobbies: []string{
			"Dev",
			"Go",
			"Web",
			"DevOps",
		},
	}

	expected := `godump_test.User {
   Name: "yassinebenaid",
   Age: 22,
   hobbies: []string:4:4 {
      "Dev",
      "Go",
      "Web",
      "DevOps",
   },
}`
	var d godump.Dumper

	if r := d.Sprint(me); expected != r {
		t.Fatalf("unexpected result by Dumper.Sprint : `%s`", r)
	}

	if r := d.Sprintln(me); expected+"\n" != r {
		t.Fatalf("unexpected result by Dumper.Sprintln : `%s`", r)
	}

	var buf bytes.Buffer

	if err := d.Fprint(&buf, me); err != nil {
		t.Fatalf("unexpected error by Dumper.Fprint : `%s`", err)
	} else if expected != buf.String() {
		t.Fatalf("unexpected result by Dumper.Fprint : `%s`", buf.String())
	}

	buf.Reset()
	if err := d.Fprintln(&buf, me); err != nil {
		t.Fatalf("unexpected error by Dumper.Fprintln : `%s`", err)
	} else if expected+"\n" != buf.String() {
		t.Fatalf("unexpected result by Dumper.Fprintln : `%s`", buf.String())
	}
}

type X int

func (X) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("foobar")
}

func TestDumperFprintReturnsAWriteErrorIfEncountered(t *testing.T) {
	var d godump.Dumper

	var x X

	if err := d.Fprint(x, nil); err == nil {
		t.Fatalf("unexpected nil error returned by Dumper.Fprint")
	} else if err.Error() != "dumper error: encountered unexpected write error, foobar" {
		t.Fatalf("unexpected error by Dumper.Fprint : `%s`", err.Error())
	}

	if err := d.Fprintln(x, nil); err == nil {
		t.Fatalf("unexpected nil error returned by Dumper.Fprintln")
	} else if err.Error() != "dumper error: encountered unexpected write error, foobar" {
		t.Fatalf("unexpected error by Dumper.Fprintln : `%s`", err.Error())
	}
}

func checkFromFeed(t *testing.T, result []byte, feedPath string) {
	t.Helper()

	expectedOutput, err := os.ReadFile(feedPath)
	if err != nil {
		t.Fatal(err)
	}

	resultLines := bytes.Split(result, []byte("\n"))
	expectedLines := bytes.Split(expectedOutput, []byte("\n"))

	if len(resultLines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(resultLines))
	}

	for i, line := range expectedLines {
		if string(line) != string(resultLines[i]) {
			t.Fatalf(`mismatch at line %d:
--- "%s" (%d)
+++ "%s" (%d)`, i+1, line, len(line), resultLines[i], len(resultLines[i]))
		}
	}
}
