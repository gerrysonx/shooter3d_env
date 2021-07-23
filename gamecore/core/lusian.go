package core

import (
	"github.com/ungerik/go3d/vec3"
)

type PrivateGrid struct {
	has_enemy        bool
	last_view_scan   float64
	last_person_scan float64
}

type PrivateMap struct {
	grids [1000 * 1000]PrivateGrid
}

type Lusian struct {
	Hero
	action_type       []int
	last_inference    float64
	inference_gap     float64
	private_map       *PrivateMap
	now_route_target  uint8
	clockwise         bool
	enemyLastPosition vec3.T
	enemyLastSeenTime float64
}

func (hero *Lusian) HandleAICommand(gap_time float64) {
	pos := hero.Position()
	// Check milestone distance
	targetPos := hero.TargetPos()
	// Todo: remain z value
	targetPos[2] = pos[2]

	if GameInst.var_player_train_state.CampCTF == 1 {
		if CheckWithinFlagArea(&GameInst, &pos) {
			GameInst.isSecuring += 2
		}
	}

	//Check if the hero is in Flag area

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
		//hero.MoveTowards(gap_time, targetPos)
		Chase(hero, targetPos, gap_time)
	}
}

func (hero *Lusian) MoveTowards(gap_time float64, pos_enemy vec3.T) {
	pos := hero.Position()
	dir := vec3.Sub(&pos_enemy, &pos)
	dir.Normalize()
	hero.SetDirection(dir)

	dir = dir.Scaled(float32(gap_time))
	dir = dir.Scaled(float32(hero.speed))
	newPos := vec3.Add(&pos, &dir)
	hero.SetPosition(newPos)
}

func (hero *Lusian) Tick(gap_time float64) {
	game := &GameInst

	if game.ManualCtrlEnemy {

		hero.ManualCtrl(gap_time)
		return
	} else {

		hero.HandleAICommand(gap_time)
		return
	}
	now_seconds := game.LogicTime

	//Get NN model Input from below
	if (hero.last_inference + hero.inference_gap) < float64(now_seconds) {
		return
	} else {
		panic("Not supposed to be here")
	}

}

//Deprecated in Self-play Mode
/*
func (hero *Lusian) PathSearching(gap_time float64) {

	// Here we need to take blockage into account
	isEnemyNearby, enemy := CheckEnemyInFrustum(hero.Camp(), hero)
	game := &GameInst
	//LogStr(fmt.Sprintf("enemy last position: %v", hero.enemyLastPosition))
	LogStr(fmt.Sprintf("Oppo ViewDir: %v", hero.Viewdir()))
	if isEnemyNearby {
		pos_enemy := enemy.Position()
		// Sometimes we cannot attack enemy that viewable to us
		canAttack := CanAttackEnemy(hero, &pos_enemy)

		if canAttack {
			NormalAttackEnemy(hero, enemy)
		} else {
			//hero.MoveTowards(gap_time, pos_enemy)
			//Chase(hero, pos_enemy, gap_time)
			//LogStr("Chase in sight")

		}
		hero.enemyLastPosition = enemy.Position()
		hero.enemyLastSeenTime = game.LogicTime
	} else {
		pos := hero.Position()
		if vec3.Distance(&pos, &hero.enemyLastPosition) < 5 || game.LogicTime-hero.enemyLastSeenTime < 1 {
			view_dir := hero.Viewdir()
			length := math.Sqrt(float64(view_dir[0]*view_dir[0] + view_dir[1]*view_dir[1]))
			theta := math.Atan2(float64(view_dir[1]), float64(view_dir[0]))
			view_dir[0] = float32(math.Cos(theta+0.2) * length)
			view_dir[1] = float32(math.Sin(theta+0.2) * length)
			hero.SetViewdir(view_dir)
		} else {
			hero.MoveTowards(gap_time, hero.enemyLastPosition)
		}
	}
}*/

//Some Pre-written Logic, Deprecated in Self-play Mode
/*
	if need_move {
			dir.Normalize()
			hero.SetDirection(dir)

			dir = dir.Scaled(float32(gap_time))
			dir = dir.Scaled(float32(hero.speed))
			newPos := vec3.Add(&pos, &dir)
			hero.SetPosition(newPos)
			// dist_after_move := vec3.Distance(&pos_enemy, &pos)
			//	LogStr(fmt.Sprintf("Lusian need move toward enemy:%v, attack range:%v, dist:%v, dist_after_move:%v, time:%v",
			//		enemy.GetId(), enemy.AttackRange(), dist, dist_after_move, game.LogicTime))
		}


	// Check enemy and self distance
	// dist := vec3.Distance(&pos_enemy, &pos)
	dir_a := enemy.Position()
	dir_b := hero.Position()
	var dir vec3.T
	need_move := true
	if true { // dist > enemy.AttackRange()+35
		if canAttack == false {
			// March towards enemy
			dir = vec3.Sub(&dir_a, &dir_b)
			//		LogStr(fmt.Sprintf("Lusian need move toward enemy:%v, attack range:%v, dist:%v time:%v", enemy.GetId(), enemy.AttackRange(), dist, game.LogicTime))
		} else {
			need_move = false
		}
	} else {
		// March to the opposite direction of enemy

		has_clear_dir, clear_dir := GetEnemyClearDir(hero.Camp(), &pos)
		if has_clear_dir {
			dir = clear_dir
		} else {
			dir = vec3.Sub(&dir_b, &dir_a)
		}
		//	LogStr(fmt.Sprintf("Lusian need move toward enemy:%v, attack range:%v, dist:%v, has_clear_dir:%v, time:%v",
		//		enemy.GetId(), enemy.AttackRange(), dist, has_clear_dir, game.LogicTime))

	}*/
