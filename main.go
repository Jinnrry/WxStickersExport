package main

import (
	"errors"
	"fmt"
	"howett.net/plist"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	archives := []string{}
	path := u.HomeDir + "/Library/Containers/com.tencent.xinWeChat/Data/Library/Application Support/com.tencent.xinWeChat/2.0b4.0.9"
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.Name() == "fav.archive" {
			archives = append(archives, path)
		}
		return nil
	})
	if len(archives) == 0 {
		panic("Not Found")
	}

	for group, file := range archives {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		var data map[string]any
		err = plist.NewDecoder(f).Decode(&data)
		if err != nil {
			panic(err)
		}
		for i, item := range data["$objects"].([]any) {
			err = os.MkdirAll(fmt.Sprintf("./imgs/%d", group), 0777)
			if err != nil && !errors.Is(err, os.ErrExist) {
				panic(err)
			}
			str, succ := item.(string)
			if succ {
				url, err := url.ParseRequestURI(str)
				if err == nil {
					fmt.Println(url.String())
					resp, err := http.Get(url.String())
					if err != nil {
						fmt.Println(err)
						continue
					}
					defer resp.Body.Close()
					content, err := io.ReadAll(resp.Body)
					if err != nil {
						fmt.Println(err)
						continue
					}
					contentType := http.DetectContentType(content[:512])
					types := strings.Split(contentType, "/")
					err = os.WriteFile(fmt.Sprintf("./imgs/%d/%d.%s", group, i, types[1]), content, 0644)
					if err != nil {
						panic(err)
					}
				}

			}
		}

	}

}
