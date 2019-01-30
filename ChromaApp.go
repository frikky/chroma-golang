package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ChromaColor struct {
	Red   int64
	Green int64
	Blue  int64
	Hex   uint32
}

// TODO - Add hex value so you don't need checkRGB() conversion
var (
	RED     ChromaColor = ChromaColor{Red: 255, Green: 0, Blue: 0}
	LIME    ChromaColor = ChromaColor{Red: 0, Green: 255, Blue: 0}
	BLUE    ChromaColor = ChromaColor{Red: 0, Green: 0, Blue: 255}
	BLACK   ChromaColor = ChromaColor{Red: 0, Green: 0, Blue: 0}
	WHITE   ChromaColor = ChromaColor{Red: 255, Green: 255, Blue: 255}
	YELLOW  ChromaColor = ChromaColor{Red: 255, Green: 255, Blue: 0}
	CYAN    ChromaColor = ChromaColor{Red: 0, Green: 255, Blue: 255}
	MAGENTA ChromaColor = ChromaColor{Red: 255, Green: 0, Blue: 255}
	SILVER  ChromaColor = ChromaColor{Red: 192, Green: 192, Blue: 192}
	GREY    ChromaColor = ChromaColor{Red: 128, Green: 128, Blue: 128}
	MAROON  ChromaColor = ChromaColor{Red: 128, Green: 0, Blue: 0}
	OLIVE   ChromaColor = ChromaColor{Red: 128, Green: 128, Blue: 0}
	GREEN   ChromaColor = ChromaColor{Red: 0, Green: 128, Blue: 0}
	PURPLE  ChromaColor = ChromaColor{Red: 128, Green: 0, Blue: 128}
	TEAL    ChromaColor = ChromaColor{Red: 0, Green: 128, Blue: 128}
	NAVY    ChromaColor = ChromaColor{Red: 0, Green: 0, Blue: 128}
)

func checkHex(curColor ChromaColor) (ChromaColor, error) {
	//newHex, err := strconv.Atoi(curColor.Hex)
	//if err == nil {
	//	curColor.Blue = int64(newHex & 255)
	//	curColor.Green = int64((newHex >> 8) & 255)
	//	curColor.Red = int64((newHex >> 16) & 255)
	//	return curColor, nil
	//}

	//var curHex string
	//if strings.HasPrefix(curColor.Hex, "#") {
	//	curHex = curColor.Hex[1:]
	//} else if strings.HasPrefix(curColor.Hex, "0x") {
	//	curHex = curColor.Hex[2:]
	//} else {
	//	return curColor, errors.New("Value is not hex formatted.")
	//}

	//newHex, err = strconv.Atoi(curHex)
	//if err != nil {
	//	return curColor, err
	//}

	//curColor.Hex = curHex
	//curColor.HexInt = strconv.(newHex, 10)
	curColor.Blue = int64(curColor.Hex & 255)
	curColor.Green = int64((curColor.Hex >> 8) & 255)
	curColor.Red = int64((curColor.Hex >> 16) & 255)

	return curColor, nil

}

func checkRGB(curColor ChromaColor) (ChromaColor, error) {
	if curColor.Red > 255 || curColor.Red < 0 {
		return curColor, errors.New("Red out of range. Must be between 0 and 255.")
	}
	if curColor.Green > 255 || curColor.Green < 0 {
		return curColor, errors.New("Green out of range. Must be between 0 and 255.")
	}
	if curColor.Blue > 255 || curColor.Blue < 0 {
		return curColor, errors.New("Blue out of range. Must be begtween 0 and 255.")
	}

	curstr := fmt.Sprintf("%02x%02x%02x",
		uint8(curColor.Blue),
		uint8(curColor.Green),
		uint8(curColor.Red),
	)

	fmt.Println(curstr)

	n, err := strconv.ParseUint(curstr, 16, 32)
	if err != nil {
		return curColor, err
	}

	curColor.Hex = uint32(n)
	if err != nil {
		return curColor, err
	}

	return curColor, nil
}

type author struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type appInfo struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Author          author   `json:"author"`
	DeviceSupported []string `json:"device_supported"`
	Category        string   `json:"category"`
}

type version struct {
	Core    string `json:"core"`
	Device  string `json:"device"`
	Version string `json:"version"`
}

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

type keyboard struct {
	MaxRow    int64
	MaxColumn int64
	Keys      map[string]string
	Uri       string
	Layoutset bool
	CurLayout [][]int64
}

type chromalink struct {
	MaxLED int64
	Uri    string
}

type mousepad struct {
	MaxLED int64
	Uri    string
}

type headset struct {
	MaxLED int64
	Uri    string
}

type mouse struct {
	MaxRow    int64
	MaxColumn int64
	Uri       string
}

type keypad struct {
	MaxRow    int64
	MaxColumn int64
	Uri       string
}

// Apikey string
type appData struct {
	Uri        string       `json:"uri,omitempty"`
	Client     *http.Client `json:"-,omitempty"`
	Keyboard   keyboard     `json:"-,omitempty"`
	Keypad     keypad       `json:"-,omitempty"`
	Chromalink chromalink   `json:"-,omitempty"`
	Mousepad   mousepad     `json:"-,omitempty"`
	Mouse      mouse        `json:"-,omitempty"`
	Headset    headset      `json:"-,omitempty"`
	SessionID  int          `json:"sessionid,omitempty"`
}

var httpClient *http.Client

// createHTTPClient for connection re-use
func createHTTPClient() appData {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	curApp := appData{
		Client: httpClient,
	}

	return curApp
}

func (app appData) GetVersion() (version, error) {
	// HTTPS source host issue
	//url := "https://localhost:54236/razer/chromasdk"
	url := "http://localhost:54235/razer/chromasdk"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return version{}, err
	}

	resp, err := app.Client.Do(req)
	if err != nil {
		return version{}, err
	}

	// Get body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return version{}, err
	}

	curVersion := version{}
	err = json.Unmarshal(body, &curVersion)
	if err != nil {
		return version{}, err
	}

	return curVersion, nil
}

// Sets appdefinitions
func App(curData appInfo) (appData, error) {
	client := createHTTPClient()
	baseUrl := "http://localhost:54235/razer/chromasdk"

	jsonByte, err := json.Marshal(curData)
	if err != nil {
		return appData{}, err
	}

	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonByte))
	if err != nil {
		return appData{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return appData{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return appData{}, err
	}

	err = json.Unmarshal(body, &client)
	if err != nil {
		return appData{}, err
	}

	// Start subprocesses
	// FIXME - Not sure how to stop these
	go client.Heartbeat()
	client = client.setKeyboard()
	client = client.setKeypad()
	client = client.setChromalink()
	client = client.setMouse()
	client = client.setMousepad()
	client = client.setHeadset()

	return client, nil
}

func (app appData) setHeadset() appData {
	app.Headset.MaxLED = 2
	app.Headset.Uri = fmt.Sprintf("%s/chromalink", app.Uri)
	return app
}

func (app appData) setChromalink() appData {
	app.Chromalink.MaxLED = 6
	app.Chromalink.Uri = fmt.Sprintf("%s/chromalink", app.Uri)
	return app
}

func (app appData) setMousepad() appData {
	app.Mousepad.MaxLED = 4
	app.Mousepad.Uri = fmt.Sprintf("%s/mousepad", app.Uri)
	return app
}

func (app appData) setMouse() appData {
	app.Mouse.MaxRow = 9
	app.Mouse.MaxColumn = 7
	app.Mouse.Uri = fmt.Sprintf("%s/mouse", app.Uri)
	return app
}

func (app appData) setKeypad() appData {
	app.Keypad.MaxRow = 4
	app.Keypad.MaxColumn = 5
	app.Keypad.Uri = fmt.Sprintf("%s/keypad", app.Uri)
	return app
}

func (app appData) setKeyboard() appData {
	app.Keyboard.MaxColumn = 22
	app.Keyboard.MaxRow = 6
	app.Keyboard.Uri = fmt.Sprintf("%s/keyboard", app.Uri)
	// Setting up enums lol
	app.Keyboard.Keys = map[string]string{
		"ESC": "0x0001",
		"F1":  "0x0002",
		"F2":  "0x0003",
		"F3":  "0x0004",
		"F4":  "0x0005",
		"F5":  "0x0006",
		"F6":  "0x0007",
		"F7":  "0x0008",
		"F8":  "0x0009",
		"F9":  "0x000A",
		"F10": "0x000B",
		"F11": "0x000C",
		"F12": "0x000D",
		"A":   "0x0302",
		"B":   "0x0407",
		"C":   "0x0405",
		"D":   "0x0304",
		"E":   "0x0204",
	}
	app.Keyboard.Layoutset = false
	return app
}
func (app mousepad) SetNone() error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_NONE"}`)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func (app headset) SetNone() error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_NONE"}`)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func (app chromalink) SetNone() error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_NONE"}`)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func (app keypad) SetNone() error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_NONE"}`)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app mouse) SetNone() error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_NONE"}`)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app keyboard) SetNone() error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_NONE"}`)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app keypad) SetStatic(curColor int64) error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_STATIC", "param": {"color": %d}}`, curColor)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app headset) SetStatic(curColor int64) error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_STATIC", "param": {"color": %d}}`, curColor)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app mousepad) SetStatic(curColor int64) error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_STATIC", "param": {"color": %d}}`, curColor)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app mouse) SetStatic(curColor int64) error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_STATIC", "param": {"color": %d}}`, curColor)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app chromalink) SetStatic(curColor int64) error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_STATIC", "param": {"color": %d}}`, curColor)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

type curgrid struct {
	Effect string    `json:"effect"`
	Param  [][]int64 `json:"param"`
}

//app.keyboard.keyboardgrid = [][]int64
// type keyboardgrid struct
func (app keyboard) SetCustomGrid(arr [][]int64) (keyboard, error) {
	arrlen := int64(len(arr))

	// Might crash if arr[0] doesn't exist? dunno if that can even happen
	if arrlen != app.MaxRow {
		errorstring := fmt.Sprintf("Invalid length of array. Expecting int64[%d][%d], got int64[%d][x]", app.MaxRow, app.MaxColumn, arrlen)
		return app, errors.New(errorstring)
	}

	arronelen := int64(len(arr[0]))
	if arronelen != app.MaxColumn {
		errorstring := fmt.Sprintf("Invalid length of array. Expecting int64[%d][%d], got int64[%d][%d]", app.MaxRow, app.MaxColumn, arrlen, arronelen)
		return app, errors.New(errorstring)
	}

	// Can make this a session in keyboard setup?
	grid := &curgrid{
		Effect: "CHROMA_CUSTOM",
		Param:  arr,
	}

	gridMarshal, err := json.Marshal(grid)
	if err != nil {
		fmt.Println(err)
		return app, err
	}

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer(gridMarshal))

	if err != nil {
		fmt.Println(err)
		return app, err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return app, err
	}

	app.CurLayout = arr
	app.Layoutset = true
	return app, nil
}

// FIXME - has functionality for any key, but missing e.g. setting a specific key
// Is there an API to check where this key would be located?
// Should take a key with int value
func (app keyboard) SetGrid(color uint32) error {
	arr := make([][]int64, app.MaxRow)
	for i, _ := range arr {
		arr[i] = make([]int64, app.MaxColumn)
		for j, _ := range arr[i] {
			arr[i][j] = int64(color)
		}
	}

	// Can make this a session in keyboard setup?
	grid := &curgrid{
		Effect: "CHROMA_CUSTOM",
		Param:  arr,
	}

	gridMarshal, err := json.Marshal(grid)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(gridMarshal))

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer(gridMarshal))

	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	ret, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(ret, err)
	return nil
}

func (app keyboard) SetStatic(curColor uint32) error {
	jsonString := fmt.Sprintf(`{"effect": "CHROMA_STATIC", "param": {"color": %d}}`, curColor)

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Finds location of a key on the keyboard to be parsed in keysetting
// Basically a shitty hex parser lol
func (app keyboard) GetKeyLocation(curkey string) (int, int, error) {
	val := ""
	curkey = strings.ToUpper(curkey)
	for key, value := range app.Keys {
		if key == curkey {
			val = value
			break
		}
	}

	if val == "" {
		return 0, 0, errors.New(fmt.Sprintf("Key %s doesn't exist.", curkey))
	}

	row, err := strconv.ParseUint(val[2:4], 16, 32)
	if err != nil {
		return 0, 0, err
	}

	column, err := strconv.ParseUint(val[4:6], 16, 32)
	if err != nil {
		return 0, 0, err
	}

	return int(row), int(column), nil
}

func (app keyboard) SetKey(key string, color uint32) (keyboard, error) {
	// NEED TO CHECK WHETHER A KEY IS ALREADY SET :)
	row, column, err := app.GetKeyLocation(key)
	if err != nil {
		return app, err
	}

	arr := make([][]int64, app.MaxRow)
	if app.Layoutset == false {
		for i, _ := range arr {
			arr[i] = make([]int64, app.MaxColumn)
			if i != row {
				continue
			}

			for j, _ := range arr[i] {
				if j != column {
					continue
				}

				// Seems to count from 0 (22 rows), but 1 is first
				arr[i][j] = 0x160F00
			}
		}
	} else {
		// Load layout and stuff
		// ONLY change the one we have
		arr = app.CurLayout
		for i, _ := range arr {
			//arr[i] = make([]int64, app.MaxColumn)
			if i != row {
				continue
			}

			for j, _ := range arr[i] {
				if j != column {
					continue
				}

				// Seems to count from 0 (22 rows), but 1 is first
				arr[i][j] = 0x160F00
			}
		}

		app.CurLayout = arr
	}

	//Can make this a session in keyboard setup?
	grid := &curgrid{
		Effect: "CHROMA_CUSTOM",
		Param:  arr,
	}

	gridMarshal, err := json.Marshal(grid)
	if err != nil {
		fmt.Println(err)
		return app, err
	}

	req, err := http.NewRequest("PUT", app.Uri, bytes.NewBuffer(gridMarshal))

	if err != nil {
		fmt.Println(err)
		return app, err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return app, err
	}

	app.Layoutset = true
	app.CurLayout = arr
	// set next and previous
	fmt.Printf("Set key %s to %d\n", key, color)
	return app, nil
}

// FIXME - create same with getnextkey
//	row, column, err := app.GetKeyLocation(key)
//	if err != nil {
//		return app, err
//	}

// Public?
func (app keyboard) findKey(inrow int, incolumn int) string {
	//val := ""
	//curkey = strings.ToUpper(curkey)
	for key, value := range app.Keys {
		row, err := strconv.ParseUint(value[2:4], 16, 32)
		if err != nil {
			fmt.Println(err)
		}

		column, err := strconv.ParseUint(value[4:6], 16, 32)
		if err != nil {
			fmt.Println(err)
		}

		if row == uint64(inrow) && column == uint64(incolumn) {
			return key
		}
	}

	return ""
}

// Finds based on string
// Dependandt on key being created
func (app keyboard) GetNextKey(key string) string {
	row, column, err := app.GetKeyLocation(key)
	if err != nil {
		return ""
	}

	newrow := 0
	newcolumn := 0
	//fmt.Println(row, column)
	if int64(column) != app.MaxColumn && int64(row) != app.MaxRow {
		// Current row, next column
		newrow, newcolumn = row, column+1
	} else if int64(column) == app.MaxColumn && int64(row) != app.MaxRow {
		// Next row, position 0
		newrow, newcolumn = row+1, column
	} else {
		// First row, first thingy
		newrow = 0
		newcolumn = 0
	}

	newkey := app.findKey(newrow, newcolumn)
	if newkey == "" {
		fmt.Println("empty newkey")
	}

	// Find key based on string
	return newkey
}

func (app keyboard) GetNextKeyPosition(row int64, column int64) (int64, int64) {
	if column != app.MaxColumn && row != app.MaxRow {
		// Current row, next column
		return row, column + 1
	} else if column == app.MaxColumn && row != app.MaxRow {
		// Next row, position 0
		return row + 1, 0
	}

	// First row, first thingy
	return 0, 0
}

func (app appData) Heartbeat() {
	for {
		newUri := fmt.Sprintf("%s/heartbeat", app.Uri)

		req, err := http.NewRequest("PUT", newUri, nil)
		if err != nil {
			fmt.Printf("Error in heartbeat Tick: %s\n", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Client.Do(req)
		if err != nil {
			fmt.Printf("Error in heartbeat Tick (2): %s\n", err)
		}

		defer resp.Body.Close()
		// body, err
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error in heartbeat body conversion (2): %s\n", err)
		}
		//fmt.Println(string(body))

		//fmt.Println("Sleeping for a second")
		time.Sleep(1 * time.Second)
	}
}
