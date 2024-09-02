# Define the table header and rows
header = ["FileName", "Count", "Time (ns/op)", "Memory (B/op)", "Allocations (allocs/op)"]
rows = [
    ["-2024-07-25_11-33-40.txt", 16716253, 71.09, 8, 1]
]

# Function to format the table
def format_markdown_table(header, rows):
    # Determine the width of each column
    col_widths = [max(len(str(item)) for item in col) for col in zip(header, *rows)]
    
    # Create the header row
    header_row = "| " + " | ".join(f"{header[i]:{col_widths[i]}}" for i in range(len(header))) + " |"
    
    # Create the separator row
    separator_row = "|-" + "-|-".join(f"{'-' * col_widths[i]}" for i in range(len(header))) + "-|"
    
    # Create the data rows
    data_rows = "\n".join(
        "| " + " | ".join(f"{str(row[i]):{col_widths[i]}}" for i in range(len(row))) + " |"
        for row in rows
    )
    
    # Combine all parts into the final table
    return f"{header_row}\n{separator_row}\n{data_rows}"

# Format and print the Markdown table
markdown_table = format_markdown_table(header, rows)
print(markdown_table)
