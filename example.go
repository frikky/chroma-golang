package main

import (
	"fmt"
	"log"
	"time"
)

//func driten(t int, f float64) int {
//	hore := float64(math.Sin(float64(2*math.Pi*t/255) + f))
//	fmt.Println(hore)
//	return asd
//}

func main() {
	//for i := 0; i < 255; i++ {
	//	driten(i, 0)
	//}

	// Register an app
	curAuthor := author{
		Name:    "Frikky",
		Contact: "@frikkylikeme",
	}

	curData := appInfo{
		Title:           "Hello",
		Description:     "Golang test",
		Author:          curAuthor,
		DeviceSupported: []string{"keyboard"},
		Category:        "application",
	}

	app, err := App(curData)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	// getnext

	key := "ESC"
	app.Keyboard, err = app.Keyboard.SetKey(key, 0xffffff)
	for {
		key = app.Keyboard.GetNextKey(key)
		app.Keyboard, err = app.Keyboard.SetKey(key, 0xffffff)
		fmt.Println(key)
		time.Sleep(1 * time.Second)
	}

	//key = "ESC"
	//row, column, err := app.GetKeyLocation(key)
	//if err != nil {
	//	fmt.Printf("Break with row %d and column %d\n", row, column)
	//	break
	//}
	//app.Keyboard, err = app.Keyboard.SetKey(key, 0xffffff)
	//for {
	//	row, column = app.Keyboard.GetNextKeyPosition(row, column)
	//	//fmt.Println(row, column)

	//	//app.Keyboard, err = app.Keyboard.SetKey(key, 0xffffff)
	//	////func (app keyboard) GetNextKeyPosition(row int, column int) {
	//	//fmt.Println(key)
	//	time.Sleep(10 * time.Second)
	//}

	// FIXME - Reset keys whenever setnone is set AKA
	//app.Keyboard, err = app.Keyboard.SetKey("ESC", 0xffffff)
	//app.Keyboard, err = app.Keyboard.SetKey("f10", 0xffffff)
	//time.Sleep(2 * time.Second)
	//app.Keyboard, err = app.Keyboard.SetKey("f11", 0xffffff)
	//time.Sleep(2 * time.Second)

	//app.Keyboard, err = app.Keyboard.SetKey("f12", 0xffffff)
	//fmt.Println(err)

	//arr := make([][]int64, 6)
	//for i, _ := range arr {
	//	arr[i] = make([]int64, 22)
	//	for j, _ := range arr[i] {
	//		arr[i][j] = 0xff0000
	//	}
	//}

	//ret := app.Keyboard.SetCustomGrid(arr)
	//fmt.Println(ret)
	//time.Sleep(10 * time.Second)

	//for i := 0; i < 255; i++ {
	//	//color, err := checkRGB(
	//	//	ChromaColor{
	//	//		Red:   int64(0),
	//	//		Green: int64(i),
	//	//		Blue:  int64(0),
	//	//	},
	//	color, err := checkRGB(RED)

	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}

	//	// Can do this to them all?
	//	ret := app.Keyboard.SetStatic(color.Hex)
	//	fmt.Println(ret)
	//	time.Sleep(1 * time.Millisecond)
	//}

	// Custom grid
	// func (app appData) setKeyboard() appData {
	// func (app keyboard) SetStatic(curColor uint32) error {
}
