import gym
import time
import subprocess
import os
import json
import numpy as np
import copy
from gym import error, spaces, utils
from gym.utils import seeding

ONE_HERO_FEATURE_SIZE = 7
BATTLE_FIELD_SIZE = 1000.0
MAIN_ACTION_DIMS = 3
MOVE_DIMS = 8
SKILL_DIMS = 3
VISION_YAW_DIMS = 3
VISION_PITCH_DIMS = 3
DEPTH_MAP_SIZE = 30

def FindHeroSkills(hero_cfg_file_path, hero_id):
    hero_cfg_file = '{}/{}.json'.format(hero_cfg_file_path, hero_id)
    hero_skills = []
    with open(hero_cfg_file, 'r') as file_handle:
        hero_cfg = json.load(file_handle)
        hero_skills = hero_cfg['skills']
        return hero_skills

def GetSkillTypes(skill_cfg_file_path, hero_skills):
    skill_dir_type_check = []
    for skill_id in hero_skills:
        skill_dir_type = False
        skill_id = int(skill_id)
        if -1 != skill_id:
            skill_cfg_file = '{}/{}.json'.format(skill_cfg_file_path, skill_id)
            with open(skill_cfg_file, 'r') as file_handle:
                skill_cfg = json.load(file_handle)
                if 0 == skill_cfg['type']:
                    skill_dir_type = True
        skill_dir_type_check.append(skill_dir_type)
    return skill_dir_type_check    

def InitMetaConfig(scene_id):
    global BATTLE_FIELD_SIZE
    global DEPTH_MAP_SIZE

    obs_size = 0
    self_hero_count = 0
    oppo_hero_count = 0

    dir_skill_mask = []
    try:
        # Load train self heroes skill masks
        root_folder = os.path.split(os.path.abspath(__file__))[0]
        
        
        cfg_file_path = '{}/../../../gamecore/cfg'.format(root_folder)
        training_map_file = '{}/maps/{}.json'.format(cfg_file_path, scene_id)
        hero_cfg_file_path = '{}/heroes'.format(cfg_file_path)
        skill_cfg_file_path = '{}/skills'.format(cfg_file_path)
        map_dict = None
        with open(training_map_file, 'r') as file_handle:
            map_dict = json.load(file_handle)

        BATTLE_FIELD_SIZE = map_dict['Width']
        DEPTH_MAP_SIZE = map_dict['DepthMapSize']

        for hero_id in map_dict['SelfHeroes']:
            hero_skills = FindHeroSkills(hero_cfg_file_path, hero_id)
            hero_skill_types = GetSkillTypes(skill_cfg_file_path, hero_skills)
            dir_skill_mask.append(hero_skill_types)
            
        self_hero_count = len(map_dict['SelfHeroes'])
        oppo_hero_count = len(map_dict['OppoHeroes'])
        obs_size = (oppo_hero_count + self_hero_count) * ONE_HERO_FEATURE_SIZE

    except Exception as ex:
        pass        
    pass

    return (self_hero_count, obs_size), (self_hero_count, MAIN_ACTION_DIMS, MOVE_DIMS, SKILL_DIMS, VISION_YAW_DIMS, VISION_PITCH_DIMS), self_hero_count, oppo_hero_count, dir_skill_mask

class ValorantEnvInfo(dict):
    def __init__(self, step_idx = 0):
        self.step_idx = step_idx
        pass

class ValorantMultiPlayerEnv(gym.Env):
    metadata = {'render.modes':['human']}
    def restart_proc(self):
    #    print('moba_env restart_proc is called.')
        self.done = False
        
        self.state = None
        self.depth = None
        self.reward = 0
        self.info = ValorantEnvInfo()
        self.last_state = None
        self.step_idx = 0                
        pass

    def __init__(self):
        # Create process, and communicate with std
        is_train = True
        root_folder = os.path.split(os.path.abspath(__file__))[0]

        my_env = os.environ.copy()
        is_train = True #my_env['moba_env_is_train'] == 'True'
        scene_id = 10 #my_env['moba_env_scene_id']    
        
        manual_str = '-manual_enemy=false'
        if not is_train:
            manual_str = '-manual_enemy=true'
        
        my_env['TF_CPP_MIN_LOG_LEVEL'] = '3'
        gamecore_file_path = '{}/../../../gamecore/gamecore'.format(root_folder)
        self.proc = subprocess.Popen([gamecore_file_path, 
                                '-render=false', '-gym_mode=true', '-debug_log=false', '-slow_tick=false', 
                                '-multi_player=true', '-scene={}'.format(scene_id), manual_str],
                                stdin=subprocess.PIPE,
                                stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE, env=my_env)

        # We shall parse the scene config file to get the input space and action space
        state_shape, action_shape, self_hero_count, oppo_hero_count, hero_dir_skill_mask = InitMetaConfig(scene_id)

        self.observation_space = spaces.Box(low = -0.5, high = 0.5, shape=state_shape, dtype=np.float32)

        # Action space omits the Tackle/Catch actions, which are useful on defense
        self.action_space = spaces.Box(low = -1, high = 65536, shape=action_shape, dtype=np.int32)  
        self.action_space.low = -1
        self.action_space.high = 65535
        self.self_hero_count = self_hero_count
        self.oppo_hero_count = oppo_hero_count

        self.restart_proc()
        self.step_idx = -1
        self.hero_dir_skill_mask = hero_dir_skill_mask
        
        print('moba_env initialized.')
        pass

    def fill_state(self, state, depth, json_data):
        # View from the point view of each self
        self_hero_count = int(json_data['SelfHeroCount'])

        for _idx in range(DEPTH_MAP_SIZE * DEPTH_MAP_SIZE):
            depth[0][_idx] = json_data['SelfHeroDepthMap'][0][_idx]
        
        for hero_idx in range(self_hero_count):        
            feature_idx = 0    
            state[hero_idx][feature_idx] = (float(json_data['SelfHeroPosX'][hero_idx]) / BATTLE_FIELD_SIZE)
            feature_idx += 1
            state[hero_idx][feature_idx] = (float(json_data['SelfHeroPosY'][hero_idx]) / BATTLE_FIELD_SIZE)
            feature_idx += 1
            state[hero_idx][feature_idx] = (float(json_data['SelfHeroPosZ'][hero_idx]) / 250) - 1
            feature_idx += 1
            state[hero_idx][feature_idx] = (float(json_data['SelfHeroHealth'][hero_idx]) / float(json_data['SelfHeroHealthFull'][hero_idx]))
            feature_idx += 1      
            state[hero_idx][feature_idx] = float(json_data['SelfHeroDirX'][hero_idx])
            feature_idx += 1
            state[hero_idx][feature_idx] = float(json_data['SelfHeroDirY'][hero_idx])
            feature_idx += 1
            state[hero_idx][feature_idx] = float(json_data['SelfHeroDirZ'][hero_idx])
            feature_idx += 1
            for _id_0 in range(self_hero_count):
                if _id_0 == hero_idx:
                    continue
                state[hero_idx][feature_idx] = (float(json_data['SelfHeroPosX'][_id_0]) / BATTLE_FIELD_SIZE)
                feature_idx += 1
                state[hero_idx][feature_idx] = (float(json_data['SelfHeroPosY'][_id_0]) / BATTLE_FIELD_SIZE)
                feature_idx += 1
                state[hero_idx][feature_idx] = (float(json_data['SelfHeroPosZ'][_id_0]) / 250) - 1
                feature_idx += 1
                state[hero_idx][feature_idx] = (float(json_data['SelfHeroHealth'][_id_0]) / float(json_data['SelfHeroHealthFull'][_id_0]))
                feature_idx += 1
                state[hero_idx][feature_idx] = float(json_data['SelfHeroDirX'][hero_idx])
                feature_idx += 1
                state[hero_idx][feature_idx] = float(json_data['SelfHeroDirY'][hero_idx])
                feature_idx += 1
                state[hero_idx][feature_idx] = float(json_data['SelfHeroDirZ'][hero_idx])
                feature_idx += 1

            oppo_hero_count = int(json_data['OppoHeroCount'])
            for _id_1 in range(oppo_hero_count):            
                state[hero_idx][feature_idx] = (float(json_data['OppoHeroPosX'][_id_1]) / BATTLE_FIELD_SIZE)
                feature_idx += 1
                state[hero_idx][feature_idx] = (float(json_data['OppoHeroPosY'][_id_1]) / BATTLE_FIELD_SIZE)
                feature_idx += 1
                state[hero_idx][feature_idx] = (float(json_data['OppoHeroPosZ'][_id_1]) / 250) - 1
                feature_idx += 1
                state[hero_idx][feature_idx] = (float(json_data['OppoHeroHealth'][_id_1]) / float(json_data['OppoHeroHealthFull'][_id_1]))
                feature_idx += 1
                state[hero_idx][feature_idx] = float(json_data['OppoHeroDirX'][hero_idx])
                feature_idx += 1
                state[hero_idx][feature_idx] = float(json_data['OppoHeroDirY'][hero_idx])
                feature_idx += 1
                state[hero_idx][feature_idx] = float(json_data['OppoHeroDirZ'][hero_idx])
                feature_idx += 1
            pass

        pass

    def get_hp_remain_reward(self):
        total_self_hero_hp = 0
        for hero_idx in range(self.self_hero_count):
            total_self_hero_hp += self.state[hero_idx][2]

        return total_self_hero_hp

    def get_harm_reward(self):
        harm_reward = 0.01#0.00002
        hero_death_penalty = 0.2
        total_harm_reward = 0
        for hero_idx in range(self.self_hero_count):
            if self.state[hero_idx][2] < self.last_state[hero_idx][2]:
                delta_hp = self.last_state[hero_idx][2] - self.state[hero_idx][2]
                total_harm_reward -= harm_reward * delta_hp

                if self.state[hero_idx][2] <= -0.5:
                    total_harm_reward -= hero_death_penalty

        start_feature_idx = self.self_hero_count * ONE_HERO_FEATURE_SIZE
        for hero_idx in range(self.oppo_hero_count):
            if self.state[0][start_feature_idx + hero_idx * ONE_HERO_FEATURE_SIZE + 2] < self.last_state[0][start_feature_idx + hero_idx * ONE_HERO_FEATURE_SIZE + 2]:
                delta_hp = self.last_state[0][start_feature_idx + hero_idx * ONE_HERO_FEATURE_SIZE + 2] - self.state[0][start_feature_idx + hero_idx * ONE_HERO_FEATURE_SIZE + 2]
                total_harm_reward += harm_reward * delta_hp


        return total_harm_reward

    def step(self, total_actions):
        #time.sleep(1)
        # action is 4-dimension vector
        player_count = '{}\n'.format(self.self_hero_count)
        self.proc.stdin.write(player_count.encode('utf-8'))
        self.proc.stdin.flush() 
        self.step_idx += 1
        for hero_idx in range(self.self_hero_count):
            action = total_actions[hero_idx]
            action_encode = 0
            isarr = isinstance(action[0], np.ndarray)
            action_0 = 0
            if isarr:
                action_0 = action[0][0]                
            else:
                action_0 = action[0]

            action_encode = (action_0 << 16)

            move_dir_encode = 0
            action_1 = 0
            isarr = isinstance(action[1], np.ndarray)
            if isarr:
                action_1 = action[1][0] 
            else:
                action_1 = action[1]

            if action_1 != -1:
                move_dir_encode = (action_1 << 12)    

            skill_dir_encode = 0
            action_2 = 0
            isarr = isinstance(action[2], np.ndarray)
            if isarr:
                action_2 = action[2][0] 
            else:
                action_2 = action[2]

            if action[2] != -1:
                skill_dir_encode = (action_2 << 8)
            
            vision_dir_encode = 0
            action_3 = 0
            isarr = isinstance(action[3], np.ndarray)
            if isarr:
                action_3 = action[3][0] 
            else:
                action_3 = action[3]

            if action[3] != -1:
                vision_dir_encode = (action_3 << 4)

            action_4 = 0
            isarr = isinstance(action[4], np.ndarray)
            if isarr:
                action_4 = action[4][0] 
            else:
                action_4 = action[4]

            if action[4] != -1:
                vision_dir_encode = vision_dir_encode + (action_4 << 2)

            encoded_action_val = (self.step_idx << 20) + action_encode + move_dir_encode + skill_dir_encode + vision_dir_encode
            self.proc.stdin.write('{}\n'.format(encoded_action_val).encode())
            self.proc.stdin.flush()
            #print(encoded_action_val)
            
        # Wait for response.
        # Parse the state
        while True:
            json_str = self.proc.stdout.readline()
            #print('When step, recv game process response {}'.format(json_str))
            if json_str == None or len(json_str) == 0:
                print('json_str == None or len(json_str) == 0')
                self.done = True
                self.reward = 0
                return [self.state, self.depth], self.reward, self.done, self.info

            try:
                str_json = json_str.decode("utf-8")
                
                parts = str_json.split('@')
                if int(parts[0]) != self.step_idx:
                    print('Step::We expect:{}, while getting {}'.format(self.step_idx, int(parts[0])))
                    continue

                jobj = json.loads(parts[1])
                #print("Spawn X: ", jobj["SelfHeroPosX"], " Y: ", jobj["SelfHeroPosY"])
                self.last_state[...] = self.state[...]
                self.fill_state(self.state, self.depth, jobj)

                if jobj['SelfWin'] != 0:
                    self.done = True
                    if 2 == jobj['SelfWin']:
                        self.reward = 0
                    elif -1 == jobj['SelfWin']:
                        self.reward = -1
                    else:
                        self.reward = 1#jobj['SelfWin']

                    # Add remain hp as reward
                    hp_remain_reward_coef = 0
                    hp_reward = self.get_hp_remain_reward()
                    self.reward += hp_remain_reward_coef * hp_reward
                    # self.reward += self.state[0][2]
                else:                    
                    self.reward = 0
                    harm_reward = self.get_harm_reward()

                    self.reward += harm_reward

                    self.done = False
                

                break
            except:
                print('Parsing json failed, terminate this game.')
                self.done = True
                self.reward = 0
                self.info.step_idx = self.step_idx
                return [self.state, self.depth], self.reward, self.done, self.info

        self.info.step_idx = self.step_idx
        
        return [self.state, self.depth], self.reward, self.done, self.info


    def reset(self):
        if 0 == self.step_idx:
            return [self.state, self.depth]
        self.restart_proc()
        # To avoid deadlocks: careful to: add \n to output, flush output, use
        # readline() rather than read()
        #self.proc.stdout.readline()
        self.proc.stdin.write(b'36864\n')
        self.proc.stdin.flush()

        while True:
            json_str = self.proc.stdout.readline()

            if json_str == None or len(json_str) == 0:
                continue
                raise Exception('When resetting env, json_str == None or len(json_str) == 0')

            try:
                str_json = json_str.decode("utf-8")
                parts = str_json.split('@')
                if int(parts[0]) != self.step_idx:
                    print('Reset::We expect:{}, while getting {}'.format(self.step_idx, int(parts[0])))
                    continue

                jobj = json.loads(parts[1])

                # Do the info check. 
                initial_self_hero_count = int(jobj['SelfHeroCount'])
                initial_oppo_hero_count = int(jobj['OppoHeroCount'])
                assert initial_self_hero_count == self.self_hero_count
                assert initial_oppo_hero_count == self.oppo_hero_count

                state_feature_count = (self.self_hero_count + self.oppo_hero_count) * ONE_HERO_FEATURE_SIZE

                self.state = np.zeros((self.self_hero_count, state_feature_count))
                self.depth = np.zeros((1, DEPTH_MAP_SIZE * DEPTH_MAP_SIZE))
                
                self.fill_state(self.state, self.depth, jobj)                
                self.last_state = copy.deepcopy(self.state)
                break

            except:
                print('When resetting env, parsing json failed.')
                continue
    
        return [self.state, self.depth]


    def render(self, mode='human'):
        pass

    def close(self):
        self.proc.stdin.close()
        self.proc.terminate()
        self.proc.wait(timeout=0.2)            
        pass
