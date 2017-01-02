package main

// #include <stdlib.h>
import "C"

import (
	"github.com/go-vgo/robotgo"
	"strings"
	"time"
	"encoding/base64"
	"bytes"
	"compress/zlib"
	"unsafe"
	"math"
)

func main() {
	count := 0

	lastRMax := 0
	fishing := true //alternatively, catching
	for {
		if isMinecraft() {
			//fmt.Println("Minecraft!")

			if count == 15 {
				img := robotgo.CaptureScreen(800, 625, 100, 75)
				robotgo.SaveBitmap(img, "test.png")
			}
			//count++

			if fishing {
				data, _ := screenshotAsColorData()

				//count := int64(0)
				//r := int64(0)
				//g := int64(0)
				//b := int64(0)
				rMax := 0
				//gMax := 0
				//bMax := 0
				for i := 0; i < len(data); i += 3 {
					//count++
					//b += int64(data[i])
					//g += int64(data[i+1])
					//r += int64(data[i+2])
					rMax = Max(rMax, int(data[i + 2]))
					//gMax = Max(gMax, int(data[i+1]))
					//bMax = Max(bMax, int(data[i]))
				}
				//ravg := Round(float64(r) / float64(count), .5, 1)
				//gavg := Round(float64(g) / float64(count), .5, 1)
				//bavg := Round(float64(b) / float64(count), .5, 1)
				//fmt.Printf("Averages: (r, g, b) => (%.1f, %.1f, %.1f)\n", ravg, gavg, bavg)
				//fmt.Printf("maxes   : (r, g, b) => (%d, %d, %d)\n", rMax, gMax, bMax)

				//if the percent change is greater than 30%, catch the fish
				delta := lastRMax - rMax
				change := math.Abs(float64(delta)) / float64(lastRMax)
				if change * 100 > 35 {
					rightClick()
					fishing = false
				} else {
					lastRMax = rMax
				}
				//fmt.Printf("Change: %.1f\n", change * 100)
			}
			if fishing {
				time.Sleep(200 * time.Millisecond)
			} else {
				//fmt.Println("Returning to fishing...")
				time.Sleep(200 * time.Millisecond)
				rightClick()
				time.Sleep(3 * time.Second)
				fishing = true
			}
		} else {
			//fmt.Println("Not minecraft")
			time.Sleep(1 * time.Second)
		}
	}
}

func isMinecraft() bool {
	title := robotgo.GetTitle()
	return strings.HasPrefix(strings.ToLower(title), "minecraft ")
}

func rightClick() {
	robotgo.MouseClick("right", false)
}

func screenshotAsColorData() ([]byte, error) {
	imgRef := robotgo.CaptureScreen(800, 625, 100, 75)
	stringBitmap := robotgo.TostringBitmap(imgRef)
	defer C.free(unsafe.Pointer(stringBitmap))
	imageData := C.GoString((*C.char)(stringBitmap))
	decoded, _ := base64.StdEncoding.DecodeString(imageData[9:])
	reader, _ := zlib.NewReader(bytes.NewReader(decoded))
	buffer := bytes.Buffer{}
	buffer.ReadFrom(reader)
	return buffer.Bytes(), nil
}

func Round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
