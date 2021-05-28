package core

import (
	//"./nn"
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
	//nn.Model
	action_type      uint8
	last_inference   float64
	inference_gap    float64
	private_map      *PrivateMap
	now_route_target uint8
	clockwise        bool
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

func (hero *Lusian) PathSearching(gap_time float64) {

	// Here we need to take blockage into account
	isEnemyNearby, enemy := CheckEnemyInFrustum(hero.Camp(), hero)
	if isEnemyNearby {
		pos_enemy := enemy.Position()
		// Sometimes we cannot attack enemy that viewable to us
		canAttack := CanAttackEnemy(hero, &pos_enemy)

		if canAttack {
			NormalAttackEnemy(hero, enemy)
		} else {
			hero.MoveTowards(gap_time, pos_enemy)
		}
	}
	/*
		game := &GameInst
		pos := hero.Position()
			else {
				// Search for enemy, loop eight directons
				// Calculate the distance between now pos and target pos
				target_pos := game.BattleField.Route[hero.now_route_target]
				dist := vec3.Distance(&target_pos, &pos)
				if dist < 10 {
					if hero.clockwise && hero.now_route_target == uint8(len(game.BattleField.Route)-1) {
						hero.clockwise = false
					} else {
						if !hero.clockwise && hero.now_route_target == 0 {
							hero.clockwise = true
						}
					}

					if hero.clockwise {
						hero.now_route_target += 1
					} else {
						hero.now_route_target -= 1
					}

				} else {
					hero.MoveTowards(gap_time, target_pos)
				}
			}
	*/
}

func (hero *Lusian) Tick(gap_time float64) {
	game := &GameInst

	if game.ManualCtrlEnemy {
		CalculateViewDepth(hero)

		hero.ManualCtrl(gap_time)
		return
	} else {
		hero.PathSearching(gap_time)
		return
	}

	pos := hero.Position()

	isEnemyNearby, enemy := CheckEnemyNearby(hero.Camp(), hero.ViewRange(), &pos)
	if isEnemyNearby {
		pos_enemy := enemy.Position()

		canAttack := CanAttackEnemy(hero, &pos_enemy)

		if canAttack {
			NormalAttackEnemy(hero, enemy)
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

		}

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

		return
	} else {
		//	panic("Not supposed to be here")
	}
}
