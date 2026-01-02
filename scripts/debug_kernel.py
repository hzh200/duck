import os
import shutil
from lib.platform import *
from lib.build import clean, Mode

kernel_name = "kernel" + selection_by_platform(".exe", "")

def build_kernel():
    os.environ["CGO_ENABLED"] = "1"
    os.chdir(KERNEL_PATH)
    system(f"go build -o {kernel_name}")
    shutil.move(f"{KERNEL_PATH}/{kernel_name}", f"{BUILD_PATH}/{Mode.DEBUG.value}/{kernel_name}")

if __name__ == "__main__":
    clean(Mode.DEBUG.value)
    build_kernel()
    system(f"{BUILD_PATH}/{Mode.DEBUG.value}/{kernel_name} -mode dev -port 9000 -dbPath {BUILD_PATH}/{Mode.DEBUG.value}/duck.db")
