package main

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

var errImg = []byte{'G', 'I', 'F', '8', '7', 'a',
	0x01, 0x00, 0x01, 0x00, 0x80, 0x00, 0x00, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0x21, 0xf9, 0x04,
	0x01, 0x0a, 0x00, 0x01, 0x00, 0x2c, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02,
	0x02, 0x4c, 0x01, 0x00, 0x3b, 0x0a}

func writeErr(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(errImg)
}

func readConfig(path string) (w, h int, ext string) {
	base := filepath.Base(path)
	ext = filepath.Ext(path)
	base = strings.TrimSuffix(base, ext)
	size := strings.Split(base, "x")
	if len(size) >= 2 {
		w, _ = strconv.Atoi(size[0])
		h, _ = strconv.Atoi(size[1])
	}
	if w <= 0 {
		w = 150
	}
	if h <= 0 {
		h = 150
	}
	return
}

func readContent(path string) string {
	dir := filepath.Dir(path)
	content := filepath.Base(dir)
	if content == "/" {
		return ""
	}
	content = fmtUnescape(content)
	return content
}

func fmtUnescape(s string) string {
	us, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return us
}

func qrGen(w http.ResponseWriter, r *http.Request) {
	//init
	width, height, ext := readConfig(r.URL.Path)
	content := readContent(r.URL.Path)

	// generat
	qrcode, err := qr.Encode(content, qr.L, qr.Auto)
	if err != nil {
		writeErr(w)
		return
	}
	qrcode, err = barcode.Scale(qrcode, width, height)
	if err != nil {
		writeErr(w)
		return
	}

	// output
	w.Header().Set("Content-Type", mime.TypeByExtension(ext))

	switch ext {
	case "gif":
		err = gif.Encode(w, qrcode, &gif.Options{NumColors: 2})
	case "jpg":
		err = jpeg.Encode(w, qrcode, &jpeg.Options{Quality: 1})
	default:
		err = png.Encode(w, qrcode)
	}
	if err != nil {
		writeErr(w)
	}
}
