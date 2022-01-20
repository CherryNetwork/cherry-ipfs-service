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

func upload_to_ipfs(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}

	files := form.File["files"]
	if len(files) != 0 {
		sh := shell.NewShell(os.Getenv("IPFS_GATEWAY"))

		var cids []string
		for _, file := range files {

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

			cids = append(cids, cid)
		}

		c.JSON(200, &gin.H{
			"cid": cids,
		})
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
