import json
import subprocess
import os
from pathlib import Path

from google import genai
from google.genai import types

from schema import Analysis


client = genai.Client(api_key=os.environ["GEMINI_API_KEY"])


def run(cmd):
    return subprocess.check_output(cmd).decode("utf-8")


readme = Path("README.md").read_text()

prompt = Path("scripts/prompts/analyze.md").read_text()

diff = run(["git", "diff", "HEAD~1", "HEAD"])

changed_files = run([
    "git",
    "diff",
    "--name-only",
    "HEAD~1",
    "HEAD"
])

commit_message = run([
    "git",
    "log",
    "-1",
    "--pretty=%B"
])

tree = run([
    "git",
    "ls-tree",
    "-r",
    "--name-only",
    "HEAD"
])


user_prompt = f"""
README
-----------------------
{readme}

COMMIT MESSAGE
-----------------------
{commit_message}

CHANGED FILES
-----------------------
{changed_files}

REPOSITORY TREE
-----------------------
{tree}

GIT DIFF
-----------------------
{diff}
"""


response = client.models.generate_content(
    model="gemini-2.5-flash",
    contents=[
        prompt,
        user_prompt,
    ],
    config=types.GenerateContentConfig(
        temperature=0,
        response_mime_type="application/json",
        response_schema=Analysis,
    ),
)


analysis = response.parsed

print()

print("=" * 40)
print("README ANALYSIS")
print("=" * 40)

print("Update :", analysis.update)
print("Reason :", analysis.reason)

if analysis.sections:
    print("\nSections")
    for s in analysis.sections:
        print("-", s)

print("=" * 40)