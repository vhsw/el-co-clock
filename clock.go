package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	tts "github.com/hegedustibor/htgo-tts"
)

var (
	store        *gtk.ListStore
	selectedPath *gtk.TreePath
	timeFormat   = "15:04:05"
)

func loadConfig() {
	file, err := os.Open("settings.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		disabled := strings.HasPrefix(line, "#")
		if disabled {
			line = strings.TrimPrefix(line, "#")
		}
		line = strings.TrimSpace(line)
		split := strings.SplitN(line, ": ", 2)
		store.Set(store.Append(),
			[]int{0, 1, 2, 3},
			[]interface{}{0, !disabled, split[0], split[1]})

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func speak(text string) {
	speech := tts.Speech{Folder: "audio", Language: "ru"}
	speech.Speak(text)
}

func update() {
	for {
		time.Sleep(1 * time.Second)
		store.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
			timeVal, _ := store.GetValue(iter, 2)
			gTimeVal, _ := timeVal.GoValue()
			timeSting := gTimeVal.(string)
			if time.Now().Format(timeFormat) == timeSting {
				msgVal, _ := store.GetValue(iter, 3)
				gMsgVal, _ := msgVal.GoValue()
				go func() { speak("Beep-bop. " + timeSting); speak(gMsgVal.(string)) }()
				return true
			}
			return false
		})
	}
}

func onEnableToggled(cell *gtk.CellRendererToggle, pathString string) {
	it, _ := store.GetIterFromString(pathString)
	valueObj, _ := store.GetValue(it, 1)
	value, _ := valueObj.GoValue()
	bValue, _ := value.(bool)
	store.SetValue(it, 1, !bValue)
}

func onTimeEdited(cell *gtk.CellRendererText, pathString string, text string) {
	_, err := time.Parse("15:04:05", text)
	if err != nil {
		println(err.Error())
		return
	}

	it, _ := store.GetIterFromString(pathString)
	store.SetValue(it, 2, text)

}

func onMessageEdited(cell *gtk.CellRendererText, pathString string, text string) {
	it, _ := store.GetIterFromString(pathString)
	store.SetValue(it, 3, text)
	go speak(text)
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
		now := time.Now()
		store.Set(it, []int{0, 1, 2, 3}, []interface{}{1, true, now.Format(timeFormat), "Enter message..."})
	})
	removeObj, _ := builder.GetObject("remove")
	removeButton := removeObj.(*gtk.Button)
	selectObj, _ := builder.GetObject("selection")
	selection := selectObj.(*gtk.TreeSelection)
	selection.Connect("changed", func() {
		model, iter, ok := selection.GetSelected()
		if ok {
			removeButton.SetSensitive(true)
			path, _ := model.(*gtk.TreeModel).GetPath(iter)
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
	go loadConfig()
	win.ShowAll()
	go update()
	gtk.Main()
}
