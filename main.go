package main

import (
	"myProjects/picUploader/COS"
	"os"
)

func main() {
	name := os.Args[1]
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cos := COS.NewCOS("AKID6AC8xhmuMiUCoGm0EPWXRF1fBvVvlLem", "GLFxTUJ2ayCPfkhT4ThjZE5lj9Wny3IW", "blog-1255588246.cos.ap-shanghai.myqcloud.com")
	err = cos.PutObject(f)
	if err != nil {
		panic(err)
	}
}
