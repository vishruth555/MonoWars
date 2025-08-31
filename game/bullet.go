package game

const MaxBulletSpeed = 0.5 //max is 1 due to swap tile constraints

type Bullet struct {
	Id       int
	X        float32
	Y        float32
	Dir      Direction
	isActive bool
}

func NewBullet(id int, x float32, y float32, dir Direction) *Bullet {

	playerCenterX := int(x + 0.4)
	playerCenterY := int(y + 0.4)

	bulletX := float32(playerCenterX) + 0.5
	bulletY := float32(playerCenterY) + 0.5

	return &Bullet{
		Id:       id,
		X:        bulletX,
		Y:        bulletY,
		Dir:      dir,
		isActive: true,
	}
}

func (bullet *Bullet) Move(g *Game) bool {
	switch bullet.Dir {

	case Left:
		bullet.X -= MaxBulletSpeed

	case Right:
		bullet.X += MaxBulletSpeed

	case Up:
		bullet.Y -= MaxBulletSpeed

	case Down:
		bullet.Y += MaxBulletSpeed

	}

	xPos := int(bullet.X)
	yPos := int(bullet.Y)

	if g.TileMap[yPos][xPos] == 0 {
		bullet.isActive = false
		return false
	}
	if g.TileMap[yPos][xPos] == bullet.Id {
		g.SwapTile(xPos, yPos)
	}
	return true
}
