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

def remove_test_release_tags():
    """Remove all tags with '-test-release' suffix locally and remotely."""
    print("Finding all tags with '-test-release' suffix...")
    tags_to_remove = run_command("git tag --list '*-test-release'").splitlines()

    if not tags_to_remove:
        print("No tags found with '-test-release' suffix.")
        return

    for tag in tags_to_remove:
        print(f"Removing tag: {tag}")
        remove_tag(tag)

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
        f.write(f"Test release of {package_dir} module\n\n")
        f.write("Features added:\n")
        f.write("- Test feature #1\n")
        f.write("- Test feature #2\n")

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
    parser = argparse.ArgumentParser(description="Remove Git tags, add a release file, and push changes.")
    parser.add_argument("-tag", required=False, help="The Git tag to remove and re-add.")
    parser.add_argument("-package", required=False, help="The package directory to create the release file in.")
    parser.add_argument("-remove-test-tags", required=False, action="store_true", help="Remove all Git tags with '-test-release' suffix.")

    args = parser.parse_args()

    if args.remove_test_tags:
        remove_test_release_tags()
    elif args.tag and args.package:
        remove_tag(args.tag)
        add_release_file(args.package, args.tag)
        add_tag(args.tag)
        push_changes()
    else:
        print("You must provide either '-remove-test-tags' or '-tag' and '-package' options.")

if __name__ == "__main__":
    main()
