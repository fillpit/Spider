package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadAll(path string) []*os.File {
	var all_file []*os.File
	finfo, _ := ioutil.ReadDir(path)
	for _ ,x := range finfo {
		real_path := path + "/" + x.Name()
		//fmt.Println(x.Name()," ",x.Size())
		if x.IsDir() {
			fmt.Println(x.Name()," ",x.Size())
			all_file = append(all_file,ReadAll(real_path)...)
		}else {
			file, _ := os.Open(real_path)

			all_file = append(all_file,file)
		}
	}
	return all_file
}
