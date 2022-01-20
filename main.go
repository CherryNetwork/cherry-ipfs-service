package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

type File struct {
	Filename string `json:"filename"`
	Cid      string `json:"cid"`
	Size     uint64 `json:"size"`
}

func upload_to_ipfs(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}

	files := []File{}

	rev_files := form.File["files"]
	if len(rev_files) != 0 {
		sh := shell.NewShell(os.Getenv("IPFS_GATEWAY"))

		for _, file := range rev_files {

			opened, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}

			var b bytes.Buffer
			if _, err := io.Copy(&b, opened); err != nil {
				log.Fatal(err)
			}

			cid, err := sh.Add(strings.NewReader(b.String()))
			if err != nil {
				log.Fatal(err)
			}

			filesize, err := sh.FilesStat(c.Request.Context(), "/ipfs/"+cid)
			if err != nil {
				log.Fatal(err)
			}

			files = append(files, File{
				Filename: file.Filename,
				Cid:  cid,
				Size: filesize.Size,
			})
		}

		c.JSON(200, files)
	} else {
		c.JSON(400, &gin.H{
			"error": "Please provide atleast one file",
		})
	}

}

func main() {
	r := gin.Default()

	r.POST("/upload", upload_to_ipfs)

	r.Run()
}
