#!/usr/bin/env python3
import re
import sys
import subprocess
import argparse
from pathlib import Path

# --- Colors (Charmtone) ---
# Using 24-bit ANSI color codes
def hex_to_rgb(hex_color):
    hex_color = hex_color.lstrip('#')
    return tuple(int(hex_color[i:i+2], 16) for i in (0, 2, 4))

def color_text(text, hex_color, bold=False):
    r, g, b = hex_to_rgb(hex_color)
    code = f"\033[38;2;{r};{g};{b}m"
    reset = "\033[0m"
    style = "\033[1m" if bold else ""
    return f"{code}{style}{text}{reset}"

# Colors from crush_colors.txt
CHARPLE = "#6B50FF"  # Primary
DOLLY   = "#FF60FF"  # Secondary
GUAC    = "#12C78F"  # Success
SRIRACHA= "#EB4268"  # Error
ZEST    = "#E8FE96"  # Warning/Accent
MALIBU  = "#00A4FF"  # Info
ASH     = "#DFDBDD"  # FgBase
SMOKE   = "#BFBCC8"  # FgHalfMuted

def print_info(msg):
    print(f"{color_text('ℹ', MALIBU, True)} {color_text(msg, ASH)}")

def print_success(msg):
    print(f"{color_text('✔', GUAC, True)} {color_text(msg, GUAC)}")

def print_error(msg):
    print(f"{color_text('✘', SRIRACHA, True)} {color_text(msg, SRIRACHA)}")
    sys.exit(1)

def print_step(msg):
    print(f"\n{color_text('::', CHARPLE, True)} {color_text(msg, DOLLY, True)}")

def run_cmd(cmd, shell=False):
    try:
        subprocess.check_output(cmd, shell=shell, stderr=subprocess.STDOUT)
    except subprocess.CalledProcessError as e:
        print_error(f"Command failed: {cmd}\nOutput: {e.output.decode()}")

def get_current_version(file_path):
    content = file_path.read_text()
    match = re.search(r'version\s*=\s*"([^"]+)"', content)
    if match:
        return match.group(1)
    print_error("Could not find version in pyproject.toml")

def bump_version_string(version, part='patch'):
    major, minor, patch = map(int, version.split('.'))
    if part == 'major':
        major += 1
        minor = 0
        patch = 0
    elif part == 'minor':
        minor += 1
        patch = 0
    else:
        patch += 1
    return f"{major}.{minor}.{patch}"

def update_file(file_path, old_ver, new_ver):
    content = file_path.read_text()
    new_content = content.replace(f'version = "{old_ver}"', f'version = "{new_ver}"')
    file_path.write_text(new_content)

def main():
    parser = argparse.ArgumentParser(description="Bump version and tag release")
    parser.add_argument('part', choices=['major', 'minor', 'patch'], default='patch', nargs='?', help="Version part to bump")
    args = parser.parse_args()

    pyproject_file = Path("pyproject.toml")
    if not pyproject_file.exists():
        print_error("pyproject.toml not found!")

    # 1. Get current version
    current_ver = get_current_version(pyproject_file)
    new_ver = bump_version_string(current_ver, args.part)
    
    print_step(f"Bumping version: {color_text(current_ver, SMOKE)} → {color_text(new_ver, ZEST, True)}")

    # 2. Update pyproject.toml
    update_file(pyproject_file, current_ver, new_ver)
    print_success(f"Updated pyproject.toml to {new_ver}")

    # Also update setup.py if it has version hardcoded
    setup_file = Path("setup.py")
    if setup_file.exists():
        content = setup_file.read_text()
        if f'version="{current_ver}"' in content:
            content = content.replace(f'version="{current_ver}"', f'version="{new_ver}"')
            setup_file.write_text(content)
            print_success(f"Updated setup.py to {new_ver}")
    
    # 3. Git commit
    print_step("Committing changes...")
    run_cmd(["git", "add", "pyproject.toml", "setup.py"])
    run_cmd(["git", "commit", "-m", f"Bump version to {new_ver}"])
    print_success("Committed changes")

    # 4. Git tag
    tag_name = f"v{new_ver}"
    print_step(f"Tagging release {color_text(tag_name, DOLLY, True)}...")
    run_cmd(["git", "tag", tag_name])
    print_success(f"Created tag {tag_name}")

    # 5. Push
    print_step("Pushing to GitHub...")
    print_info("Pushing commit and tag...")
    run_cmd(["git", "push", "origin", "main"]) # Assuming main branch
    run_cmd(["git", "push", "origin", tag_name])
    
    print("\n" + color_text("✨ Release Initiated! ✨", CHARPLE, True))
    print(f"The {color_text('release.yml', MALIBU)} workflow should now be running on GitHub.")

if __name__ == "__main__":
    main()

