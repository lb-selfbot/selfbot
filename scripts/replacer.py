#!/usr/bin/env python3
# filepath: replacer.py

import sys
import os
import re

def transform_file_contents(lines, package_name):
    """
    Transform the old file contents.
    Captures the command category via the RegisterCommands function.
    Performs string replacements on commands:
     - Replaces dgc with router.
     - Changes Usage to Arguments and strips "[p]".
     - Converts []string{...} to []router.CommandFlag{...} and
       CmdFlagDangerous to router.CmdFlagDangerous.
     - Changes handler references from a. (old) to c. (new).
    """
    # Captures the struct name (category) from the RegisterCommands method.
    func_pattern = r'^func\s+\(\w+\s+\*(\w+)Commands\)\s+RegisterCommands'
    category_pattern = r'^\s*router\.StartCategory\("([^"]+)",\s*"([^"]+)"\)'
    cmd_start_pattern = r'^\s*router\.RegisterCmd\(&dgc\.Command\{\s*$'
    cmd_end_pattern = r'^\s*\}\)\s*$'
    old_struct_pattern = r'^type\s+\w+Commands\s+struct\s*\{\s*\}'

    top_section = []
    bottom_section = []
    struct_name = None

    found_register_func = False
    finished_register_func = False

    cat_name = None
    cat_desc = None
    commands = []
    cur_cmd_block = []
    collecting_cmd_block = False

    for line in lines:
        stripped = line.rstrip('\n')

        # Replace the package declaration with the given package name
        if stripped.startswith('package '):
            top_section.append(f'package {package_name}\n')
            continue

        # Skip empty struct definitions
        if re.match(old_struct_pattern, stripped):
            continue

        # Detect the start of the old RegisterCommands function and capture struct name
        m = re.match(func_pattern, stripped)
        if m:
            struct_name = m.group(1)
            found_register_func = True
            continue

        # End of the RegisterCommands function block
        if found_register_func and stripped.startswith('}'):
            finished_register_func = True
            found_register_func = False
            continue

        # Inside the RegisterCommands function block
        if found_register_func and not finished_register_func:
            if re.match(category_pattern, stripped):
                m = re.match(category_pattern, stripped)
                cat_name, cat_desc = m.groups()
            elif re.match(cmd_start_pattern, stripped):
                collecting_cmd_block = True
                cur_cmd_block = []
            elif re.match(cmd_end_pattern, stripped) and collecting_cmd_block:
                collecting_cmd_block = False
                commands.append(cur_cmd_block)
                cur_cmd_block = []
            elif collecting_cmd_block:
                cur_cmd_block.append(stripped)
            continue

        # Lines outside of the RegisterCommands block
        if not finished_register_func:
            top_section.append(line)
        else:
            bottom_section.append(line)

    # If we didn't capture a struct name, use package_name capitalized as fallback.
    if not struct_name:
        struct_name = package_name.capitalize()

    # Build the new struct, constructor, and Register method
    new_code = []
    new_code.append(f'type {struct_name}Commands struct {{ services *services.Services }}\n')
    new_code.append(f'func New{struct_name}Commands(s *services.Services) *{struct_name}Commands {{\n')
    new_code.append(f'\treturn &{struct_name}Commands{{services: s}}\n')
    new_code.append('}\n')
    new_code.append(f'func (c *{struct_name}Commands) Register() {{\n')
    new_code.append('\tcmdRouter := c.services.Discord.GetRouter()\n\n')

    if cat_name and cat_desc and commands:
        new_code.append('\tcmdRouter.RegisterCategory(&router.Category{\n')
        new_code.append(f'\t\tName:        "{cat_name}",\n')
        new_code.append(f'\t\tDescription: "{cat_desc}",\n')
        new_code.append('\t}, []*router.Command{\n')
        for cmd_block in commands:
            new_code.append('\t\t{\n')
            for line_cmd in cmd_block:
                # Replace occurrences of "dgc" with "router"
                line_cmd = line_cmd.replace('dgc', 'router')
                # Change "Usage:" to "Arguments:" and remove "[p]"
                if 'Usage:' in line_cmd:
                    line_cmd = line_cmd.replace('Usage:', 'Arguments:')
                    line_cmd = line_cmd.replace('[p]', '')
                # Convert flag declarations from []string{...} to []router.CommandFlag{...}
                if 'Flags:' in line_cmd and '[]string{' in line_cmd:
                    line_cmd = line_cmd.replace('[]string{', '[]router.CommandFlag{')
                # Change "CmdFlagDangerous" to "router.CmdFlagDangerous"
                line_cmd = line_cmd.replace('CmdFlagDangerous', 'router.CmdFlagDangerous')
                # Update handler references from a. to c.
                if 'Handler:' in line_cmd and 'a.' in line_cmd:
                    line_cmd = line_cmd.replace('Handler:     a.', 'Handler:     c.')
                new_code.append(f'\t\t\t{line_cmd}\n')
            new_code.append('\t\t},\n')
        new_code.append('\t})\n')

    new_code.append('}\n\n')

    # Merge all sections; ensure a single package line remains
    final_lines = []
    found_package_line = False
    for ln in top_section:
        if ln.strip().startswith('package '):
            if not found_package_line:
                found_package_line = True
                final_lines.append(f'package {package_name}\n')
            continue
        final_lines.append(ln)
    final_lines.append(''.join(new_code))
    final_lines.extend(bottom_section)
    return final_lines

def main():
    if len(sys.argv) < 2:
        print("Usage: replacer.py <directory path>")
        sys.exit(1)

    base_dir = sys.argv[1]

    # Walk through subdirectories in base_dir
    for root, dirs, files in os.walk(base_dir):
        # Process each file ending in "-old.go"
        for filename in files:
            if filename.endswith("-old.go"):
                old_path = os.path.join(root, filename)
                # New filename is the same but without "-old"
                new_filename = filename.replace("-old.go", ".go").lower()
                new_path = os.path.join(root, new_filename)

                # Delete the existing .go file if it exists
                if os.path.exists(new_path):
                    os.remove(new_path)

                # Read the old file
                with open(old_path, 'r') as f:
                    old_lines = f.readlines()

                # Use the directory's name as the package name
                package_name = os.path.basename(root).lower()

                # Transform file contents
                new_contents = transform_file_contents(old_lines, package_name)

                # Write the transformed contents back to the old file
                with open(old_path, 'w') as f:
                    f.writelines(new_contents)

                # Rename the file from old to new
                os.rename(old_path, new_path)
                print(f"Processed: {old_path} -> {new_path}")

    os.system('gofmt -w .')

if __name__ == "__main__":
    main()