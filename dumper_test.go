package godump

import (
	"bytes"
	"os"
	"reflect"
	"testing"
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

		Nil: nil,
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

	var d dumper
	d.dump(reflect.ValueOf(node))

	checkFromFeed(t, d.buf, "./testdata/primitives.txt")
}

func TestCanDumpStructes(t *testing.T) {

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

		Ref *Node
	}

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
	}

	node.Inline.Field2.Field2.Field2 = node.Inline.Field2.Field2

	node.Typed.Field2 = &node.Inline.Field2
	node.Ref = &node

	var d dumper
	d.dump(reflect.ValueOf(node))

	checkFromFeed(t, d.buf, "./testdata/structs.txt")
}

func TestCanDumpPrivateStructes(t *testing.T) {

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

		Ref *node
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
	}

	n.inline.field2.field2.field2 = n.inline.field2.field2

	n.typed.field2 = &n.inline.field2
	n.Ref = &n

	var d dumper
	d.dump(reflect.ValueOf(n))

	checkFromFeed(t, d.buf, "./testdata/private-structs.txt")
}

func TestCanDumpPrivateStructesWhenPrivateFieldsDumpingIsEnabled(t *testing.T) {

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

		Ref *node
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
	}

	n.inline.field2.field2.field2 = n.inline.field2.field2

	n.typed.field2 = &n.inline.field2
	n.Ref = &n

	var d dumper
	d.dumpPrivateFields = true
	d.dump(reflect.ValueOf(n))

	checkFromFeed(t, d.buf, "./testdata/private-structs-dumped.txt")
}

func checkFromFeed(t *testing.T, result []byte, feed_path string) {
	t.Helper()

	expectedOutput, err := os.ReadFile(feed_path)
	if err != nil {
		t.Fatal(err)
	}

	r_lines := bytes.Split(result, []byte("\n"))
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
