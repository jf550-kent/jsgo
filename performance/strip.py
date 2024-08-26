import os

def extract_lines_from_file(file_path):
    with open(file_path, 'r') as file:
        for line in file:
            print(line)
            # if line.startswith('pkg:') or line.startswith('Bench'):
            #     print(line.strip())

# def process_directory(directory_path):
#     for root, dirs, files in os.walk(directory_path):
#         for file in files:
#             file_path = os.path.join(root, file)
#             extract_lines_from_file(file_path)

if __name__ == "__main__":
    # process_directory("./2024-08-16_19-46-21.txt")
    current_directory = os.getcwd() + "/performance/2024-08-16_19-46-21.txt"
    with open(current_directory, 'r') as file:
      line = file.readline()
      while line:
          if line.startswith('pkg:') or line.startswith('Bench'):
            print(line.strip())
          line = file.readline()
