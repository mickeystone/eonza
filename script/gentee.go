// Copyright 2020 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package script

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentee/gentee"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/vm"
	"gopkg.in/yaml.v2"
)

type FileLines struct {
	File    *os.File
	Scanner *bufio.Scanner
}

func YamlToMap(in string) (*core.Map, error) {
	var (
		tmp map[string]string
		ret interface{}
		err error
	)
	if err = yaml.Unmarshal([]byte(in), &tmp); err != nil {
		return nil, err
	}
	ret, err = gentee.Go2GenteeType(tmp, `map.str`)
	if err != nil {
		return nil, err
	}
	return ret.(*core.Map), nil
}

func CopyName(rt *vm.Runtime, fname string) (string, error) {
	var (
		err error
		i   int
	)
	if !filepath.IsAbs(fname) {
		fname, err = filepath.Abs(fname)
		if err != nil {
			return ``, err
		}
	}
	dir := filepath.Dir(fname)
	base := strings.SplitN(filepath.Base(fname), `.`, 2)
	if len(base) == 1 {
		base = append(base, ``)
	} else {
		base[1] = `.` + base[1]
	}
	for {
		i++
		exist, err := vm.ExistFile(rt, fname)
		if err != nil {
			return ``, err
		}
		if exist == 0 {
			break
		}
		fname = filepath.Join(dir, fmt.Sprintf("%s (%d)%s", base[0], i, base[1]))
	}
	return fname, nil
}

func CloseLines(flines *FileLines) error {
	return flines.File.Close()
}

func GetLine(flines *FileLines) string {
	return flines.Scanner.Text()
}

func ReadLines(filename string) (*FileLines, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	return &FileLines{File: file, Scanner: scanner}, nil
}

func ScanLines(flines *FileLines) int64 {
	if flines.Scanner.Scan() {
		return 1
	}
	return 0
}

/*
// Subbuf(buf, int, int) buf
func Subbuf(buf *core.Buffer, off int64, size int64) (*core.Buffer, error) {
	if off < 0 || off+size > int64(len(buf.Data)) {
		return nil, fmt.Errorf(vm.ErrorText(core.ErrInvalidParam))
	}
	ret := core.NewBuffer()
	ret.Data = append(ret.Data, buf.Data[off:off+size]...)
	return ret, nil
}
*/
