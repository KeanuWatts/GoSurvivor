package Util

import (
	"embed"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

//go:embed Assets/*

var Assets embed.FS

func LoadSprite(pic string) *pixel.Sprite {
	//reemovee Util/ from the path
	pic = pic[5:]
	file, err := Assets.Open(pic)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	picData := pixel.PictureDataFromImage(img)
	picSprite := pixel.NewSprite(picData, picData.Bounds())
	return picSprite
}

func DrawText(win *pixelgl.Window, str string, pos pixel.Vec, size float64) {
	// Initialize a basic Atlas and Text
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	//place it Center
	//move it to the left by half the length of the string times the size
	pos = pos.Sub(pixel.Vec{X: float64(len(str)) * size * 3.5, Y: 0})
	pos = pos.Sub(pixel.Vec{X: 0, Y: size * 5})
	TextDisplay := text.New(pos, basicAtlas)
	TextDisplay.WriteString(str)
	TextDisplay.Draw(win, pixel.IM.Scaled(pos, size))
}

func DrawTextInBox(win *pixelgl.Window, str string, pos pixel.Vec, size pixel.Vec, Color color.Color, thickness float64, fontSize float64) {
	DrawHollowRect(win, pos, size, Color, thickness, false)
	DrawText(win, str, pos, fontSize)
}

func DrawHollowRect(win *pixelgl.Window, pos pixel.Vec, size pixel.Vec, color color.Color, thickness float64, fill bool) {
	imd := imdraw.New(nil)
	imd.Color = color
	Top := pos.Add(pixel.Vec{X: 0, Y: size.Y / 2}).Y
	Bottom := pos.Add(pixel.Vec{X: 0, Y: -size.Y / 2}).Y
	Left := pos.Add(pixel.Vec{X: -size.X / 2, Y: 0}).X
	Right := pos.Add(pixel.Vec{X: size.X / 2, Y: 0}).X
	imd.Push(pixel.Vec{X: Left, Y: Top},
		pixel.Vec{X: Right, Y: Top},
		pixel.Vec{X: Right, Y: Bottom},
		pixel.Vec{X: Left, Y: Bottom},
		pixel.Vec{X: Left, Y: Top})
	if fill {
		imd.Polygon(0)
	} else {
		imd.Line(thickness)
	}

	imd.Draw(win)

}
func GetRandomEntries(list []interface{}, n int) []interface{} {
	rand.Seed(time.Now().UnixNano())
	result := make([]interface{}, n)
	if len(list) == 0 {
		return result
	}
	selected := make(map[int]bool)
	for i := 0; i < n; i++ {
		index := rand.Intn(len(list))
		for selected[index] {
			index = rand.Intn(len(list))
		}
		selected[index] = true
		result[i] = list[index]
	}
	return result
}

type IntList []int

func (list IntList) Contains(element int) bool {
	for i := range list {
		if list[i] == element {
			return true
		}
	}
	return false
}
