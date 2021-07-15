# coding=utf-8
import sys, os
import time
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
scene_id = 10

import subprocess


if __name__=='__main__':
    TIMESTEPS_PER_ACTOR_BATCH = 2048
    root_folder = os.path.split(os.path.abspath(__file__))[0]

    horizon_total = 8192*224
    horizon_per_worker = 8192*8
    worker_count = horizon_total / horizon_per_worker

    train_full_path = '{}/ppo_trainer_local.py'.format(root_folder)
    generate_full_path = '{}/ppo_generator_local.py'.format(root_folder)
    my_env = os.environ.copy()
    my_env['TF_CPP_MIN_LOG_LEVEL'] = '3'
    
    global_step = 0
    scene_id = 10
    
    time_stamp = int(time.time()*1000)
    train_log_full_path = '{}/../summary_log_distributed/train_{}.log'.format(root_folder, time_stamp)
    file_handle = open(train_log_full_path, 'wb')
    subprocess.Popen(['python', train_full_path, 
    '{}'.format(worker_count), 
    '{}'.format(global_step), 
    '{}'.format(scene_id)], env=my_env)

    for i in range(worker_count):
        worker_log_full_path = '{}/../summary_log_distributed/worker_{}_{}.log'.format(root_folder, time_stamp, i)
        file_handle = open(worker_log_full_path, 'wb')        
        subprocess.Popen(['python', generate_full_path, 
        '{}'.format(horizon_per_worker), 
        '{}'.format(i), 
        '{}'.format(global_step),
        '{}'.format(scene_id)], stdout=file_handle, stderr=file_handle, env=my_env, bufsize=1)
        
    while True:
        time.sleep(10)
