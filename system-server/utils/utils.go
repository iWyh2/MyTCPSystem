package utils

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// TimeStr 返回当前时间
func TimeStr() string {
	t := time.Now()
	return t.Format("06/01/02 15:04")
}

// Clear 清屏函数
func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// Banner 打印banner
func Banner() {
	fileData := []byte{
		32, 95, 95, 95, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 95, 95,
		32, 32, 32, 32, 95, 32, 32, 32, 32, 32, 95, 95, 95, 32, 32, 13,
		10, 40, 95, 41, 32, 92, 32, 32, 32, 32, 32, 32, 32, 32, 47, 32,
		47, 32, 32, 32, 124, 32, 124, 32, 32, 32, 124, 95, 95, 32, 92,
		32, 13, 10, 32, 95, 32, 92, 32, 92, 32, 32, 47, 92, 32, 32, 47,
		32, 47, 32, 32, 32, 95, 124, 32, 124, 95, 95, 32, 32, 32, 32,
		41, 32, 124, 13, 10, 124, 32, 124, 32, 92, 32, 92, 47, 32, 32,
		92, 47, 32, 47, 32, 124, 32, 124, 32, 124, 32, 39, 95, 32, 92,
		32, 32, 47, 32, 47, 32, 13, 10, 124, 32, 124, 32, 32, 92, 32,
		32, 47, 92, 32, 32, 47, 124, 32, 124, 95, 124, 32, 124, 32, 124,
		32, 124, 32, 124, 47, 32, 47, 95, 32, 13, 10, 124, 95, 124, 32,
		32, 32, 92, 47, 32, 32, 92, 47, 32, 32, 92, 95, 95, 44, 32, 124,
		95, 124, 32, 124, 95, 124, 95, 95, 95, 95, 124, 13, 10, 32, 32,
		32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 95, 95, 47,
		32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 13, 10, 32,
		32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 95, 95,
		95, 47, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 13, 10,
	}
	fmt.Println(string(fileData))
}

// BannerByFile 打印banner
func BannerByFile(filename string) {
	fileData, _ := os.ReadFile(filename)
	fmt.Println(string(fileData))
}

// PrintProgressBar 打印进度条
func PrintProgressBar() {
	total := 100 // Total number of steps
	for i := 0; i <= total; i++ {
		width := 50 // Width of the progress bar
		progress := float64(i) / float64(total)
		bar := int(progress * float64(width))

		fmt.Printf("\r[")
		for i := 0; i < width; i++ {
			if i < bar {
				fmt.Print("=")
			} else if i == bar {
				fmt.Print(">")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Printf("] %d%%", int(progress*100))
		time.Sleep(50 * time.Millisecond) // Simulate work by sleeping
	}
}
