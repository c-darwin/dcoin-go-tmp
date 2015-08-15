// +build android

package main

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/dcoin"

	"image"
	"log"
	"math"
	"time"

	_ "image/png"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/config"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
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
)

func main() {

	dir:= C.GoString(C.JGetTmpDir2());
	fmt.Println("dir111()>::", dir)

	//go func(dir string) {
	//  dcoin.Start(dir)
	//}(dir)
	go dcoin.Start(dir)

	app.Main(func(a app.App) {
		var c config.Event
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
				case config.Event:
				c = e
				case paint.Event:
				onPaint(c)
				a.EndPaint(e)
			}
		}
	})
}

func onPaint(c config.Event) {
	if scene == nil {
		loadScene()
	}
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, c)
	debug.DrawFPS(c)
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

	var n *sprite.Node

	n = newNode()
	eng.SetSubTex(n, texs[texBooks])
	eng.SetTransform(n, f32.Affine{
		{36, 0, 0},
		{0, 36, 0},
	})

	n = newNode()
	eng.SetSubTex(n, texs[texFire])
	eng.SetTransform(n, f32.Affine{
		{72, 0, 144},
		{0, 72, 144},
	})

	n = newNode()
	n.Arranger = arrangerFunc(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		// TODO: use a tweening library instead of manually arranging.
		t0 := uint32(t) % 120
		if t0 < 60 {
			eng.SetSubTex(n, texs[texGopherR])
		} else {
			eng.SetSubTex(n, texs[texGopherL])
		}

		u := float32(t0) / 120
		u = (1 - f32.Cos(u*2*math.Pi)) / 2

		tx := 18 + u*48
		ty := 36 + u*108
		sx := 36 + u*36
		sy := 36 + u*36
		eng.SetTransform(n, f32.Affine{
			{sx, 0, tx},
			{0, sy, ty},
		})
	})
}

const (
	texBooks = iota
	texFire
	texGopherR
	texGopherL
)

func loadTextures() []sprite.SubTex {
	a, err := asset.Open("logo-big.png")
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

	return []sprite.SubTex{
		texBooks:   sprite.SubTex{t, image.Rect(4, 71, 132, 182)},
		texFire:    sprite.SubTex{t, image.Rect(330, 56, 440, 155)},
		texGopherR: sprite.SubTex{t, image.Rect(152, 10, 152+140, 10+90)},
		texGopherL: sprite.SubTex{t, image.Rect(162, 120, 162+140, 120+90)},
	}
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
