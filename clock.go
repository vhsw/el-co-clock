package main

import (
	"time"

	"github.com/gotk3/gotk3/gtk"
)

func update(store *gtk.ListStore, n *int) {
	for {
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	builder, _ := gtk.BuilderNew()

	builder.AddFromFile("clock.glade")

	// treeObj, _ := builder.GetObject("tree")
	// tree := treeObj.(*gtk.TreeView)

	storeObj, _ := builder.GetObject("liststore")
	store := storeObj.(*gtk.ListStore)
	addObj, _ := builder.GetObject("add")
	addButton := addObj.(*gtk.Button)
	n := 1
	addButton.Connect("clicked", func() {
		it := store.Append()
		n++
		store.Set(it, []int{0, 1}, []interface{}{1, true})
	})
	winObj, _ := builder.GetObject("main")
	win := winObj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.ShowAll()
	go update(store, &n)
	gtk.Main()
}
