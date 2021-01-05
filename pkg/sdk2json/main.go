package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	RegParam = `, `

	RegDllName = `req\.dll: (?P<DLL>[\w]+\.dll)`
)

func regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func regSubMatchToMapString(regEx, s string) (paramsMap map[string]string) {

	r := regexp.MustCompile(regEx)
	match := r.FindStringSubmatch(s)

	paramsMap = make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}


func removeAnnotations(apiPrototype string) string {
	apiPrototype = strings.Replace(apiPrototype, "_Must_inspect_result_", "", -1)
	apiPrototype = strings.Replace(apiPrototype, "_Success_(return != 0 && return < nBufferLength)", "", -1)
	apiPrototype = strings.Replace(apiPrototype, "_Success_(return != 0 && return < cchBuffer)", "", -1)
	apiPrototype = strings.Replace(apiPrototype, "_Success_(return != FALSE)", "", -1)
	apiPrototype = strings.Replace(apiPrototype, "_Ret_maybenull_", "", -1)
	apiPrototype = strings.Replace(apiPrototype, "_Post_writable_byte_size_(dwSize)", "", -1)
	apiPrototype = strings.Replace(apiPrototype, "__out_data_source(FILE)", "", -1)
	apiPrototype = strings.Replace(apiPrototype, " OPTIONAL", "", -1)


	return apiPrototype
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func standardize(s string) string {
	if strings.HasPrefix(s, "BOOLAPI") {
		s = strings.Replace(s, "BOOLAPI", "BOOL WINAPI", -1)
	} else if strings.HasPrefix(s, "INTERNETAPI_(HINTERNET)") {
		s = strings.Replace(s, "INTERNETAPI_(HINTERNET)", "HINTERNET WINAPI", -1)
	} else if strings.HasPrefix(s, "INTERNETAPI_(DWORD)") {
		s = strings.Replace(s, "INTERNETAPI_(DWORD)", "DWORD WINAPI", -1)
	} else if strings.HasPrefix(s, "STDAPI") {
		s = strings.Replace(s, "STDAPI", "HRESULT WINAPI", -1)
	}
	return s
}

// WriteStrSliceToFile writes a slice of string line by line to a file.
func WriteStrSliceToFile(filename string, data []string) (int, error) {
	// Open a new file for writing only
	file, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Create a new writer.
	w := bufio.NewWriter(file)
	nn := 0
	for _, s := range data {
		n, _ := w.WriteString(s + "\n")
		nn += n
	}

	w.Flush()
	return nn, nil
}

// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {

	var (
		part   []byte
		prefix bool
	)

	// Start by getting a file descriptor over the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// SliceContainsStringReverse returns if slice contains substring
func SliceContainsStringReverse(a string, list []string) bool {
	for _, b := range list {
		b = " " + b
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func getDLLName(file, apiname, sdkpath string) (string, error) {
	cat := strings.TrimSuffix(filepath.Base(file), ".h")
	functionName := "nf-" + cat + "-" + strings.ToLower(apiname) + ".md"
	mdFile := path.Join(sdkpath, "sdk-api-src", "content", cat, functionName)
	mdFileContent, err := utils.ReadAll(mdFile)
	if err != nil {
		log.Printf("Failed to find file: %s", mdFile)
		return "", err
	}
	m := regSubMatchToMapString(RegDllName, string(mdFileContent))
	return strings.ToLower(m["DLL"]), nil
}

func unique(slice []string) []string {
	encountered := map[string]int{}
	diff := []string{}

	for _, v := range slice {
		encountered[v] = encountered[v] + 1
	}

	for _, v := range slice {
		if encountered[v] == 1 {
			diff = append(diff, v)
		}
	}
	return diff
}

func main() {

	// Parse arguments.
	// C:\Program Files (x86)\Windows Kits\10\Include\10.0.19041.0\
	sdkumPath := flag.String("sdk", "", "The path to the windows sdk directory")
	// https://github.com/MicrosoftDocs/sdk-api
	sdkapiPath := flag.String("sdk-api", "sdk-api", "The path to the sdk-api docs directory")
	hookapisPath := flag.String("hookapis", "hookapis.txt", "The path to a a text file which define which APIs to trace, new line separated.")
	printretval := flag.Bool("printretval", false, "Print return value type for each API")

	printanno := flag.Bool("printanno", false, "Print list of annotation values")
	minify := flag.Bool("minify", false, "Mininify json")

	flag.Parse()

	if *sdkumPath == "" {
		flag.Usage()
		os.Exit(0)
	}

	if !Exists(*sdkumPath) {
		log.Fatal("sdk directory does not exist")
	}

	if !Exists(*sdkapiPath) {
		log.Fatal("sdk-api directory does not exist")
	}
	if !Exists(*hookapisPath) {
		log.Fatal("hookapis.txt does not exists")
	}

	// Read the list of APIs we are interested to keep.
	wantedAPIs, err := readLines(*hookapisPath)
	if err != nil {
		log.Fatalln(err)
	}
	if len(wantedAPIs) == 0 {
		log.Fatal("hookapis.txt is empty")
	}

	files, err := utils.WalkAllFilesInDir(*sdkumPath)
	if err != nil {
		log.Fatalln(err)
	}

	m := make(map[string]map[string]API)
	var winStructsRaw []string
	var winStructs []Struct

	parsedAPI := 0
	for _, file := range files {

		var prototypes []string

		file = strings.ToLower(file)

		if !strings.HasSuffix(file, "\\fileapi.h") &&
			!strings.HasSuffix(file, "\\processthreadsapi.h") &&
			!strings.HasSuffix(file, "\\winreg.h") &&
			!strings.HasSuffix(file, "\\bcrypt.h") &&
			!strings.HasSuffix(file, "\\winbase.h") &&
			!strings.HasSuffix(file, "\\urlmon.h") &&
			!strings.HasSuffix(file, "\\memoryapi.h") &&
			!strings.HasSuffix(file, "\\tlhelp32.h") &&
			!strings.HasSuffix(file, "\\debugapi.h") &&
			!strings.HasSuffix(file, "\\handleapi.h") &&
			!strings.HasSuffix(file, "\\winsvc.h") &&
			!strings.HasSuffix(file, "\\libloaderapi.h") &&
			!strings.HasSuffix(file, "\\sysinfoapi.h") &&
			!strings.HasSuffix(file, "\\winuser.h") &&
			!strings.HasSuffix(file, "\\winhttp.h") &&
			!strings.HasSuffix(file, "\\wininet.h") {
			continue
		}

		// Read Win32 include API headers.
		data, err := utils.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
		}

		// Start parsing all struct in header file.
		a, b := getAllStructs(string(data))
		winStructsRaw = append(winStructsRaw, a...)
		winStructs = append(winStructs, b...)

		// Grab all API prototypes
		// 1. Ignore: FORCEINLINE
		r := regexp.MustCompile(RegAPIs)
		matches := r.FindAllString(string(data), -1)

		for _, v := range matches {
			prototype := removeAnnotations(v)
			prototype = standardizeSpaces(prototype)
			prototype = standardize(prototype)
			prototypes = append(prototypes, prototype)

			// Only parse APIs we want to hook.
			if !SliceContainsStringReverse(prototype, wantedAPIs) {
				continue
			}

			// Parse the API prototype.
			papi := parseAPI(prototype)

			// Find which DLL this API belongs to. Unfortunately, the sdk does
			// not give you this information, we look into the sdk-api markdown
			// docs instead. (Normally, we could have parsed everything from
			// the md files, but they are missing the parameters type!)
			dllname, err := getDLLName(file, papi.Name, *sdkapiPath)
			if err != nil {
				continue
			}
			if _, ok := m[dllname]; !ok {
				m[dllname] = make(map[string]API)
			}
			m[dllname][papi.Name] = papi
			parsedAPI++
		}

		if len(prototypes) > 0 {
			// Write raw prototypes to a text file.
			WriteStrSliceToFile("prototypes-"+filepath.Base(file)+".inc", prototypes)
		}
	}

	if len(m) > 0 {
		// Marshall and write to json file.
		data, _ := json.MarshalIndent(m, "", " ")
		utils.WriteBytesFile("apis.json", bytes.NewReader(data))
	}

	var foundAPIs []string
	if *printretval {
		for dll, v := range m {
			log.Printf("DLL: %s\n", dll)
			log.Println("====================")
			for api, vv := range v {
				log.Printf("API: %s:%s() => %s\n", vv.CallingConvention, api, vv.ReturnValueType)
				if !utils.StringInSlice(api, wantedAPIs) {
					log.Printf("Not found")
				}
				foundAPIs = append(foundAPIs, api)
			}
		}
	}

	log.Printf("Parsed API count: %d, Hooked API Count: %d", parsedAPI, len(wantedAPIs))
	log.Print(unique(append(wantedAPIs, foundAPIs...)))

	// Write struct results
	WriteStrSliceToFile("winstructs.h", winStructsRaw)

	b, _ := json.MarshalIndent(winStructs, "", " ")
	utils.WriteBytesFile("structs.json", bytes.NewReader(b))

	if *printanno || *minify {
		data, err := utils.ReadAll("apis.json")
		if err != nil {
			log.Fatalln(err)
		}
		apis := make(map[string]map[string]API)
		err = json.Unmarshal(data, &apis)
		if err != nil {
			log.Fatalln(err)
		}

		if *printanno {
			var annotations []string
			var types []string
			for _, v := range apis {
				for _, vv := range v {
					for _, param := range vv.Params {
						if !utils.StringInSlice(param.Annotation, annotations) {
							annotations = append(annotations, param.Annotation)
							// log.Println(param.Annotation)
						}

						if !utils.StringInSlice(param.Type, types) {
							types = append(types, param.Type)
							log.Println(param.Type)
						}
					}
				}
			}
		}


		if *minify {
			data, _ := json.Marshal(minifyAPIs(apis))
			utils.WriteBytesFile("mini-apis.json", bytes.NewReader(data))
		}
		os.Exit(0)
	}
}
