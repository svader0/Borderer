package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/AllenDang/giu"
	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
	i "github.com/svader0/Image-Kit-Util"
)

var (
	originalImage image.Image = nil
	workingImage  *image.RGBA = nil
	texture       *g.Texture  = nil
	borderSize    int32       = 10
	borderColor   color.RGBA  = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

const (
	maxDisplayWidth  float32 = 600
	maxDisplayHeight float32 = 400
)

// loadImageAction loads the image and stores it without modifying.
func loadImageAction() {
	filename, err := dialog.File().Load()
	if err != nil {
		fmt.Println("Error loading file:", err)
		return
	}

	img, err := i.LoadImage(filename)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	if img == nil {
		fmt.Println("No image data loaded")
		return
	}

	originalImage = img
	updateImage() // Initial update
}

func saveImageAction() {
	if workingImage != nil {
		filename, err := dialog.File().SetStartFile("output.png").Filter("PNG Image", "png").Filter("JPEG Image", "jpg").Filter("GIF Image", "gif").Filter("All Files", "*").Title("Save Image").Save()
		if err != nil {
			fmt.Println("Error saving file", err)
			return
		} else {
			err = i.SaveImage(filename, workingImage)
			if err != nil {
				fmt.Println("Error saving image", err)
				return
			}
		}
	}
}

func updateImage() {
	if originalImage != nil {
		workingImage = i.AddBorder(originalImage, int(borderSize), borderColor)

		if texture != nil {
			texture = nil
		}
		g.NewTextureFromRgba(workingImage, func(t *giu.Texture) {
			texture = t
		})
	}
}

func loop() {
	g.SingleWindow().Layout(
		g.Label("Instructions:\n1. Load an image\n2. Adjust the border size and color\n3. Save the image to your computer."),
		g.Row(
			g.Button("Load Image").OnClick(loadImageAction),
			g.Button("Save Image").OnClick(saveImageAction),
		),
		g.Row(
			g.Label("Border Size:"),
			g.SliderInt(&borderSize, 1, 100).Size(200).OnChange(updateImage),
		),
		g.Row(
			g.Label("Border Color:"),
			g.ColorEdit("Border color", &borderColor).Size(300).OnChange(updateImage),
		),

		g.Custom(func() {
			if texture != nil && workingImage != nil {
				// Calculate scaling factor
				imgWidth := float32(workingImage.Rect.Size().X)
				imgHeight := float32(workingImage.Rect.Size().Y)

				var scaleFactor float32
				if scaleFactor > 1 {
					scaleFactor = 1
				} else {
					scaleFactor = min(maxDisplayWidth/imgWidth, maxDisplayHeight/imgHeight)
				}

				// Display scaled image
				g.Image(texture).Size(imgWidth*scaleFactor, imgHeight*scaleFactor).Build()
			}
		}),
	)
}

func main() {
	wnd := g.NewMasterWindow("Image Borderer", 800, 600, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
