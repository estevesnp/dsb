package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}

	foobar1 := &String{Value: "foobar"}
	foobar2 := &String{Value: "foobar"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("hello strings with the same content have different hash keys")
	}

	if foobar1.HashKey() != foobar2.HashKey() {
		t.Errorf("foobar strings with the same content have different hash keys")
	}

	if hello1.HashKey() == foobar1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestStringHashCache(t *testing.T) {
	s := "foo"
	str := &String{Value: s}

	if _, ok := stringHashCache[s]; ok {
		t.Fatalf("string hash cache for value %q should be empty", s)
	}

	str.HashKey()

	if _, ok := stringHashCache[s]; !ok {
		t.Fatalf("string hash cache for value %q should be filled", s)
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}

	ten1 := &Integer{Value: 10}
	ten2 := &Integer{Value: 10}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("one integers with the same content have different hash keys")
	}

	if ten1.HashKey() != ten2.HashKey() {
		t.Errorf("ten integers with the same content have different hash keys")
	}

	if one1.HashKey() == ten1.HashKey() {
		t.Errorf("integers with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}

	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("true booleans with the same content have different hash keys")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("false booleans with the same content have different hash keys")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("booleans with different content have same hash keys")
	}
}
