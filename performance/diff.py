import matplotlib.pyplot as plt
import numpy as np
import glob
import os


directory = os.getcwd() + "/performance/"

count_y = []
time_y = []
memory_y = []
allocations_y = []
fileName = []

def log(s):
  return np.log10(float(s))

for file_path in glob.glob(os.path.join(directory, '*.txt')):
    with open(file_path, 'r') as file:
        filename = os.path.basename(file_path)
        fileName.append(filename)
        for line in file:
          if line.startswith('BenchmarkExample-4'):
            result = line.split()
            header = { "Name": 0, "Count": 1, "Time": 2, "Memory": 4, "Allocations": 6 }
            count = result[header["Count"]]
            time = float(result[header["Time"]])
            memory = result[header["Memory"]]
            alloc = result[header["Allocations"]]

            count_y.append(count)
            time_y.append(time)
            memory_y.append(memory)
            allocations_y.append(alloc)
            break

time_changes = np.diff(time_y)

dot_changes = (time_changes / time_y[:-1])

plt.plot(dot_changes, marker='o', linestyle='-', color='b', label='Data')

# Adding labels and title
plt.xlabel('X Axis')
plt.title('Simple Plot')

# Adding grid
plt.grid(True)

# Adding legend
plt.legend()

# Displaying the plot
plt.show()