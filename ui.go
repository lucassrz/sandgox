package main

import (
	"bytes"
	"github.com/ebitenui/ebitenui"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"time"
)

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func setupUI(g *Game) {
	res := newResources()
	var buttons = make([]*widget.Button, 0)
	var elements = []buttonData{
		{"Sand", Sand},
		{"Water", Water},
		{"Air", Air},
		{"Metal", Metal},
		{"Black Hole", BlackHole},
		{"Water Generator", WaterGenerator}}

	for _, el := range elements {
		buttons = append(buttons, createButton(g, res, el.label, el.cellType))
	}

	slider := widget.NewSlider(
		widget.SliderOpts.MinMax(0, 20),
		widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.SliderOpts.Images(res.sliderImage, res.buttonImage),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			g.brushSize = args.Current
			isChangingBrush = true
			changingBrushTime = time.Now().Second()
		}),
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
	)

	checkboxShowOnlyUpdated := widget.NewCheckbox(
		widget.CheckboxOpts.ButtonOpts(
			widget.ButtonOpts.WidgetOpts(
				// Set the location of the checkboxShowOnlyUpdated
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				// Set the minimum size of the checkboxShowOnlyUpdated
				widget.WidgetOpts.MinSize(30, 30),
			),
			// Set the background images - idle, hover, pressed
			widget.ButtonOpts.Image(res.buttonImage),
		),
		widget.CheckboxOpts.Image(res.checkboxImage),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			onlyShowUpdatedCells = args.State == widget.WidgetChecked
		}),
	)

	checkboxUpdateAllCells := widget.NewCheckbox(
		widget.CheckboxOpts.ButtonOpts(
			widget.ButtonOpts.WidgetOpts(
				// Set the location of the checkboxShowOnlyUpdated
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				// Set the minimum size of the checkboxShowOnlyUpdated
				widget.WidgetOpts.MinSize(30, 30),
			),
			// Set the background images - idle, hover, pressed
			widget.ButtonOpts.Image(res.buttonImage),
		),
		widget.CheckboxOpts.Image(res.checkboxImage),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			updateAllCells = args.State == widget.WidgetChecked
		}),
	)
	checkboxOnlyOneColor := widget.NewCheckbox(
		widget.CheckboxOpts.ButtonOpts(
			widget.ButtonOpts.WidgetOpts(
				// Set the location of the checkboxShowOnlyUpdated
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				// Set the minimum size of the checkboxShowOnlyUpdated
				widget.WidgetOpts.MinSize(30, 30),
			),
			// Set the background images - idle, hover, pressed
			widget.ButtonOpts.Image(res.buttonImage),
		),
		widget.CheckboxOpts.Image(res.checkboxImage),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			onlyOneColor = args.State == widget.WidgetChecked
		}),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		),
		), widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)))

	for _, button := range buttons {
		buttonContainer.AddChild(button)
	}

	buttonContainer.AddChild(slider)
	buttonContainer.AddChild(checkboxShowOnlyUpdated)
	buttonContainer.AddChild(checkboxUpdateAllCells)
	buttonContainer.AddChild(checkboxOnlyOneColor)

	brushButtonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		),
		), widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)))

	buttonContainer.AddChild(brushButtonContainer)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(
				widget.Insets{
					Left:   screenWidth,
					Right:  0,
					Top:    0,
					Bottom: 0,
				},
			),
		)),
	)
	rootContainer.AddChild(buttonContainer)
	g.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}

func createButton(g *Game, res *resources, label string, cellType CellType) *widget.Button {
	return widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text(label, res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.selectedCellType = cellType
		}),
	)
}

type buttonData struct {
	label    string
	cellType CellType
}

func newResources() *resources {
	idle := image2.NewNineSliceColor(color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff})
	hover := image2.NewNineSliceColor(color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff})
	pressed := image2.NewNineSliceColor(color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff})
	sliderImage := image2.NewNineSliceColor(color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff})
	uncheckedImage := ebiten.NewImage(20, 20)
	uncheckedImage.Fill(color.White)

	checkedImage := ebiten.NewImage(20, 20)
	checkedImage.Fill(color.NRGBA{82, 114, 255, 255})

	greyedImage := ebiten.NewImage(20, 20)
	greyedImage.Fill(color.NRGBA{255, 0, 0, 255})
	font, _ := loadFont(10)
	return &resources{
		buttonImage: &widget.ButtonImage{
			Idle:    idle,
			Hover:   hover,
			Pressed: pressed,
		},
		checkboxImage: &widget.CheckboxGraphicImage{
			Unchecked: &widget.GraphicImage{
				Idle: uncheckedImage,
			},
			Checked: &widget.GraphicImage{
				Idle: checkedImage,
			},
			Greyed: &widget.GraphicImage{
				Idle: greyedImage,
			},
		},
		textColor: &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		},
		sliderImage: &widget.SliderTrackImage{
			Idle:  sliderImage,
			Hover: sliderImage,
		},
		font: font,
		padding: widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		},
	}
}
