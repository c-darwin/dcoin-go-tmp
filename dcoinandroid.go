// +build android

package main

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/dcoin"
	"image"
	"log"
	"time"
	_ "image/png"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/config"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"fmt"
)

/*
#include <stdio.h>
#include <stdlib.h>

char* JGetTmpDir2() {
	return getenv("TMPDIR");
}
*/
import "C"


var (
	startTime = time.Now()
	eng       = glsprite.Engine()
	scene     *sprite.Node
	cfg config.Event
)

func main() {

	//dir:= C.GoString(C.JGetTmpDir2());
	dir := C.GoString(C.getenv(C.CString("FILESDIR")))
	fmt.Println("dir111()>::", dir)
	fmt.Println("dir122()>::", C.GoString(C.getenv(C.CString("FILESDIR"))))
	fmt.Println("dir122()>::", C.GoString(C.getenv(C.CString("TMPDIR"))))

	//go func(dir string) {
	//  dcoin.Start(dir)
	//}(dir)
	go dcoin.Start(dir)

	app.Main(func(a app.App) {

		for e := range a.Events() {
			fmt.Println("e:", e)
			switch e := app.Filter(e).(type) {
				case config.Event:
				cfg = e
				case paint.Event:
				onPaint(cfg)
				a.EndPaint(e)
			}
		}
	})
}

func onPaint(c config.Event) {
	loadScene()
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, c)
}

func newNode() *sprite.Node {
	n := &sprite.Node{}
	eng.Register(n)
	scene.AppendChild(n)
	return n
}

func loadScene() {
	texs := loadTextures()
	scene = &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	//var n *sprite.Node
	new_w := float32(cfg.WidthPt)
	new_h := new_w*1.77
	if float32(cfg.WidthPt)/float32(cfg.HeightPt) > 1 {
		new_w = float32(cfg.WidthPt)
		new_h = float32(cfg.WidthPt)*0.5625

	}
	n := newNode()
	eng.SetSubTex(n, texs)
	eng.SetTransform(n, f32.Affine{
		{new_w, 0, 0},
		{0, new_h, 0},
	})
}

const (
	texBooks = iota
	texFire
	texGopherR
	texGopherL
)

func loadTextures() sprite.SubTex {
	imgPath := "mobile.png"
	w := 1080
	h := 1920
	if float32(cfg.WidthPt)/float32(cfg.HeightPt) > 1 {
		imgPath = "mobile-landscape.png"
		w = 1920
		h = 1080
	}
	a, err := asset.Open(imgPath)

	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	img, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(img)
	if err != nil {
		log.Fatal(err)
	}

	return sprite.SubTex{t, image.Rect(0, 0, w, h)}
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
