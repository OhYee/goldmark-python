import matplotlib.pyplot as plt
import matplotlib
import sys

matplotlib.use("svg")

y = [1, 2, 3, 4, 5]
x = [5, 4, 3, 2, 1]

plt.plot(x, y)
plt.savefig(sys.stdout)
