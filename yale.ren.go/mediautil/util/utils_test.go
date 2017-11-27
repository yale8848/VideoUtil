// Create by Yale 2017/11/27 16:37
package util

import (
	"testing"
)

func Test_IteratorFiles(t *testing.T) {
	ll := NewListExt()
	IteratorFiles("C:\\YaleSoftFiles\\WorkSpace\\IntelliJ\\EasyStudyCard\\deploy", ll, func(s string) bool {
		return true
	})
	for e := ll.Front(); e != nil; e = e.Next() {
		_, ret := GetCurrentFileDir(e.Value.(string))
		if ret {
			GetFileName(e.Value.(string))
			//fmt.Println()
			//CreateDir(v + FILE_SEP + "AA")
		}
	}
	t.Log("")
}
