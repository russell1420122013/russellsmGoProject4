package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"time"
)

//go:embed Go_sprites/*
var EmbeddedAssets embed.FS
var textWidget *widget.Text
var rootContainer *widget.Container
var comKillCount int
var game Game

const (
	GameWidth   = 700
	GameHeight  = 700
	PlayerSpeed = 5
)

type Sprite struct {
	pict *ebiten.Image
	xloc int
	yloc int
	dX   int
	dY   int
}

type Game struct {
	player  Sprite
	com     []Sprite
	score   int
	drawOps ebiten.DrawImageOptions
	AppUI   *ebitenui.UI
}

func (g *Game) Update() error {
	g.AppUI.Update()
	processPlayerInput(g)
	i := 0
	for range g.com {
		checkPos(g, i)
		i++
	}
	displayScore(g)
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	g.drawOps.GeoM.Reset()
	g.drawOps.GeoM.Translate(float64(g.player.xloc), float64(g.player.yloc))
	screen.DrawImage(g.player.pict, &g.drawOps)
	for _, comSprite := range g.com {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Translate(float64(comSprite.xloc), float64(comSprite.yloc))
		screen.DrawImage(comSprite.pict, &g.drawOps)
	}
	g.AppUI.Draw(screen)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}
func main() {
	rand.NewSource(time.Now().UnixNano())
	println(rand.Intn(600))
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("project 3")
	game = Game{score: 0}
	game.player = Sprite{
		pict: loadPNGImageFromEmbedded("player.png"),
		xloc: 200,
		yloc: 300,
		dX:   0,
		dY:   0,
	}
	for i := 0; i < 10; i++ {
		game.com = append(game.com, Sprite{
			pict: loadPNGImageFromEmbedded("com.png"),
			xloc: rand.Intn(600),
			yloc: rand.Intn(600),
			dX:   0,
			dY:   0,
		})
	}
	game.AppUI = MakeUIWindow()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
}
func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	rootContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
	)
	textInfo := widget.TextOptions{}.Text("Score: 0", basicfont.Face7x13, color.White)
	idle, err := loadImageNineSlice("button_idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("button_hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("button_pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("button_disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Press for more to come", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),

		widget.ButtonOpts.ClickedHandler(makeMoreCom),
	)
	rootContainer.AddChild(button)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)
	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("Go_sprites")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("Go_sprites/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}
func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}
func processPlayerInput(theGame *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		theGame.player.dY = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		theGame.player.dY = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		theGame.player.dY = 0
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		theGame.player.dX = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		theGame.player.dX = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		theGame.player.dX = 0
	}
	theGame.player.yloc += theGame.player.dY
	if theGame.player.yloc <= 0 {
		theGame.player.dY = 0
		theGame.player.yloc = 0
	} else if theGame.player.yloc+theGame.player.pict.Bounds().Size().Y > GameHeight {
		theGame.player.dY = 0
		theGame.player.yloc = GameHeight - theGame.player.pict.Bounds().Size().Y
	}
	theGame.player.xloc += theGame.player.dX
	if theGame.player.xloc <= 0 {
		theGame.player.dX = 0
		theGame.player.xloc = 0
	} else if theGame.player.xloc+theGame.player.pict.Bounds().Size().X > GameWidth {
		theGame.player.dX = 0
		theGame.player.xloc = GameWidth - theGame.player.pict.Bounds().Size().X
	}

}
func checkPos(theGame *Game, comNum int) {
	if theGame.player.xloc >= theGame.com[comNum].xloc &&
		theGame.player.yloc >= theGame.com[comNum].yloc &&
		theGame.player.xloc <= theGame.com[comNum].pict.Bounds().Size().X+theGame.com[comNum].xloc &&
		theGame.player.yloc <= theGame.com[comNum].pict.Bounds().Size().Y+theGame.com[comNum].yloc {
		for theGame.com[comNum].yloc <= GameHeight {
			theGame.com[comNum].yloc += 1
		}
		theGame.score += 100
		comKillCount += 1
	} else if theGame.player.xloc+theGame.player.pict.Bounds().Size().X <=
		theGame.com[comNum].pict.Bounds().Size().X+theGame.com[comNum].xloc &&
		theGame.player.yloc+theGame.player.pict.Bounds().Size().Y <=
			theGame.com[comNum].pict.Bounds().Size().Y+theGame.com[comNum].yloc &&
		theGame.player.xloc+theGame.player.pict.Bounds().Size().X >=
			theGame.com[comNum].xloc &&
		theGame.player.yloc+theGame.player.pict.Bounds().Size().Y >=
			theGame.com[comNum].yloc {
		for theGame.com[comNum].yloc <= GameHeight {
			theGame.com[comNum].yloc += 100
		}
		theGame.score += 100
		comKillCount += 1
	}
}
func displayScore(theGame *Game) {
	message := fmt.Sprintf("Score: %d", theGame.score)
	textWidget.Label = message
	fmt.Println(theGame.score)
}
func makeMoreCom(args *widget.ButtonClickedEventArgs) {
	for i := 0; i < 10; i++ {
		game.com = append(game.com, Sprite{
			pict: loadPNGImageFromEmbedded("com.png"),
			xloc: rand.Intn(600),
			yloc: rand.Intn(600),
			dX:   0,
			dY:   0,
		})
	}
}
