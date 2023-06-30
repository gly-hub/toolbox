package captcha

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"testing"
)

func TestAddFontFromBytes(t *testing.T) {
	captcha := New()

	err1 := captcha.AddFontFromBytes([]byte("3fgh"))
	assert.NotNil(t, err1)

	fontData, _ := Asset("default.ttf")
	err2 := captcha.AddFontFromBytes(fontData)
	assert.Nil(t, err2)

	err3 := captcha.SetFont("test.ttf")
	assert.NotNil(t, err3)
}

func TestNew(t *testing.T) {
	captcha := New(Option{
		Size:         &image.Point{X: 128, Y: 64},
		FrontColors:  []color.Color{color.White},
		BkgColors:    []color.Color{color.Black},
		FontPath:     []string{"default.ttf"},
		DisturbLevel: NORMAL,
	})

	err1 := captcha.SetFont("default.ttf")
	assert.Nil(t, err1)
	captcha.SetSize(128, 64)
	captcha.SetDisturbance(MEDIUM)
	captcha.SetFrontColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})

	img1, str := captcha.Create(6, ALL)
	assert.NotEqual(t, str, "")
	assert.NotNil(t, img1)

	img2 := captcha.CreateCustom("3fgh")
	assert.NotNil(t, img2)

	//http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
	//	img, str := cap.Create(6, ALL)
	//	png.Encode(w, img)
	//	println(str)
	//})
	//
	//http.ListenAndServe(":8085", nil)
}
