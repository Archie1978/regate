package drivers

import (
	"fmt"
	"reflect"
)

// Liste of driver remote protocol
var listDriver map[string]DriverRP

type DriverRP interface {
	New() DriverRP
	Process(msg interface{})
	Start(chanelWebSocket chan interface{}, numeroSession int, urlString string)
	Close()

	GetCodeJavascript() (content string)
}

func init() {
	listDriver = make(map[string]DriverRP)
}

func AddDriver(d DriverRP) {
	if reflect.ValueOf(d).Kind() != reflect.Ptr {
		return
	}
	var nameStruct = reflect.TypeOf(d).Elem().Name()
	listDriver[nameStruct] = d
}

func GetDriver(nameDriver string) (DriverRP, error) {
	if d, ok := listDriver[nameDriver]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("Driver not found into %v", ListDriver())
}

func ListDriver() (list []string) {
	list = make([]string, 0)
	for key, _ := range listDriver {
		list = append(list, key)
	}
	return
}

func GetAllCodeJavascript() string {
	ret := ""
	for name, driver := range listDriver {
		ret += "/*** " + name + " ***/\n"
		ret += driver.GetCodeJavascript()
		ret += "\n"
	}
	return ret
}
