import subprocess
import argparse

default_tag = "v0.1.0-test-alpha-release"

def run_command(command):
    """Run a shell command and capture its output."""
    try:
        result = subprocess.run(command, shell=True, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result.stdout.decode('utf-8').strip()
    except subprocess.CalledProcessError as e:
        print(f"Command '{command}' failed with error: {e.stderr.decode('utf-8').strip()}")
        return None

def remove_tag(tag):
    """Remove the specified Git tag locally and remotely, if it exists."""
    # Check if the tag exists locally
    print(f"Checking if tag '{tag}' exists locally...")
    tag_exists = run_command(f"git tag -l {tag}")

    if tag_exists:
        print(f"Tag '{tag}' found. Removing locally...")
        run_command(f"git tag -d {tag}")
    else:
        print(f"Tag '{tag}' does not exist locally. Continuing...")

    # Attempt to remove the tag remotely
    print(f"Removing tag '{tag}' from remote (if exists)...")
    run_command(f"git push origin :refs/tags/{tag}")

def add_tag(tag):
    """Create a new tag and push it."""
    print(f"Creating and adding new tag '{tag}' locally...")
    run_command(f"git tag {tag}")

    print(f"Pushing new tag '{tag}' to remote...")
    run_command(f"git push origin {tag}")

def push_changes():
    """Push changes to remote and push all tags."""
    print("Pushing changes to remote with --no-verify and --force...")
    run_command("git push --no-verify --force")

    print("Pushing all tags to remote...")
    run_command("git push --tags")

def main():
    parser = argparse.ArgumentParser(description="Remove a Git tag, add it if necessary, and push changes.")
    parser.add_argument("-tag", required=False, default=default_tag, help="The Git tag to remove and re-add.")

    args = parser.parse_args()

    # Remove the specified tag if it exists
    remove_tag(args.tag)

    # Add (or re-add) the tag
    add_tag(args.tag)

    # Push remaining changes and all tags
    push_changes()

if __name__ == "__main__":
    main()
