package core

import (
	"math"
	"os"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/ungerik/go3d/vec3"
)

var LogHandle *os.File

func GetCallStack() string {
	var buffer [2048 * 1024]byte
	size := runtime.Stack(buffer[0:], true)
	crashed_stack := string(buffer[0:size])
	return crashed_stack
}

func LogStr(log string) {
	if LogHandle == nil {
		return
	}
	LogHandle.WriteString(log)

	if GameInst.LogoutStack {
		LogHandle.WriteString("\n")
		LogHandle.WriteString(GetCallStack())
	}

	LogHandle.WriteString("\n")
	LogHandle.Sync()
}

func LogBytes(log []byte) {
	if LogHandle == nil {
		return
	}
	LogHandle.Write(log)
	LogHandle.WriteString("\n")
	LogHandle.Sync()
}

func remove(s []Buff, i int) []Buff {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func CanAttackEnemy(unit BaseFunc, enemy_pos *vec3.T) bool {
	unit_pos := unit.Position()
	dist := vec3.Distance(enemy_pos, &unit_pos)
	if dist < unit.AttackRange() {
		return true
	}

	return false

}

func CheckUnitOnDir(position *vec3.T, dir *vec3.T) (bool, BaseFunc) {
	game := &GameInst

	for _, v := range game.BattleUnits {
		if v.Attackable() == false {
			continue
		}

		unit_pos := v.Position()
		unit_dir := unit_pos.Sub(position)
		angle := vec3.Angle(unit_dir, dir)
		if angle < float32(0.3) {
			return true, v
		}
	}

	return false, nil
}

func GetEnemyClearDir(my_camp int32, position *vec3.T) (bool, vec3.T) {
	game := &GameInst
	total_dir := vec3.T{0, 0, 0}

	for _, v := range game.BattleUnits {
		if v.Attackable() == false {
			continue
		}

		if v.Camp() != my_camp {
			unit_pos := v.Position()
			unit_dir := unit_pos.Sub(position)
			unit_dir.Normalize()
			total_dir = *total_dir.Add(unit_dir)
		}
	}

	total_dir.Normalize()
	total_dir[0] = -total_dir[0]
	total_dir[1] = -total_dir[1]

	return false, total_dir
}

func CheckEnemyOnDirAngle(my_camp int32, position *vec3.T, dir *vec3.T, angle_threshold float32) (bool, BaseFunc) {
	game := &GameInst

	for _, v := range game.BattleUnits {
		if v.Attackable() == false {
			continue
		}

		unit_pos := v.Position()
		unit_dir := unit_pos.Sub(position)
		angle := vec3.Angle(unit_dir, dir)
		if angle < float32(angle_threshold) && v.Camp() != my_camp {
			return true, v
		}
	}

	return false, nil
}

func CheckEnemyOnDir(my_camp int32, position *vec3.T, dir *vec3.T) (bool, BaseFunc) {
	game := &GameInst

	for _, v := range game.BattleUnits {
		if v.Attackable() == false {
			continue
		}

		unit_pos := v.Position()
		unit_dir := unit_pos.Sub(position)
		angle := vec3.Angle(unit_dir, dir)
		if angle < float32(0.3) && v.Camp() != my_camp {
			return true, v
		}
	}

	return false, nil
}

func CheckEnemyOnDirWithinDist(my_camp int32, position *vec3.T, dir *vec3.T, dist_thres float32) (bool, BaseFunc) {
	game := &GameInst

	for _, v := range game.BattleUnits {

		if v.Attackable() == false {
			continue
		}
		unit_pos := v.Position()
		dist := vec3.Distance(position, &unit_pos)
		unit_dir := unit_pos.Sub(position)
		angle := vec3.Angle(unit_dir, dir)
		if angle < float32(0.3) && v.Camp() != my_camp && dist < dist_thres {
			return true, v
		}
	}

	return false, nil
}

func CheckSegmentThroughObject(points []vec3.T, actor *StaticUnit) bool {
	through := false
	through, _ = actor.CheckBeenThrough(points[0], points[1])

	return through
}

func max_of_arr(a []float32) float32 {
	max_val := -float32(math.MaxFloat32)
	for _, v := range a {
		if v > max_val {
			max_val = v
		}
	}

	return max_val
}

func min_of_arr(a []float32) float32 {
	min_val := float32(math.MaxFloat32)
	for _, v := range a {
		if v < min_val {
			min_val = v
		}
	}

	return min_val
}

func rearr_minmax(min float32, max float32) (float32, float32) {
	if max < min {
		return max, min
	} else {
		return min, max
	}
}

func CalculateViewDepth(f0 BaseFunc) {
	game := &GameInst

	view_depth := f0.ViewDepth()
	start_point := f0.Position()
	view_frustum := f0.ViewFrustum()
	y_gap := vec3.Sub(&view_frustum[0], &view_frustum[3])
	x_gap := vec3.Sub(&view_frustum[2], &view_frustum[3])

	for _idx := 0; _idx < game.DepthMapSize; _idx += 1 {
		x_ratio := float32(_idx+1) / float32(game.DepthMapSize)
		x_offset := x_gap
		x_offset.Scale(x_ratio)

		for _idy := 0; _idy < game.DepthMapSize; _idy += 1 {
			y_ratio := float32(_idy+1) / float32(game.DepthMapSize)
			y_offset := y_gap
			y_offset.Scale(y_ratio)
			end_point := view_frustum[3]
			end_point.Add(&x_offset).Add(&y_offset)
			nearest_depth := float32(1.0)
			for _, v := range game.BattleField.Props {
				_collide, distance, _ := v.GetNearestCollidePoint(start_point, end_point)
				if _collide {
					if distance < nearest_depth {
						nearest_depth = distance
					}
				}
			}

			view_depth[_idx*game.DepthMapSize+_idy] = nearest_depth
		}
	}

}

func CheckSegmentThroughAABB(points []vec3.T, actor *StaticUnit) bool {

	a := points[0][0]
	b := points[1][0]
	m := actor.BB.Xmin
	txmin := (m - b) / (a - b)

	a = points[0][0]
	b = points[1][0]
	m = actor.BB.Xmax
	txmax := (m - b) / (a - b)

	a = points[0][1]
	b = points[1][1]
	m = actor.BB.Ymin
	tymin := (m - b) / (a - b)

	a = points[0][1]
	b = points[1][1]
	m = actor.BB.Ymax
	tymax := (m - b) / (a - b)

	a = points[0][2]
	b = points[1][2]
	m = actor.BB.Zmin
	tzmin := (m - b) / (a - b)

	a = points[0][2]
	b = points[1][2]
	m = actor.BB.Zmax
	tzmax := (m - b) / (a - b)

	txmin, txmax = rearr_minmax(txmin, txmax)
	tymin, tymax = rearr_minmax(tymin, tymax)
	tzmin, tzmax = rearr_minmax(tzmin, tzmax)

	enter_latest := max_of_arr([]float32{txmin, tymin, tzmin})
	leave_earliest := min_of_arr([]float32{txmax, tymax, tzmax})
	if leave_earliest >= enter_latest {
		return true
	}

	return false
}

func CheckIfSeperated(hero BaseFunc, enemy BaseFunc) bool {
	game := &GameInst

	var segment []vec3.T
	for _, v := range game.BattleField.Props {

		segment = []vec3.T{hero.Position(), enemy.Position()}
		through_aabb := CheckSegmentThroughAABB(segment, v)
		if through_aabb {
			through_object := CheckSegmentThroughObject(segment, v)
			if through_object {
				return true
			}
		}

	}
	return false
}

func CheckEnemyInFrustum(camp int32, hero BaseFunc) (bool, BaseFunc) {
	view_dir := hero.Direction()
	if view_dir[0] == 0 {
		return false, nil
	}

	game := &GameInst

	position := hero.Position()
	dist := float32(0)
	min_dist := hero.ViewRange()
	var min_dist_enemy BaseFunc
	min_dist_enemy = nil
	for _, v := range game.BattleUnits {
		if v.Attackable() == false {
			continue
		}

		unit_pos := v.Position()
		if v.Camp() != camp && v.Health() > 0 {
			dist = vec3.Distance(&position, &unit_pos)
			if dist < min_dist {
				enemy_dir := vec3.Sub(&unit_pos, &position)
				view_dir = hero.Direction()
				enemy_angle := vec3.Angle(&view_dir, &enemy_dir)
				if enemy_angle < mgl32.DegToRad(hero.Fov()/2) {
					// Check if the two objects are not seperated by Static objects
					seperated := CheckIfSeperated(hero, v)
					if !seperated {
						min_dist = dist
						min_dist_enemy = v
					}
				}
			}
		}
	}

	if min_dist_enemy != nil {
		return true, min_dist_enemy
	}

	// fmt.Println("->CheckEnemyNearby")
	return false, nil
}

func CheckEnemyNearby(camp int32, radius float32, position *vec3.T) (bool, BaseFunc) {
	game := &GameInst
	dist := float32(0)
	min_dist := radius
	var min_dist_enemy BaseFunc
	min_dist_enemy = nil
	for _, v := range game.BattleUnits {
		if v.Attackable() == false {
			continue
		}

		unit_pos := v.Position()
		dist = vec3.Distance(position, &unit_pos)
		if (dist < min_dist) && (v.Camp() != camp) && (v.Health() > 0) {
			min_dist = dist
			min_dist_enemy = v
		}
	}

	if min_dist_enemy != nil {
		return true, min_dist_enemy
	}

	// fmt.Println("->CheckEnemyNearby")
	return false, nil
}

func SelectFirstNonHeroEnemy(camp int32, radius float32, position *vec3.T) (bool, BaseFunc) {
	game := &GameInst
	dist := float32(0)
	var fall_back_hero_enemy BaseFunc
	fall_back_hero_enemy = nil

	for _, v := range game.BattleUnits {
		if v.Attackable() == false {
			continue
		}

		unit_pos := v.Position()
		dist = vec3.Distance(position, &unit_pos)
		if (dist < radius) && (v.Camp() != camp) && (v.Health() > 0) {
			if IsCreep(v) {
				return true, v
			}

			if fall_back_hero_enemy == nil {
				fall_back_hero_enemy = v
			}
		}
	}

	// fmt.Println("->CheckEnemyNearby")
	return fall_back_hero_enemy != nil, fall_back_hero_enemy
}

func InitWithCamp(battle_unit BaseFunc, camp int32) {
	const_num := float32(0.707106781)
	extent := battle_unit.Extent()
	if camp == 0 {
		battle_unit.SetPosition(vec3.T{1, 999, extent[2]})
		battle_unit.SetCamp(0)
		battle_unit.SetDirection(vec3.T{const_num, -const_num})
	} else {
		battle_unit.SetPosition(vec3.T{999, 1, extent[2]})
		battle_unit.SetCamp(1)
		battle_unit.SetDirection(vec3.T{-const_num, const_num})
	}
}

func InitHeroWithCamp(battle_unit BaseFunc, camp int32, pos_x float32, pos_y float32) {
	InitWithCamp(battle_unit, camp)
	extent := battle_unit.Extent()
	battle_unit.SetPosition(vec3.T{pos_x, pos_y, extent[2]})

	hero_unit, ok := battle_unit.(HeroFunc)
	if ok {
		hero_unit.SetTargetPos(pos_x, pos_y)
	}

	battle_unit.SetDirection(vec3.T{0, 0, 0})
	battle_unit.ClearAllBuff()
	battle_unit.ClearLastSkillUseTime()
}

func ConvertNum2Dir(action_code int) (dir vec3.T) {
	offset_x := float32(0)
	offset_y := float32(0)
	const_val := float32(200)

	switch action_code {
	case 0: // do nothing
		offset_x = float32(const_val)
		offset_y = float32(const_val)
	case 1:
		offset_x = float32(-const_val)
		offset_y = float32(-const_val)
	case 2:
		offset_x = float32(0)
		offset_y = float32(-const_val)
	case 3:
		offset_x = float32(const_val)
		offset_y = float32(-const_val)
	case 4:
		offset_x = float32(-const_val)
		offset_y = float32(0)
	case 5:
		offset_x = float32(const_val)
		offset_y = float32(0)
	case 6:
		offset_x = float32(-const_val)
		offset_y = float32(const_val)
	case 7:
		offset_x = float32(0)
		offset_y = float32(const_val)
	}

	dir[0] = offset_x
	dir[1] = offset_y
	return dir
}
