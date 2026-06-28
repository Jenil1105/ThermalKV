import json
import subprocess
from pathlib import Path

from openai import OpenAI


client = OpenAI()


def run(cmd):
    return subprocess.check_output(cmd).decode("utf-8")


# Read README
readme = Path("README.md").read_text(encoding="utf-8")

# Read prompt
prompt = Path("scripts/prompts/analyze.md").read_text()

# Git diff
diff = run(["git", "diff", "HEAD~1", "HEAD"])

# Changed files
changed_files = run([
    "git",
    "diff",
    "--name-only",
    "HEAD~1",
    "HEAD",
])

# Commit message
commit_message = run([
    "git",
    "log",
    "-1",
    "--pretty=%B",
])

# Repository tree
try:
    tree = run(["git", "ls-tree", "-r", "--name-only", "HEAD"])
except:
    tree = ""

user_prompt = f"""
README

{readme}

------------------------

COMMIT MESSAGE

{commit_message}

------------------------

CHANGED FILES

{changed_files}

------------------------

REPOSITORY TREE

{tree}

------------------------

GIT DIFF

{diff}
"""

response = client.responses.create(
    model="gpt-4.1",
    input=[
        {
            "role": "system",
            "content": prompt,
        },
        {
            "role": "user",
            "content": user_prompt,
        },
    ],
)

text = response.output_text.strip()

print("\n===== AI RESPONSE =====\n")
print(text)

try:
    result = json.loads(text)

    print("\nDecision :", result["update"])
    print("Reason   :", result["reason"])

    if result["sections"]:
        print("Sections :")
        for section in result["sections"]:
            print(" -", section)

except Exception:
    print("AI returned invalid JSON.")