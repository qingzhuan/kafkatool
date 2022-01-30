package images

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"kafkatool/common"
	"kafkatool/config"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	// 消防通道占用报警图片路径和图片
	FireEscapeAlarmBasePath = config.Config.FireEscape.FireBasePath
	FireEscapeAlarmJpg      = config.Config.FireEscape.JpgPath
	FireEscapeAlarmPgm      = config.Config.FireEscape.PgmPath
	RandomImagePath 		= config.Config.GroundRandomImagePath
)

var FireEscapeImageQueue = make(chan DetectFile, 10)
var GroundImageQueue = make(chan DetectFile, 500)

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

func GetRandomImage(path string) (detectFile DetectFile) {
	imageList := GetImageList(path)
	lenght := len(imageList)
	if lenght < 1 {
		log.Printf("从%s获取图片异常\n", path)
		return
	}
	rand.Seed(time.Now().UnixNano())
	intn := rand.Intn(lenght)
	detectFile.JPG = imageList[intn]
	detectFile.PGM = imageList[intn] + ".pgm"
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

func ProduceFireEscapeJpgImage() {
	jpgPath := filepath.Join(config.Config.FireEscape.FireBasePath, "jpg")
	if !common.DirIsExist(jpgPath) {
		_ = os.Mkdir(jpgPath, 0644)
	}

	for {
		// 消防通道报警的图片
		for i := 0; i < config.Config.FireEscape.CarImageContinueTime; i++ {
			fileName := common.MD5Value(strconv.Itoa(time.Now().Nanosecond())) + ".jpg"
			dstName := filepath.Join(jpgPath, fileName)
			readFile, err := ioutil.ReadFile(config.Config.FireEscape.JpgPath) // 能产生报警的图片
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
			detectFile.PGM = config.Config.FireEscape.PgmPath
			detectFile.JPG = dstName

			FireEscapeImageQueue <- detectFile
		}
		// 其他随机图片，用于智能巡查掉报警
		for i := 0; i < config.Config.FireEscape.GroundRandomImageContinueTime; i++ {
			detectFile := <- GroundImageQueue
			FireEscapeImageQueue <- detectFile
		}
	}
}

func ProduceGroundJpgImage() {

	for {
		groundJpgImagePath := config.Config.GroundRandomImagePath
		detectFile := GetRandomImage(groundJpgImagePath)
		GroundImageQueue <- detectFile
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
