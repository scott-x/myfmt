package main

import (
	"fmt"
	"github.com/scott-x/myfmt/db"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	dir          string
	err          error
	subDirToSkip = ".git"
	files        []string
	goRe         = regexp.MustCompile(`.*\.go$`)
	NUM          int
)

func init() {
	dir, err = os.Getwd()
}

func main() {
	if _, err = os.Stat(path.Join(dir, "go.mod")); err != nil {
		log.Println("please move to golang root directory first...")
		return
	}

	err = walk(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		res := goRe.FindString(file)
		if len(res) > 0 {
			//log.Println(file)
			if db.IsItemExist(file) {
				if db.IsMd5Same(file) {
					continue
				}
				format(file, &NUM)

				err = db.UpdateRecordViaPth(file)
				if err != nil {
					log.Printf("UpdateRecordViaPth error: %s\n", strings.TrimPrefix(file, dir+"/"))
				}
			} else {
				format(file, &NUM)
				err = db.Record(file)
				if err != nil {
					log.Printf("Record error: %s\n", strings.TrimPrefix(file, dir+"/"))
				}
			}
		}
	}

	if NUM == 0 {
		log.Printf("[myfmt]: %d file\n", 0)
	} else if NUM == 1 {
		log.Printf("[myfmt]: %d file\n", 1)
	} else {
		log.Printf("[myfmt]: %d files\n", NUM)
	}

}

func walk(pth string) error {
	return filepath.Walk(pth, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			//fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		//fmt.Printf("visited file or dir: %q\n", path)
		files = append(files, path)
		return nil
	})
}

func format(file string, num *int) {
	(*num)++
	err = exec.Command("go", "fmt", file).Run()
	if err != nil {
		log.Printf("fmt error: %s\n", strings.TrimPrefix(file, dir+"/"))
	}
	log.Println(strings.TrimPrefix(file, dir+"/"))
}