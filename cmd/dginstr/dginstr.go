// dginstr is a command-line utility to check and convert CSV instruction definitions into
// files usable by various emulators.

// Copyright (C) 2017,2018,2019,2021 Steve Merrony

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	maxTypes   = 20
	maxFormats = 40
	maxInstrs  = 500
	instrAttrs = 7
)

var (
	// command arguments
	actionFlag  = flag.String("action", "", "specify operation to perform ie. 'checkbits', 'makeada', or 'makego'")
	cpuTypeFlag = flag.String("cputype", "", "Target CPU type for output <nova|eclipse|mv>")
	csvFlag     = flag.String("csv", "", "CSV file to source data from")
	goFlag      = flag.String("go", "", "Go filename for output")
	adaFlag     = flag.String("ada", "", "Ada file for output")

	typesList    [maxTypes]string
	formatsList  [maxFormats]string
	formatCounts map[string]int
	instrsTable  [maxInstrs][]string

	// headers = [...]string{"#", "Mnem", "Bits", "BitMask", "Len", "Instruction Format", "Instruction Type"}

	numTypes, numFormats, numInstrs int
	genNova, genEclipse, genMV      bool
)

func main() {
	flag.Parse()

	if *csvFlag == "" {
		log.Fatalln("ERROR: Must specify source CSV file with -csv=<csvfile> argument")
	}
	if *actionFlag == "" {
		log.Fatalln("ERROR: Must specify action with -action=<action> argument i.e. checkbits, makeada, or makego")
	}

	switch *actionFlag {
	case "checkbits":
		if loadCSV() {
			checkBits()
		}
	case "makeada":
		if *adaFlag == "" {
			log.Fatalln("ERROR: Must specify Ada file for output with -ada=<adafile> argument")
		}
		if *cpuTypeFlag == "" {
			log.Fatalln("ERROR: Must specify DG -cpuType when generating Ada file")
		}
		switch *cpuTypeFlag {
		case "nova":
			genNova = true
		case "eclipse":
			genNova, genEclipse = true, true
		case "mv":
			genNova, genEclipse, genMV = true, true, true
		default:
			log.Fatalln("ERROR: cpuType must be one of nova, eclipse or mv")
		}
		if loadCSV() {
			exportAda()
		}
	case "makego":
		if *goFlag == "" {
			log.Fatalln("ERROR: Must specify Go file for output with -go=<gofile> argument")
		}
		if *cpuTypeFlag == "" {
			log.Fatalln("ERROR: Must specify DG -cpuType when generating Go file")
		}
		switch *cpuTypeFlag {
		case "nova":
			genNova = true
		case "eclipse":
			genNova, genEclipse = true, true
		case "mv":
			genNova, genEclipse, genMV = true, true, true
		default:
			log.Fatalln("ERROR: cpuType must be one of nova, eclipse or mv")
		}
		if loadCSV() {
			exportGo()
		}
	default:
		log.Fatalln("ERROR: No such action")
	}
}

func loadCSV() bool {

	csvFile, err := os.Open(*csvFlag)
	if err != nil {
		log.Fatalf("ERROR: Could not open CSV file %v", err)
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	line, err := csvReader.Read()
	if err != nil {
		log.Fatalf("ERROR: Could not read CSV file %v", err)
	}
	if line[0] != ";Types" {
		log.Printf("Error: expecting <;Types> got <%s>\n", line[0])
		return false
	}

	// reset data counts
	numTypes = 0
	numInstrs = 0
	numInstrs = 0
	formatCounts = make(map[string]int)

	numTypes = 0
	for {
		line, _ = csvReader.Read()
		if line[0] == ";" {
			break
		}
		if (genNova && strings.Contains(line[0], "NOVA")) ||
			(genEclipse && strings.Contains(line[0], "ECLIPSE")) ||
			(genMV && strings.Contains(line[0], "EAGLE")) {
			typesList[numTypes] = line[0]
			numTypes++
		}
	}

	line, _ = csvReader.Read()
	if line[0] != ";Formats" {
		log.Printf("Error: expecting <;Formats> got <%s>\n", line[0])
		return false
	}

	numFormats = 0
	for {
		line, _ = csvReader.Read()
		if line[0] == ";" {
			break
		}
		formatsList[numFormats] = line[0]
		formatCounts[line[0]] = 0
		//log.Printf("Loading format #%d: %s\n", numFormats, line[0])
		numFormats++
	}

	line, _ = csvReader.Read()
	if line[0] != ";Instructions" {
		log.Printf("Error: expecting <;Instructions> got <%s>\n", line[0])
		return false
	}

	numInstrs = 0
	for {
		line, _ = csvReader.Read()
		if line[0] == ";" {
			break
		}
		if *actionFlag == "checkbits" ||
			(genNova && strings.Contains(line[5], "NOVA")) ||
			(genEclipse && strings.Contains(line[5], "ECLIPSE")) ||
			(genMV && strings.Contains(line[5], "EAGLE")) {
			row := make([]string, 7)
			for c := 0; c < instrAttrs; c++ {
				row[c] = line[c]
			}
			instrsTable[numInstrs] = row
			formatCounts[line[4]]++
			numInstrs++
		}
	}

	csvFile.Close()
	fmt.Printf("%d instruction definitions read from CSV\n", numInstrs)
	return true
}

func exportAda() bool {
	adaFile, err := os.Create(*adaFlag)
	if err != nil {
		log.Println(err)
		return false
	}
	adaWriter := bufio.NewWriter(adaFile)

	fmt.Fprintf(adaWriter, "-- Code generated by dginstr.go; DO NOT EDIT.\n\n")
	fmt.Fprintf(adaWriter, `-- Permission is hereby granted, free of charge, to any person obtaining a copy
-- of this software and associated documentation files (the "Software"), to deal
-- in the Software without restriction, including without limitation the rights
-- to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
-- copies of the Software, and to permit persons to whom the Software is
-- furnished to do so, subject to the following conditions:
-- The above copyright notice and this permission notice shall be included in
-- all copies or substantial portions of the Software.

-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
-- AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
-- OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
-- THE SOFTWARE.

with Ada.Strings.Unbounded; use Ada.Strings.Unbounded;

package CPU.Instructions is

`)

	fmt.Fprintln(adaWriter, "-- Instruction Classes")
	fmt.Fprintln(adaWriter, "type Instr_Class_T is (")
	for t := 0; t < numTypes; t++ {
		if t > 0 {
			fmt.Fprintln(adaWriter, ",")
		}
		fmt.Fprintf(adaWriter, "\t%s", typesList[t])
	}

	fmt.Fprintf(adaWriter, "\n);\n\n-- Instruction Formats\n")
	fmt.Fprintln(adaWriter, "type Instr_Format_T is (")
	first := true
	for f := 0; f < numFormats; f++ {
		if formatCounts[formatsList[f]] > 0 {
			if first {
				first = false
			} else {
				fmt.Fprintln(adaWriter, ",")
			}
			fmt.Fprintf(adaWriter, "\t%s", formatsList[f])
		}
	}

	fmt.Fprintf(adaWriter, "\n);\n\n-- Instruction Mnemonic Consts\n")
	fmt.Fprintln(adaWriter, "type Instr_Mnemonic_T is (")
	for i := 0; i < numInstrs; i++ {
		if i > 0 {
			fmt.Fprintln(adaWriter, ",")
		}
		fmt.Fprintf(adaWriter, "   I_%s", instrsTable[i][0])
	}

	fmt.Fprintf(adaWriter, "\n);\n")
	fmt.Fprintf(adaWriter,
		`
type Instr_Char_Rec is record
   Mnemonic    : Unbounded_String;
   Bits        : Word_T;
   Mask        : Word_T;
   Instr_Len   : Positive;
   Instr_Fmt   : Instr_Format_T;
   Instr_Class : Instr_Class_T;
   Disp_Offset : Natural;
end record;

type Instructions is array (Instr_Mnemonic_T range Instr_Mnemonic_T'Range) of Instr_Char_Rec;

Instruction_Set : constant Instructions :=
(
`)

	for i := 0; i < numInstrs; i++ {
		fmt.Fprintf(adaWriter, "I_%s => (To_Unbounded_String(\"%s\"), 16#%s#, 16#%s#, %s, %s, %s, %s)",
			instrsTable[i][0],
			instrsTable[i][0],
			instrsTable[i][1][2:],
			instrsTable[i][2][2:],
			instrsTable[i][3],
			instrsTable[i][4],
			instrsTable[i][5],
			instrsTable[i][6])
		if i < numInstrs-1 {
			fmt.Fprintf(adaWriter, ",\n")
		}
	}

	fmt.Fprintf(adaWriter, "\n);\n\nend CPU.Instructions;\n")
	adaWriter.Flush()
	adaFile.Close()
	fmt.Println("Ada file written")
	return true
}
func exportGo() bool {
	goFile, err := os.Create(*goFlag)
	if err != nil {
		log.Println(err)
		return false
	}
	goWriter := bufio.NewWriter(goFile)

	fmt.Fprintf(goWriter, "// Code generated by dginstr.go; DO NOT EDIT.\n\n")
	fmt.Fprintf(goWriter, `// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package mvcpu

`)

	fmt.Fprintf(goWriter, "// Instruction Types\nconst (\n")
	fmt.Fprintf(goWriter, "\t%s = iota\n", typesList[0])
	for t := 1; t < numTypes; t++ {
		fmt.Fprintf(goWriter, "\t%s\n", typesList[t])
	}

	fmt.Fprintf(goWriter, ")\n\n// Instruction Formats\nconst (\n")
	fmt.Fprintf(goWriter, "\t%s = iota\n", formatsList[0])
	for f := 1; f < numFormats; f++ {
		if formatCounts[formatsList[f]] > 0 {
			fmt.Fprintf(goWriter, "\t%s\n", formatsList[f])
		}
	}

	fmt.Fprintf(goWriter, ")\n\n// Instruction Mnemonic Consts\nconst (\n")
	fmt.Fprintf(goWriter, "\tinstr%s = iota\n", instrsTable[0][0])
	for i := 1; i < numInstrs; i++ {
		fmt.Fprintf(goWriter, "\tinstr%s\n", instrsTable[i][0])
	}

	fmt.Fprintf(goWriter, ")\n\n// InstructionsInit initialises the instruction characterstics for each instruction\n")
	fmt.Fprintf(goWriter, "func InstructionsInit() {\n")

	for i := 0; i < numInstrs; i++ {
		fmt.Fprintf(goWriter, "\tinstructionSet[instr%s] = instrChars{\"%s\", %s, %s, %s, %s, %s, %s}\n",
			instrsTable[i][0],
			instrsTable[i][0],
			instrsTable[i][1],
			instrsTable[i][2],
			instrsTable[i][3],
			instrsTable[i][4],
			instrsTable[i][5],
			instrsTable[i][6])
	}

	fmt.Fprintf(goWriter, "}\n")
	goWriter.Flush()
	goFile.Close()
	fmt.Println("Go file written")
	return true
}

// checkBits tests every instruction to ensure that (at least) all set bits are covered by the
// associated bit mask and that patterns are unique
func checkBits() {
	errors := 0
	for i := 0; i < numInstrs; i++ {
		bitsUint, _ := strconv.ParseUint(instrsTable[i][1], 0, 16)
		maskUint, _ := strconv.ParseUint(instrsTable[i][2], 0, 16)
		diff := bitsUint ^ maskUint // XOR
		if (diff&bitsUint) != 0 || len(instrsTable[i][1]) != 6 || len(instrsTable[i][2]) != 6 {
			errors++
			fmt.Printf("Bitmasking error in  %s\n", instrsTable[i][0])
		}
	}
	if errors == 0 {
		fmt.Printf("No bitmasking errors detected in %d instructions\n", numInstrs)
	}
	errors = 0
	insMap := make(map[string]string)
	for i := 0; i < numInstrs; i++ {
		if _, already := insMap[instrsTable[i][1]]; already {
			errors++
			fmt.Printf("Bit pattern for %s is a duplicate for %s\n", instrsTable[i][0], insMap[instrsTable[i][1]])
		} else {
			insMap[instrsTable[i][1]] = instrsTable[i][0]
		}
	}
	if errors == 0 {
		fmt.Printf("No duplication errors detected in %d instructions\n", len(insMap))
	}
}
