import networkx as nx
import numpy as np
import string
from networkx.drawing.nx_agraph import graphviz_layout
from numpy import genfromtxt
import matplotlib
import matplotlib.pyplot as plt
import csv

matplotlib.use("Agg")

A = genfromtxt('graph/matrix.csv', delimiter=',')

dt = [('len', float)]
A = A.view(dt)
G = nx.from_numpy_matrix(A)

datafile = open('graph/path.csv', 'r')
myreader = csv.reader(datafile)

r = [int(i) for i in list(myreader)[0]]
routes=[r]
edges = []
for r in routes:
    route_edges = [(r[n],r[n+1],r[n]) for n in range(len(r)-1)]
    G.add_nodes_from(r)
    G.add_weighted_edges_from(route_edges)
    edges.append(route_edges)


pos = nx.spring_layout(G)
nx.draw_networkx_nodes(G,pos=pos)
nx.draw_networkx_labels(G,pos=pos)
colors = ['r']
linewidths = [5]
for ctr, edgelist in enumerate(edges):
    nx.draw_networkx_edges(G,pos=pos,edgelist=edgelist,edge_color = colors[ctr], width=linewidths[ctr])

nx.draw(G, pos)
plt.savefig("graph/graph.png")
