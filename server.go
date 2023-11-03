package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/huichen/sego"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
)

var (
	segmenter = sego.Segmenter{}
)

func init() {
}

func main() {
	/**
	* 可以载入多个词典文件，文件名用","分隔，排在前面的词典优先载入分词，当一个分词同时出现在多个词典里时，前面载入的优先
	* 词典的格式为（每个分词一行），频率越大，优先级越高，词性说明，见README.MD
	* 分词文本 频率 词性
	*/
	dict := flag.String("dict", "../dict/dictionary.txt,../dict/user.txt", "词典文件")
	bindaddr := flag.String("b", "0.0.0.0:80", "listen port")
	flag.Parse()
	_file, _ := exec.LookPath(os.Args[0])
	_pwd, _ := path.Split(_file)
	os.Chdir(_pwd)
	runtime.GOMAXPROCS(runtime.NumCPU())
	segmenter.LoadDictionary(*dict)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/words", func(c *gin.Context) {
		text := c.DefaultQuery("key", "")
		segments := segmenter.Segment([]byte(text))
		words := sego.SegmentsToSlice(segments, false)
		c.JSON(200, gin.H{
			"message": "pong",
			"words":   words,
		})
	})
	r.GET("/reload", func(c *gin.Context) {
		segmenter.Close()
		segmenter.LoadDictionary(*dict)
		c.JSON(200, gin.H{
			"message": "pong",
			"reload":  true,
		})
	})
	log.Println("notice: bind addr:%v", *bindaddr)
	err := r.Run(*bindaddr) // listen and serve on 0.0.0.0:8080
	log.Println("error: %v", err)
}
