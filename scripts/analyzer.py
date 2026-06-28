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

prompt = Path("scripts/prompts/analyze.md").read_text()

from utils.context import build_context

context = build_context()


user_prompt = f"""
README
-----------------------
{context["readme"]}

COMMIT MESSAGE
-----------------------
{context["commit"]}

REPOSITORY TREE
-----------------------
{context["tree"]}

GIT DIFF
-----------------------
{context["diff"]}

CHANGED FILES
-----------------------
{chr(10).join(context["changed_files"])}

SOURCE FILES
-----------------------

"""
for file in context["files"]:
    user_prompt += f"""

FILE: {file["path"]}

{file["content"]}

--------------------------------------------
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