from enum import Enum
import os
import shutil
from lib.platform import BUILD_PATH

class Mode(Enum):
    DEBUG = "debug"
    RELEASE = "release"

def clean(mode: Mode):
    target_path = os.path.abspath(f"{BUILD_PATH}/{mode}")

    if not os.path.exists(target_path):
        os.makedirs(target_path)

    for sub_name in os.listdir(target_path):
        sub_path = os.path.join(target_path, sub_name)
        if os.path.isdir(sub_path):
            shutil.rmtree(sub_path)
        else:
            os.remove(sub_path)
