import shutil
import subprocess
import os

from clean import clean

platform = ""

if os.name == "nt":
    platform = "windows"
elif os.name == "posix":
    platform = "linux"
else:
    exit()

mode = "debug"

clean(mode)

os.environ["CGO_ENABLED"] = "1"

kernel_name = "kernel" + ".exe" if platform == "windows" else ""

os.chdir("src/kernel")
subprocess.run(f"go build -o {kernel_name}", shell=True, check=True)
os.chdir("../..")
shutil.move(f"./src/kernel/{kernel_name}", f"./build/{mode}/{kernel_name}")

os.chdir(f"build/{mode}")
subprocess.run(f"{kernel_name} -mode dev -port 9000 -dbPath ./build/debug/duck.db", shell=True, check=True)
