import glob
import os

directory = os.getcwd()

requiredFileName = ["2024-08-05_22-23-57_NULL.txt", "2024-08-02_18-55-50_Assignment.txt", "2024-08-06_14-22-05.txt", "2024-08-06_15-28-32_Object_Key_Declaration.txt", "2024-08-15_14-35-56.txt", "2024-08-07_20-18-17_Bit_Operation.txt", "2024-08-02_19-21-02.txt ", "2024-07-28_14-07-40.txt ", "2024-08-15_17-35-58.txt", "2024-08-11_20-34-52.txt", "2024-08-15_14-33-19.txt"] 

def isNeededFile(fileName):
    for f in requiredFileName:
        if f == fileName:
            return True
    return False


headerTable = ["FileName", "Count", "Time (ns/op)", "Memory (B/op)", "Allocations (allocs/op)"]

records = []
for file_path in glob.glob(os.path.join(directory, '*.txt')):
    with open(file_path, 'r') as file:
        filename = os.path.basename(file_path)
        print(filename)
        if isNeededFile(filename):
            for line in file:
                if line.startswith('BenchmarkExample-4'):
                    result = line.split()
                    header = { "Name": 0, "Count": 1, "Time": 2, "Memory": 4, "Allocations": 6 }
                    record = [filename]
                    record.append(result[header["Count"]])
                    record.append(result[header["Time"]])
                    record.append(result[header["Memory"]])
                    record.append(result[header["Allocations"]])
                    records.append(record)
                    break

def format_markdown_table(header, rows):
    col_widths = [max(len(str(item)) for item in col) for col in zip(header, *rows)]
    header_row = "| " + " | ".join(f"{header[i]:{col_widths[i]}}" for i in range(len(header))) + " |"
    separator_row = "|-" + "-|-".join(f"{'-' * col_widths[i]}" for i in range(len(header))) + "-|"
    data_rows = "\n".join(
        "| " + " | ".join(f"{str(row[i]):{col_widths[i]}}" for i in range(len(row))) + " |"
        for row in rows
    )
    
    # Combine all parts into the final table
    return f"{header_row}\n{separator_row}\n{data_rows}"

print(format_markdown_table(headerTable, records))