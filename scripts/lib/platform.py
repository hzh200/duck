import os
import subprocess

BASE_PATH = os.path.abspath(os.path.join(os.path.dirname(__file__), "../.."))

SRC_PATH = os.path.join(BASE_PATH, "src")
APP_PATH = os.path.join(SRC_PATH, "app")
KERNEL_PATH = os.path.join(SRC_PATH, "kernel")
SCRIPTS_PATH = os.path.join(BASE_PATH, "scripts")
BUILD_PATH = os.path.join(BASE_PATH, "build")

def selection_by_platform(str1: str, str2: str) -> str:
    if os.name == "posix":
        return str1
    else: # "nt"
        return str2

def combine_path(*args: tuple[str]) -> str:
    return os.path.join(*args)

def system(cmd: str) -> None:
    subprocess.run(cmd, shell=True, check=True)
