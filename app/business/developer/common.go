package developer

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func CheckIsAddController(modelname, path string) {
	filePath := filepath.Join("app/", modelname, "/controller.go")
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			os.Create(filePath)
			err := CopyFileContents("resource/developer/codetpl/go/controller.gos", filePath)
			if err != nil {
				panic(err)
			}
		}
	}
	con_path := "gofly/app/" + path
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	ishase := false
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), con_path) {
			ishase = true
		}
		result += string(a) + "\n"
	}
	if ishase == false {
		addstr := "	_ \"" + con_path + "\""
		datestr := strings.ReplaceAll(result, ")", addstr)
		result = datestr + ")\n"
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
	fw.Close()
}

func CheckApiRemoveController(modelname, path string) {
	filePath := filepath.Join("app/", modelname, "/controller.go")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}
	con_path := "gofly/app/" + path
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), con_path) {
			// result += datestr + "\n"
			continue
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
	fw.Close()
}

func CheckIsAddAppController(modelname string) {
	filePath := filepath.Join("app/controller.go")
	con_path := "gofly/app/" + modelname
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Print(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	ishase := false
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), con_path) {
			ishase = true
		}
		result += string(a) + "\n"
	}
	if ishase == false {
		addstr := "	_ \"" + con_path + "\""
		datestr := strings.ReplaceAll(result, ")", addstr)
		result = datestr + ")\n"
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
	fw.Close()
}

func CheckApiRemoveAppController(modelname string) {
	filePath := filepath.Join("app/controller.go")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}
	con_path := "gofly/app/" + modelname
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), con_path) {
			continue
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
	fw.Close()
}

func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func CopyAllDir(targetPath string, destPath string) error {
	err := filepath.Walk(targetPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(destPath, path[len(targetPath):])
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return copyFile(path, destPath)
	})
	return err
}

func copyFile(srcFile, destFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}
	err = dest.Sync()
	if err != nil {
		return err
	}
	return nil
}
