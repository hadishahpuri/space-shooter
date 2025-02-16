package main

import (
	"log"
	"fmt"
	"image/color"
	"time"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 400
	screenHeight = 600
	playerSpeed  = 5
)

type Bullet struct {
	x, y float64
}

type Enemy struct {
	x, y float64
}

type Game struct {
	playerX float64
	bullets []Bullet
	enemies []Enemy
    score int
    gameOver bool
	lastSpawnTime time.Time // Track last enemy spawn time
	lastShootTime time.Time // Track last bullet shoot time
    enemyRespawnGap int64
}

func (g *Game) Update() error {
    // Stop updates if game over
	if g.gameOver {
        if ebiten.IsKeyPressed(ebiten.KeyR) {
	        g.ResetGame()
        }
		return nil
	}

    g.SetDifficultyLevel()
    g.HandleMovements()
    g.HandleBulletFires()
    g.HandleBulletPositions()
    g.HandleEnemiesMovements()
    g.DestroyEnemyOnShot()
    g.AddNewEnemy()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the player spaceship (rectangle for now)
	ebitenutil.DrawRect(screen, g.playerX, screenHeight-50, 40, 10, color.White)

    // Draw bullets
    for _, b := range g.bullets {
	    ebitenutil.DrawRect(screen, b.x, b.y, 5, 10, color.White)
    }

    // Draw enemies
    for _, e := range g.enemies {
	    ebitenutil.DrawRect(screen, e.x, e.y, 40, 40, color.RGBA{255, 0, 0, 255})
    }

    ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))

	if g.gameOver {
		gameOverText := "GAME OVER\nPress R to Restart"

		// Get the font width and height for the text
		fontFace := basicfont.Face7x13
		textWidth := text.BoundString(fontFace, gameOverText).Dx()
		textHeight := fontFace.Metrics().Height.Ceil()

		// Calculate position to center the text
		xPos := (screenWidth - textWidth) / 2
		yPos := (screenHeight - textHeight) / 2

		// Draw the "Game Over" message at the calculated position
		text.Draw(screen, gameOverText, fontFace, xPos, yPos, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) HandleMovements() {
	// Move player left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.playerX > 0 {
		g.playerX -= playerSpeed
	}
	// Move player right
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.playerX < screenWidth-40 {
		g.playerX += playerSpeed
	}
}

func (g *Game) HandleBulletFires() {
    // Fire bullet when Space key is pressed
    if time.Since(g.lastShootTime).Milliseconds() >= 300 && ebiten.IsKeyPressed(ebiten.KeySpace) {
	    g.bullets = append(g.bullets, Bullet{x: g.playerX + 15, y: screenHeight - 60})
        g.lastShootTime = time.Now()
    }
}

func (g *Game) HandleBulletPositions() {
    // Update bullet positions
    for i := range g.bullets {
	    g.bullets[i].y -= 8 // Move bullets up
    }
}

func (g *Game) HandleEnemiesMovements() {
    // Move enemies downward
    for i := range g.enemies {
	    g.enemies[i].y += 1

		// 🔴 Check if an enemy reaches the bottom (GAME OVER)
		if g.enemies[i].y > screenHeight-40 {
			g.gameOver = true // Set game over state
            break
		}
    }
}

func (g *Game) SpawnEnemy() {
	enemyX := float64(rand.Intn(screenWidth - 40)) // Random X position
	g.enemies = append(g.enemies, Enemy{x: enemyX, y: 50}) // Add new enemy
}

func (g *Game) InitEnemies() {
	for i := 0; i < 5; i++ {
		g.enemies = append(g.enemies, Enemy{x: float64(i * 80), y: 50})
	}
}

func (g *Game) AddNewEnemy() {
    // ⏳ Enemy Respawn Every 1 Second
	if time.Since(g.lastSpawnTime).Milliseconds() >= g.enemyRespawnGap {
		g.SpawnEnemy() // Call function to spawn a new enemy
		g.lastSpawnTime = time.Now() // Reset spawn timer
	}
}

func (g *Game) ResetGame() {
	g.playerX = screenWidth / 2
	g.bullets = []Bullet{}
	g.enemies = []Enemy{}
	g.lastSpawnTime = time.Now()
	g.gameOver = false
}

func (g *Game) SetDifficultyLevel() {
    switch {
	    case g.score < 100:
		    g.enemyRespawnGap = 1000
	    case g.score >= 100 && g.score < 200:
		    g.enemyRespawnGap = 900
	    case g.score >= 200 && g.score < 300:
		    g.enemyRespawnGap = 800
	    case g.score >= 300 && g.score < 400:
		    g.enemyRespawnGap = 700
	    case g.score >= 400 && g.score < 500:
		    g.enemyRespawnGap = 500
	    default:
		    g.enemyRespawnGap = 300
	}
}

func (g *Game) DestroyEnemyOnShot() {
    // 🔥 Collision Detection (Bullets vs. Enemies)
	for bi := 0; bi < len(g.bullets); bi++ {
		b := g.bullets[bi]

		for ei := 0; ei < len(g.enemies); ei++ {
			e := g.enemies[ei]

			if b.x > e.x && b.x < e.x+40 && b.y > e.y && b.y < e.y+40 {
				// Remove the enemy
				g.enemies = append(g.enemies[:ei], g.enemies[ei+1:]...)

				// Remove the bullet
				g.bullets = append(g.bullets[:bi], g.bullets[bi+1:]...)

                // Update user score
                g.score += 10
				// Break out of the loop to prevent index errors
				break
			}
		}
	}
}

func main() {
	game := &Game{
        playerX: screenWidth / 2,
        lastSpawnTime: time.Now(),
        enemyRespawnGap: 1000, 
    }

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Space Shooter in Go")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
