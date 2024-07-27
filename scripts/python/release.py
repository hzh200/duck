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

mode = "release"

clean(mode)
