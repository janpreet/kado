import re
import subprocess

def get_current_branch():
    try:
        branch = subprocess.check_output(["git", "rev-parse", "--abbrev-ref", "HEAD"]).strip().decode('utf-8')
        return branch
    except subprocess.CalledProcessError as e:
        print(f"Error getting current branch: {e}")
        return None

def read_version(version_file):
    with open(version_file, 'r') as file:
        return file.read().strip()

def write_version(version_file, version):
    with open(version_file, 'w') as file:
        file.write(version + '\n')

def bump_version(version, part):
    major, minor, patch = map(int, version.split('.'))
    if part == 'major':
        major += 1
        minor = 0
        patch = 0
    elif part == 'minor':
        minor += 1
        patch = 0
    elif part == 'patch':
        patch += 1
    return f"{major}.{minor}.{patch}"

def main():
    version_file = 'VERSION'
    current_version = read_version(version_file)
    branch = get_current_branch()
    
    if not branch:
        print("Could not determine the current git branch.")
        return
    
    if re.match(r'^feature/', branch):
        new_version = bump_version(current_version, 'major')
    elif re.match(r'^fix/', branch):
        new_version = bump_version(current_version, 'minor')
    else:
        new_version = bump_version(current_version, 'patch')
    
    write_version(version_file, new_version)
    print(f"Bumped version from {current_version} to {new_version}")

if __name__ == "__main__":
    main()
