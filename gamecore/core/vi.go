package core

import (
	"fmt"

	"github.com/ungerik/go3d/vec3"
)

type Vi struct {
	Hero
}

func (hero *Vi) HandleAICommand(gap_time float64) {
	pos := hero.Position()
	// Check milestone distance
	targetPos := hero.TargetPos()
	// Todo: remain z value
	targetPos[2] = pos[2]

	//Check if the hero is in Flag area
	if CheckWithinFlagArea(&GameInst, &pos) {
		LogStr(fmt.Sprintf("Securing"))
		GameInst.isSecuring = 1
	}

	if GameInst.TimeInitSecuring != -1 {
		return
	}

	dist2tar := vec3.Distance(&pos, &targetPos)

	pos_ready := false
	if dist2tar < 5 {
		// Already the last milestone
		pos_ready = true
	}

	//LogStr(fmt.Sprintf("self ViewDir: %v", hero.Viewdir()))

	isEnemyNearby, enemy := CheckEnemyInFrustum(hero.Camp(), hero)
	if isEnemyNearby {
		// Check if time to make hurt
		NormalAttackEnemy(hero, enemy)
	}
	if !pos_ready {
		Chase(hero, targetPos, gap_time)
	}
}

func (hero *Vi) Tick(gap_time float64) {
	// Used for AI control
	game := &GameInst
	targetPos := hero.TargetPos()

	//CalculateViewDepth(hero)
	//LogStr(fmt.Sprintf("target pos: %v, %v", targetPos[0], targetPos[1]))
	hero.HandleAICommand(gap_time)

	//	hero.ManualCtrl(gap_time)
	return

	if game.ManualCtrlEnemy {
		hero.ManualCtrl(gap_time)
		return
	}

	pos := hero.Position()
	// Check milestone distance
	//targetPos := hero.TargetPos()
	dist := vec3.Distance(&pos, &targetPos)
	pos_ready := false
	if dist < 4 {
		// Already the last milestone
		pos_ready = true
	}

	isEnemyNearby, enemy := CheckEnemyNearby(hero.Camp(), hero.AttackRange(), &pos)
	if isEnemyNearby && pos_ready {
		// Check if time to make hurt
		NormalAttackEnemy(hero, enemy)
	} else {
		if pos_ready {
			// Do nothing
		} else {
			LogStr("Possible 2")
			Chase(hero, targetPos, gap_time)
		}
	}
}
