package griblib

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func Test_read_section0(t *testing.T) {

	test_data := Section0{Discipline: 2,
		Edition:       3,
		MessageLength: 4,
		Indicator:     5,
		Reserved:      6,
	}

	section0, readError := ReadSection0(toIoReader(test_data))

	if readError != nil {
		t.Fatal(readError)
	}

	if test_data != section0 {
		t.Error("Deserialized struct is not equal to original struct")
	}
}
func Test_read_section1(t *testing.T) {

	test_data := Section1{
		LocalTablesVersion:   1,
		MasterTablesVersion:  2,
		OriginatingCenter:    3,
		OriginatingSubCenter: 4,
		ProductionStatus:     5,
		ReferenceTime: Time{
			Day:    1,
			Hour:   2,
			Minute: 3,
			Month:  4,
			Second: 5,
			Year:   6,
		},
		ReferenceTimeSignificance: 7,
		Type: 8,
	}

	section1, readError := ReadSection1(toIoReader(test_data))

	if readError != nil {
		t.Fatal(readError)
	}

	if test_data != section1 {
		t.Error("Deserialized section1 struct is not equal to original struct")
	}
}

// create a reader from a struct for testing purposes
func toIoReader(data interface{}) (reader io.Reader) {
	var bin_buf bytes.Buffer

	binary.Write(&bin_buf, binary.BigEndian, data)

	reader = bytes.NewReader(bin_buf.Bytes())

	return

}
