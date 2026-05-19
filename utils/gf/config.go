package gf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UpConfFieldData(path string, parameter map[string]interface{}) error {
	file_path := filepath.Join(path, "/resource/config.yml")
	f, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	var is_hose = false
	for {
		is_hose = false
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		for keys, Val := range parameter {
			if strings.Contains(string(a), keys) {
				is_hose = true
				datestr := strings.ReplaceAll(string(a), string(a), fmt.Sprintf("     %v: %v\n", keys, Val))
				result += datestr
			}
		}
		if !is_hose {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
	return nil
}
