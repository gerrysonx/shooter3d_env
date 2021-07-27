import numpy as np
from numpy.core.numeric import Inf

class Node:
    def __init__(self, id, pos, adj) -> None:
        self.id = id
        self.pos = pos
        self.adjNodes = adj
        self.dist = -1


class DistanceHelper:
    nodes = [Node(0, (0.2, -0.27),   [1]), 
             Node(1, (0.27, -0.27),  [2]), 
             Node(2, (0.27, -0.2),   [3]), 
             Node(3, (0.2, -0.2),    []),
             Node(4, (-0.2, -0.27),   [5]),
             Node(5, (-0.27, -0.27), [6]),
             Node(6, (-0.27, -0.2),  [7]),
             Node(7, (-0.2, -0.2),   [])]
    

    def __init__(self) -> None:
        self.calcDist()

    @staticmethod
    def calcDist():
        for node in DistanceHelper.nodes:
            node.dist = DistanceHelper.dfs(node)

    @staticmethod
    def dfs(node):
        if node.id in [3, 7]:
            return 0
        tmp = Inf
        for i in node.adjNodes:
            tmp = min(tmp, DistanceHelper.dfs(DistanceHelper.nodes[i]) + np.sqrt((node.pos[0]-DistanceHelper.nodes[i].pos[0])**2 + (node.pos[1]-DistanceHelper.nodes[i].pos[1])**2))
        return tmp

    @staticmethod
    def findMinDist(pos):
        tmp = Inf
        t_id = -1
        if pos[0] < 0.2 and pos[0] > 0 and pos[1] < -0.2:
            t_id = 0
        elif pos[0] > 0.2 and pos[1] < -0.27:
            t_id = 1
        elif pos[0] > 0.2 and pos[1] > -0.27 and pos[1] < -0.2:
            t_id = 2
        elif pos[0] > 0 and pos[1] > -0.2:
            t_id = 3

        elif pos[0] > -0.2 and pos[0] < 0 and pos[1] < -0.2:
            t_id = 4
        elif pos[0] < -0.2 and pos[1] < -0.27:
            t_id = 5
        elif pos[0] < -0.2 and pos[1] > -0.27 and pos[1] < -0.2:
            t_id = 6
        elif pos[0] < 0 and pos[1] > -0.2:
            t_id = 7

        return np.sqrt((pos[0] - DistanceHelper.nodes[t_id].pos[0])**2 + (pos[1] - DistanceHelper.nodes[t_id].pos[1])**2) + DistanceHelper.nodes[t_id].dist if t_id != -1 else 0

'''
disthelper = DistanceHelper()
for node in disthelper.nodes:
    print(node.dist)
print(disthelper.findMinDist((.2,-.2)))
print(np.sqrt(0.23**2 + 0.23**2))'''