package main

import (
	"github.com/paked/engi"
	"log"
	"math/rand"
)

var (
	W Game
)

type Game struct {
	engi.World
}

func (game Game) Preload() {
	engi.Files.Add("guy", "data/icon.png")
	engi.Files.Add("rock", "data/rock.png")
	engi.Files.Add("font", "data/font.go")
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)
	game.AddSystem(&engi.RenderSystem{})
	game.AddSystem(&engi.CollisionSystem{})
	game.AddSystem(&DeathSystem{})
	game.AddSystem(&FallingSystem{})
	game.AddSystem(&ControlSystem{})
	game.AddSystem(&RockSpawnSystem{})

	guy := engi.NewEntity([]string{"RenderSystem", "ControlSystem", "RockSpawnSystem", "CollisionSystem", "DeathSystem"})
	texture := engi.Files.Image("guy")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "guy")
	collisionMaster := engi.CollisionMasterComponent{}

	width := texture.Width() * render.Scale.X
	height := texture.Height() * render.Scale.Y

	space := engi.SpaceComponent{engi.Point{(engi.Width() - width) / 2, (engi.Height() - height) / 2}, width, height}
	guy.AddComponent(&render)
	guy.AddComponent(&space)
	guy.AddComponent(&collisionMaster)

	game.AddEntity(guy)
}

type ControlSystem struct {
	*engi.System
}

func (control *ControlSystem) New() {
	control.System = &engi.System{}
}

func (control ControlSystem) Name() string {
	return "ControlSystem"
}

func (control *ControlSystem) Update(entity *engi.Entity, dt float32) {
	var space *engi.SpaceComponent
	if !entity.GetComponent(&space) {
		return
	}

	speed := 400 * dt

	if engi.Keys.KEY_A.Down() {
		space.Position.X -= speed
	}

	if engi.Keys.KEY_D.Down() {
		space.Position.X += speed
	}

	if engi.Keys.KEY_W.Down() {
		space.Position.Y -= speed
	}

	if engi.Keys.KEY_S.Down() {
		space.Position.Y += speed
	}
}

type RockSpawnSystem struct {
	*engi.System
}

func (rock RockSpawnSystem) Name() string {
	return "RockSpawnSystem"
}

func (rock *RockSpawnSystem) New() {
	rock.System = &engi.System{}
}

func (rock *RockSpawnSystem) Update(entity *engi.Entity, dt float32) {
	if rand.Float32() < .96 {
		return
	}

	position := engi.Point{0, -32}
	position.X = rand.Float32() * (engi.Width())
	W.AddEntity(NewRock(position))
}

func NewRock(position engi.Point) *engi.Entity {
	rock := engi.NewEntity([]string{"RenderSystem", "FallingSystem", "CollisionSystem", "SpeedSystem"})
	texture := engi.Files.Image("rock")
	render := engi.NewRenderComponent(texture, engi.Point{4, 4}, "rock")
	space := engi.SpaceComponent{position, texture.Width() * render.Scale.X, texture.Height() * render.Scale.Y}
	rock.AddComponent(&render)
	rock.AddComponent(&space)
	return rock
}

type FallingSystem struct {
	*engi.System
}

func (fs *FallingSystem) New() {
	fs.System = &engi.System{}
	engi.Mailbox.Listen("CollisionMessage", fs)
}

func (fs FallingSystem) Name() string {
	return "FallingSystem"
}

func (fs FallingSystem) Update(entity *engi.Entity, dt float32) {
	var space *engi.SpaceComponent
	if !entity.GetComponent(&space) {
		return
	}
	space.Position.Y += 200 * dt
}

type DeathSystem struct {
	*engi.System
}

func (ds *DeathSystem) New() {
	ds.System = &engi.System{}
	engi.Mailbox.Listen("CollisionMessage", ds)
}

func (ds DeathSystem) Name() string {
	return "DeathSystem"
}

func (fs DeathSystem) Update(entity *engi.Entity, dt float32) {

}

func (fs DeathSystem) Receive(message engi.Message) {
	collision, isCollision := message.(engi.CollisionMessage)
	if isCollision {
		log.Println(collision, message)
		log.Println("DEAD")
	}
}

func main() {
	log.Println("[Game] Says hello, written in github.com/paked/engi + Go")
	W = Game{}
	engi.Open("Stream Game", 800, 800, false, &W)
}
