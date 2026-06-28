from pathlib import Path
import subprocess


def run(cmd):
    return subprocess.check_output(cmd).decode("utf-8")


def read_file(path, max_chars=15000):
    try:
        text = Path(path).read_text(encoding="utf-8")
        return text[:max_chars]
    except:
        return None

def get_changed_files():
    out = run([
        "git",
        "diff",
        "--name-only",
        "HEAD~1",
        "HEAD"
    ])

    return [f for f in out.splitlines() if f]


def build_context():
    context = {}

    context["readme"] = read_file("README.md")

    context["diff"] = run([
        "git",
        "diff",
        "HEAD~1",
        "HEAD"
    ])

    context["commit"] = run([
        "git",
        "log",
        "-1",
        "--pretty=%B"
    ])

    context["tree"] = run([
        "git",
        "ls-tree",
        "-r",
        "--name-only",
        "HEAD"
    ])

    context["changed_files"] = get_changed_files()

    context["files"] = []

    for file in context["changed_files"]:
        text = read_file(file)

        if text:
            context["files"].append({
                "path": file,
                "content": text
            })

    return context