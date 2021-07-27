# coding=utf-8
"""
Some of code copy from
https://github.com/openai/baselines/blob/master/baselines/ppo1/pposgd_simple.py
"""

from paramiko import file
from EloHelper import EloHelper
import sys
import copy
import numpy as np

from collections import deque
import gym, os, glob
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
import tensorflow as tf
import random, cv2
import time
import math
import pickle
import json
from ppo import GetDataGeneratorAndTrainer
import boto3
from boto3.session import Session
import clear

mybucket = "haytham_traindata__datamingplatform"

access_key = "vkimgpull-datamingplatform-f98d2d53"

secret_key = "vkimgpull-datamingplatform-6ab40a41"

url = "http://shlightspeedrgw.cephrados.so.db:7480"

session = Session(access_key, secret_key)

s3_client = session.client('s3', endpoint_url=url)

s3 = session.resource('s3', endpoint_url=url)

bucket = None

g_step = 0

g_log_file_name = None
g_save_pb_model = True
g_new_training = True

TIMESTEPS_PER_ACTOR_BATCH = 0

def init():
    global bucket
    bucket = s3.Bucket(mybucket)
    try:
        os.mkdir("./ckpt/model_in_train")
    except:
        pass
    try:
        os.mkdir('../summary_log_distributed')
    except:
        pass
    try:
        os.mkdir('../distribute_collected_train_data')
    except:
        pass
    if g_new_training:
        try:
            os.system('rm -r ../distribute_collected_train_data/*')
        except:
            pass
        try:
            os.system('rm -r ./model/*')
        except:
            pass
        clear.clear.clear_model(mybucket)

    
def log_out(str_log):
    global g_log_file_name
    if g_log_file_name == None:
        root_folder = os.path.split(os.path.abspath(__file__))[0]
        g_log_file_name = '{}/../summary_log_distributed/train_at_{}.log'.format(root_folder, int(time.time()*100000))

    _handle = open(g_log_file_name, 'a')
    _handle.write(str_log)
    _handle.write('\n')
    _handle.close()

    pass


def get_one_step_data(timestep, work_thread_count):
    root_folder = os.path.split(os.path.abspath(__file__))[0]
    ob, ac, std_atvtg, tdlamret, lens, rets, unclipped_rets, news, depth, hidden_states, model_idx, wins, camps = [], [], [], [], [], [], [], [], [], [], [], [], []

    data_folder_path = '{}/../distribute_collected_train_data'.format(root_folder)
    i = 0
    while i < 15:
        for root, _, files in os.walk(data_folder_path):
            for file_name in files:
                full_file_name = '{}/{}'.format(root, file_name)
                if file_name == 'end.flag':
                    break
                else:
                    while True:
                        if os.path.exists('{}/end.flag'.format(root)):
                            break   
                        time.sleep(0.1)
                with open(full_file_name, 'rb') as file_handle:
                    _seg = pickle.load(file_handle)
                    ob.extend(_seg["ob"])
                    ac.extend(_seg["ac"])
                    std_atvtg.extend(_seg["std_atvtg"])
                    tdlamret.extend(_seg["tdlamret"])
                    lens.extend(_seg["ep_lens"])
                    rets.extend(_seg["ep_rets"])
                    unclipped_rets.extend(_seg["ep_unclipped_rets"])
                    news.extend(_seg["new"])
                    depth.extend(_seg["depth"])
                    model_idx.extend(_seg["model_idx"])
                    wins.extend(_seg["wins"])
                    camps.extend(_seg["camps"])
                    if 'hidden_states' in _seg:
                        hidden_states.extend(_seg["hidden_states"])
                os.remove(full_file_name)
                os.remove('{}/end.flag'.format(root))
                
                i += 1
        print('{} segments uploaded already, waiting...'.format(i))
        time.sleep(5)
    print('Successfully collected {} files, data size:{} from {}.'.format(i, len(ob), timestep))
    seg = {}
    seg["ob"] = np.array(ob)
    seg["ac"] = np.array(ac)
    seg["std_atvtg"] = np.array(std_atvtg)
    seg["tdlamret"] = np.array(tdlamret)
    seg["ep_lens"] = np.array(lens)
    seg["ep_rets"] = np.array(rets)
    seg["new"] = np.array(news)
    seg["ep_unclipped_rets"] = np.array(unclipped_rets)
    seg["hidden_states"] = np.array(hidden_states)
    seg['depth'] =  np.array(depth)
    seg["model_idx"] = np.array(model_idx)
    seg["wins"] = np.array(wins)
    seg["camps"] = np.array(camps)
    return seg


def learn(scene_id, num_steps):
    root_folder = os.path.split(os.path.abspath(__file__))[0]
    global g_step
    agent, _, session = GetDataGeneratorAndTrainer(scene_id, TIMESTEPS_PER_ACTOR_BATCH)

    saver = tf.train.Saver(max_to_keep=0)
    max_rew = -10000
    base_step = g_step
    #saver = tf.train.import_meta_graph('{}/../ckpt/mnist.ckpt-191.meta'.format(root_folder))
    for obj in bucket.objects.filter(Prefix = 'ckpt/model_in_train'):
        bucket.download_file(obj.key, './'+ obj.key)
    model_file = '{}/ckpt/model_in_train/mnist.ckpt-0'.format(root_folder)
    try:
        saver.restore(session, model_file)
        print('restored file successfully:{}'.format(model_file))
    except:
        print('restoring file failed, continue')
    
    while True:
        for obj in bucket.objects.filter(Prefix = 'model'):
            if obj.key == 'model/model_list.json':
                bucket.download_file(obj.key, './'+ obj.key)
        if os.path.isfile('./model/model_list.json'):
            print('model list found, proceed')
            break
        else:
            print('model list not found, retry')
            time.sleep(5)

    for timestep in range(num_steps):
        #g_step = base_step + timestep
        seg = get_one_step_data(timestep, g_data_generator_count)

        if ((sum(seg['wins']) / len(seg['wins'])) + 1) / 2 > EloHelper.updateThreshold:
            score, model_info_t = EloHelper.getEloScore(seg['model_idx'], seg['wins'])
            if g_save_pb_model:
                idx = str(int(model_info_t['Model'][-1]) + 1).zfill(4)
                tf.saved_model.simple_save(session,
                        "./model/model_{}".format(idx),
                        inputs={"input_state":agent.multi_s, "input_depth":agent.multi_d, "input_camp":agent.multi_c},
                        outputs={"output_policy_0": agent.a_policy_new[0][0], "output_policy_1": agent.a_policy_new[0][1], "output_policy_2": agent.a_policy_new[0][2], 
                        "output_policy_3": agent.a_policy_new[0][3], "output_policy_4": agent.a_policy_new[0][4], 
                        "output_value":agent.value[0]}) 
                saver.save(session,'ckpt/model_{}/mnist.ckpt'.format(idx), global_step=g_step)
                for root, _, files in os.walk('{}/ckpt/model_{}'.format(root_folder, idx)):
                    for file_name in files:
                        full_file_name = '{}/{}'.format(root, file_name)
                        response = s3.meta.client.upload_file(full_file_name, mybucket, 'ckpt/model_{}/{}'.format(idx, file_name), ExtraArgs={'ACL':'public-read'})
                        if response != None:
                            print(response)
                for root, _, files in os.walk('{}/model/model_{}'.format(root_folder, idx)):
                    for file_name in files:
                        index = root.find('model_{}'.format(idx))
                        response = s3.meta.client.upload_file('{}/{}'.format(root, file_name), mybucket, 'model/{}/{}'.format(root[index:], file_name), ExtraArgs={'ACL':'public-read'})
                        if response != None:
                            print(response)
                model_info_t["Model"].append(idx)
                model_info_t["Score"].append(int(round(score)))
                EloHelper.writeModelList(model_info_t)
                response = s3.meta.client.upload_file('{}/model_list.json'.format(EloHelper.modelPath), mybucket, 'model/model_list.json', ExtraArgs={'ACL':'public-read'})
                if response != None:
                    print(response)

        entropy, kl_distance = agent.learn_one_traj(timestep, seg)

        max_rew = max(max_rew, np.max(agent.unclipped_rewbuffer))

        saver.save(session, '{}/ckpt/model_in_train/mnist.ckpt'.format(root_folder), global_step = g_step)
        for root, _, files in os.walk('{}/ckpt/model_in_train'.format(root_folder)):
            for file_name in files:
                full_file_name = '{}/{}'.format(root, file_name)
                response = s3.meta.client.upload_file(full_file_name, mybucket, 'ckpt/model_in_train/{}'.format(file_name), ExtraArgs={'ACL':'public-read'})
                if response != None:
                    print(response)
        
        str_log = 'Timestep:{}\tEpLenMean:{}\tEpRewMean:{}\tUnClippedEpRewMean:{}\tMaxUnClippedRew:{}\tEntropy:{}\tKL_distance:{}\tWinning Rate:{}'.format(timestep,
        '%.3f'%np.mean(agent.lenbuffer),
        '%.3f'%np.mean(agent.rewbuffer),
        '%.3f'%np.mean(agent.unclipped_rewbuffer),
        max_rew,
        '%.3f'%entropy,
        '%.8f'%kl_distance, 
        '%.3f'%(((sum(seg['wins']) / len(seg['wins'])) + 1) / 2))
        log_out(str_log)

        print(str_log)

if __name__=='__main__':
    scene_id = 10
    g_step = 0
    g_data_generator_count = 1

    TIMESTEPS_PER_ACTOR_BATCH = 2048
    if len(sys.argv) > 1:
        TIMESTEPS_PER_ACTOR_BATCH = int(sys.argv[1])


    if len(sys.argv) > 2:
        mybucket = sys.argv[2]

    if len(sys.argv) > 3:
        g_step = int(sys.argv[3])

    if len(sys.argv) > 4:
        scene_id = int(sys.argv[4])

    my_env = os.environ
    my_env['moba_env_is_train'] = 'True'
    my_env['moba_env_scene_id'] = '{}'.format(scene_id)

    init()
    learn(scene_id, num_steps=1000)