package log4g

import (
	"time"
	"os"
)

var file_detect_ticker = time.NewTicker(time.Minute)
var detectFiles = make(map[string]*file_changed_info)
var running = false

type file_changed_info struct {
	filepath string
	mapping interface{}
	mod_time time.Time
	notifyFunc func(filepath string, mapping interface{}) error
}

func start() {
	go func() {
		for {
			<- file_detect_ticker.C
			for filename,fileinfo := range detectFiles {
				if fi, err := os.Stat(filename); err == nil {
					if !fileinfo.mod_time.Equal(fi.ModTime()) {
						fileinfo.mod_time = fi.ModTime()
						fileinfo.notifyFunc(fileinfo.filepath, fileinfo.mapping)
					}
				}
			}
		}
	}()
}

func AddFileChangedListener(filepath string, mapping interface{}, notifyFunc func(filepath string, mapping interface{}) error) error {
	if !running {
		start()
		running = true
	}
	fileinfo := new (file_changed_info)
	fileinfo.filepath = filepath
	fileinfo.mapping = mapping
	fileinfo.notifyFunc = notifyFunc
	if fi, err := os.Stat(filepath); err == nil {
		fileinfo.mod_time = fi.ModTime()
	} else {
		return err
	}
	detectFiles[filepath] = fileinfo
	return notifyFunc(filepath, mapping)
}

func RemoveFileChangedListener(filepath string) {
	delete(detectFiles, filepath)
}