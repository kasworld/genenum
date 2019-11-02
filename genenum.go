// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func makeGenComment() string {
	return fmt.Sprintf("// Code generated by \"%s %s\"\n",
		filepath.Base(os.Args[0]), strings.Join(os.Args[1:], " "))
}

// loadEnumWithComment load list of enum + comment
func loadEnumWithComment(filename string) ([][]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	rtn := make([][]string, 0)
	rd := bufio.NewReader(fd)
	for {
		line, err := rd.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) != 0 && line[0] != '#' {
			s2 := strings.SplitN(line, " ", 2)
			if len(s2) == 1 {
				s2 = append(s2, "")
			}
			rtn = append(rtn, s2)
		}
		if err != nil { // eof
			break
		}
	}
	return rtn, nil
}

// saveTo save go source with format, saved file may need goimport
func saveTo(outdata *bytes.Buffer, buferr error, outfilename string) error {
	if buferr != nil {
		fmt.Printf("fail %v %v\n", outfilename, buferr)
		return buferr
	}
	src, err := format.Source(outdata.Bytes())
	if err != nil {
		fmt.Println(outdata)
		fmt.Printf("fail %v %v\n", outfilename, err)
		return err
	}
	if werr := ioutil.WriteFile(outfilename, src, 0644); werr != nil {
		fmt.Printf("fail %v %v\n", outfilename, werr)
		return werr
	}
	fmt.Printf("goimports -w %v\n", outfilename)
	return nil
}

var (
	typename    = flag.String("typename", "", "enum typename")
	basedir     = flag.String("basedir", "", "base directory of enumdata, gen code ")
	packagename = flag.String("packagename", "", "load basedir/packagename.enum")
	genstats    = flag.Bool("genstats", false, "generate stats package")
)

func main() {
	flag.Parse()

	if *typename == "" {
		fmt.Println("typename not set")
	}
	if *packagename == "" {
		fmt.Println("packagename not set")
	}
	if *basedir == "" {
		fmt.Println("base dir not set")
	}

	os.MkdirAll(path.Join(*basedir, *packagename), os.ModePerm)

	enumdatafile := path.Join(*basedir, *packagename+".enum")
	enumdata, err := loadEnumWithComment(enumdatafile)
	if err != nil {
		fmt.Printf("fail to load %v %v\n", enumdatafile, err)
		return
	}

	buf, err := buildEnumCode(*packagename, *typename, enumdata)
	saveTo(buf, err, path.Join(*basedir, *packagename, *packagename+"_gen.go"))

	if *genstats {
		os.MkdirAll(path.Join(*basedir, *packagename+"_stats"), os.ModePerm)
		buf, err = buildStatsCode(*packagename, *typename)
		saveTo(buf, err, path.Join(*basedir, *packagename+"_stats", *packagename+"_stats_gen.go"))
	}
}

func buildEnumCode(
	pkgname string, typename string, enumdata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s
		import "fmt"
		type %v uint8
	`, pkgname, typename)

	fmt.Fprintf(&buf, "const (\n")
	for i, v := range enumdata {
		if i == 0 {
			fmt.Fprintf(&buf, "%v %v = iota // %v \n", v[0], typename, v[1])
		} else {
			fmt.Fprintf(&buf, "%v // %v\n", v[0], v[1])
		}
	}
	fmt.Fprintf(&buf, `
	%[1]s_Count int = iota 
	)`, typename)

	fmt.Fprintf(&buf, `
	var _%[1]s2string = map[%[1]s]string{
	`, typename)

	for _, v := range enumdata {
		fmt.Fprintf(&buf, "%v : \"%v\", \n", v[0], v[0])
	}
	fmt.Fprintf(&buf, "\n}\n")
	fmt.Fprintf(&buf, `
	func (e %[1]s) String() string {
		if s, exist := _%[1]s2string[e]; exist {
			return s
		}
		return fmt.Sprintf("%[1]s%%d", uint8(e))
	}
	`, typename)

	fmt.Fprintf(&buf, `
	var _string2%[1]s = map[string]%[1]s{
	`, typename)

	for _, v := range enumdata {
		fmt.Fprintf(&buf, "\"%v\" : %v, \n", v[0], v[0])
	}
	fmt.Fprintf(&buf, "\n}\n")
	fmt.Fprintf(&buf, `
	func  String2%[1]s(s string) (%[1]s, bool) {
		v, b :=  _string2%[1]s[s]
		return v,b
	}
	`, typename)

	return &buf, nil
}

func buildStatsCode(pkgname string, typename string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
	package %[1]s_stats
	import (
		"bytes"
		"fmt"
		"html/template"
		"net/http"
	)
	`, pkgname, typename)

	fmt.Fprintf(&buf, `
	type %[2]sStat [%[1]s.%[2]s_Count]int
	func (es *%[2]sStat) String() string {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%[1]s[")
		for i, v := range es {
			fmt.Fprintf(&buf,
				"%%v:%%v ",
				%[1]s.%[2]s(i), v)
		}
		buf.WriteString("]")
		return buf.String()
	}
	func (es *%[2]sStat) Inc(act %[1]s.%[2]s) {
		es[act]++
	}
	func (es *%[2]sStat) Add(act %[1]s.%[2]s, v int) {
		es[act]+=v
	}
	func (es *%[2]sStat) Get(act %[1]s.%[2]s) int {
		return es[act]
	}
	
	func (es *%[2]sStat) ToWeb(w http.ResponseWriter, r *http.Request) error {
		tplIndex, err := template.New("index").Funcs(IndexFn).Parse(%[3]c
		<html>
		<head>
		<title>%[2]s stat Info</title>
		</head>
		<body>
		<table border=1 style="border-collapse:collapse;">%[3]c +
			HTML_tableheader +
			%[3]c{{range $i, $v := .}}%[3]c +
			HTML_row +
			%[3]c{{end}}%[3]c +
			HTML_tableheader +
			%[3]c</table>
	
		<br/>
		</body>
		</html>
		%[3]c)
		if err != nil {
			return err
		}
		if err := tplIndex.Execute(w, es); err != nil {
			return err
		}
		return nil
	}
	
	func Index(i int) string {
		return %[1]s.%[2]s(i).String()
	}
	
	var IndexFn = template.FuncMap{
		"%[2]sIndex": Index,
	}
	
	const (
		HTML_tableheader = %[3]c<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>%[3]c
		HTML_row = %[3]c<tr>
		<td>{{%[2]sIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		%[3]c
	)
	`, pkgname, typename, '`')

	return &buf, nil
}
