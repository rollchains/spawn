import os
import re
import subprocess
import time

import requests

# current file location
curr_dir = os.path.dirname(os.path.realpath(__file__))
parent_dir = os.path.dirname(curr_dir)

# Commands which block the process
BLOCKING_START_COMMANDS = ["local-ic start", "make testnet", "make sh-testnet"]
ignore_commands = ["gh repo create"]

START_PID = -1
DEBUGGING = False

def main():
    demos = os.path.join(parent_dir, "versioned_docs", "version-v0.50.x", "03-demos")
    with open(os.path.join(demos, "03-tokenfactory.md"), "r") as f:
        parse_docs(f.read())

def poll_for_start(API_URL: str, pid=-1, waitSeconds=300):
    for i in range(waitSeconds + 1):
        try:
            requests.get(API_URL)
            return
        except Exception:
            # if i % 5 == 0:
            print(f"waiting for server to start (iter:{i}) ({API_URL})")
            time.sleep(1)

    if pid != -1:
        print(f"Killing process with pid for failing to start: {pid}")
        os.system(f"kill {pid}")

    raise Exception("Server did not start")

def clean_lines(text: str) -> list[str]:
    # remove comment lines
    sec = [line for line in text.split("\n") if not line.startswith("#")]

    for l in sec:
        l = l.strip()
        # if l starts with any in ignore-commands, remove the line from the section
        if any([l.startswith(cmd) for cmd in ignore_commands]):
            sec.remove(l)

    # remove blank lines
    return [line for line in sec if line.strip() != ""]

def get_env_variables_from_lines(lines: list[str]) -> dict[str, str]:
    env_vars = {}
    for line in lines:
        # DENOM=factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
        match = re.match(r"([A-Z_]+)=(.*)", line)
        if match:
            env_vars[match.group(1)] = match.group(2)

    return env_vars

def parse_docs(text: str):
    global START_PID

    # split the text by ```bash
    sections = text.split("```bash")
    for section in sections[1:]: # first section has nothing we want
        end = section.find("```")
        sec = clean_lines(section[:end])

        if len(sec) == 0:
            # print("debugging: empty (commands are ignored / there is nothing here)")
            continue

        secStr = "\n".join(sec)

        print("="*20, "\n", secStr, "\n", "="*20)

        # run sec as the os terminal
        if not DEBUGGING:
            # if sec contains anything from within BLOCKING_START_COMMANDS, it should run in its
            # own terminal

            # TODO: keep up with these globally and unset at end? so we dont polute test
            envs = get_env_variables_from_lines(sec)
            for k, v in envs.items():
                os.environ[k] = v
                print(f"{k}={v}")

            if any([cmd in sec for cmd in BLOCKING_START_COMMANDS]):
                START_PID = subprocess.Popen(secStr, shell=True).pid
                print(f"Started process with pid: {START_PID}")
                poll_for_start("http://127.0.0.1:26657", START_PID, waitSeconds=60)
            else:
                os.system(secStr)
                time.sleep(2)

            # input("--- Press Enter to continue...")

    # TODO: verification logic here
    # Kill / cleanup the instance (direct from pid would be nice, or just a killall <binary>)

if __name__ == '__main__':
    main()
