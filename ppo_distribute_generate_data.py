# coding=utf-8
"""
Some of code copy from
https://github.com/openai/baselines/blob/master/baselines/ppo1/pposgd_simple.py
"""

import sys
import copy
import numpy as np

from collections import deque
import gym, os
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
import tensorflow as tf
import random, cv2
import time
import math
import pickle
import json
import shutil
import random
from utils import LoadModel
from ppo import GetDataGeneratorAndTrainer
from tensorflow.python.lib.io import file_io
import boto3
from boto3.session import Session
from paramiko import SSHClient
from scp import SCPClient
import paramiko
from paramiko.py3compat import decodebytes
import base64


mybucket = "haytham_traindata__datamingplatform"

access_key = "vkimgpull-datamingplatform-f98d2d53"

secret_key = "vkimgpull-datamingplatform-6ab40a41"

url = "http://shlightspeedrgw.cephrados.so.db:7480"

session = Session(access_key, secret_key)
s3_client = session.client('s3', endpoint_url=url)
s3 = session.resource('s3', endpoint_url=url)
bucket = None
TRAINER_IP = None


g_step = 0
g_worker_id = 0
TIMESTEPS_PER_ACTOR_BATCH = 0

def init():
    global bucket
    global TRAINER_IP
    bucket = s3.Bucket(mybucket)
    try:
        os.mkdir('../distribute_collected_train_data')
    except:
        pass
    try:
        os.system('rm -r ../distribute_collected_train_data/*')
    except:
        pass
    try:
        os.system('rm -r ./model/*')
    except:
        pass
    for obj in bucket.objects.filter(Prefix = 'model'):
        if not os.path.isfile('./{}'.format(obj.key)):
            try:
                t_dir = obj.key.split('/')
                t_dir = '/'.join(t_dir[:-1])
                os.mkdir("./{}".format(t_dir))
            except:
                pass
            bucket.download_file(obj.key, './'+ obj.key)
        if obj.key == 'model/model_list.json':
            bucket.download_file(obj.key, './'+ obj.key)
    while True:
        for obj in bucket.objects.filter(Prefix = 'ip_addr'):
            TRAINER_IP = obj.key.split('/')[-1]
            print("host IP address: {}".format(TRAINER_IP))
            return
        print("host IP address not found")
        time.sleep(5)
        


def dump_generated_data_2_file(file_name, seg):
    with open(file_name, 'wb') as file_handle:
        pickle.dump(seg, file_handle)
        pass

def delete_outdated_folders(step):
    # Delete folders within range [0, step)
    root_folder = os.path.split(os.path.abspath(__file__))[0]
    if step < 1:
        return
    for idx in range(step):
        try:
            data_folder_path = '{}/../distribute_collected_train_data/{}'.format(root_folder, idx)
            shutil.rmtree(data_folder_path)
        except:
            pass
    pass

def generate_data(scene_id):

    root_folder = os.path.split(os.path.abspath(__file__))[0]
    _, data_generator, session = GetDataGeneratorAndTrainer(scene_id, TIMESTEPS_PER_ACTOR_BATCH)
    _step = g_step

    saver = tf.train.Saver(max_to_keep=1)

    sshClient = SSHClient()
    sshClient.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    sshClient.connect(TRAINER_IP)
    sshClient.exec_command('mkdir {}/../distribute_collected_train_data/{}'.format(root_folder, g_worker_id))
    scpInst = SCPClient(sshClient.get_transport())

    while True:
        while True:
            for obj in bucket.objects.filter(Prefix = 'ckpt/model_in_train'):
                try:
                    t_dir = obj.key.split('/')
                    t_dir = '/'.join(t_dir[:-1])
                    os.mkdir("./{}".format(t_dir))
                except :
                    pass
                bucket.download_file(obj.key, './'+ obj.key)
            model_file = '{}/ckpt/model_in_train/mnist.ckpt-0'.format(root_folder)

            try:
                saver.restore(session, model_file)
                print('restore file success:{}'.format(model_file))
                break
            except:
                print('restore file failed:{}, continue to try...'.format(model_file))
                time.sleep(30)
                continue

        while True:
            for obj in bucket.objects.filter(Prefix = 'model'):
                if not os.path.isfile('./{}'.format(obj.key)):
                    try:
                        t_dir = obj.key.split('/')
                        t_dir = '/'.join(t_dir[:-1])
                        os.mkdir("./{}".format(t_dir))
                    except:
                        pass
                    # if not os.path.exists(os.path.dirname('./'+ obj.key)):
                    #     os.makedirs(os.path.dirname('./'+ obj.key))
                    bucket.download_file(obj.key, './'+ obj.key)
                if obj.key == 'model/model_list.json':
                    bucket.download_file(obj.key, './'+ obj.key)
            if os.path.isfile('./model/model_list.json'):
                print('model list found, proceed')
                break
            else:
                print('model list not found, retry')
                time.sleep(5)
        
        seg = data_generator.get_one_step_data()

        data_folder_path = '{}/../distribute_collected_train_data/{}'.format(root_folder, g_worker_id)
        while True:
            if os.path.exists(data_folder_path):
                break
            else:
                try:
                    os.mkdir(data_folder_path)
                    # Need to delete all previous folders
                    delete_outdated_folders(_step)
                    break
                except:
                    print('mkdir {} failed, but we caught the exception.'.format(data_folder_path))
                    continue

        
        data_file_name = 'distribute_collected_train_data/{}/seg.data'.format(g_worker_id)
        dump_generated_data_2_file('{}/../{}'.format(root_folder, data_file_name), seg)
        print('Timestep:{}, generated data:{}.'.format(_step, data_file_name))
        fileExists = True
        while fileExists:
            _,_,stderr = sshClient.exec_command('ls -l {}/../distribute_collected_train_data/\
                                                {}/end.flag'.format(root_folder, g_worker_id))
            for line in stderr:
                if len(line) > 0:
                    fileExists = False
                    break
            time.sleep(1)
        scpInst.put('{}/../{}'.format(root_folder, data_file_name), '{}/../{}'.format(root_folder, data_file_name))
        sshClient.exec_command('touch {}/../distribute_collected_train_data/{}/end.flag'.format(root_folder, g_worker_id))
        _step += 1

if __name__=='__main__':
 
    g_step = 0
    scene_id = 10
    g_worker_id = random.randint(0,100000)
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
    generate_data(scene_id)
