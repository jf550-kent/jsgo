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
            time = result[header["Time"]]
            memory = result[header["Memory"]]
            alloc = result[header["Allocations"]]

            count_y.append(log(count))
            time_y.append(time)
            memory_y.append(memory)
            allocations_y.append(alloc)
            break
figure, axis = plt.subplots(2, 2, figsize=(12, 4))

axis[0, 0].plot(count_y, marker='o', linestyle='-', color='b')
axis[0, 0].set_title("Count")

axis[0, 1].plot(time_y, marker='o', linestyle='-', color='b')
axis[0, 1].set_title("Time taken for each operation")

# For Tangent Function
axis[1, 0].plot(memory_y, marker='o', linestyle='-', color='b')
axis[1, 0].set_title("Mo")

# For Tanh Function
axis[1, 1].plot(allocations_y, marker='o', linestyle='-', color='b')
axis[1, 1].set_title("Tanh Function")

for ax in axis.flatten():
    ax.grid(True)

# Adjust layout to prevent overlap
plt.tight_layout()

# Display the plot
plt.show()


import numpy as np

# Sample data
values = [100, 150, 120, 180]

# Compute the differences
changes = np.diff(values)

# Compute the percentage changes
percentage_changes = (changes / values[:-1]) 

print("Values:", values)
print("Changes:", changes)
print("Percentage Changes:", percentage_changes)

