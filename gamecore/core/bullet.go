package core

import (
	"github.com/ungerik/go3d/vec3"
)

type Bullet struct {
	BaseInfo
}

var bullet_template Bullet

func (bullet *Bullet) Init(a ...interface{}) BaseFunc {
	return bullet
}

func (bullet *Bullet) Tick(gap_time float64) {
	game := &GameInst

	now_seconds := game.LogicTime
	if (bullet.LastAttackTime() + bullet.AttackFreq()) < float64(now_seconds) {
		bullet.SetHealth(0)
		return
	}

	pos := bullet.Position()
	//LogStr(fmt.Sprintf("Bullet: %v", pos))
	isEnemyNearby, enemy := CheckEnemyNearby(bullet.Camp(), bullet.AttackRange(), &pos)
	if isEnemyNearby && enemy.GetId() != 0 {
		// Check if time to make hurt
		enemy.DealDamage(bullet.Damage())
		//LogStr("Bullet: Damage Dealed")
		bullet.SetHealth(0)

	} else {
		// Check if the bullet is within a building
		within := false
		//LogStr(fmt.Sprintf("Prop:"))
		for _, v := range game.BattleField.Props {
			//LogStr(fmt.Sprintf("Prop: %v", v))
			within = v.CheckWithin(pos)
			//LogStr(fmt.Sprintf("Prop Within: %v", within))
			if within {
				bullet.SetHealth(0)
				//LogStr(fmt.Sprintf("Bullet: In a prop"))
				return
			}
		}
		// March towards target direction
		dir := bullet.Direction()
		dir = dir.Scaled(float32(gap_time))
		dir = dir.Scaled(float32(bullet.speed))
		newPos := vec3.Add(&pos, &dir)
		bullet.SetPosition(newPos)
	}
}
