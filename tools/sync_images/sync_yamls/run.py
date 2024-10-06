import subprocess
import time


# Function to execute a batch of commands
def execute_batch(commands, start, batch_size=5):
    for i in range(start, min(start + batch_size, len(commands))):
        command = commands[i]
        try:
            print(f"Executing command {i+1}/{len(commands)}: {command}")
            subprocess.run(command, shell=True, check=True)
        except subprocess.CalledProcessError as e:
            print(f"Command failed with error: {e}")
        time.sleep(1)  # Adding a delay to avoid overloading the system

# Read commands from file
with open('sync_images.sh', 'r') as file:
    commands = file.readlines()


# Strip any extra whitespace/newlines
commands = [cmd.strip() for cmd in commands]


# Execute commands in batches of 4
batch_size = 4
for start in range(0, len(commands), batch_size):
    execute_batch(commands, start, batch_size)
    print("Batch executed, waiting before next batch...")
    time.sleep(5)  # Adding a delay between batches to avoid overloading the system