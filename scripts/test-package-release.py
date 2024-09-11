import subprocess
import sys
import argparse

def run_command(command):
    """Run a shell command and capture its output."""
    try:
        result = subprocess.run(command, shell=True, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result.stdout.decode('utf-8').strip()
    except subprocess.CalledProcessError as e:
        print(f"Command '{command}' failed with error: {e.stderr.decode('utf-8').strip()}")
        sys.exit(1)

def remove_tag(tag):
    """Remove the specified Git tag locally and remotely."""
    print(f"Removing tag '{tag}' locally...")
    run_command(f"git tag -d {tag}")

    print(f"Removing tag '{tag}' from remote...")
    run_command(f"git push origin :refs/tags/{tag}")

def push_changes():
    """Push changes to remote and push all tags."""
    print("Pushing changes to remote with --no-verify and --force...")
    run_command("git push --no-verify --force")

    print("Pushing all tags to remote...")
    run_command("git push --tags")

def main():
    parser = argparse.ArgumentParser(description="Remove a Git tag and push changes.")
    parser.add_argument("-tag", required=True, help="The Git tag to remove.")

    args = parser.parse_args()

    remove_tag(args.tag)
    push_changes()

if __name__ == "__main__":
    main()
