// Create by Yale 2017/11/27 16:20
package util

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

var FILE_SEP = "/"

func init() {
	if runtime.GOOS == "windows" {
		FILE_SEP = "\\"
	}
}

var fileSuffix = []string{".mp4", ".mp3", ".dat", ".flv"}

func IsMediaFile(fullName string) bool {
	fs := strings.ToLower(path.Ext(fullName))
	for _, v := range fileSuffix {
		if fs == v {
			return true
		}
	}
	return false
}
func IteratorFiles(path string, list *ListExt, filter func(string) bool) {
	dirs, error := ioutil.ReadDir(path)
	if error != nil {
		fmt.Println("读取文件失败")
		return
	}
	for _, v := range dirs {
		if !v.IsDir() {
			p := path + FILE_SEP + v.Name()
			if filter(p) {
				list.PushBack(p)
			}
		} else {
			IteratorFiles(path+FILE_SEP+v.Name(), list, filter)
		}
	}
}
func ClearDir(path string, sub bool) {
	dirs, error := ioutil.ReadDir(path)
	if error != nil {
		fmt.Println("读取文件失败")
		return
	}
	for _, v := range dirs {
		if !v.IsDir() {
			os.Remove(path + FILE_SEP + v.Name())
		} else {
			if sub {
				ClearDir(path+FILE_SEP+v.Name(), sub)
			}
		}
	}
}
func List2ArrayString(list *list.List) []string {

	items := make([]string, list.Len())
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		items[i] = e.Value.(string)
		i++
	}
	return items
}
func CreateDir(p string) {
	os.Mkdir(p, os.ModePerm)
}
func GetCurrentFileDir(p string) (string, bool) {
	pos := strings.LastIndex(p, FILE_SEP)
	if pos != -1 {
		return p[0:pos], true
	}
	return "", false
}

type ListExt struct {
	list.List
}

func NewListExt() *ListExt {
	list := &ListExt{}
	return list
}
func (list *ListExt) ToString() {
	for e := list.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
func (list *ListExt) SetByIndex(pos int, v interface{}) {
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		if i == pos {
			e.Value = v
		}
		i++
	}
}
func (list *ListExt) GetByIndex(pos int) interface{} {
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		if i == pos {
			return e.Value
		}
		i++
	}
	return nil
}
func (list *ListExt) Iterator(callback func(interface{}, int)) {
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		callback(e.Value, i)
		i++
	}

}

func GetFileNameWithSuffix(fullName string) string {
	pos := strings.LastIndex(fullName, FILE_SEP)
	if pos != -1 {
		return fullName[pos+1:]
	}
	return ""
}

func GetFileNameOnly(fullName string) string {
	pos := strings.LastIndex(fullName, FILE_SEP)
	if pos != -1 {
		p := fullName[pos+1:]
		e := path.Ext(p)
		return strings.TrimSuffix(p, e)
	}
	return ""
}
func IsAudio(p string) bool {
	ex := strings.ToLower(path.Ext(p))
	if ex == ".mp3" {
		return true
	}
	return false
}
