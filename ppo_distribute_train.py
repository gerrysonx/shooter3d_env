# coding=utf-8
"""
Some of code copy from
https://github.com/openai/baselines/blob/master/baselines/ppo1/pposgd_simple.py
"""

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

mybucket = "haytham_traindata__datamingplatform"

access_key = "vkimgpull-datamingplatform-f98d2d53"

secret_key = "vkimgpull-datamingplatform-6ab40a41"

url = "http://shpublicrgw.cephrados.so.db:7480"

session = Session(access_key, secret_key)

s3_client = session.client('s3', endpoint_url=url)

s3 = session.resource('s3', endpoint_url=url)

bucket = s3.Bucket('haytham_traindata__datamingplatform')

g_step = 0

g_log_file_name = None

TIMESTEPS_PER_ACTOR_BATCH = 0

jobid = 0

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
    ob, ac, std_atvtg, tdlamret, lens, rets, unclipped_rets, news, depth, hidden_states = [], [], [], [], [], [], [], [], [], []
    # Enumerate data files under folder
    files = glob.glob('../distribute_collected_train_data/*')
    for f in files:
        os.remove(f)
    data_folder_path = '{}/../distribute_collected_train_data'.format(root_folder, timestep)
    collected_all_data_files = False
    while not collected_all_data_files:
        objectlist =  bucket.objects.filter(Prefix = 'distribute_collected_train_data')
        i = 0
        for obj in objectlist:
            try:
                bucket.download_file(obj.key, '../distribute_collected_train_data/{}.data'.format(i))
                i += 1
            except:
                continue
        for root, _, files in os.walk(data_folder_path):
            '''
            if len(files) < work_thread_count:
                print('Already has {} files, waiting for the worker thread generate more data files.'.format(len(files)))
                break'''

            for file_name in files:
                full_file_name = '{}/{}'.format(root, file_name)
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
                    if 'hidden_states' in _seg:
                        hidden_states.extend(_seg["hidden_states"])
            collected_all_data_files = True
            break
        time.sleep(10)
    print('Successfully collected {} files, data size:{} from {}.'.format(len(files), len(ob), timestep))
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
    return seg


def learn(scene_id, num_steps):
    root_folder = os.path.split(os.path.abspath(__file__))[0]
    global g_step
    agent, _, session = GetDataGeneratorAndTrainer(scene_id, TIMESTEPS_PER_ACTOR_BATCH)

    saver = tf.train.Saver(max_to_keep=1)
    max_rew = -10000
    base_step = g_step
    #saver = tf.train.import_meta_graph('{}/../ckpt/mnist.ckpt-191.meta'.format(root_folder))
    for obj in bucket.objects.filter(Prefix = 'ckpt'):
        bucket.download_file(obj.key, '../'+ obj.key)
    model_file = '{}/../ckpt/mnist.ckpt-0'.format(root_folder)
    try:
        saver.restore(session, model_file)
        print('restore file success:{}'.format(model_file))
    except:
        print('restore file failed, continue')
    objectlist =  bucket.objects.filter(Prefix = 'distribute_collected_train_data')
    for obj in objectlist:
        s3.Object(mybucket, obj.key).delete()
    for timestep in range(num_steps):
        #g_step = base_step + timestep
        seg = get_one_step_data(timestep, g_data_generator_count)

        entropy, kl_distance = agent.learn_one_traj(timestep, seg)

        max_rew = max(max_rew, np.max(agent.unclipped_rewbuffer))

        saver.save(session, '{}/../ckpt/mnist.ckpt'.format(root_folder), global_step = g_step)
        for root, _, files in os.walk('{}/../ckpt/'.format(root_folder)):
            for file_name in files:
                full_file_name = '{}/{}'.format(root, file_name)
                print(full_file_name)
                print(file_name)
                response = s3.meta.client.upload_file(full_file_name, mybucket, 'ckpt/{}'.format(file_name))
                print(response)
        str_log = 'Timestep:{}\tEpLenMean:{}\tEpRewMean:{}\tUnClippedEpRewMean:{}\tMaxUnClippedRew:{}\tEntropy:{}\tKL_distance:{}'.format(timestep,
        '%.3f'%np.mean(agent.lenbuffer),
        '%.3f'%np.mean(agent.rewbuffer),
        '%.3f'%np.mean(agent.unclipped_rewbuffer),
        max_rew,
        '%.3f'%entropy,
        '%.8f'%kl_distance)
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
        jobid = int(sys.argv[2])

    if len(sys.argv) > 3:
        g_step = int(sys.argv[3])

    if len(sys.argv) > 4:
        scene_id = int(sys.argv[4])

    my_env = os.environ
    my_env['moba_env_is_train'] = 'True'
    my_env['moba_env_scene_id'] = '{}'.format(scene_id)

    learn(scene_id, num_steps=500)