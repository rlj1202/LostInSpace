package lostinspace

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

type Texture interface {
	Bind(uint32)
	GetSize() (int32, int32)
}

type Texture2D struct {
	width  int32
	height int32
	id     uint32
}

type Texture2DArray struct {
	width  int32
	height int32
	id     uint32
}

func NewTexture2D(imgFile *os.File) *Texture2D {
	tex := new(Texture2D)

	gl.GenTextures(1, &tex.id)
	gl.BindTexture(gl.TEXTURE_2D, tex.id)

	rawBytes, width, height, err := rawBytes(imgFile)
	if err != nil {
		panic(err)
	}
	tex.width = width
	tex.height = height

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rawBytes))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	return tex
}

// Sizes of given images have to be in given width and height.
func NewTexture2DArray(width, height int32, imgFiles []*os.File) *Texture2DArray {
	tex := new(Texture2DArray)
	tex.width = width
	tex.height = height

	gl.GenTextures(1, &tex.id)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, tex.id)

	gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, 1, gl.RGBA8, width, height, int32(len(imgFiles)))
	for i, imgFile := range imgFiles {
		rawBytes, imgWidth, imgHeight, err := rawBytes(imgFile)
		if err != nil {
			panic(err)
		}

		gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, 0, 0, 0, int32(i), imgWidth, imgHeight, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rawBytes))
	}

	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	return tex
}

func (tex *Texture2D) Bind(i uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + i)
	gl.BindTexture(gl.TEXTURE_2D, tex.id)
}

func (tex *Texture2DArray) Bind(i uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + i)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, tex.id)
}

func (tex *Texture2D) GetSize() (int32, int32) {
	return tex.width, tex.height
}

func (tex *Texture2DArray) GetSize() (int32, int32) {
	return tex.width, tex.height
}

func (tex *Texture2D) Destroy() {
	if tex.id != 0 {
		gl.DeleteTextures(1, &tex.id)
		tex.id = 0
	}
}

func (tex *Texture2DArray) Destroy() {
	if tex.id != 0 {
		gl.DeleteTextures(1, &tex.id)
		tex.id = 0
	}
}

// return: (rawBytes, width, height, error)
func rawBytes(imgFile *os.File) ([]uint8, int32, int32, error) {
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, 0, 0, err
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, 0, 0, fmt.Errorf("Unsupported stride.")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba.Pix, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), nil
}
