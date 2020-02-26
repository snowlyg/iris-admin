package files

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// 创建临时文件
func CreateTemFile(name string, file multipart.File) (string, error) {
	f, err := ioutil.TempFile("upload/tmp/", "*_"+name)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

// 调用os.MkdirAll递归创建文件夹
func CreateFile(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

//  判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) // os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//  WriteFile writes the contents of the output buffer to a file
func WriteFile(filename string, output []byte) error {
	return ioutil.WriteFile(filename, output, 0666)
}

//  AppendFile writes the contents of the output buffer to a file
func AppendFile(filename string, output []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(output); err != nil {
		return err
	}
	return nil
}

// 压缩文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	重置配置文件路径
	该方法目前仅用于测试 ，有些多余
*/
func GetAbsPath(confPath string) string {
	getwd, err := os.Getwd()
	if err != nil {
		color.Red(fmt.Sprintf("Getwd err %v",err) )
	}

	end := filepath.Base(getwd)

	if end != "IrisAdminApi" {
		return filepath.Join(filepath.Dir(getwd) , confPath)
	}

	return confPath
}
