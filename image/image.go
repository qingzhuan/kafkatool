package images

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"kafkatool/common"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// 消防通道占用报警图片路径和图片
	FireEscapeAlarmBasePath = "/mnt/data/test_data/car"
	FireEscapeAlarmJpg      = FireEscapeAlarmBasePath + "/image_1643158302.308963776.jpg"
	FireEscapeAlarmPgm      = FireEscapeAlarmBasePath + "/image_1643158302.308963776.pgm"
)

var ImageQueue = make(chan DetectFile, 100)

type DetectFile struct {
	PGM string
	JPG string
}

func GetImageList(path string) (fileList []string) {
	if path == "" {
		path = "."
	}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("GetImageList err:", err)
		return
	}
	for _, item := range dir {
		if !item.IsDir() && strings.HasSuffix(item.Name(), ".jpg") {
			if path == "." {
				pwd, _ := os.Getwd()
				path = pwd + "/"

			}
			fileList = append(fileList, strings.Join([]string{path, item.Name()}, ""))
		}
	}
	return
}

func GetRandImage(path string) (detectFile DetectFile){
	imageList := GetImageList(path)
	length := len(imageList)
	if length < 1 {
		log.Printf("从%s获取图片异常\n", path)
		return
	}
	rand.Seed(time.Now().UnixNano())
	intn := rand.Intn(length)
	detectFile.PGM = imageList[intn] + ".pgm"
	detectFile.JPG = imageList[intn]
	return
}

func GetImageWidthAndHeight(name string) (width, height int) {
	width, height = 1920, 1080
	file, err := os.Open(name)

	if err != nil {
		log.Println("open err", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	config, _, err := image.DecodeConfig(file)
	log.Println("DecodeConfig err", err)
	height = config.Height
	width = config.Width
	return
}

func ProduceJpgImage() {
	path := FireEscapeAlarmBasePath
	// 定时删除生成的文件
	go ClearImage(path + "/jpg")

	for {
		dstName := path + "/jpg/" + common.MD5Value(strconv.Itoa(time.Now().Nanosecond())) + ".jpg"
		readFile, err := ioutil.ReadFile(FireEscapeAlarmJpg) // 能产生报警的图片
		if err != nil {
			log.Println("读取文件失败，err：", err)
			continue
		}
		err = ioutil.WriteFile(dstName, readFile, 0644)
		if err != nil {
			log.Println("复制文件失败，err：", err)
			continue
		}
		var detectFile DetectFile
		detectFile.PGM = FireEscapeAlarmPgm
		detectFile.JPG = dstName

		ImageQueue <- detectFile
	}
}

func ClearImage(path string) {

	for {
		infos, err := ioutil.ReadDir(path)
		if err != nil {
			log.Println("ClearImage", err)
			return
		}
		if len(infos) < 50 {
			return
		}

		for _, info := range infos {
			if info.ModTime().Add(time.Minute * 10).Unix() < time.Now().Unix() && strings.HasSuffix(info.Name(), ".jpg") {
				fileName := path + "/" + info.Name()
				err := os.Remove(fileName)
				if err != nil {
					log.Println("remove file err:", err)
					continue
				}
				log.Println("删除成功", fileName)

			}
		}
		time.Sleep(10 * time.Minute)

	}
}
