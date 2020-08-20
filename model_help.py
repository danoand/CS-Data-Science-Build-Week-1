# Imports
from math import sqrt

# calc_vct_dist calculates the Euclidean space between two vectors
def calc_vct_dist(vct1, vct2):
    dist = 0.0      # initialize a distance variable
    for idx in range(len(vct1)-1):
        dist = dist + (vct1[idx] - vct2[idx])**2   # sum the squared elm wise diffs

    return sqrt(dist)

# Determine most similar neighbors (closes neighbors in the space)
def det_nghbrs(space, tst_vctr, nnghbrs):
    dists = []

    # Iterate through each space vector
    #   and determine the distance from the given vector
    for vctr in space:
        # Construct a list of distances
        dists.append(calc_vct_dist(tst_vctr, vctr))

    dists.sort()

    nghbors = []

    for i in range(nnghbrs):
        nghbors.append(dists[i])

    return nghbors

# Test distance function
dataset = [[2.7810836,2.550537003,0],
	[1.465489372,2.362125076,0],
	[3.396561688,4.400293529,0],
	[1.38807019,1.850220317,0],
	[3.06407232,3.005305973,0],
	[7.627531214,2.759262235,1],
	[5.332441248,2.088626775,1],
	[6.922596716,1.77106367,1],
	[8.675418651,-0.242068655,1],
	[7.673756466,3.508563011,1]]

my_neighbors = det_nghbrs(dataset, dataset[0], 3)
for tmp_nbr in my_neighbors:
    print(tmp_nbr)


