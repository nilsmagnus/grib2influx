package griblib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

type Message struct {
	Section0 Section0
	Section1 Section1
	Section2 Section2
	Section3 Section3
	Section4 Section4
	Section5 Section5
	Section6 Section6
	Section7 Section7
}

type Options struct {
	Operation               string    `json:"operation"`
	Discipline              int       `json:"discipline"`
	DataExport              bool      `json:"dataExport"`
	Category                int       `json:"category"`
	Filepath                string    `json:"filePath"`
	ReduceFilePath          string    `json:"reduceFilePath"`
	ExportType              int       `json:"exportType"`
	MaximumNumberOfMessages int       `json:"maximumNumberOfMessages"`
	GeoFilter               GeoFilter `json:"geoFilter"`
}

const (
	Grib                 = 0x47524942
	EndSectionLength     = 926365495
	SupportedGribEdition = 2
)

// ReadMessages reads all message from gribFile
func ReadMessages(gribFile io.Reader, options Options) (messages []Message, err error) {

	for {
		message, messageErr := ReadMessage(gribFile)
		if messageErr != nil {
			if strings.Contains(messageErr.Error(), "EOF") {
				return messages, nil
			} else {
				fmt.Println("Error when parsing a message, ", messageErr.Error())
				return messages, err
			}
		} else {
			messages = append(messages, message)
			if len(messages) >= int(options.MaximumNumberOfMessages) {
				return messages, nil
			}
		}
	}
}

func ReadMessage(gribFile io.Reader) (message Message, err error) {

	section0, headError := ReadSection0(gribFile)

	if headError != nil {
		return message, headError
	}

	fmt.Sprintln("section 0 length is ", unsafe.Sizeof(section0))
	messageBytes := make([]byte, section0.MessageLength-16)

	numBytes, readError := gribFile.Read(messageBytes)

	if readError != nil {
		fmt.Println("Error reading message")
		return message, readError
	}

	if numBytes != int(section0.MessageLength-16) {
		fmt.Println("Did not read full message")
	}

	return readMessage(bytes.NewReader(messageBytes), section0)

}

func readMessage(gribFile io.Reader, section0 Section0) (message Message, err error) {

	message.Section0 = section0
	for {

		// pre-parse section head to decide which struct use
		sectionHead, headErr := ReadSectionHead(gribFile)
		if headErr != nil {
			fmt.Println("Error reading header", headErr.Error())
			return message, headErr
		}

		switch sectionHead.Number {

		case 1:
			message.Section1, err = ReadSection1(gribFile)
		case 2:
			message.Section2, err = ReadSection2(gribFile, sectionHead.ContentLength())
		case 3:
			message.Section3, err = ReadSection3(gribFile)
		case 4:
			message.Section4, err = ReadSection4(gribFile)
		case 5:
			message.Section5, err = ReadSection5(gribFile)
		case 6:
			message.Section6, err = ReadSection6(gribFile, sectionHead.ContentLength())
		case 7:
			message.Section7, err = ReadSection7(gribFile, sectionHead.ContentLength(), message.Section5.DataTemplate)
		case 8:
			// end-section, return
			return message, nil
		default:
			err = fmt.Errorf("Unknown section number %d  (Something bad with parser or files)", sectionHead.Number)
		}
		if err != nil {
			return message, err
		}
	}
}

type Section0 struct {
	Indicator     uint32 `json:"indicator"`
	Reserved      uint16 `json:"reserved"`
	Discipline    uint8  `json:"discipline"`
	Edition       uint8  `json:"edition"`
	MessageLength uint64 `json:"messageLength"`
}

func ReadSection0(reader io.Reader) (section0 Section0, err error) {
	section0.Indicator=255
	err = binary.Read(reader, binary.BigEndian, &section0)
	if err != nil {
		return section0, err
	}

	if section0.Indicator == Grib {
		if section0.Edition != SupportedGribEdition {
			return section0, fmt.Errorf("Unsupported  grib edition %d", section0.Edition)
		}
	}

	return

}

type Section interface {
	SectionNumber() uint8
	String() string
}

type SectionHead struct {
	ByteLength uint32 `json:"byteLength"`
	Number     uint8  `json:"number"`
}

func ReadSectionHead(section io.Reader) (head SectionHead, err error) {
	var length uint32
	err = binary.Read(section, binary.BigEndian, &length)
	if err != nil {
		return head, fmt.Errorf("Read of Length failed: %s", err.Error())
	}
	if length == EndSectionLength {
		return SectionHead{
			ByteLength: 4,
			Number:     8,
		}, nil
	}
	var sectionNumber uint8
	err = binary.Read(section, binary.BigEndian, &sectionNumber)
	if err != nil {
		return head, err
	}

	return SectionHead{
		ByteLength: length,
		Number:     sectionNumber,
	}, nil
}

func (s SectionHead) SectionNumber() uint8 {
	return s.Number
}

func (s SectionHead) ContentLength() int {
	return int(s.ByteLength) - binary.Size(s)
}

func (s SectionHead) String() string {
	return fmt.Sprint("Section ", s.Number)
}

type Time struct {
	Year   uint16 `json:"year"`   // year
	Month  uint8  `json:"month"`  // month + 1
	Day    uint8  `json:"day"`    // day
	Hour   uint8  `json:"hour"`   // hour
	Minute uint8  `json:"minute"` // minute
	Second uint8  `json:"second"` // second
}

type Section1 struct {
	OriginatingCenter         uint16 `json:"ooriginatingCenter"`
	OriginatingSubCenter      uint16 `json:"originatingSubCenter"`
	MasterTablesVersion       uint8  `json:"masterTablesVersion"`
	LocalTablesVersion        uint8  `json:"localTablesVersion"`
	ReferenceTimeSignificance uint8  `json:"referenceTimeSignificance"`
	ReferenceTime             Time   `json:"referenceTime"`
	ProductionStatus          uint8  `json:"productionStatus"`
	Type                      uint8  `json:"type"`
}

func ReadSection1(f io.Reader) (section Section1, err error) {
	return section, binary.Read(f, binary.BigEndian, &section)
}

type Section2 struct {
	LocalUse []uint8 `json:"localUse"`
}

func ReadSection2(f io.Reader, len int) (section Section2, err error) {
	section.LocalUse = make([]uint8, len)
	return section, read(f, &section.LocalUse)
}

type Section3 struct {
	Source                   uint8       `json:"source"`
	DataPointCount           uint32      `json:"datapointCount"`
	PointCountOctets         uint8       `json:"pointCountOctets"`
	PointCountInterpretation uint8       `json:"pointCountInterpretation"`
	TemplateNumber           uint16      `json:"templateNumber"`
	Definition               interface{} `json:"definition"`
}

func (s Section3) String() string {
	return fmt.Sprint("Point count: ", s.DataPointCount, " Definition: ", ReadGridDefinitionTemplateNumber(int(s.TemplateNumber)))
}

func ReadSection3(f io.Reader) (section Section3, err error) {

	err = read(f, &section.Source, &section.DataPointCount, &section.PointCountOctets, &section.PointCountInterpretation, &section.TemplateNumber)
	if err != nil {
		return section, err
	}

	section.Definition, err = ReadGrid(f, section.TemplateNumber)
	return section, err
}

type Section4 struct {
	CoordinatesCount                uint16   `json:"coordinatesCount"`
	ProductDefinitionTemplateNumber uint16   `json:"productDefinitionTemplateNumber"`
	ProductDefinitionTemplate       Product0 `json:"productDefinitionTemplate"` // FIXME
	Coordinates                     []byte   `json:"coordinates"`
}

func ReadSection4(f io.Reader) (section Section4, err error) {
	err = read(f, &section.CoordinatesCount, &section.ProductDefinitionTemplateNumber)
	if err != nil {
		return section, err
	}

	switch section.ProductDefinitionTemplateNumber {
	case 0:
		err = read(f, &section.ProductDefinitionTemplate)
	default:
		return section, fmt.Errorf("Category definition template number %d not implemented yet", section.ProductDefinitionTemplateNumber)
	}

	if err != nil {
		return section, err
	}

	section.Coordinates = make([]byte, section.CoordinatesCount)

	return section, read(f, &section.Coordinates)
}

type Section5 struct {
	PointsNumber       uint32 `json:"pointsNumber"`
	DataTemplateNumber uint16 `json:"dataTemplateNumber"`
	DataTemplate       Data3  `json:"dataTemplate"` // FIXME
}

func ReadSection5(f io.Reader) (section Section5, err error) {
	err = read(f, &section.PointsNumber, &section.DataTemplateNumber, &section.DataTemplate)
	if err != nil {
		return section, err
	}

	if section.DataTemplateNumber != 3 {
		return section, fmt.Errorf("Template number not supported: %d", section.DataTemplateNumber)
	}
	//f.Seek(int64(length - 11), 1);

	//fmt.Println(section.DataTemplate)

	return section, nil
}

type Section6 struct {
	BitmapIndicator uint8  `json:"bitmapIndicator"`
	Bitmap          []byte `json:"bitmap"`
}

func ReadSection6(f io.Reader, length int) (section Section6, err error) {
	section.Bitmap = make([]byte, length-1)

	return section, read(f, &section.BitmapIndicator, &section.Bitmap)
}

type Section7 struct {
	Data []int64 `json:"data"`
}

func ReadSection7(f io.Reader, length int, template Data3) (section Section7, err error) {
	section.Data = ParseData3(f, length, &template) // 5 is the length of (octet 1-5)
	return section, err
}

func read(reader io.Reader, data ...interface{}) (err error) {
	for _, what := range data {
		err = binary.Read(reader, binary.BigEndian, what)
		if err != nil {
			return err
		}
	}
	return nil
}
