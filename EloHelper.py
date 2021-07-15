import numpy as np
import os
import json


class EloHelper:
    eloScoreBase = 1000
    eloIncrementStep = 16
    modelPath = "./model"
    updateThreshold = 0.6


    @staticmethod
    def getModelIdx():
        with open("{}/model_list.json".format(EloHelper.modelPath),'r') as load_f:
            load_dict = json.load(load_f)
            print("load_dict:", load_dict)
            return int(load_dict['Model'][-1])
    
    @staticmethod
    def writeModelList(model_info):
        with open("{}/model_list.json".format(EloHelper.modelPath),"w") as f:
            json.dump(model_info,f)

    @staticmethod
    def getScoreIncrement(ra, rb, sa):
        ea = 1 / (1 + np.power(10,(rb-ra)/400))
        # eb = 1 / (1 + np.power(10,(ra-rb)/400))
        return EloHelper.eloIncrementStep * (sa - ea)

    @staticmethod
    def getEloScore(model_idx, wins):
        with open("{}/model_list.json".format(EloHelper.modelPath),'r') as model_file:
            model_info_t = json.load(model_file)
        model_score = dict()
        for _idx in range(len(model_info_t["Model"])):
            model_score[model_info_t["Model"][_idx]] = model_info_t["Score"][_idx]

        score = EloHelper.eloScoreBase
        for _idx in range(len(wins)):
            score_increment = 0
            for _model_idx in model_idx[_idx]:
                score_increment += EloHelper.getScoreIncrement(score, model_score[_model_idx],  (wins[_idx] + 1) / 2)
            score += score_increment / len(model_idx[_idx])
        return score, model_info_t
