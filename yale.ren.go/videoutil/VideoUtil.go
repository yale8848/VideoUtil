// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/yale8848/gffmpeg"
	"io/ioutil"
	"log"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const OUTPUT = "output"

type envModel struct {
	walk.ListModelBase
	items []string
}

type myMainWindow struct {
	*walk.MainWindow
	lb    *walk.ListBox
	label *walk.Label
	model *envModel
	btn   *walk.PushButton
}

var mainWindow = &myMainWindow{}

func (m *envModel) ItemCount() int {
	return len(m.items)
}

func (m *envModel) Value(index int) interface{} {
	return m.items[index]
}
func destorytemp(path string) {

	filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}

		return nil
	})

}
func exe1(src string, dest string, chanel chan bool) {
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println(dir)
	//11025 22050 44100
	//ffmpeg.exe -i L1.mp3 -ar 22050 -ac 1 1.mp3

	//cmd:=exec.Command("ffmpeg.exe","-i 1-1.mp4 -b:v 100k -r 15 -bufsize 100k -x264opts keyint=25 1-1_.mp4")
	//cmd:=exec.Command("ffmpeg.exe","-i",src, "-b:v", "100k", "-r", "15", "-bufsize", "100k" ,"-x264opts", "keyint=25", dest)
	cmd := exec.Command(dir+"\\"+"ffmpeg.exe", "-i", src, "-b:v", "800k", "-bufsize", "800k", "-x264opts", "keyint=25", dest)
	//cmd := exec.Command(dir+"\\"+"ffmpeg.exe", "-i", src, "-x264opts", "keyint=25", dest)
	//cmd:=exec.Command("ffmpeg.exe","-i",src,"-x264opts","keyint=25",dest)
	//cmd:=exec.Command("t.exe",src,dest)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", cmd.Stdout)
	chanel <- true

}
func exe(src string, dest string, chanel chan bool) {
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	ff, e := gffmpeg.NewGFFmpeg(dir + "\\" + "ffmpeg.exe")
	ff.SetDebug(true)
	if e != nil {
		panic(e)
	}

	getBitRate := func(src int) int {
		if src <= 500 {
			return src
		}
		ns := src / 2
		return ns
	}
	bd := gffmpeg.NewBuilder().SrcPath(src)
	info := ff.Set(bd).GetMediaInfo()
	if info != nil {

		if getBitRate(info.BitRate) != info.BitRate {
			bd.BitRate(getBitRate(info.BitRate)).
				BufSize(getBitRate(info.BitRate))
		}
		bd.KeyInt(25).Threads(4).
			DistPath(dest)
		ff.Set(bd).Start(nil)
	}
	chanel <- true

}

func (w *myMainWindow) set_listitem(i int, text string) {
	w.model.items[i] = text
	w.lb.SetModel(w.model)
}
func (w *myMainWindow) start_convert(base string, items []string) {
	channel := make(chan bool)
	output := base + "\\" + OUTPUT
	remove_dir(output)
	os.Mkdir(output, os.ModePerm)
	for i, v := range items {
		w.label.SetText(fmt.Sprintf(" [正在处理]    %s    (%d/%d)", v, i+1, len(items)))
		w.set_listitem(i, " [正在处理...]            "+v)
		go exe(add_basePath(base, v), add_basePath(output, v), channel)
		<-channel
		w.set_listitem(i, " 处理完毕...           "+v)
	}
	w.label.SetText("全部处理完毕")
	fmt.Println(output)
	cmd := exec.Command("explorer.exe", output)
	cmd.Start()
}
func add_basePath(base string, path string) string {
	return base + "\\" + path
}
func get_fileNam(base string, path string) string {
	return strings.Replace(path, base+"\\", "", -1)
}
func remove_dir(path string) {
	dirs, error := ioutil.ReadDir(path)
	if error != nil {
		fmt.Println("读取文件失败")
		return
	}
	for _, v := range dirs {
		if !v.IsDir() {
			os.Remove(path + "\\" + v.Name())
		}
	}
}
func (w *myMainWindow) lb_showFileList(path string) {
	dirs, error := ioutil.ReadDir(path)
	if error != nil {
		fmt.Println("读取文件失败")
		return
	}
	items := make([]string, 0)
	for _, v := range dirs {
		if !v.IsDir() {
			if strings.Contains(strings.ToLower(v.Name()), ".mp4") {
				items = append(items, v.Name())
			}
		}
	}
	m := &envModel{items: items}
	w.lb.SetModel(m)
	w.model = m
	w.btn.SetEnabled(false)
	go w.start_convert(path, items)

}

func (w *myMainWindow) fd_choseForlder() error {
	dlg := new(walk.FileDialog)

	//dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.Title = "Select an Image"
	dlg.ShowReadOnlyCB = true

	if ok, err := dlg.ShowBrowseFolder(w); err != nil {
		return err
	} else if !ok {
		return nil
	}
	w.lb_showFileList(dlg.FilePath)
	return nil
}

func main() {
	mw := MainWindow{
		AssignTo: &mainWindow.MainWindow,
		Title:    "视频处理(导学号)V1.3",
		MinSize:  Size{600, 600},
		Layout:   VBox{},
		Children: []Widget{

			Label{
				Text: "处理完后的视频都在本目录的output目录里，最好在本地把每一个视频打开都看看，备份原视频",
			},
			PushButton{
				AssignTo: &mainWindow.btn,
				Text:     "选取视频目录",
				OnClicked: func() {
					mainWindow.fd_choseForlder()
				},
			},
			Label{
				Text:     "",
				AssignTo: &mainWindow.label,
			},
			ListBox{
				AssignTo: &mainWindow.lb,
			},
		},
	}

	if _, err := mw.Run(); err != nil {
		log.Fatal(err)
	}
}
