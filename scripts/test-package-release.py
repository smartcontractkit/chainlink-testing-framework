# This script is used to verify release pipeline end-to-end
# Usage:
# python ./scripts/test-package-release.py -tag k8s-test-runner/v1.999.0-test-release -package ./k8s-test-runner
import subprocess
import argparse
import os

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

def add_release_file(package_dir, tag):
    """Change directory to the package and create a release file."""
    # Extract the version part of the tag
    version_part = tag.split('/')[-1]
    filename = f".changeset/{version_part}.md"

    # Change directory to the package
    os.chdir(package_dir)
    print(f"Changed directory to {package_dir}. Creating file {filename}...")

    # Write example data to the release file
    with open(filename, 'w') as f:
        f.write(f"Initial release of {package_dir} test runner\n\n")
        f.write("Features added:\n")
        f.write("- One\n")
        f.write("- Two\n")

    # Add and commit the new file
    run_command(f"git add {filename}")
    print("Yubikey signature might be required if the script is hanging for too long check your Yubikey...")
    commit_message = f"Test release commit {version_part}"
    run_command(f"git commit -m '{commit_message}' --no-verify")
    print(f"Committed the release file: {filename}")

def add_tag(tag):
    """Create a new tag and push it."""
    print(f"Creating and adding new tag '{tag}' locally...")
    run_command(f"git tag {tag}")

    print(f"Pushing new tag '{tag}' to remote...")
    run_command(f"git push origin :refs/tags/{tag}")

def push_changes():
    """Push changes to remote and push all tags."""
    print("Pushing changes to remote with --no-verify and --force...")
    run_command("git push --no-verify --force")

    print("Pushing all tags to remote...")
    run_command("git push --tags")

def main():
    parser = argparse.ArgumentParser(description="Remove a Git tag, add a release file, and push changes.")
    parser.add_argument("-tag", required=True, help="The Git tag to remove and re-add.")
    parser.add_argument("-package", required=True, help="The package directory to create the release file in.")

    args = parser.parse_args()

    # Remove the specified tag if it exists
    remove_tag(args.tag)

    # Add release file to the package directory
    add_release_file(args.package, args.tag)

    # Add (or re-add) the tag
    add_tag(args.tag)

    # Push remaining changes and all tags
    push_changes()

if __name__ == "__main__":
    main()
