import json
import os
from pathlib import Path

from google import genai
from google.genai import types

from schema import ReadmeUpdate
from utils.context import build_context

client = genai.Client(
    api_key=os.environ["GEMINI_API_KEY"]
)

analysis = json.loads(
    Path("analysis.json").read_text()
)

if not analysis["update"]:
    print("README update not required.")
    exit(0)

readme = Path("README.md").read_text()

prompt = Path(
    "scripts/prompts/edit.md"
).read_text()


context = build_context()

user_prompt = f"""
README
-----------------------
{readme}

Sections To Update
-----------------------
{", ".join(analysis["sections"])}

Reason
-----------------------
{analysis["reason"]}

Repository Tree
-----------------------
{context["tree"]}

Git Diff
-----------------------
{context["diff"]}

Changed Files
-----------------------
{chr(10).join(context["changed_files"])}

Source Files
-----------------------
"""

for file in context["files"]:
    user_prompt += f"""

FILE:
{file["path"]}

{file["content"]}

------------------------------------------
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
        response_schema=ReadmeUpdate,
    ),
)

updated = response.parsed.readme

if updated.strip() == readme.strip():
    print("README unchanged.")
    exit(0)

Path("README.md").write_text(
    updated,
    encoding="utf-8"
)

print("README updated successfully.")