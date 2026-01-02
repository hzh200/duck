import shutil
from debug_kernel import build_kernel
from lib.platform import *
from lib.build import clean, Mode

def debug():
    shutil.copyfile(f"{APP_PATH}/index.html", f"{BUILD_PATH}/{Mode.DEBUG.value}/index.html")
    system(f"npx tailwindcss -i {APP_PATH}/interfaces/styles/globals.css -o {BUILD_PATH}/{Mode.DEBUG.value}/globals.css")
    system(f"npx webpack --config ./webpack.config.{Mode.DEBUG.value}.js")

    # subprocess.run(["go", f"build -C ./src/kernel -o ../../build/{mode}/kernel"], shell=True, check=True)

    build_kernel()

    system(f"npx tsc {SRC_PATH}/utils/preload.ts --outdir {BUILD_PATH}/{Mode.DEBUG.value}")
    system("npx electron . --enable-logging --inspect")

if __name__ == "__main__":
    clean(Mode.DEBUG.value)
    debug()
