// Example program that uses blakjack/webcam library
// for working with V4L2 devices.
// The application reads frames from device and writes them to stdout
// If your device supports motion formats (e.g. H264 or MJPEG) you can
// use it's output as a video stream.
// Example usage: go run stdout_streamer.go | vlc -
package main

import (
	//"io/ioutil"
	//"io"
	//"encoding/base64"
	"fmt"
	"image"
	"log"
	"os"
	"sort"
	//"strings"
	"bytes"
	_ "image/jpeg"

	"github.com/blackjack/webcam"
)

func readChoice(s string) int {
	var i int
	for true {
		print(s)
		_, err := fmt.Scanf("%d\n", &i)
		if err != nil || i < 1 {
			println("Invalid input. Try again")
		} else {
			break
		}
	}
	return i
}

type FrameSizes []webcam.FrameSize

func (slice FrameSizes) Len() int {
	return len(slice)
}

//For sorting purposes
func (slice FrameSizes) Less(i, j int) bool {
	ls := slice[i].MaxWidth * slice[i].MaxHeight
	rs := slice[j].MaxWidth * slice[j].MaxHeight
	return ls < rs
}

//For sorting purposes
func (slice FrameSizes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func main() {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	format_desc := cam.GetSupportedFormats()
	var formats []webcam.PixelFormat
	for f := range format_desc {
		formats = append(formats, f)
	}

	println("Available formats: ")
	for i, value := range formats {
		fmt.Fprintf(os.Stderr, "[%d] %s\n", i+1, format_desc[value])
	}

	choice := readChoice(fmt.Sprintf("Choose format [1-%d]: ", len(formats)))
	format := formats[choice-1]

	fmt.Fprintf(os.Stderr, "Supported frame sizes for format %s\n", format_desc[format])
	frames := FrameSizes(cam.GetSupportedFrameSizes(format))
	sort.Sort(frames)

	for i, value := range frames {
		fmt.Fprintf(os.Stderr, "[%d] %s\n", i+1, value.GetString())
	}
	choice = readChoice(fmt.Sprintf("Choose format [1-%d]: ", len(frames)))
	size := frames[choice-1]

	f, w, h, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Fprintf(os.Stderr, "Resulting image format: %s (%dx%d)\n", format_desc[f], w, h)
	}

	println("Press Enter to start streaming")
	fmt.Scanf("\n")
	err = cam.StartStreaming()
	if err != nil {
		panic(err.Error())
	}

	timeout := uint32(5) //5 seconds
	//for {
		err = cam.WaitForFrame(timeout)

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			//continue
		default:
			panic(err.Error())
		}

		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			//fmt.Println("len: ", len(frame))
			//d1 := frame
			/*
				err := ioutil.WriteFile("/tmp/foto.jpg", d1, 0644)
				if err != nil {
					panic(err)
				}
			*/
			//fmt.Println("Done")

			//reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
			reader := bytes.NewReader(frame)
			m, _, err := image.Decode(reader)
			if err != nil {
				log.Fatal(err)
			}
			bounds := m.Bounds()

			// Calculate a 16-bin histogram for m's red, green, blue and alpha components.
			//
			// An image's bounds do not necessarily start at (0, 0), so the two loops start
			// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
			// likely to result in better memory access patterns than X first and Y second.
			var histogram [16][4]int
			
			fmt.Println("Rango Y: ",bounds.Min.Y, " - ", bounds.Max.Y )
			fmt.Println("Rango X: ",bounds.Min.X, " - ", bounds.Max.X )
			
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := m.At(x, y).RGBA()
					// A color's RGBA method returns values in the range [0, 65535].
					// Shifting by 12 reduces this to the range [0, 15].
					histogram[r>>12][0]++
					histogram[g>>12][1]++
					histogram[b>>12][2]++
					histogram[a>>12][3]++
				}
			}

			// Print the results.
			fmt.Printf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
			for i, x := range histogram {
				fmt.Printf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
			}

			//os.Stdout.Write(frame)
			//os.Stdout.Sync()
		} else if err != nil {
			panic(err.Error())
		}
	//}
}
