import os
import shutil
import sys

def clean(mode=""):
    if mode != "release" and mode != "debug":
        exit()

    target_path = os.path.abspath(f"./build/{mode}")

    if not os.path.exists(target_path):
        os.makedirs(target_path)

    for sub_name in os.listdir(target_path):
        sub_path = os.path.join(target_path, sub_name)
        if os.path.isdir(sub_path):
            shutil.rmtree(sub_path)
        else:
            os.remove(sub_path)

if __name__ == "_main__":
    mode = sys.argv[1]
    clean(mode)
