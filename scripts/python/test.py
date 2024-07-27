import shutil
import subprocess
import os

os.chdir("src/kernel")
subprocess.run(f"go test duck/kernel/server/... -v", shell=True, check=True)
os.chdir("../..")
os.chdir("src/kernel")
subprocess.run(f"go test duck/kernel/persistence/... -v", shell=True, check=True)
