package bone

import (
	"testing"
)

func TestDecodeBool(t *testing.T) {
	d := Decoder{}
	if d.Accept(0x20) != nil {
		t.Fatal("unable to accept: false")
	}
	if b, err := d.Values[0].Bool(); err != nil || b != false {
		t.Fatal("failed to decode: false")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if b, err := d.Values[1].Bool(); err != nil || b != true {
		t.Fatal("failed to decode: true")
	}
}

func TestDecodeTuple(t *testing.T) {
	d := Decoder{}
	if d.Accept(0xA0) != nil {
		t.Fatal("unable to accept: tuple 1")
	}
	if d.Accept(0x20) != nil {
		t.Fatal("unable to accept: false")
	}
	if d.Accept(0xFF) != nil {
		t.Fatal("unable to accept: level")
	}
	if d.Accept(0xA0) != nil {
		t.Fatal("unable to accept: tuple 1")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if d.Accept(0xA0) != nil {
		t.Fatal("unable to accept: tuple 1")
	}
	if d.Accept(0xA0) != nil {
		t.Fatal("unable to accept: tuple 1")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if d.Accept(0xB0) != nil {
		t.Fatal("unable to accept: tuple 2")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if b, err := d.Values[0].Values[0].Bool(); err != nil || b != false {
		t.Fatal("failed to decode: false")
	}
	if b, err := d.Values[1].Values[0].Bool(); err != nil || b != true {
		t.Fatal("failed to decode: true")
	}
	if l := d.Values[1].Level; l != 1 {
		t.Fatal("wrong level")
	}
	if b, err := d.Values[2].Values[0].Values[0].Bool(); err != nil || b != true {
		t.Fatal("failed to decode: true")
	}
	if b, err := d.Values[3].Values[0].Bool(); err != nil || b != true {
		t.Fatal("failed to decode: true")
	}
	if b, err := d.Values[3].Values[1].Bool(); err != nil || b != true {
		t.Fatal("failed to decode: true")
	}
	if b, err := d.Values[4].Bool(); err != nil || b != true {
		t.Fatal("failed to decode: true")
	}
}

func TestDecodeInts(t *testing.T) {
	d := Decoder{}
	if d.Accept(0x10) != nil {
		t.Fatal("unable to accept: 0")
	}
	if c := d.Values[0].Code; c != 0x10 {
		t.Fatal("failed to decode: 0")
	}
}

func TestDecodeLists(t *testing.T) {
	d := Decoder{}
	if d.Accept(0xF0) != nil {
		t.Fatal("unable to accept: list")
	}
	if d.Accept(0x00) != nil {
		t.Fatal("unable to accept: end list")
	}
	if d.Accept(0xF0) != nil {
		t.Fatal("unable to accept: list")
	}
	if d.Accept(0xF0) != nil {
		t.Fatal("unable to accept: list")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if d.Accept(0x00) != nil {
		t.Fatal("unable to accept: end list")
	}
	if d.Accept(0x21) != nil {
		t.Fatal("unable to accept: true")
	}
	if d.Accept(0x00) != nil {
		t.Fatal("unable to accept: end list")
	}
	if d.Values[0].Code != 0xF0 {
		t.Fatal("failed to decode: list")
	}
	if d.Values[1].Code != 0xF0 {
		t.Fatal("failed to decode: list")
	}
	if len(d.Values[1].Values) != 2 {
		t.Fatal("wrong length")
	}
	if len(d.Values[1].Values[0].Values) != 1 {
		t.Fatal("wrong length")
	}
	if d.Values[1].Values[0].Values[0].Code != 0x21 {
		t.Fatal("failed to decode: true")
	}
	if d.Values[1].Values[1].Code != 0x21 {
		t.Fatal("failed to decode: true")
	}
}
