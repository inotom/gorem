package main

import (
	"flag"
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"html/template"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	// ImageJpeg は出力画像のタイプが jpg
	ImageJpeg = iota
	// ImagePng は出力画像のタイプが png
	ImagePng
	// ImageGif は出力画像のタイプが gif
	ImageGif
)

// ImageType は出力画像タイプの型
type ImageType int

// ImageProp は出力画像のプロパティを保持する
type ImageProp struct {
	Width     int
	Height    int
	TextColor color.RGBA
	FillColor color.RGBA
	Text      string
	FontSize  float64
	HasSize   bool
	Type      ImageType
}

// StringAt はテキストが中心に表示される x, y 座標の値を返す
func (ip *ImageProp) StringAt(gc *draw2dimg.GraphicContext) (float64, float64) {
	left, top, right, bottom := gc.GetStringBounds(ip.Text)
	x := (float64(ip.Width) - (right - left)) / 2
	var fontHeight float64
	if !ip.HasSize {
		fontHeight = (bottom - top) / 4
	}
	y := (float64(ip.Height) / 2) + fontHeight
	return x, y
}

// PropAt は画像サイズ(幅x高さ)のテキストが中心に表示される x, y 座標の値を返す
func (ip *ImageProp) PropAt(text string, gc *draw2dimg.GraphicContext) (float64, float64) {
	left, top, right, bottom := gc.GetStringBounds(text)
	x := (float64(ip.Width) - (right - left)) / 2
	y := (float64(ip.Height) / 2) + ((bottom - top) * 2)
	return x, y
}

// makeImage は ImageProp のプロパティから画像(image.RGBA)を生成して返すl
func makeImage(ip *ImageProp) *image.RGBA {
	// ベース画像を作成
	img := image.NewRGBA(image.Rect(0, 0, ip.Width, ip.Height))
	gc := draw2dimg.NewGraphicContext(img)

	// 背景を描画
	gc.SetFillColor(ip.FillColor)
	draw2dkit.Rectangle(gc, 0, 0, float64(ip.Width), float64(ip.Height))
	gc.Fill()

	// フォント設定
	draw2d.SetFontFolder("./fonts")
	gc.SetFontData(draw2d.FontData{
		Name:   "mplus",
		Family: draw2d.FontFamilyMono,
	})
	// フォント色
	gc.SetFillColor(ip.TextColor)
	// フォントサイズ
	gc.SetFontSize(ip.FontSize)

	// テキストを配置
	x, y := ip.StringAt(gc)
	gc.FillStringAt(ip.Text, x, y)

	// 画像サイズ(幅x高さ)を表示
	if ip.HasSize {
		gc.SetFontSize(ip.FontSize * 0.8)
		propText := fmt.Sprintf("%d x %d", ip.Width, ip.Height)
		x, y := ip.PropAt(propText, gc)
		gc.FillStringAt(propText, x, y)
	}

	gc.Close()

	return img
}

func makeImageHandler(w http.ResponseWriter, r *http.Request) {
	// ImageProp の初期化
	ip := ImageProp{
		Width:     180,
		Height:    120,
		TextColor: color.RGBA{0xff, 0xff, 0xff, 0xff},
		FillColor: color.RGBA{0xcc, 0xcc, 0xcc, 0xff},
		Text:      "",
		FontSize:  14,
		HasSize:   false,
		Type:      ImageJpeg,
	}

	// URI 引数を取得して ImageProp の値を更新する
	query := r.URL.Query()
	// 幅
	if newWidth, err := strconv.Atoi(query.Get("w")); err == nil {
		if newWidth > 0 {
			ip.Width = newWidth
		}
	}
	// 高さ
	if newHeight, err := strconv.Atoi(query.Get("h")); err == nil {
		if newHeight > 0 {
			ip.Height = newHeight
		}
	}
	// フォントサイズ
	if newFs, err := strconv.ParseFloat(query.Get("fs"), 64); err == nil {
		if newFs > 0 {
			ip.FontSize = newFs
		}
	}
	// テキスト
	if newText := query.Get("s"); len(newText) > 0 {
		ip.Text = newText
	}
	// 画像サイズ(幅x高さ)の表示
	ip.HasSize = (query.Get("p") == "1")
	// 画像タイプ
	switch query.Get("t") {
	case "png":
		ip.Type = ImagePng
	case "gif":
		ip.Type = ImageGif
	}

	// ImageProp のプロパティを元に画像を作成し、レスポンスに出力する
	img := makeImage(&ip)
	switch ip.Type {
	case ImagePng:
		png.Encode(w, img)
	case ImageGif:
		gif.Encode(w, img, nil)
	default:
		jpeg.Encode(w, img, nil)
	}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	if err := t.templ.Execute(w, nil); err != nil {
		log.Fatalf("templateHandler.ServeHTTP:", err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解析

	http.Handle("/", &templateHandler{filename: "home.html"})
	http.HandleFunc("/lorem", makeImageHandler)

	// Web サーバを開始
	log.Println("Webサーバを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("ListenAndServe:", err)
	}
}
