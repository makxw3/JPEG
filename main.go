package main

import (
	"bytes"
	"fmt"
	"jp/rd"
	"os"
)

func main() {
	const FILEPATH = "./cat.jpg"
	file, err := os.Open(FILEPATH)
	if err != nil {
		fmt.Printf("Error in  Opening the file %s\n", FILEPATH)
		os.Exit(1)
	}
	rdr := rd.Get(file)
	// Make Sure that the first byte is FF and the second byte is D8
	fb, _ := rdr.Read()
	sb, _ := rdr.Read()
	if fb != 0xFF && sb != SOI {
		fmt.Printf("Error! The file %s is not a valid JPEG file\n", FILEPATH)
		os.Exit(1)
	}
	fmt.Printf("The file %s is a valid JPEG file\n", FILEPATH)
	// TODO: Move on to reading the relevant markers
	for {
		fb, _ := rdr.Read()
		sb, _ := rdr.Read()
		if fb != 0xFF {
			fmt.Printf("Error! Expected a Marker but got the byte 0x%X\n", fb)
			os.Exit(1)
		}
		if sb >= APP0 && sb <= APP15 {
			readAPPN(rdr)
		}
		if sb == DQT {
			readQuantizationTable(rdr)
		}
	}
	file.Close()
}

func readAPPN(rdr *rd.Reader) {
	fmt.Printf("**** START -> Reading APPN Marker ****\n")
	// Read the next 16 bits that represent the size of the payload in the marker
	fb, _ := rdr.Read()
	sb, _ := rdr.Read()
	// TODO: Learn about Big and Small Endian
	length := int((fb << 8) + sb)
	fmt.Printf("Payload Length --> #%d\n", length-2)
	rdr.ReadBuffer(length - 2)
	fmt.Printf("**** END -> Reading APPN Marker ****\n\n")
}

func printTable(tb []byte) {
	var out bytes.Buffer
	out.WriteString("[")
	for a := range tb {
		if a != 0 && a%8 == 0 {
			out.WriteString("\n")
		}
		out.WriteString(fmt.Sprintf("%v ", tb[a]))
	}
	out.WriteString("]")
	fmt.Printf("%s\n", out.String())
}

func readQuantizationTable(rdr *rd.Reader) [][]byte {
	fmt.Printf("**** START -> Reading Quantization Table ****\n")
	fb, _ := rdr.Read()
	sb, _ := rdr.Read()
	length := int((fb << 8) + sb)
	length -= 2
	tables := [][]byte{}
	for {
		if length <= 0 {
			break
		}
		tableInfo, _ := rdr.Read()
		length -= 1
		fmt.Printf("Table Info --> %X\n", tableInfo)
		tableID := int(tableInfo & 0x0F)
		// Table ID should be >= 0 && <= 3
		if tableID > 3 {
			fmt.Printf("Error! Invalid TableID #%d\n", tableID)
			os.Exit(1)
		}
		fmt.Printf("Table ID --> %X\n", tableID)
		// Get the upper-nibble -> 8/16 bit table
		upperN := int(tableInfo >> 4)
		// If the (upperN != 0) This implies that the table contains 16 bit values
		if upperN != 0 {
			tb := []byte{}
			tb = make([]byte, 64)
			for a := 0; a < 64; a++ {
				fb, _ := rdr.Read()
				sb, _ := rdr.Read()
				tb[a] = (fb << 8) + sb
			}
			length -= 128
			tables = append(tables, tb)
			printTable(tb)
		} else {
			tb := []byte{}
			tb = make([]byte, 64)
			for a := 0; a < 64; a++ {
				fb, _ := rdr.Read()
				tb[a] = (fb)
			}
			length -= 64
			tables = append(tables, tb)
			printTable(tb)
		}
	}
	if length != 0 {
		fmt.Printf("Expected length = 0 but got #%d instead\n", length)
		os.Exit(1)
	}
	return tables
}

// The Marker Bytes
const (
	SOI  = 0xD8 // Start of Image
	SOF0 = 0xC0 // Start of frame for baseline DCT
	SOF2 = 0xC2 // Start of frame for progressive DCT
	DHT  = 0xC4 // Define Huffman Table(s)
	DQT  = 0xDB // Define Quantization Table(s)
	DRI  = 0xDD // Define Restart Interval
	SOS  = 0xDA // Start of Scan
	/** RSTn **/
	RST0 = 0xD0
	RST1 = 0xD1
	RST2 = 0xD2
	RST3 = 0xD3
	RST4 = 0xD4
	RST5 = 0xD5
	RST6 = 0xD6
	RST7 = 0xD7
	/** APPn **/
	APP0  = 0xE0
	APP1  = 0xE1
	APP2  = 0xE2
	APP3  = 0xE3
	APP4  = 0xE4
	APP5  = 0xE5
	APP6  = 0xE6
	APP7  = 0xE7
	APP8  = 0xE8
	APP9  = 0xE9
	APP10 = 0xEA
	APP11 = 0xEB
	APP12 = 0xEC
	APP13 = 0xED
	APP14 = 0xEE
	APP15 = 0xEF
	/****/
	COM = 0xFE // Comments
	EOI = 0xD9
)
