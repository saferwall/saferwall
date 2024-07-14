# Saferwall Query Language

`ql` package implements a parsing and transpiling of _search terms or modifiers_ into Couchbase SQL++.

## API and Workflow

The package API is rather simple.

```go

package main

import ql "github.com/saferwall/internal/ql"

func main() {

    exampleQuery := "type:pe and size:1mb+  and ls:2021-03-01T00:00:00"
    n1ql := ql.NewParser().Parse(exampleQuery).Compile()
    fmt.Println(n1ql)
}

```

```sh

$ go run example.go

SELECT * from bucket where filetype == "pe" AND filesize > 300 AND lastSubmission = 2021-03-01T00:00:00

```



## Query Language Reference

MODIFIER

### Modifiers

TODO: Add full list of supported modifiers.

### Literals

* Hexadecimal strings.
* Integers.

```
"type": "file"
"md5": "6468ee100d88c71d55dfdcf4e30f991e"
"sha1": "5c520d2d7dc4c9e5d536d3aff998185657d40ac8"
"sha256": "b102ed1018de0b7faea37ca86f27ba3025c0c70f28417ac3e9ef09d32617f801"
"sha512": "41913eb5adaab42c7ebff547421c0faedede5a3356cb2aa8b92ab20320f73766101056853f450435281cf31e7f32603c62fbd88fa3a680b19abda5d8cc9a98ae"
"ssdeep": "768:QzG3EG0IUJrd6dQar/MjfW33AMar6q3Fu:QKEG4Jx6Ky/Mjo3AMa13U"
"crc32": "0x7017fca6"
"size": 32768
"magic": "PE32+ executable (GUI) x86-64, for MS Windows"
"trid": [

```

### File Information

The following keywords are valid search terms for all file types.

* ```size```
* ```type```
* ```fs```
* ```ls```
* ```positives```
* ```name```
* ```tag```
* ```meta```

### Filetype Specific Keywords

The following keywords are valid search terms for binary files only.

* ```section```
* ```imports```
* ```exports```

### Operators

You can apply operators on different modifiers to build precise search terms.

* ```and``` : boolean AND operation, both modifiers must be satisified.
* ```or```  : boolean OR operation, only a single modifier needs to be satisfied.

## Supported Modifiers

| Modifier  | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| size      | Filters the files to be returned according to size. The size can be specified in bytes (default), kilobytes or megabytes. Trailing plus or minus sign will retrieve those files with a size, respectively, larger than or smaller than the one provided.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| type      | Filters the type of file to be returned (i.e. magic signature). Example: type:elf. This is the full list of available file type literals: Executables: peexe, pedll, neexe, nedll, mz, msi, com, coff, elf, krnl, rpm, linux, macho. Internet: html, xml, flash, fla, iecookie, bittorrent, email, outlook, cap. Phones&tablets: symbian, palmos, wince, android, iphone. Images: jpeg, emf, tiff, gif, png, bmp, gimp, indesign, psd, targa, xws, dib, jng, ico, fpx, eps, svg. Video&audio: ogg, flc, fli, mp3, flac, wav, midi, avi, mpeg, qt, asf, divx, flv, wma, wmv, rm, mov, mp4, 3gp. Documents: text, pdf, ps, doc, docx, rtf, ppt, pptx, xls, xlsx, odp, ods, odt, hwp, gul, ebook, latex. Bundles: isoimage, zip, gzip, bzip, rzip, dzip, 7zip, cab, jar, rar, mscompress, ace, arc, arj, asd, blackhole, kgb. Code: script, php, python, perl, ruby, c, cpp, java, shell, pascal, awk, dyalog, fortran, java-bytecode. Apple: apple, mac, applesingle, appledouble, machfs, appleplist, maclib. Miscellaneous: lnk, ttf, rom. |
| fs        | Filters the files to be returned according to the first submission date.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| ls        | Filters the files to be returned according to the last submission date.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| positives | Filters the files to be returned according to the number of antivirus vendors that detected it upon scanning with VirusTotal. It allows you to specify larger than or smaller than values.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| name      | Returns the files submitted to VirusTotal with a file name that contains the literal provided.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
