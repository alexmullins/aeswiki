// Program grabs the AES SBox Wiki page and checks whether the published
// SBox/Inverse Table matches the published SBox/Inverse byte values.

package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PreType int

const PreCount = 4

const (
	SBoxTable PreType = iota
	SBoxBytes
	SBoxInvTable
	SBoxInvBytes
)

func (p PreType) String() string {
	switch p {
	case 0:
		return "SBoxTable"
	case 1:
		return "SBoxBytes"
	case 2:
		return "SBoxInvTable"
	case 3:
		return "SBoxInvBytes"
	default:
		return "Unknown"
	}
}

var (
	url   = "https://en.wikipedia.org/wiki/Rijndael_S-box"
	prees = make(map[PreType]string)
)

func parseSBoxTable(table string) []byte {
	// Ignore 1st and 2nd rows.
	// Loop over each subsequent line:
	// 		1. Split on '|' and ignore everything before
	//		2. Split on ' ' and convert each element to a
	tr := strings.NewReader(table)
	scanner := bufio.NewScanner(tr)

	b := make([]byte, 0)

	i := 0
	for scanner.Scan() {
		if i < 2 {
			// Skip over the first 2 rows
			i++
			continue
		}
		rawLine := scanner.Text()
		rawSplit := strings.Split(rawLine, "|")
		if len(rawSplit) != 2 {
			return nil
		}
		byteLine := strings.Replace(rawSplit[1], " ", "", -1)
		bytes, err := hex.DecodeString(byteLine)
		if err != nil || len(bytes) != 16 {
			return nil
		}
		b = append(b, bytes...)
	}

	if err := scanner.Err(); err != nil {
		return nil
	}
	return b
}

func parseSBoxBytes(table string) []byte {
	// Ignore 1st and 2nd rows
	// Skip last row '};'
	// Every other row:
	// 		1. Trimspace
	//		2. Replace ',' -> ''
	//		3. Replace ' ' -> ''
	//		3. Replace '0x' -> ''

	tr := strings.NewReader(table)
	scanner := bufio.NewScanner(tr)

	b := make([]byte, 0)

	i := 0

	for scanner.Scan() {
		if i < 2 {
			// Skip the first 2 rows
			i++
			continue
		}
		line := scanner.Text()
		if strings.Contains(line, "};") {
			// Skip the last line
			continue
		}
		line = strings.TrimSpace(line)
		line = strings.Replace(line, ",", "", -1)
		line = strings.Replace(line, " ", "", -1)
		line = strings.Replace(line, "0x", "", -1)

		bytes, err := hex.DecodeString(line)
		if err != nil || len(bytes) != 16 {
			return nil
		}

		b = append(b, bytes...)
	}

	if err := scanner.Err(); err != nil {
		return nil
	}
	return b
}

func main() {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	preElems := doc.Find("#mw-content-text > pre")

	if preElems.Length() != PreCount {
		panic("Did not find enough elements on the page.")
	}

	preElems.Each(func(i int, sel *goquery.Selection) {
		prees[PreType(i)] = sel.Text()
	})

	tableB := parseSBoxTable(prees[SBoxTable])
	bytesB := parseSBoxBytes(prees[SBoxBytes])

	if !bytes.Equal(tableB, bytesB) {
		panic("SBoxTable and SBoxBytes do not match")
	}

	iTableB := parseSBoxTable(prees[SBoxInvTable])
	iBytesB := parseSBoxBytes(prees[SBoxInvBytes])

	if !bytes.Equal(iTableB, iBytesB) {
		panic("SBoxInvTable and SBoxInvBytes do not match.")
	}

	fmt.Println("Everything matches.")
}
