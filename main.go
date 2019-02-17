// Geo project main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/fogleman/gg"
	"github.com/paulmach/go.geojson"
	"gopkg.in/yaml.v2"
)

type StyleStruct struct {
	Rb        float64 `yaml:"rb"`
	Gb        float64 `yaml:"gb"`
	Bb        float64 `yaml:"bb"`
	Ab        float64 `yaml:"ab"`
	Rf        float64 `yaml:"rf"`
	Gf        float64 `yaml:"gf"`
	Bf        float64 `yaml:"bf"`
	Af        float64 `yaml:"af"`
	LineWidth float64 `yaml:"lw"`
	Layer     int     `yaml:"layer"`
	Zoom      int     `yaml:"zoom"`
	Type      string  `yaml:"type"`
}

type StyleCollection struct {
	styles []StyleStruct
}

var Input string
var Output string
var Style string
var X int
var Y int
var scale float64
var shiftX float64
var shiftY float64

func init() {
	flag.StringVar(&Input, "i", "./Geo.txt", "Path to Input file")
	flag.StringVar(&Output, "o", "./Png.png", "Path to Output file")
	flag.StringVar(&Style, "s", "./Style.yaml", "Path to Style file")
	flag.IntVar(&X, "x", 5120, "Size x")
	flag.IntVar(&Y, "y", 5120, "Size y")
	/*flag.IntVar(&X, "x", 1366, "Size x")
	flag.IntVar(&Y, "y", 1024, "Size y")*/
}

func main() {
	flag.Parse()
	b, s, err := ReadFiles(Input, Style)
	if err != nil {
		fmt.Println("readfiles")
		fmt.Println(err.Error)
		return
	}
	var styles []StyleStruct
	var zoom = 1
	styleMap := make(map[string][]StyleStruct)
	err = yaml.Unmarshal(s, &styles)
	if err != nil {
		fmt.Println("marshal")
		fmt.Println(err.Error)
		return
	}
	for _, s := range styles {
		styleMap[s.Type] = append(styleMap[s.Type], s)
	}
	//fmt.Println(styleMap)
	scale = 25
	shiftX = float64(X / 25)
	shiftY = float64(Y / 5)
	geoFeature, err := geojson.UnmarshalFeatureCollection(b)
	con := gg.NewContext(X, Y)
	con.InvertY()

	myPosition := "Volga Federal District"
	for _, g := range geoFeature.Features {
		//fmt.Println(g.Type)
		switch g.Geometry.Type {
		case "Point":
			con.DrawPoint(g.Geometry.Point[0], g.Geometry.Point[1], 1)
			con.SetRGB255(0, 0, 0)
			con.Fill()
		case "LineString":
			for i := 0; i < len(g.Geometry.LineString)-1; i++ {
				x1 := g.Geometry.LineString[i][0]
				y1 := g.Geometry.LineString[i][1]
				x2 := g.Geometry.LineString[i+1][0]
				y2 := g.Geometry.LineString[i+1][1]
				con.SetRGBA(0, 0, 0, 1)
				con.SetLineWidth(2)
				con.DrawLine(x1+float64(X/3.0), y1+float64(Y/3.0), x2+float64(X/3.0), y2+float64(Y/3.0))
				con.Stroke()
			}

		case "Polygon":
			for _, geom := range g.Geometry.Polygon {
				firstPoint, x0, y0 := true, 0.0, 0.0
				for _, point := range geom {
					x := point[0]
					y := point[1]
					if firstPoint {
						x0 = x*scale + shiftX
						y0 = y*scale + shiftY
						firstPoint = false
					}
					con.LineTo(x, y)

				}
				con.LineTo(x0, y0)

				setPolyParams(con, 1, 0, 0, 1, 0, 0, 0, 0.5, 0.5)
			}
		case "MultiPolygon":
			for _, geom := range g.Geometry.MultiPolygon {
				for _, poly := range geom {
					firstPoint, x0, y0 := true, 0.0, 0.0
					for _, point := range poly {
						x := 0.0
						if point[0] > 0 {
							x = point[0]*scale + shiftX
						} else {
							x = (point[0]+360)*scale + shiftX
						}
						y := point[1]*scale + shiftY
						if firstPoint {
							con.MoveTo(x, y)
							x0 = x
							y0 = y
							firstPoint = false
						}
						con.LineTo(x, y)
					}
					con.LineTo(x0, y0)
					mpolStyle := styleMap["MultiPolygon"]
					var style StyleStruct
					//fmt.Println(g.Properties["name:en"])
					if g.Properties["name:en"] == myPosition {
						for _, s := range mpolStyle {
							if s.Layer == 2 {
								style = s
								break
							}
						}
					} else {
						for _, s := range mpolStyle {
							if s.Zoom == zoom {
								style = s
								break
							}
						}
						zoom = zoom
					}
					setPolyParams(con, style.Rb, style.Gb, style.Bb, style.Ab, style.Rf, style.Gf, style.Bf, style.Af, style.LineWidth)
				}
			}
		}
	}
	err = con.SavePNG(Output)
	if err != nil {
		return
	}
}

func ReadFiles(fnamejs string, fnamest string) (js []byte, st []byte, err error) {
	js, err = ioutil.ReadFile(fnamejs)
	if err != nil {
		return
	}
	st, err = ioutil.ReadFile(fnamest)
	return
}

func setPolyParams(con *gg.Context, rb float64, gb float64, bb float64, ab float64, rf float64, gf float64, bf float64, af float64, lw float64) {
	con.SetRGBA(rf, gf, bf, af)
	con.FillPreserve()
	con.SetRGBA(rb, gb, bb, ab)
	con.SetLineWidth(lw)
	con.Stroke()
}
