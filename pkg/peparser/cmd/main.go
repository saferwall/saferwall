package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	// "encoding/hex"
	// "github.com/donutloop/toolkit/debugutil"

	peparser "github.com/saferwall/saferwall/pkg/peparser"
	"github.com/saferwall/saferwall/pkg/peparser/pedumper"
)

func dump(filename string) {
	err := pedumper.Dump(filename)
	if err != nil {
		fmt.Println(filename, err)
	}
}

func parse(filename string) {
	fmt.Println("Processing: ", filename)
	pe, err := peparser.Open(filename)
	if err != nil {
		// log.Printf("Error while opening file: %s, reason: %s", filename, err)
		return
	}

	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "%s, Exception: %v\n", filename, err)
		}
	}()

	err = pe.Parse()
	if err != nil {
		fmt.Println(filename, err)
	}

	// if err == nil {
	// 	if pe.IsDLL() {
	// 		log.Print("File is DLL")
	// 	}
	// 	if pe.IsDriver() {
	// 		log.Print("File is Driver")
	// 	}
	// 	if pe.IsEXE() {
	// 		log.Print("File is Exe")
	// 	}
	// }

	// if len(pe.Anomalies) > 0 {
	// 	fmt.Printf("Anomalies found while parsing %s\n", filename)
	// 	for _, anomaly := range pe.Anomalies {
	// 		fmt.Println(anomaly)
	// 	}
	// }
	// for _, s := range pe.Sections {
	// 	fmt.Println(s.NameString(), pe.PrettySectionFlags(s.Characteristics))
	// }

	// fmt.Println()
	// fmt.Println(hex.EncodeToString(pe.Authentihash()))
	// pe.GetAnomalies()
	// fmt.Println(debugutil.PrettySprint(pe.DosHeader))
	// fmt.Println(debugutil.PrettySprint(pe.NtHeader))
	// fmt.Println(debugutil.PrettySprint(pe.FileHeader))
	// fmt.Println(pe.PrettyImageFileCharacteristics())
	// fmt.Println(pe.PrettyDllCharacteristics())
	// fmt.Println(pe.Checksum())

	// fmt.Print()
	// fmt.Println(debugutil.PrettySprint(pe.BoundImports))

	// for _, imp := range pe.Imports {
	// 	log.Println(imp.Name)
	// 	log.Println("=============================================")
	// 	for _, function := range imp.Functions {
	// 		hint := fmt.Sprintf("%X", function.Hint)
	// 		offset := fmt.Sprintf("%X", function.Offset)

	// 		log.Printf("%s, hint: 0x%s, thunk: 0x%s", function.Name, hint, offset)
	// 	}
	// 	log.Println("=============================================")

	// }

	pe.Close()

}

func main() {
	var searchDir string

	if len(os.Args) > 1 {
		searchDir = os.Args[1]

	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		searchDir = currentDir + string(os.PathSeparator) + "bin"
	}

	log.Printf("Processing directory %s", searchDir)

	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && !strings.HasSuffix(path, ".xml") &&
			!strings.HasSuffix(path, ".bat") && !strings.HasSuffix(path, ".js") &&
			!strings.HasSuffix(path, ".chm") && !strings.HasSuffix(path, ".jar") &&
			!strings.HasSuffix(path, ".cmd") && !strings.HasSuffix(path, ".ps1") {
			fileList = append(fileList, path)
		}
		return nil
	})

	for _, file := range fileList {
		parse(file)
	}
}
