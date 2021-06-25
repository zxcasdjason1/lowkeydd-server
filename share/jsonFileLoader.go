package share

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// jsonFileReader
func JSONFileLoader(path string, jsonStruct interface{}) []byte {
	f, err := os.Open(fmt.Sprintf("%s", path))
	if err != nil {
		log.Printf("%s 讀取失敗", path)
		panic(err)
	}
	log.Printf("%s 讀取成功", path)
	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	json.Unmarshal(data, &jsonStruct)
	return data
}
