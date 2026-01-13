import sys
import shutil
from lib.platform import *
from lib.build import clean, Mode

kernel_name = "kernel" + selection_by_platform(windows_str=".exe", unix_str="")

def build_app():
    shutil.copyfile(f"{APP_PATH}/index.html", f"{BUILD_PATH}/{Mode.DEBUG.value}/index.html")
    system(f"npx tailwindcss -i {APP_PATH}/interfaces/styles/globals.css -o {BUILD_PATH}/{Mode.DEBUG.value}/globals.css")
    system(f"npx webpack --config {BASE_PATH}/webpack.config.{Mode.DEBUG.value}.js")
    system(f"npx tsc {SRC_PATH}/preload.ts --outdir {BUILD_PATH}/{Mode.DEBUG.value}")

def build_kernel():
    os.environ["CGO_ENABLED"] = "1"
    # os.chdir(KERNEL_PATH)
    # system(f"go build -o {kernel_name}")
    # shutil.move(f"{KERNEL_PATH}/{kernel_name}", f"{BUILD_PATH}/{Mode.DEBUG.value}/{kernel_name}")
    system(f"go build -C {KERNEL_PATH} -o {BUILD_PATH}/{Mode.DEBUG.value}/{kernel_name}")
    
def run_kernel():
    system(f"{BUILD_PATH}/{Mode.DEBUG.value}/{kernel_name} -mode dev -port 9000 -dbPath {BUILD_PATH}/{Mode.DEBUG.value}/duck.db")

def run_electron():
    system("npx electron . --enable-logging --inspect")

if __name__ == "__main__":
    if len(sys.argv) == 1:
        clean(Mode.DEBUG.value)
        build_app()
        build_kernel()
        run_electron()
        exit()

    if len(sys.argv) > 2 or sys.argv[1] not in ("kernel", "app"):
        print("Usage: python debug.py (kernel/app)")

    clean(Mode.DEBUG.value)
    if sys.argv[1] == "kernel":
        build_kernel()
        run_kernel()
    else:
        build_app()
