package main

import(
	"log"
	"fmt"
	"math/rand"
	"time"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const RESOLUTION_X_Y = 400
const CHANCE_TO_BE_ALIVE = 5

type Game struct{
	board [][]int
	generation int
}

var game *Game



func main(){
	game = initGame()
	randBoardState(game)
	if err := ebiten.Run(updateLoop, RESOLUTION_X_Y, RESOLUTION_X_Y, 2, "Conway's Game of Life"); err != nil{
		log.Fatal(err)
	}
	//fmt.Println("Hello")

}

func initGame() *Game{
	board := make([][]int, RESOLUTION_X_Y)
	for index, _ := range board{
		board[index] = make([]int, RESOLUTION_X_Y)
	}

	return &Game{board: board, generation: 1}
}

func randBoardState(g *Game){
	rand.Seed(time.Now().UnixNano())
	for x := 0; x < RESOLUTION_X_Y; x++{
		for y := 0; y < RESOLUTION_X_Y; y++{
			if rand.Intn(CHANCE_TO_BE_ALIVE) == 1 {
				g.board[x][y] = 1
			}
		}
	}
}

func updateLoop(screen *ebiten.Image) error{
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft){
		x, y := ebiten.CursorPosition()
		createLivingCels(x, y, game)
	}

	if ebiten.IsDrawingSkipped(){
		return nil
	}

	screen.Fill(color.RGBA{0,0,0,0xff})
	background, _ := ebiten.NewImage(RESOLUTION_X_Y, RESOLUTION_X_Y, ebiten.FilterDefault)
	game = checkRules(game)
	draw(game, background)
	screen.DrawImage(background, &ebiten.DrawImageOptions{})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Generation: %v", game.generation))
	return nil

}

func draw(g *Game, background *ebiten.Image) {
	for x := 0; x < RESOLUTION_X_Y; x++ {
		for y := 0; y < RESOLUTION_X_Y; y++ {
			if g.board[x][y] == 1 {
				ebitenutil.DrawRect(background, float64(x), float64(y), 1, 1, color.White)
			}
		}
	}
}

func createLivingCels(x int, y int, g *Game) *Game{
	x = clamp(x, 0, RESOLUTION_X_Y - 1)
	y = clamp(y, 0, RESOLUTION_X_Y - 1)
	topX, topY := x, clamp(y+1, 0, RESOLUTION_X_Y - 1)
	leftX, leftY := clamp(x-1, 0, RESOLUTION_X_Y - 1), y
	bottomX, bottomY :=x, clamp(y-1, 0, RESOLUTION_X_Y -1)
	rightX, rightY := clamp(x+1, 0, RESOLUTION_X_Y -1), y

	g.board[x][y] = 1
	g.board[topX][topY] = 1
	g.board[leftX][leftY] = 1
	g.board[bottomX][bottomY] = 1
	g.board[rightX][rightY] = 1

	return g
}

func checkRules(g *Game) *Game{
	gameAfterCheck := initGame()
	gameAfterCheck.generation = g.generation + 1

	for x := 0; x < RESOLUTION_X_Y; x++{
		for y := 0; y < RESOLUTION_X_Y; y++{
			neighbors := countNeighbors(x, y, g)
			live := g.board[x][y] == 1
			// Any live cell with fewer than two live neighbors dies, as if by underpopulation
			if live && neighbors < 2 {
				gameAfterCheck.board[x][y] = 0
			}
			// Any live cell with two or three live neighbors lives on to the next generation
			if live && (neighbors == 2 || neighbors == 3) {
				gameAfterCheck.board[x][y] = 1
			}
			// Any live cell with more than three live neighbors dies, as if by overpopulation
			if live && neighbors > 3 {
				gameAfterCheck.board[x][y] = 0
			}
			// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction
			if !live && neighbors == 3 {
				gameAfterCheck.board[x][y] = 1
			}
		}
	}
	return gameAfterCheck
}

func countNeighbors(x int, y int, g *Game) int{
	neighbors := 0
	if y+1 < RESOLUTION_X_Y && g.board[x][y+1] == 1 { // top
		neighbors += 1
	}
	if y+1 < RESOLUTION_X_Y && x+1 < RESOLUTION_X_Y && g.board[x+1][y+1] == 1 { // top right
		neighbors += 1
	}
	if x+1 < RESOLUTION_X_Y && g.board[x+1][y] == 1 { // right
		neighbors += 1
	}
	if x+1 < RESOLUTION_X_Y && y-1 > 0 && g.board[x+1][y-1] == 1 { // bottom right
		neighbors += 1
	}
	if y-1 > 0 && g.board[x][y-1] == 1 { // bottom
		neighbors += 1
	}
	if x-1 > 0 && y-1 > 0 && g.board[x-1][y-1] == 1 { // bottom left
		neighbors += 1
	}
	if x-1 > 0 && g.board[x-1][y] == 1 { // left
		neighbors += 1
	}
	if x-1 > 0 && y+1 < RESOLUTION_X_Y && g.board[x-1][y+1] == 1 { // top left
		neighbors += 1
	}
	return neighbors
}

func clamp(value int, minimum int, maximum int) int{
	if value < minimum{
		return minimum
	}else if(value> maximum){
		return maximum
	}
	return value
}

