package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

type Captcha struct {
	frontColors  []color.Color
	bkgColors    []color.Color
	disturbLevel DisturbLevel
	fonts        []*truetype.Font
	size         image.Point
}

type StrType int

const (
	NUM   StrType = iota // 数字
	LOWER                // 小写字母
	UPPER                // 大写字母
	ALL                  // 全部
	CLEAR                // 去除部分易混淆的字符
)

type DisturbLevel int

const (
	NORMAL DisturbLevel = 4
	MEDIUM DisturbLevel = 8
	HIGH   DisturbLevel = 16
)

type Option struct {
	Size         *image.Point
	FrontColors  []color.Color
	BkgColors    []color.Color
	FontPath     []string
	DisturbLevel DisturbLevel
}

var defaultOpt = Option{
	Size:         &image.Point{X: 128, Y: 64},
	FrontColors:  []color.Color{color.Black},
	BkgColors:    []color.Color{color.White},
	FontPath:     []string{"default.ttf"},
	DisturbLevel: NORMAL,
}

func New(opts ...Option) *Captcha {
	var opt = defaultOpt

	if len(opts) > 0 {
		if opts[0].Size != nil && opts[0].Size.X > 64 && opts[0].Size.Y > 32 {
			opt.Size = opts[0].Size
		}
		if len(opts[0].FrontColors) > 0 {
			opt.FrontColors = opts[0].FrontColors
		}

		if len(opts[0].BkgColors) > 0 {
			opt.BkgColors = opts[0].BkgColors
		}

		if len(opts[0].FontPath) > 0 {
			opt.FontPath = opts[0].FontPath
		}

		if opts[0].DisturbLevel != 0 {
			opt.FontPath = opts[0].FontPath
		}
	}

	c := &Captcha{
		disturbLevel: opt.DisturbLevel,
		size:         *opt.Size,
		frontColors:  opt.FrontColors,
		bkgColors:    opt.BkgColors,
	}
	_ = c.SetFont(opt.FontPath...)

	return c
}

// AddFont 添加一个字体
func (c *Captcha) AddFont(path string) error {
	var (
		fontData []byte
		err      error
		font     *truetype.Font
	)
	if path == "default.ttf" {
		fontData, err = Asset("default.ttf")
	} else {
		fontData, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}
	}

	font, err = freetype.ParseFont(fontData)
	if err != nil {
		return err
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

// AddFontFromBytes allows to load font from slice of bytes, for example, load the font packed by https://github.com/jteeuwen/go-bindata
func (c *Captcha) AddFontFromBytes(contents []byte) error {
	font, err := freetype.ParseFont(contents)
	if err != nil {
		return err
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

// SetFont 设置字体 可以设置多个
func (c *Captcha) SetFont(paths ...string) error {
	for _, v := range paths {
		if err := c.AddFont(v); err != nil {
			return err
		}
	}
	return nil
}

func (c *Captcha) SetDisturbance(d DisturbLevel) {
	if d > 0 {
		c.disturbLevel = d
	}
}

func (c *Captcha) SetFrontColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.frontColors = c.frontColors[:0]
		for _, v := range colors {
			c.frontColors = append(c.frontColors, v)
		}
	}
}

func (c *Captcha) SetBkgColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.bkgColors = c.bkgColors[:0]
		for _, v := range colors {
			c.bkgColors = append(c.bkgColors, v)
		}
	}
}

func (c *Captcha) SetSize(w, h int) {
	if w < 48 {
		w = 48
	}
	if h < 20 {
		h = 20
	}
	c.size = image.Point{X: w, Y: h}
}

func (c *Captcha) randFont() *truetype.Font {
	return c.fonts[rand.Intn(len(c.fonts))]
}

// 绘制背景
func (c *Captcha) drawBkg(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//填充主背景色
	bgColorIndex := ra.Intn(len(c.bkgColors))
	bkg := image.NewUniform(c.bkgColors[bgColorIndex])
	img.FillBkg(bkg)
}

// 绘制噪点
func (c *Captcha) drawNoises(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	dLen := int(c.disturbLevel)
	// 绘制干扰斑点
	for i := 0; i < dLen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		r := ra.Intn(size.Y/20) + 1
		colorIndex := ra.Intn(len(c.frontColors))
		img.DrawCircle(x, y, r, i%4 != 0, c.frontColors[colorIndex])
	}

	// 绘制干扰线
	for i := 0; i < dLen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := ra.Intn(size.Y) * o
		h := ra.Intn(size.Y/10) * o
		colorIndex := ra.Intn(len(c.frontColors))
		img.DrawLine(x, y, x+w, y+h, c.frontColors[colorIndex])
		colorIndex++
	}

}

// 绘制文字
func (c *Captcha) drawString(img *Image, str string) {

	if c.fonts == nil {
		panic("没有设置任何字体")
	}
	tmp := NewImage(c.size.X, c.size.Y)

	// 文字大小为图片高度的 0.6
	fSize := int(float64(c.size.Y) * 0.6)
	// 用于生成随机角度
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 文字之间的距离
	// 左右各留文字的1/4大小为内部边距
	padding := fSize / 4
	gap := (c.size.X - padding*2) / (len(str))

	// 逐个绘制文字到图片上
	for i, char := range str {
		// 创建单个文字图片
		// 以文字为尺寸创建正方形的图形
		str := NewImage(fSize, fSize)
		// str.FillBkg(image.NewUniform(color.Black))
		// 随机取一个前景色
		colorIndex := r.Intn(len(c.frontColors))

		//随机取一个字体
		font := c.randFont()
		str.DrawString(font, c.frontColors[colorIndex], string(char), float64(fSize))

		// 转换角度后的文字图形
		rs := str.Rotate(float64(r.Intn(40) - 20))
		// 计算文字位置
		s := rs.Bounds().Size()
		left := i*gap + padding
		top := (c.size.Y - s.Y) / 2
		// 绘制到图片上
		draw.Draw(tmp, image.Rect(left, top, left+s.X, top+s.Y), rs, image.ZP, draw.Over)
	}
	if c.size.Y >= 48 {
		// 高度大于48添加波纹 小于48波纹影响用户识别
		tmp.distortTo(float64(fSize)/10, 200.0)
	}

	draw.Draw(img, tmp.Bounds(), tmp, image.ZP, draw.Over)
}

// Create 生成一个验证码图片
func (c *Captcha) Create(num int, t StrType) (*Image, string) {
	if num <= 0 {
		num = 4
	}
	dst := NewImage(c.size.X, c.size.Y)

	c.drawBkg(dst)
	c.drawNoises(dst)

	str := string(c.randStr(num, int(t)))
	c.drawString(dst, str)

	return dst, str
}

func (c *Captcha) CreateCustom(str string) *Image {
	if len(str) == 0 {
		str = "unknown"
	}
	dst := NewImage(c.size.X, c.size.Y)
	c.drawBkg(dst)
	c.drawNoises(dst)
	c.drawString(dst, str)
	return dst
}

var fontKinds = [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}
var letters = []byte("34578acdefghjkmnpqstwxyABCDEFGHJKMNPQRSVWXY")

// 生成随机字符串
// size 个数 kind 模式
func (c *Captcha) randStr(size int, kind int) []byte {
	iKind, result := kind, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			iKind = rand.Intn(3)
		}
		scope, base := fontKinds[iKind][0], fontKinds[iKind][1]
		result[i] = uint8(base + rand.Intn(scope))
		// 不易混淆字符模式：重新生成字符
		if kind == 4 {
			result[i] = letters[rand.Intn(len(letters))]
		}
	}
	return result
}
