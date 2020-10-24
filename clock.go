package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"

	htgotts "github.com/hegedustibor/htgo-tts"
)

type item struct {
	id      *gtk.TreeIter
	time    string
	message string
}

func update() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(store.IterNChildren(nil))
	}
}

var (
	store *gtk.ListStore
)

func onEnableToggled(cell *gtk.CellRendererToggle, pathString string) {
	it, _ := store.GetIterFromString(pathString)
	valueObj, _ := store.GetValue(it, 1)
	value, _ := valueObj.GoValue()
	bValue, _ := value.(bool)
	store.SetValue(it, 1, !bValue)
}

func onTimeEdited(cell *gtk.CellRendererText, pathString string, text string) {
	it, _ := store.GetIterFromString(pathString)
	store.SetValue(it, 2, text)
}

func onMessageEdited(cell *gtk.CellRendererText, pathString string, text string) {
	it, _ := store.GetIterFromString(pathString)
	store.SetValue(it, 3, text)
	speech := htgotts.Speech{Folder: "audio", Language: "en"}
	speech.Speak(text)
}

func main() {
	gtk.Init(nil)
	builder, _ := gtk.BuilderNewFromFile("clock.glade")
	signals := map[string]interface{}{
		"on_enable_toggled": onEnableToggled,
		"on_time_edited":    onTimeEdited,
		"on_message_edited": onMessageEdited,
	}
	builder.ConnectSignals(signals)
	storeObj, _ := builder.GetObject("liststore")
	store = storeObj.(*gtk.ListStore)
	addObj, _ := builder.GetObject("add")
	addButton := addObj.(*gtk.Button)
	addButton.Connect("clicked", func() {
		it := store.Append()
		store.Set(it, []int{0, 1, 2, 3}, []interface{}{1, true, "", ""})
	})
	removeObj, _ := builder.GetObject("remove")
	removeButton := removeObj.(*gtk.Button)
	selectObj, _ := builder.GetObject("selection")
	selection := selectObj.(*gtk.TreeSelection)
	var selectedPath *gtk.TreePath
	selection.Connect("changed", func() {
		var iter *gtk.TreeIter
		var model gtk.ITreeModel
		var ok bool
		model, iter, ok = selection.GetSelected()
		if ok {
			removeButton.SetSensitive(true)
			path, err := model.(*gtk.TreeModel).GetPath(iter)
			if err != nil {
				log.Printf("treeSelectionChangedCB: Could not get path from model: %s\n", err)
				return
			}
			selectedPath = path
		}
	})
	removeButton.Connect("clicked", func() {
		it, _ := store.GetIter(selectedPath)
		store.Remove(it)
	})
	winObj, _ := builder.GetObject("main")
	win := winObj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.ShowAll()
	go update()
	gtk.Main()
}
