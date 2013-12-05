package main

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
)

func main() {
	z, _ := os.Open("C:\\Users\\Public\\Music\\Sample Music\\Kalimba.mp3")
	v := url.Values{}
	v.Add("artist", "Bon Jovi")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, z)
	v.Add("file", string(buf.Bytes()))
	fmt.Printf("%+v", v)
	z.Close()
}
