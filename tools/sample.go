package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ariyn/f1/packet"
	"os"
)

func main() {
	path := "/tmp/f1-data/f1-1"
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	buf := bufio.NewReader(file)

	for true {
		var isPrefix = true
		line := []byte{}
		for isPrefix {
			var _line []byte
			_line, isPrefix, err = buf.ReadLine()
			if err != nil {
				panic(err)
			}

			line = append(line, _line...)
		}

		//p := packet.Packet{}
		//err = json.Unmarshal(line, &p)
		//if err != nil {
		//	panic(err)
		//}

		position := p.Data.(map[string]interface{})["CarMotionData"].([]interface{})[p.Header.PlayerCarIndex].(map[string]interface{})["worldPosition"].(map[string]interface{})
		fmt.Printf("[%f, %f, %f],\n", position["X"], position["Y"], position["Z"])
	}
}
