package main

import (
	"github.com/gh0o0st/uploader/COS"
)

func main() {
	cos := COS.NewClient("AKID6AC8xhmuMiUCoGm0EPWXRF1fBvVvlLem", "GLFxTUJ2ayCPfkhT4ThjZE5lj9Wny3IW", "blog-1255588246.cos.ap-shanghai.myqcloud.com")
	err := cos.GetBucket()
	if err != nil {
		panic(err)
	}
}
