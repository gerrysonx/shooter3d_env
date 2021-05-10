from gym.envs.registration import register


register(
	id='valorant-multiplayer-v0',
	entry_point='gym_valorant.envs:ValorantMultiPlayerEnv',
)
