package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	// "encoding/hex"
	// "github.com/donutloop/toolkit/debugutil"

	peparser "github.com/saferwall/saferwall/pkg/peparser"
)

func prettyPrint(buff []byte) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, buff, "", "\t")
	if error != nil {
		log.Println("JSON parse error: ", error)
		return string(buff)
	}

	return string(prettyJSON.Bytes())
}

func printAnomalies(anomalies []string) {
	log.Printf("Anomalies: \n")
	for _, ano := range anomalies {
		log.Printf("         - %s\n", ano)
	}
}


func parse(filename string) {

	// fmt.Println("Processing: ", filename)
	pe, err := peparser.Open(filename)
	if err != nil {
		// log.Printf("Error while opening file: %s, reason: %s", filename, err)
		return
	}

	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			log.Printf("%s\n", filename)
			log.Printf("Exception: %v\n", err)
			if len(pe.Anomalies) > 0 {
				printAnomalies(pe.Anomalies)
			}
			log.Println("==============================================================================")
		}
	}()

	err = pe.Parse()
	if err != nil && 
		err != peparser.ErrImageOS2SignatureFound &&
		err != peparser.ErrDOSMagicNotFound &&
		err != peparser.ErrImageNtSignatureNotFound &&
		err != peparser.ErrImageNtOptionalHeaderMagicNotFound &&
		err != peparser.ErrImageBaseNotAligned &&
		err != peparser.ErrImageOS2LESignatureFound &&
		err != peparser.ErrImageVXDSignatureFound && 
		err != peparser.ErrInvalidPESize &&
		err != peparser.ErrInvalidElfanewValue {
		log.Printf("%s\n", filename)
		log.Printf("Error: %v\n", err)
		if len(pe.Anomalies) > 0 {
			printAnomalies(pe.Anomalies)
		}
		log.Println("=====================================================================")
	}

	// if err == peparser.ErrDOSMagicNotFound {
	// 	if strings.Contains(filename, "Windows 10 x64") && 
	// 		strings.HasSuffix(filename, ".dll"){
	// 		// os.Remove(filename)
	// 	}
	// }

	// var buff []byte
	// buff, err = json.Marshal(pe.RichHeader)
	// fmt.Print(prettyPrint(buff))
	pe.Close()
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

	// fmt.Println(hex.EncodeToString(pe.Authentihash()))
	// pe.GetAnomalies()
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
