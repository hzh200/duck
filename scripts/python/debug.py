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

shutil.copyfile("src/app/index.html", f"build/{mode}/index.html")
subprocess.run(f"yarn tailwindcss -i src/app/interfaces/styles/globals.css -o ./build/{mode}/globals.css", shell=True, check=True)
subprocess.run(f"yarn webpack --config ./webpack.config.{mode}.js", shell=True, check=True)

# subprocess.run(["go", f"build -C ./src/kernel -o ../../build/{mode}/kernel"], shell=True, check=True)

kernel_name = "kernel" + ".exe" if platform == "windows" else ""

os.chdir("src/kernel")
subprocess.run(f"go build -o {kernel_name}", shell=True, check=True)
os.chdir("../..")
shutil.move(f"./src/kernel/{kernel_name}", f"./build/{mode}/{kernel_name}")

subprocess.run(f"yarn tsc ./src/utils/preload.ts --outdir ./build/{mode}", shell=True, check=True)
subprocess.run("yarn electron . --enable-logging --inspect", shell=True, check=True)
