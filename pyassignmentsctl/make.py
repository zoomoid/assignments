#!/usr/local/bin/python3

from glob import glob
from os import environ, mkdir
from shutil import copy2
import subprocess
from pathlib import Path
import argparse
import configparser
import logging
import requests
import os.path

CONFIGMAP_PATH = Path(".", "./.assignments.rc")

ASSIGNMENTS_CLASS_URL = "https://gist.githubusercontent.com/zoomoid/df2f5687c59f83d32f5169927321ebfb/raw/88496bd0f8fd6d19a7063581bd0c3c255b307e12/assignments.cls"

ASSIGNMENT_TEMPLATE = """\
\documentclass{{../assignments}}
\course{{{course}}}
\group{{{group}}}
{members}
\\title{{}}
\\author{{}}
\date{{}}
\sheet{{{sheet}}}
\due{{{due}}}

\\begin{{document}}
\maketitle
\gradingtable{{}}

% \exercise[<number of points>]{{<Exercise Title>}}
% \subexercise{{<Subexercise Title>}}

\end{{document}}
"""


def has_configmap() -> bool:
    return CONFIGMAP_PATH.is_file()


def configmap_is_valid():
    config = configparser.ConfigParser()
    config.read(CONFIGMAP_PATH)
    return (
        "general" in config
        and "course" in config["general"]
        and "assignments" in config
        and "number" in config["assignments"]
    ), config


def bootstrap(args):
    if Path("./.assignments.rc").is_file():
        print("❌ Directory is already bootstrapped, clean up before trying again!")
        exit(1)
    class_path = "./assignments.cls"
    if not Path(class_path).is_file():
        print(
            f"⬇️ Downloading 'assignments' class from {ASSIGNMENTS_CLASS_URL}")
        cls = requests.get(ASSIGNMENTS_CLASS_URL)
        with open(class_path, "w") as f:
            f.write(cls.text)
    check_configmap()
    config = configparser.ConfigParser()
    config["general"] = {}
    if "course" in args and args.course:
        config["general"]["course"] = str(args.course)
    else:
        course = input("❓ Please enter the course's name: ")
        config["general"]["course"] = str(course)

    if "group" in args and args.group:
        config["general"]["group"] = str(args.group)
    else:
        group = input(
            "❓ Please enter a group name (or leave empty to use nothing): ")
        config["general"]["group"] = str(group)

    if "members" in args and args.members:
        config["members"] = {}
        for member in args.members:
            matrnr, name = transform_member_to_configmap(member)
            config["members"][matrnr] = name
    else:
        members = []
        while True:
            member_or_abort = input(
                "❓ Please enter a group member's name followed by their immatriculation number (e.g. Max Mustermann,123456), or press 'q' to move on: "
            )
            if member_or_abort == "q":
                break
            else:
                members.append(member_or_abort)
        config["members"] = {}
        for member in members:
            matrnr, name = transform_member_to_configmap(member)
            config["members"][matrnr] = name

    config["assignments"] = {}
    config["assignments"]["number"] = str(0)

    with open(CONFIGMAP_PATH, "w") as configfile:
        config.write(configfile)
    print(f"✅ Bootstrapped {config['general']['course']} course.")


def transform_member_to_configmap(member_str):
    name, matrnr = str(member_str).split(",")
    return str(matrnr.strip()), str(name.strip())


def transform_member_to_latex(matrnr, name):
    *firstname, lastname = str(name).split()
    firstname = " ".join(firstname)
    return f"\member{{{matrnr}}}{{{firstname}}}{{{lastname}}}"


def members_to_latex(config):
    member_list = list(
        zip(list(config["members"].keys()), list(config["members"].values()))
    )
    members = []
    for matrnr, name in member_list:
        members.append(transform_member_to_latex(matrnr, name))
    return "\n".join(members)


def check_configmap():
    is_valid, config = configmap_is_valid()
    if not has_configmap() and not is_valid:
        number = 1
        if "number" in args and args.number:
            number = int(args.number)
        if not "noninteractive" in args and not "number" in args:
            number = input("❓ Please provide an assignment number: ")
            if number.startswith("0") and len(number) > 1:
                number.removeprefix("0")
        config = configparser.ConfigParser()
        config["assignments"] = {}
        config["assignments"]["number"] = str(number)
        with open(CONFIGMAP_PATH, "w") as configfile:
            config.write(configfile)
            logging.info("Created configmap")
        return int(number)
    else:
        if "number" in args and args.number:
            return args.number
        return int(config["assignments"]["number"]) + 1


def write_configmap(number):
    config = configparser.ConfigParser()
    config.read(CONFIGMAP_PATH)
    config["assignments"]["number"] = str(number)
    with open(CONFIGMAP_PATH, "w") as configfile:
        config.write(configfile)


def add_leading_zero(number):
    if int(number) >= 10:
        return str(number)
    return "0" + str(int(number))


def generate_assignment(args):
    print("Generating new assignment")
    number = check_configmap()
    h = Path(".")
    dir_name = f"assignment-{add_leading_zero(number)}"
    due_date = input(
        "⏱️  When is the assignment due? (e.g.,'April 20, 2021): ")

    config = configparser.ConfigParser()
    config.read(CONFIGMAP_PATH)

    group = None
    if "group" in config["general"]:
        group = config["general"]["group"]
    if "group" in args:
        group = args.group

    members = ""
    if "members" in config:
        members = members_to_latex(config)

    course = config["general"]["course"]

    if (h / dir_name).is_dir() and not "force" in args:
        logging.error(f"❌ Directory {dir_name} already exists, exiting!")
        exit(1)

    mkdir(h / dir_name)
    mkdir(h / dir_name / "source")
    mkdir(h / dir_name / "code")

    template = ASSIGNMENT_TEMPLATE.format(
        due=due_date,
        sheet=add_leading_zero(number),
        course=course,
        group=group,
        members=members,
    )

    assignment_file = h / dir_name / "assignment.tex"

    with open(assignment_file, "w") as f:
        f.write(template)
    print(f"📜 Templated {str(assignment_file)}")
    write_configmap(number)


def run_latexmk(assignment_path):
    cmd = f"latexmk -pdf -interaction=nonstopmode -file-line-error -shell-escape -f {os.path.basename(assignment_path)}"
    dir_name = os.path.dirname(assignment_path)
    for i in range(0, args.runs):
        print(f"⏳ Running latexmk [{i+1}/3]")
        _ = subprocess.run(cmd.split(), capture_output=(
            args.quiet), cwd=dir_name)

    artifact_path = Path('.') / str(assignment_path).replace('.tex', '.pdf')
    artifact_no = str(dir_name).replace('assignment-', '')
    artifact_fn = os.path.basename(artifact_path)
    artifact_dst = Path(
        '.') / 'dist' / (artifact_fn.replace('.pdf', '') + '-' + artifact_no + '.pdf')
    if not (Path('.') / 'dist').is_dir():
        mkdir(Path('.') / 'dist')
    skip = False
    if artifact_dst.is_file():
        if 'force' in args and args.force:
            artifact_dst.unlink()
        else:
            print(
                f"📚 {artifact_dst} already exists, skipping... (provide --force to override file)")
            skip = True
    if not skip:
        copy2(artifact_path, artifact_dst)
        print(f"✅ Built {assignment_path} to {artifact_dst}")
    if not args.keep:
        clean_cmd = f"latexmk -C"
        _ = subprocess.run(clean_cmd.split(),
                           capture_output=(args.quiet), cwd=dir_name)


def build_all():
    files = list(Path('.').glob('assignment-*/assignment.tex'))
    for file in files:
        print("📄 Building {file}...".format(file=str(file)))
        run_latexmk(file)


def build_specific(number):
    file = Path('.') / ('assignment' + '-' +
                        add_leading_zero(number)) / 'assignment.tex'
    print("📄 Building {file}...".format(file=str(file)))
    run_latexmk(file)


def build(args):
    if 'all' in args and args.all:
        build_all()
    else:
        if 'number' in args and args.number:
            build_specific(int(args.number))
        else:
            build_all()


def bundle(assignment_no):
    config = configparser.ConfigParser()
    config.read(CONFIGMAP_PATH)
    members = list(config['members'].keys())
    archive_name = f"sheet_{add_leading_zero(assignment_no)}_{'_'.join(members)}.zip"
    archive_path = (Path('.') / 'dist' / archive_name)
    skip = False
    if archive_path.is_file():
        if 'force' in args and args.force:
            archive_path.unlink()
        else:
            print(
                f"📚 {archive_path} already exists, skipping... (provide --force to override file)")
            skip = True
    if not skip:
        cwd = Path('./dist')
        ano = add_leading_zero(assignment_no)
        mkdir(cwd / f"assignment-{ano}")
        cp_code_cmd = f"cp -r assignment-{ano}/code dist/assignment-{ano}"
        cp_pdf_cmd = f"cp dist/assignment-{ano}.pdf dist/assignment-{ano}/"
        _ = subprocess.run(cp_code_cmd.split(),
                           capture_output=(args.quiet), cwd=Path('.'))
        _ = subprocess.run(cp_pdf_cmd.split(),
                           capture_output=(args.quiet), cwd=Path('.'))

        code_files = list(Path(cwd / f"assignment-01").glob('code/*'))
        code_files = [str(p).replace(
            f"dist/assignment-{ano}/", "") for p in code_files]
        zip_cmd = f"zip {archive_name} assignment-{ano}.pdf {''.join(code_files)}"
        _ = subprocess.run(zip_cmd.split(), capture_output=(
            args.quiet), cwd=(cwd / f"assignment-{ano}"))
        _ = subprocess.run(
            f"mv assignment-{ano}/{archive_name} .".split(), capture_output=(args.quiet), cwd=cwd)
        _ = subprocess.run(
            f"rm -r assignment-{ano}".split(), capture_output=(args.quiet), cwd=cwd)
        print(f"🗄️ Compiled archive at {archive_path}")


def compile_all():
    files = list(Path('.').glob('dist/assignment-*.pdf'))
    for f in files:
        assignment_no = str(f).removeprefix(
            'dist/assignment-').removesuffix('.pdf')
        print("⏳ Compiling assignment {number}...".format(
            number=assignment_no))
        bundle(int(assignment_no))


def compile_specific(number):
    print("⏳ Compiling assignment {number}...".format(number=number))
    bundle(int(number))


def compile(args):
    if 'all' in args and args.all:
        compile_all()
    else:
        if 'number' in args and args.number:
            compile_specific(int(args.number))
        else:
            compile_all()


def release(args):
    tag = environ.get("CI_COMMIT_TAG")
    assignment = tag.removeprefix('assignment-')
    artifacts_id = environ.get("CI_JOB_ID")
    archive_name = str(os.path.basename(list(Path(
        './dist').glob('sheet_{assignment}_*.zip'.format(assignment=assignment)))[0]))
    pdf_name = str(os.path.basename(
        Path('./dist/assignment-{assignment}.pdf'.format(assignment=assignment))))
    print("🏷️ Git tag is {tag}, releasing assignment {assignment} with {archive_name}...".format(
        tag=tag, assignment=assignment, archive_name=archive_name))
    with open('artifacts.env', "a") as envfile:
        envfile.write("ASSIGNMENT={assignment}\n".format(
            assignment=assignment))
        envfile.write("TAG={tag}\n".format(tag=tag))
        envfile.write("ARTIFACTS_ID={artifacts_id}\n".format(
            artifacts_id=artifacts_id))
        envfile.write("ARCHIVE_NAME={archive_name}\n".format(
            archive_name=archive_name))
        envfile.write("PDF_NAME={pdf_name}\n".format(
            pdf_name=pdf_name))


HELP = """\

💡 Toolchain utility for building LaTeX assignments

   Run './make.py bootstrap' to create a new environment {needed}
   Run './make.py generate <number>' to template a new assignment
   Run './make.py build <number>' to build a specfic assignment with latexmk (add '--all' to build all assignments)
   Run './make.py compile <number>' to compile a zip file for a specific assignment (add '--all' to compile all assignments)
   Run './make.py release' inside a Gitlab CI/CD pipeline to create files required for automatic release (see manual)

   See each command's help for description of arguments

   Copyright (C) zoomoid, 2022
"""


def print_help(args):
    needs_env = Path(
        './.assignments.rc').is_file() and Path('./assignments.cls').is_file()
    print(HELP.format(needed='(not required right now)' if needs_env else ''))
    return ""


root = argparse.ArgumentParser(description="", prog="./make.py")
root.add_argument("--noninteractive", action="store_true",
                  default=False, help="Skip any user prompts")
root.set_defaults(func=print_help)

subparsers = root.add_subparsers(help="Sub-command parsers")

bootstrap_cmd = subparsers.add_parser(
    "bootstrap", help="Bootstraps an environment with LaTeX class and .assignments.rc file")
bootstrap_cmd.set_defaults(func=bootstrap)
bootstrap_cmd.add_argument("--course", help="Course name")
bootstrap_cmd.add_argument("--group", help="Group name, leave empty to omit")
bootstrap_cmd.add_argument(
    "--members",
    nargs="*",
    help="Group members, e.g., 'John Doe, 123456' 'Jane Doe, 789012'",
)

generate_cmd = subparsers.add_parser(
    "generate", help="Generate new assignments from the template. Omit to use assignments.number from .assignments.rc"
)
generate_cmd.set_defaults(func=generate_assignment)
generate_cmd.add_argument(
    "number", nargs="?", default=None, type=int, help="Assignment number")
generate_cmd.add_argument(
    "--no-increment",
    action="store_true",
    help="Skip incrementing assignments.number in .assignments.rc",
)
generate_cmd.add_argument(
    "--force", "-F", default=False, action="store_true", help="Overrides any existing assignment source files"
)

build_cmd = subparsers.add_parser(
    "build", help="Runs 'latexmk' on assignment files and collects them in dist/")
build_cmd.set_defaults(func=build)
build_cmd.add_argument("number", nargs="?", type=int,
                       help="Assignment number to build. Omit to build ALL assignments")
build_cmd.add_argument("--all", "-A", action="store_true",
                       default=False, help="Build all assignments in assignment-*/")
build_cmd.add_argument("--quiet", "-q", action="store_true",
                       default=False, help="Suppress output from latexmk subprocesses")
build_cmd.add_argument("--keep", action="store_true", default=False,
                       help="Skip latexmk -C cleaning up all files in the source directory")
build_cmd.add_argument("--force", "-F", action="store_true", default=False,
                       help="Override any existing assignments with the same name")
build_cmd.add_argument("--runs", "-r", default=3,
                       type=int, help="latexmk compiler runs")

compile_cmd = subparsers.add_parser(
    "compile", help="Bundles assignments into submittable zip archives")
compile_cmd.set_defaults(func=compile)
compile_cmd.add_argument("number", nargs="?", type=int,
                         help="Assignment number to compile. Omit to bundle ALL assignments found in dist/")
compile_cmd.add_argument("--all", "-A", action="store_true",
                         default=False, help="Compile all assignments in dist/")
compile_cmd.add_argument("--quiet", "-q", action="store_true",
                         default=False, help="Suppress output from zip subprocess")
compile_cmd.add_argument("--force", "-F", action="store_true", default=False,
                         help="Override any existing archives with the same name")

release_cmd = subparsers.add_parser(
    "release", help="Creates a pre-release file storing values required for the Gitlab Release CLI")
release_cmd.set_defaults(func=release)

args = root.parse_args()
args.func(args)
