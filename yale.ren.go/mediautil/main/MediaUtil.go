package main

import (
	"VideoUtil/yale.ren.go/mediautil/util"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

const OUTPUT = "output_dxh"

type envModel struct {
	walk.ListModelBase
	list *util.ListExt
}

type myMainWindow struct {
	*walk.MainWindow
	lb             *walk.ListBox
	label          *walk.Label
	model          *envModel
	btn            *walk.PushButton
	rootPath       string
	outputRootPath string
}

func (m *envModel) ItemCount() int {
	return m.list.Len()
}

func (m *envModel) Value(index int) interface{} {
	return m.list.GetByIndex(index)
}

var mainWindow = &myMainWindow{}

func (w *myMainWindow) showMediaFiles(path string) {
	lists := util.NewListExt()
	util.IteratorFiles(path, lists, func(p string) bool {
		return util.IsMediaFile(p) && strings.Index(p, OUTPUT) == -1
	})
	m := &envModel{list: lists}
	w.lb.SetModel(m)
	w.model = m
	go w.startConvert(lists)
}
func getOutPutPath(path string) string {
	p, ret := util.GetCurrentFileDir(path)
	if ret {
		p = p + util.FILE_SEP + OUTPUT
	}
	return p
}
func (w *myMainWindow) preDeal() {

	t := w.outputRootPath
	os.RemoveAll(w.outputRootPath)
	util.CreateDir(t)

}

func (w *myMainWindow) getOutPutFileName(sp string) string {
	p := strings.Replace(sp, w.rootPath, w.outputRootPath, -1)
	fmt.Println("+p+" + p)
	po, ret := util.GetCurrentFileDir(p)
	if ret {

	}
	fmt.Println("+po+" + po)
	n := util.GetFileNameOnly(p)
	sf := path.Ext(p)
	if strings.ToLower(sf) == ".dat" || strings.ToLower(sf) == ".flv" {
		sf = ".mp4"
	}

	util.CreateDir(po)
	return po + util.FILE_SEP + n + sf

}
func exe(src string, dest string, isAudio bool, ch chan bool) {

	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	var cmd *exec.Cmd
	if isAudio {
		cmd = exec.Command(dir+util.FILE_SEP+"ffmpeg.exe", "-i", src, "-ar", "22050", "-ac", "1", dest)
	} else {
		cmd = exec.Command(dir+util.FILE_SEP+"ffmpeg.exe", "-i", src, "-minrate", "100k", "-b:v", "150k", "-maxrate", "200k", "-bufsize", "150k", "-r", "12", "-x264opts", "keyint=15", "-ar", "22050", "-ac", "1", dest)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", cmd.Stdout)
	ch <- true

}
func (w *myMainWindow) setListitem(i int, text string) {
	w.model.list.SetByIndex(i, text)
	w.lb.SetModel(w.model)
}

func (w *myMainWindow) getOutPutText(src, out string) string {

	return " [处理完成...]     " + out
}
func (w *myMainWindow) startConvert(list *util.ListExt) {
	if list.Len() == 0 {
		w.label.SetText("没有发现多媒体内容")
		return
	}
	channel := make(chan bool)
	list.Iterator(func(v interface{}, pos int) {
		p := v.(string)
		outf := w.getOutPutFileName(p)

		w.label.SetText(fmt.Sprintf(" [正在处理]    %s    (%d/%d)", p, pos+1, list.Len()))
		w.setListitem(pos, " [正在处理...]     "+p)
		go exe(p, outf, util.IsAudio(p), channel)
		<-channel

		w.setListitem(pos, w.getOutPutText(p, outf))
	})
	w.label.SetText("全部处理完毕")
}

func (w *myMainWindow) choseFloder() error {
	dlg := new(walk.FileDialog)

	//dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.Title = "Select an Image"
	dlg.ShowReadOnlyCB = true

	if ok, err := dlg.ShowBrowseFolder(w); err != nil {
		return err
	} else if !ok {
		return nil
	}
	w.btn.SetEnabled(false)

	w.rootPath = dlg.FilePath
	w.outputRootPath = dlg.FilePath + util.FILE_SEP + OUTPUT
	w.preDeal()
	w.showMediaFiles(dlg.FilePath)

	return nil
}

func createWindow() {
	mw := MainWindow{
		AssignTo: &mainWindow.MainWindow,
		Title:    "导学号多媒体文件处理V1.1",
		MinSize:  Size{600, 600},
		Layout:   VBox{},
		Children: []Widget{

			Label{
				Text: "处理完后会在当前目录生成output_dxh，最好在本地把每一个视频打开都看看验证一下",
			},
			PushButton{
				AssignTo: &mainWindow.btn,
				Text:     "选取要处理的目录",
				OnClicked: func() {
					mainWindow.choseFloder()
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
func main() {
	createWindow()
}
