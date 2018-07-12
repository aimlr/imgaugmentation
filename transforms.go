package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/segment"
	"github.com/anthonynsimon/bild/transform"
)

var m = map[string]func(image.Image, json.RawMessage) (image.Image, string){
	"brightness":    brightness,
	"contrast":      contrast,
	"gamma":         gamma,
	"hue":           hue,
	"saturation":    saturation,
	"boxBlur":       boxBlur,
	"gaussianBlur":  gaussianBlur,
	"dilate":        dilate,
	"edgeDetection": edgeDetection,
	"emboss":        emboss,
	"erode":         erode,
	"grayscale":     grayscale,
	"invert":        invert,
	"median":        median,
	"sepia":         sepia,
	"sharpen":       sharpen,
	"sobel":         sobel,
	"unsharpMask":   unsharpMask,
	"threshold":     threshold,
	"cropIn":        cropIn,
	"flipH":         flipH,
	"flipV":         flipV,
	"shearH":        shearH,
	"shearV":        shearV,
	"translate":     translate,
	"resize":        resize,
	"rotate":        rotate,
}

type Variation struct {
	Type    string
	Suffix  *string
	Details json.RawMessage
}

func doTransforms(img image.Image, tFile *string) []ImageVariation {

	m["combine"] = combine

	var imgs = []ImageVariation{
		{
			img,
			"",
		},
	}

	if tFile == nil || *tFile == "" {
		log.Printf("No variations to be done.  Variations file is empty.")
		return imgs
	}

	b, err := ioutil.ReadFile(*tFile)
	if err != nil {
		log.Fatalf("Error reading transforms file %s: %v\n", *tFile, err)
		return nil
	}

	var variations []Variation
	err = json.Unmarshal(b, &variations)
	if err != nil {
		log.Fatalf("Error unmarshaling transforms file %s: %v\n", *tFile, err)
		return nil
	}

	// log.Printf("Variations ... \n")

	for _, v := range variations {

		givenSuffix, givenSuffixOk := getSuffix(v)
		// fmt.Printf("Given suffix: %v, %t\n", givenSuffix, givenSuffixOk)
		fn, ok := m[v.Type]
		if !ok {
			log.Fatalf("Fatal Error! No mapped function for variation type: %s\n", v.Type)
		}

		var suffix string
		newImg, defaultSuffix := fn(img, v.Details)
		if givenSuffixOk {
			suffix = givenSuffix
		} else {
			suffix = defaultSuffix
		}
		imgs = append(imgs, ImageVariation{newImg, suffix})

	}

	return imgs
}

func getObject(msg json.RawMessage, in interface{}) {
	err := json.Unmarshal(msg, &in)
	if err != nil {
		log.Fatalf("Error unmarshaling transforms file %s: %v\n", string(msg), err)
	}
}

type Brightness struct {
	Amount float64
}

func brightness(img image.Image, msg json.RawMessage) (image.Image, string) {

	var v Brightness
	getObject(msg, &v)
	newImg := adjust.Brightness(img, v.Amount)
	suffix := fmt.Sprintf("brightness-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Contrast struct {
	Amount float64
}

func contrast(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Contrast
	getObject(msg, &v)
	newImg := adjust.Contrast(img, v.Amount)
	suffix := fmt.Sprintf("contrast-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Gamma struct {
	Amount float64
}

func gamma(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Gamma
	getObject(msg, &v)
	newImg := adjust.Contrast(img, v.Amount)
	suffix := fmt.Sprintf("contrast-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Hue struct {
	Amount int
}

func hue(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Hue
	getObject(msg, &v)
	newImg := adjust.Hue(img, int(v.Amount))
	suffix := fmt.Sprintf("hue-%d", v.Amount)
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Saturation struct {
	Amount float64
}

func saturation(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Saturation
	getObject(msg, &v)
	newImg := adjust.Saturation(img, v.Amount)
	suffix := fmt.Sprintf("saturation-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type BoxBlur struct {
	Amount float64
}

func boxBlur(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v BoxBlur
	getObject(msg, &v)
	newImg := blur.Box(img, v.Amount)
	suffix := fmt.Sprintf("boxBlur-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type GaussianBlur struct {
	Amount float64
}

func gaussianBlur(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v GaussianBlur
	getObject(msg, &v)
	newImg := blur.Gaussian(img, v.Amount)
	suffix := fmt.Sprintf("gaussianBlur-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Dilate struct {
	Amount float64
}

func dilate(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Dilate
	getObject(msg, &v)
	newImg := effect.Dilate(img, v.Amount)
	suffix := fmt.Sprintf("dilate-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type EdgeDetection struct {
	Amount float64
}

func edgeDetection(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v EdgeDetection
	getObject(msg, &v)
	newImg := effect.EdgeDetection(img, v.Amount)
	suffix := fmt.Sprintf("edgeDetection-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func emboss(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := effect.Emboss(img)
	suffix := fmt.Sprintf("emboss")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Erode struct {
	Amount float64
}

func erode(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Erode
	getObject(msg, &v)
	newImg := effect.Erode(img, v.Amount)
	suffix := fmt.Sprintf("erode-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func grayscale(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := effect.Grayscale(img)
	suffix := fmt.Sprintf("grayscale")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func invert(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := effect.Invert(img)
	suffix := fmt.Sprintf("invert")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Median struct {
	Amount float64
}

func median(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Median
	getObject(msg, &v)
	newImg := effect.Median(img, v.Amount)
	suffix := fmt.Sprintf("median-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func sepia(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := effect.Sepia(img)
	suffix := fmt.Sprintf("sepia")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}
func sharpen(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := effect.Sharpen(img)
	suffix := fmt.Sprintf("sharpen")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func sobel(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := effect.Sobel(img)
	suffix := fmt.Sprintf("sobel")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type UnsharpMask struct {
	Amount  float64
	Amount2 float64
}

func unsharpMask(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v UnsharpMask
	getObject(msg, &v)
	newImg := effect.UnsharpMask(img, v.Amount, v.Amount2)
	suffix := fmt.Sprintf("unsharpMask-%s-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64), strconv.FormatFloat(v.Amount2, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Threshold struct {
	Amount uint8
}

func threshold(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Threshold
	getObject(msg, &v)
	newImg := segment.Threshold(img, v.Amount)
	suffix := fmt.Sprintf("threshold-%d", v.Amount)
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type CropIn struct {
	Top    int
	Left   int
	Bottom int
	Right  int
}

func cropIn(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v CropIn
	getObject(msg, &v)
	rect := img.Bounds()
	rect.Min.X -= v.Left
	rect.Min.Y -= v.Top
	rect.Max.X -= v.Right
	rect.Max.Y -= v.Bottom
	newImg := transform.Crop(img, rect)
	suffix := fmt.Sprintf("cropIn-%d-%d-%d-%d", v.Left, v.Top, v.Bottom, v.Right)
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func flipH(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := transform.FlipH(img)
	suffix := fmt.Sprintf("flipH")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func flipV(img image.Image, msg json.RawMessage) (image.Image, string) {
	newImg := transform.FlipV(img)
	suffix := fmt.Sprintf("flipV")
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Shear struct {
	Amount float64
}

func shearH(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Shear
	getObject(msg, &v)
	newImg := transform.ShearH(img, v.Amount)
	suffix := fmt.Sprintf("shearH-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

func shearV(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Shear
	getObject(msg, &v)
	newImg := transform.ShearV(img, v.Amount)
	suffix := fmt.Sprintf("shearY-%s", strconv.FormatFloat(v.Amount, 'f', -1, 64))
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Translate struct {
	X int
	Y int
}

func translate(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Translate
	getObject(msg, &v)
	newImg := transform.Translate(img, v.X, v.Y)
	suffix := fmt.Sprintf("translate-%d-%d", v.X, v.Y)
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Resize struct {
	PercentX int
	PercentY int
}

func resize(img image.Image, msg json.RawMessage) (image.Image, string) {
	var v Resize
	getObject(msg, &v)
	rect := img.Bounds()
	x := int(rect.Max.X * (1.0 + v.PercentX/100))
	y := int(rect.Max.Y * (1.0 + v.PercentY/100))
	newImg := transform.Resize(img, x, y, transform.NearestNeighbor)
	suffix := fmt.Sprintf("resize-%d-%d", v.PercentX, v.PercentY)
	// log.Printf("\t%s\n", suffix)
	return newImg, suffix
}

type Rotation struct {
	Degrees      float64
	ResizeBounds *bool
	PivotX       *int
	PivotY       *int
}

func rotate(img image.Image, msg json.RawMessage) (image.Image, string) {

	var v Rotation
	getObject(msg, &v)
	var opts transform.RotationOptions
	if v.ResizeBounds != nil && *v.ResizeBounds {
		opts.ResizeBounds = true
	}
	if v.PivotX != nil && v.PivotY != nil {
		opts.Pivot = &image.Point{*v.PivotX, *v.PivotY}
	}
	newImg := transform.Rotate(img, v.Degrees, &opts)
	suffix := fmt.Sprintf("rotate-%d", int(v.Degrees))
	// log.Printf("\t%s\n", suffix)

	return newImg, suffix
}

type Variations []Variation

func combine(img image.Image, msg json.RawMessage) (image.Image, string) {
	var vs Variations
	getObject(msg, &vs)

	for _, v := range vs {
		fn, ok := m[v.Type]
		if !ok {
			log.Fatalf("Fatal Error! No mapped function for variation type: %s\n", v.Type)
		}
		img, _ = fn(img, v.Details)
	}

	return img, "combine"
}

func getSuffix(in interface{}) (string, bool) {
	val := reflect.ValueOf(in).FieldByName("Suffix").Interface().(*string)
	if val == nil {
		return "", false
	}

	// fmt.Printf("Suffix value is: %s\n", *val)
	return *val, true
}
