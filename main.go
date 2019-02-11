// Geo project main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/fogleman/gg"

	"github.com/paulmach/go.geojson"
)

var Input string
var Output string
var X int
var Y int
var scale float64
var shiftX float64
var shiftY float64

func init() {
	flag.StringVar(&Input, "i", "./Geo.txt", "Path to Input file")
	flag.StringVar(&Output, "o", "./Png.png", "Path to Output file")
	flag.IntVar(&X, "x", 1366, "Size x")
	flag.IntVar(&Y, "y", 1024, "Size x")
}

func main() {
	flag.Parse()
	b, err := ReadFiles(Input)
	if err != nil {
		fmt.Println(err.Error)
		return
	}
	scale = 5
	shiftX = float64(X / 5)
	shiftY = float64(Y / 5)
	geoFeature, err := geojson.UnmarshalFeatureCollection(b)
	con := gg.NewContext(X, Y)
	con.InvertY()
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
				con.SetRGB255(255, 0, 0)
				con.SetLineWidth(2)
				con.DrawLine(x1+float64(X/3.0), y1+float64(Y/3.0), x2+float64(X/3.0), y2+float64(Y/3.0))
				//con.DrawLine(x1, y1, x2, y2)
				con.Stroke()
				//con.DrawLine(x1, y1, x2, y2)
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
						x := point[0]*scale + shiftX
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
					setPolyParams(con, 1, 0.5, 0, 1, 0.5, 0, 1, 0.2, 0.5)
				}
			}
		}
	}
	err = con.SavePNG(Output)
	if err != nil {
		return
	}
}

func ReadFiles(fnamejs string) (js []byte, err error) {
	js, err = ioutil.ReadFile(fnamejs)
	return
}

func setPolyParams(con *gg.Context, rb float64, gb float64, bb float64, ab float64, rf float64, gf float64, bf float64, af float64, lw float64) {
	con.SetRGBA(rf, gf, bf, af)
	con.FillPreserve()
	con.SetRGBA(rb, gb, bb, ab)
	con.SetLineWidth(lw)
	con.Stroke()
}
