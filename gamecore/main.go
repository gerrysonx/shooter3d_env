package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"./core"
	"./render"
	"./test"
)

func main() {
	root_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		var buffer [2048 * 1024]byte
		size := runtime.Stack(buffer[0:], true)
		crashed_stack := string(buffer[0:size])
		if strings.Contains(crashed_stack, "panic") {
			fmt.Printf("stack is:%v\n", crashed_stack)
			now := time.Now()
			crash_file_handle, _err0 := os.Create(fmt.Sprintf("%s/../crash_%v.log", root_dir, now.UnixNano()))
			if _err0 != nil {
				fmt.Println("Create crash_file_handle file failed.")
				return
			}
			crash_file_handle.WriteString(crashed_stack)
			crash_file_handle.WriteString("\n")
			crash_file_handle.Sync()
			crash_file_handle.Close()
		} else {
			fmt.Printf("no panic")
		}
	}()

	_target_frame_gap_time := flag.Float64("frame_gap", 0.015, "")
	_fix_update := flag.Bool("fix_update", true, "a bool")
	_run_render := flag.Bool("render", true, "a bool")
	_input_gap_time := flag.Float64("input_gap", 0.2, "")
	_manual_enemy := flag.Bool("manual_enemy", false, "a bool")
	_gym_mode := flag.Bool("gym_mode", true, "a bool")
	_debug_log := flag.Bool("debug_log", false, "a bool")
	_slow_tick := flag.Bool("slow_tick", false, "a bool")
	_multi_player := flag.Bool("multi_player", true, "a bool")
	_scene_id := flag.Int("scene", 10, "a bool")
	_test_code := flag.Bool("test_code", false, "a bool")

	if *_test_code {
		test.Main()
		return
	}

	file_handle, _err := os.Create(fmt.Sprintf("%s/../shooter3d_core.log", root_dir))
	if _err != nil {
		fmt.Println("Create log file failed.")
		return
	}

	defer file_handle.Close()

	if *_debug_log {
		core.LogHandle = file_handle
	} else {
		core.LogHandle = nil
	}
	//core.LogHandle = file_handle
	flag.Parse()
	core.LogStr(fmt.Sprintf("main is called"))

	core.SkillMgrInst.LoadCfgFolder(fmt.Sprintf("%s/cfg/skills", root_dir))
	core.HeroMgrInst.LoadCfgFolder(fmt.Sprintf("%s/cfg/heroes", root_dir))

	core.GameInst = core.Game{}

	core.GameInst.ManualCtrlEnemy = *_manual_enemy
	core.GameInst.MultiPlayer = *_multi_player
	core.GameInst.SceneId = *_scene_id

	core.GameInst.Init()

	if *_run_render {
		render.RendererInst = render.Renderer{}
		render.RendererInst.InitRenderEnv(&core.GameInst)
		defer render.RendererInst.Release()
	}

	now := time.Now()
	_before_tick_time := float64(now.UnixNano()) / 1e9
	_after_tick_time := _before_tick_time
	_logic_cost_time := _before_tick_time
	_gap_time := float64(0)

	_action_stamp := int(0)
	// Set target frame time gap

	if *_fix_update {
		_last_input_time := core.GameInst.LogicTime
		var action_code int   // raw action code
		var action_code_0 int // action code
		var action_code_1 int // move dir code
		var action_code_2 int // skill dir code
		var action_code_3 int // vision yaw code
		var action_code_4 int // vision pitch code
		for {
			// Process input from user
			if core.GameInst.LogicTime > _last_input_time+*_input_gap_time && *_gym_mode {
				_last_input_time = core.GameInst.LogicTime

				// Output game state to stdout
				if *_multi_player {
					game_state_str := core.GameInst.DumpVarPlayerGameState()
					// core.LogBytes(file_handle, game_state_str)
					// core.LogStr(fmt.Sprintf("Every step, logic_time:%v, _action_stamp:%d, game_state_str:%s", core.GameInst.LogicTime, _action_stamp, game_state_str))
					fmt.Printf("%d@%s\n", _action_stamp, game_state_str)
				} else {
				}

				//action_code = 4352
				if *_multi_player {
					// Input hero count.
					fmt.Scanf("%d\n", &action_code)
					// core.LogStr(fmt.Sprintf("_multi_player mode, time:%v, get action code:%v", core.GameInst.LogicTime, action_code))
					if action_code == 36864 {
						// Init game
						_action_stamp = 0
						core.GameInst.HandleMultiPlayerAction(0, 9, 0, 0)
						_last_input_time = 0
						// game_state_str := core.GameInst.DumpVarPlayerGameState()
						// core.LogStr(fmt.Sprintf("After game.init, time:%v, game state is:%s", core.GameInst.LogicTime, game_state_str))
						continue
					}

					// First action_code is player count
					player_count := action_code
					for _idx := 0; _idx < player_count; _idx += 1 {
						fmt.Scanf("%d\n", &action_code)

						_action_stamp = action_code >> 20
						//core.LogStr(fmt.Sprintf("action: %v", action_code>>2))
						action_code_0 = (action_code >> 16) & 0xf
						action_code_1 = (action_code >> 12) & 0xf
						action_code_2 = (action_code >> 8) & 0xf
						action_code_3 = (action_code >> 4) & 0x3
						action_code_4 = (action_code >> 2) & 0x3
						// core.LogStr(fmt.Sprintf("_multi_player mode, player id:%d get action code:%v, action_code_0:%v, action_code_1:%v, action_code_2:%v",
						// 	_idx, action_code, action_code_0, action_code_1, action_code_2))
						//core.LogStr(fmt.Sprintf("action 3: %v, action 4: %v", action_code_3, action_code_4))
						core.GameInst.HandleMultiPlayerAction(_idx, action_code_0, action_code_1, action_code_2)
						core.GameInst.HandleMultiPlayerVision(_idx, action_code_3, action_code_4)
						//core.GameInst.HandleMultiPlayerVision(_idx, 0, 1)

						if 9 == action_code_0 {
							// Instantly output
							_last_input_time = 0
						}
					}

				} else {
					// Input action code
					// Only handle multi player mode.

				}

				if *_slow_tick {
					gap_time_in_nanoseconds := *_target_frame_gap_time * float64(time.Second)
					time.Sleep(time.Duration(gap_time_in_nanoseconds) * 100) // was two
				}
			}

			if false == *_gym_mode {
				// Output game state to stdout
				if *_multi_player {
					game_state_str := core.GameInst.DumpVarPlayerGameState()
					// core.LogBytes(file_handle, game_state_str)
					fmt.Printf("%d@%s\n", _action_stamp, game_state_str)
				} else {

				}
				gap_time_in_nanoseconds := *_target_frame_gap_time * float64(time.Second)
				time.Sleep(time.Duration(gap_time_in_nanoseconds))
			}
			//	game_state_str := core.GameInst.DumpVarPlayerGameState()
			// core.LogBytes(file_handle, game_state_str)
			//	core.LogStr(fmt.Sprintf("BeforeTick, _action_stamp:%v, game_state_str:%s", core.GameInst.LogicTime, game_state_str))
			if !core.GameInst.Paused {
				core.GameInst.Tick(*_target_frame_gap_time)
			}

			//	game_state_str = core.GameInst.DumpVarPlayerGameState()
			// core.LogBytes(file_handle, game_state_str)
			//	core.LogStr(fmt.Sprintf("AfterTick, _action_stamp:%v, game_state_str:%s", core.GameInst.LogicTime, game_state_str))
			// Draw logic units on output image
			if *_run_render {
				render.RendererInst.Render()
			}

		}
	} else {
		for {
			_gap_time = _after_tick_time - _before_tick_time

			// fmt.Printf("------------------------------\nMain game loop, _before_tick_time is:%v\n", _before_tick_time)
			_before_tick_time = _after_tick_time
			// Do the tick
			core.GameInst.Tick(_gap_time)
			// Draw logic units on output image
			if *_run_render {
				render.RendererInst.Render()
			}

			now = time.Now()
			_logic_cost_time = (float64(now.UnixNano()) / 1e9) - _before_tick_time

			if _logic_cost_time > *_target_frame_gap_time {
			} else {
				gap_time_in_nanoseconds := (*_target_frame_gap_time - _logic_cost_time) * float64(time.Second)
				time.Sleep(time.Duration(gap_time_in_nanoseconds))
			}
			now = time.Now()
			_after_tick_time = float64(now.UnixNano()) / 1e9
			// fmt.Printf("Main game loop, _after_tick_time is:%v, _gap_time is:%v, logic cost time is:%v\n------------------------------\n", _after_tick_time, _gap_time, _logic_cost_time)
		}
	}

}
