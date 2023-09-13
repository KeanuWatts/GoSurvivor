package main

import (
	"Go_Survivor/Entities"
	"Go_Survivor/Util"
	"embed"
	"encoding/gob"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

//go:embed Util/Assets/*
var Assets embed.FS

type gameState struct {
	Center            pixel.Vec
	TimeSinceStart    time.Time
	ProjectileList    *Entities.ProjectileGroup
	EnemyList         []*Entities.EnemyGroup
	Character         Entities.Character
	Win               *pixelgl.Window
	PEPE              *pixel.Sprite
	AttackSpeed       time.Duration
	Bg                Entities.WorldTiles
	Tick              <-chan time.Time
	LastTickTimes     []time.Time
	UpgradesAvailable int
	Restart           bool
	Paused            bool
	SelectedBox       int
}

var Center = pixel.Vec{}

var timeSincStart time.Time
var projectileList *Entities.ProjectileGroup
var enemyList []*Entities.EnemyGroup
var character Entities.Character
var win *pixelgl.Window
var PEPE *pixel.Sprite
var AttackSpeed time.Duration
var bg Entities.WorldTiles
var tick <-chan time.Time
var lastTickTimes []time.Time
var UpgradesAvailable int
var restart bool
var paused bool
var selectedBox int
var volume = 0.5

func run() {
	selectedBox = 0
	paused = false
	restart = false
	timeSincStart = time.Now()
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	Center = pixel.Vec{X: cfg.Bounds.W() / 2, Y: cfg.Bounds.H() / 2}
	win, _ = pixelgl.NewWindow(cfg)
	for !win.Closed() {
		win.Clear(colornames.Black)
		win.Update()

		PEPE = Util.LoadSprite("Util/Assets/PePe.png")

		character = *Entities.NewCharacterBuilder().
			WithSpriteFile("Util/Assets/48930.png").
			WithHP(100).
			WithWin(win).
			WithSpeed(3).
			WithAttackSpeed(1).
			WithScale(0.3).
			Build()

		enemySprite := Util.LoadSprite("Util/Assets/Reddit.png")
		Spiral := Entities.NewEnemyGroupBuilder().
			WithWin(win).
			WithSprite(enemySprite).
			WithHP(100).
			WithDamage(15).
			WithSpeed(5).
			WithScale(0.02).
			WithCharacter(&character).
			WithMaxSpawnDistance(800).
			WithMinSpawnDistance(800).
			WithMaxDistance(2000).
			WithMovementType(*Entities.NewSpiralFollow()).
			WithMaxSpawns(20).
			WithName("Spiral").
			Build()

		Liniar := Entities.NewEnemyGroupBuilder().
			WithWin(win).
			WithSprite(enemySprite).
			WithHP(10).
			WithDamage(15).
			WithSpeed(2).
			WithScale(0.015).
			WithCharacter(&character).
			WithMaxSpawnDistance(1600).
			WithMinSpawnDistance(1200).
			WithMaxDistance(2000).
			WithMovementType(*Entities.NewDirectFollow()).
			WithMaxSpawns(40).
			WithName("Liniar").
			Build()

		Orth := Entities.NewEnemyGroupBuilder().
			WithWin(win).
			WithSprite(enemySprite).
			WithHP(20).
			WithDamage(20).
			WithSpeed(3).
			WithScale(0.025).
			WithCharacter(&character).
			WithMaxSpawnDistance(1600).
			WithMinSpawnDistance(1200).
			WithMaxDistance(2000).
			WithMovementType(*Entities.NewOrthogonalFollow()).
			WithMaxSpawns(40).
			WithName("Orth").
			Build()

		Diag := Entities.NewEnemyGroupBuilder().
			WithWin(win).
			WithSprite(enemySprite).
			WithHP(20).
			WithDamage(20).
			WithSpeed(3).
			WithScale(0.025).
			WithCharacter(&character).
			WithMaxSpawnDistance(1600).
			WithMinSpawnDistance(1200).
			WithMaxDistance(2000).
			WithMovementType(*Entities.NewDiagonalFollow()).
			WithMaxSpawns(40).
			WithName("Diag").
			Build()

		enemyList = []*Entities.EnemyGroup{Liniar, Orth, Diag, Spiral}

		LiniarMovemeent := *Entities.NewMovementTypeBuilder().
			WithName("linear").
			WithSpeeds([]float64{2}).
			Build()
		projectileList = Entities.NewProjectileGroupBuilder().
			WithWin(win).
			WithSpriteFile("Util/Assets/Bullet.png").
			WithScale(0.2).
			WithMovementType(LiniarMovemeent).
			WithDamage(1).
			WithChar(&character).
			WithLifeSpan(2000).
			WithAttackPerSec(0.25).
			WithName("Basic Attack").
			Build()
		bg = *Entities.NewWorldTilesBuilder().
			WithSpriteFile("Util/Assets/Town.png").
			WithWin(win).
			WithCharacter(&character).
			WithCenter(Center).
			Build()
		tick = time.Tick(16 * time.Millisecond)
		lastTickTimes = make([]time.Time, 10)
		for i := 0; i < 10; i++ {
			lastTickTimes[i] = time.Now()
		}
		UpgradesAvailable = 0
		character.AddEnemyGroup(*enemyList[0])
		character.AddProjectileGroup(*projectileList)
		LastScore := 0

		for !win.Closed() && !restart {
			if LastScore != character.GetScore() {
				score := character.GetScore()
				LastScore = score
				if isTriangularNumber(score) {
					UpgradesAvailable++
				}

				if score == 50 {
					//remove all enemies with the name "Liniar"
					character.MarkEnemiesForDeathWithName("Liniar")
					//add the Orth enemy group
					character.AddEnemyGroup(*enemyList[1])
				}
				if score == 100 {
					//remove all enemies with the name "Orth"
					character.MarkEnemiesForDeathWithName("Orth")
					//add the Diag enemy group
					character.AddEnemyGroup(*enemyList[2])
				}
				if score == 150 {
					//add back the Orth enemy group
					character.AddEnemyGroup(*enemyList[1])
				}
				if score == 200 {
					//remove all enemies with the name "Orth"
					character.MarkEnemiesForDeathWithName("Orth")
					//remove all enemies with the name "Diag"
					character.MarkEnemiesForDeathWithName("Diag")
					//add the Spiral enemy group
					character.AddEnemyGroup(*enemyList[3])
				}
				if score == 250 {
					//add back linear
					character.AddEnemyGroup(*enemyList[0])
				}
				if score == 300 {
					//add back orth and diag
					character.AddEnemyGroup(*enemyList[1])
					character.AddEnemyGroup(*enemyList[2])
				}
				//for every 50 points over 300 add another random enemy group
				if score > 300 && score%50 == 0 {
					selection := rand.Intn(len(enemyList))
					character.AddEnemyGroup(*enemyList[selection])
				}
			}
			if UpgradesAvailable > 0 {
				DrawUpgradeScreen()
			} else if paused {
				DrawPauseScreen()
			} else if character.GetHP() > 0 {
				MainGameLoop()
			} else {
				DrawGameoverTextWithScore(win, character.GetScore())
			}
			win.Update()
			<-tick // Wait until next tick
		}
		restart = false
	}
}

func isTriangularNumber(n int) bool {
	if n < 0 {
		return false
	}
	prev := 0
	next := 1
	for {
		if n == prev {
			// fmt.Println("Truee")
			return true
		}
		if n < prev {
			// fmt.Println("False")
			return false
		}
		temp := next
		next = prev + next
		prev = temp
		// fmt.Println("prev: ", prev, " n: ", n, " next: ", next)
	}
}

func DrawPauseScreen() {

	bg.Draw()
	character.Draw()
	color := colornames.Black
	color.A = 100
	Util.DrawHollowRect(win, Center, pixel.Vec{X: win.Bounds().W(), Y: win.Bounds().H()}, color, 1, true)
	Util.DrawText(win, "Press Enter to Unpause", pixel.Vec{X: Center.X, Y: Center.Y}, 5)
	if win.JustPressed(pixelgl.KeyEnter) {
		paused = false
		projectileList.ProlongLifeSpan()
	}
}

func GenerateUpgradeOptions() []string {
	//for now just return the same 3 options
	return []string{"Health", "Attack Speed", "Damage"}
}

func MainGameLoop() {
	startTime := time.Now()
	YMovement := 0.0
	XMovement := 0.0
	if win.JustPressed(pixelgl.KeyEscape) {
		paused = true
	}
	if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
		YMovement += 1
	}
	if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
		YMovement -= 1
	}
	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
		XMovement -= 1
	}
	if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
		XMovement += 1
	}

	movVec := pixel.Vec{X: XMovement, Y: YMovement}
	// if movVec.Len() > 0 {
	// 	character.Move(movVec.Unit())
	// }
	character.Move(movVec)

	win.Clear(colornames.Skyblue)
	bg.Draw()
	character.Draw()
	if character.HandleCollisions() {
		win.Clear(colornames.Red)

		rand.Seed(time.Now().UnixNano())
		if rand.Float64()*100 <= 1 {
			PEPE.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		}
	}

	character.IncrementScore(character.RemoveDead())
	character.GamestateUpdate()
	// enemyList.IncreaseDifficulty(score)
	DrawFPSAndTime(win, lastTickTimes, timeSincStart, character.GetScore())
	lastTickTimes = append([]time.Time{startTime}, lastTickTimes[:len(lastTickTimes)-1]...)
}

func DrawGameoverTextWithScore(win *pixelgl.Window, score int) {

	win.Clear(colornames.Red)
	Util.DrawText(win, "Game Over", pixel.Vec{X: Center.X - 50, Y: Center.Y}, 5)
	Util.DrawText(win, fmt.Sprintf("Score: %d", score), pixel.Vec{X: Center.X - 50, Y: Center.Y - 20}, 3)
	if win.JustPressed(pixelgl.KeyEnter) {
		restart = true
	}
}

var firstLoop = true
var UpgradeOptions []*Entities.Upgrade

func DrawUpgradeScreen() {
	win.Clear(colornames.Black)
	bg.Draw()
	character.Draw()
	if win.JustPressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyW) {
		go playAudioFile("Util/Assets/menu-selection-102220.mp3")
		selectedBox--
		if selectedBox < 0 {
			selectedBox = 2
		}
	}
	if win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeyS) {
		go playAudioFile("Util/Assets/menu-selection-102220.mp3")
		selectedBox++
		if selectedBox > 2 {
			selectedBox = 0
		}
	}
	if firstLoop {
		UpgradeOptions = character.GetAvailableUpgrades(3)
		firstLoop = false
	}
	//get a transparent grey color
	color := colornames.Black
	color.A = 100
	Util.DrawHollowRect(win, Center, pixel.Vec{X: win.Bounds().W(), Y: win.Bounds().H()}, color, 1, true)
	//split the screen into 4 Verticle sections
	//top section is just text, other 3 are boxes
	WindowHeight := win.Bounds().H()
	WindowWidth := win.Bounds().W()
	//draw the top section so it takes up the top 1/4 of the screen	and is centered
	TextPos := pixel.Vec{X: WindowWidth / 2, Y: WindowHeight * 7 / 8}
	Util.DrawText(win, "Upgrades Available", TextPos, 5)
	//draw the 3 boxes
	//draw the first box
	BoxSize := pixel.Vec{X: WindowWidth - 75, Y: WindowHeight/4 - 75}
	SelectedBoxSize := pixel.Vec{X: WindowWidth - 25, Y: WindowHeight/4 - 25}
	BoxPos := pixel.Vec{X: WindowWidth / 2, Y: WindowHeight * 5 / 8}
	for i := 0; i < len(UpgradeOptions); i++ {
		if i == selectedBox {
			Util.DrawTextInBox(win, (UpgradeOptions)[i].GetName(), BoxPos, SelectedBoxSize, colornames.White, 5, 5)
		} else {
			Util.DrawTextInBox(win, (UpgradeOptions)[i].GetName(), BoxPos, BoxSize, colornames.White, 4, 4)
		}
		BoxPos = BoxPos.Sub(pixel.Vec{X: 0, Y: WindowHeight / 4})
	}

	if win.JustPressed(pixelgl.KeyEnter) {
		UpgradesAvailable--
		character.AddUpgrade((UpgradeOptions)[selectedBox])
		(UpgradeOptions)[selectedBox].ApplyUpgrade(&character)
		selectedBox = 0
		firstLoop = true
		projectileList.ProlongLifeSpan()
	}

}

func DrawFPSAndTime(win *pixelgl.Window, movingAverageTickTime []time.Time, timeSincStart time.Time, score int) {
	var totalTickTime int64
	for i := 0; i < len(movingAverageTickTime)-1; i++ {
		totalTickTime += movingAverageTickTime[i].Sub(movingAverageTickTime[i+1]).Milliseconds()
	}
	totalTickTime /= int64(len(movingAverageTickTime) - 1)
	if totalTickTime == 0 {
		totalTickTime = 1
	}

	minutes := int(time.Since(timeSincStart).Minutes())
	seconds := int(time.Since(timeSincStart).Seconds()) % 60
	timeDisplay := fmt.Sprintf("%d:%d", minutes, seconds)
	Util.DrawText(win, fmt.Sprintf("FPS: %d", int(1000/totalTickTime)), pixel.Vec{X: win.Bounds().W() - 100, Y: win.Bounds().H() - 20}, 1)
	Util.DrawText(win, fmt.Sprintf("Time: %s", timeDisplay), pixel.Vec{X: win.Bounds().W() - 100, Y: win.Bounds().H() - 40}, 1)
	Util.DrawText(win, fmt.Sprintf("Score: %d", score), pixel.Vec{X: win.Bounds().W() - 100, Y: win.Bounds().H() - 60}, 1)
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:8080", nil))
	// }()
	defer func() {
		if r := recover(); r != nil {
			// fmt.Println("Recovered in f", r)
			// fmt.Scanln()

		}
	}()
	f, err := Assets.Open("Util/Assets/Mili - Between Two Worlds [Limbus Company].mp3")
	if err != nil {
		// log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		// log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	go func() {
		for {
			playAudioFile("Util/Assets/Mili - Between Two Worlds [Limbus Company].mp3")
		}
	}()
	pixelgl.Run(run)

}

type AudioTask struct {
	FileName string
}

func PlayAudio() {
	f, err := Assets.Open("Util/Assets/Mili - Between Two Worlds [Limbus Company].mp3")
	if err != nil {
		// log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		// log.Fatal(err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		// log.Fatal(err)
	}

	defer streamer.Close()

	// Use a non-blocking select to allow the audio to play concurrently
	for {
		speaker.Play(streamer)
		if streamer.Len() == streamer.Position() {
			return
		}
		select {}
	}
}

func playAudioFile(filePath string) {
	f, err := Assets.Open(filePath)
	if err != nil {
		// log.Fatal(err)
	}

	streamer, _, err := mp3.Decode(f)
	if err != nil {
		// log.Fatal(err)
	}
	defer streamer.Close()
	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   -1.5, //each value of 1 will double the volume(half if negative) ((double is based off the base))
		Silent:   false,
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func SaveGameState(state *gameState) error {
	time := fmt.Sprintf("%v", time.Now().Unix())
	filename := "gameState" + time + ".gob"
	fmt.Println(filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(state)
}

func LoadGameState(filename string) (*gameState, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var state gameState
	err = decoder.Decode(&state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}
