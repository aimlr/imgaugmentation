package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/font"

	"github.com/fogleman/gg"
	"github.com/zach-klippenstein/goregen"
)

type Bounds struct {
	Left, Top float64
}

type FixedContent struct {
	Content string
}

type FromRegexContent struct {
	Pattern string
}

type FromFileContent struct {
	Content  string
	FilePath string
}

type Text struct {
	ContentType   string
	Fixed         *FixedContent
	FromFile      *FromFileContent
	FromRegex     *FromRegexContent
	Ignore        bool
	Bounds        `json:"bounds"`
	TextTransform string
}

type Print struct {
	SourceImg     string
	OutFolder     string
	FontPath      string
	FontSize      float64
	OutFilePrefix string
	Rgba          string
	Texts         []Text
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var n = flag.Int("n", 10, "number of images to generate")
	var tFile = flag.String("t", "", "transforms file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()

		log.Fatalf("Cmd should be: %s [options] <path to json config file>\n", os.Args[0])
	}
	jsonPath := flag.Args()[0]

	var p Print
	bs, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v\n", jsonPath, err)
	}

	err = json.Unmarshal(bs, &p)
	if err != nil {
		log.Fatalf("Error unmarshaling data from file %s: %v\n", jsonPath, err)
	}

	//fmt.Printf("%+v", p)
	outputFolderCheck(p.OutFolder)

	img, err := gg.LoadImage(p.SourceImg)
	if err != nil {
		log.Fatalf("Error loading file %s: %v\n", p.SourceImg, err)
	}

	r, g, b, a := colors(p.Rgba)
	var ff font.Face
	if ff, err = gg.LoadFontFace(p.FontPath, p.FontSize); err != nil {
		panic(err)
	}

	if ff == nil {
		log.Fatalf("ff was nil")
	}

	var wg sync.WaitGroup
	for i := 0; i < *n; i++ {

		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			log.Printf("Working on variation %d ... \n", index)

			dc := gg.NewContextForImage(img)
			for _, t := range p.Texts {

				if t.Ignore { // might want to disable some but not delete the json
					continue
				}

				if err := dc.LoadFontFace(p.FontPath, p.FontSize); err != nil {
					panic(err)
				}
				// dc.SetFontFace(ff)

				dc.SetRGBA255(r, g, b, a)
				writeText(dc, t)
			}

			transforms := doTransforms(dc.Image(), tFile)
			for _, t := range transforms {
				fname := p.OutFilePrefix + fmt.Sprintf("-%05d-%s.png", index, t.info)
				outPath := p.OutFolder + string(os.PathSeparator) + fname
				log.Printf("Saving %s\n", fname)
				if err := savePng(outPath, t.img); err != nil {
					log.Fatalf("Error while saving png %s: %v\n", outPath, err)
				}

			}

		}(i)

		//dc.SavePNG()
	}

	wg.Wait()

}

type ImageVariation struct {
	img  image.Image
	info string
}

func savePng(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func outputFolderCheck(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		log.Fatalf("Could not create dir %s: %v\n", path, err)
	}
}

func writeText(dc *gg.Context, t Text) {
	s, err := getContent(t)
	if err != nil {
		log.Fatalf("Error getting content for %+v: %v\n", t, err)
	}
	dc.DrawString(s, t.Left, t.Top)
}

func getContent(t Text) (string, error) {
	var data string
	var err error

	switch strings.ToLower(t.ContentType) {
	case "fixed":
		data = t.Fixed.Content
	case "fromfile":
		if data, err = getRandomStringFromFile(t.FromFile.FilePath); err != nil {
			return "", err
		}
	case "fromregex":
		if data, err = regen.Generate(t.FromRegex.Pattern); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("Invalid content type: %s for text: %+v", t.ContentType, t)
	}

	transformParts := strings.Split(t.TextTransform, ",")
	for _, transform := range transformParts {
		switch strings.TrimSpace(strings.ToLower(transform)) {
		case "":
			break
		case "capitalize":
			data = strings.ToUpper(data)
		default:
			log.Fatalf("Unrecognized text transform %s in text: %+v.\n", transform, t)
		}
	}

	return data, nil

}

func getRandomStringFromFile(f string) (string, error) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		//Do something
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	i := r.Int63n(int64(len(lines)))
	return lines[i], nil
}

func colors(s string) (int, int, int, int) {
	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return 0, 0, 0, 255
	}
	if len(parts) < 4 {
		log.Fatalf(`rgba should be in the form: "int,int,int,int". Example, "225,220,230,255"`)
	}

	r, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("Invalid r value %s: %v\n", parts[0], err)
	}

	g, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatalf("Invalid g value %s: %v\n", parts[1], err)
	}

	b, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Fatalf("Invalid b value %s: %v\n", parts[2], err)
	}

	a, err := strconv.Atoi(parts[3])
	if err != nil {
		log.Fatalf("Invalid a value %s: %v\n", parts[3], err)
	}

	return r, g, b, a
}
