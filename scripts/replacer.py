#!/usr/bin/env python3
# filepath: /Users/kaan/Code/LightningBot/LB-Selfbot/scripts/replacer.py

import sys
import os
import re

def transform_file_contents(lines, package_name):
    """
    Transform the file contents according to the old replacer logic.
    We treat 'package <something>' as 'package package_name'.
    We detect an old RegisterCommands function to replace with a new one.
    We also do string replacements for flags, usage -> arguments, etc.
    """
    func_pattern = r'^func\s+\(\w+\s+\*ActivityCommands\)\s+RegisterCommands'
    category_pattern = r'^\s*router\.StartCategory\("([^"]+)",\s*"([^"]+)"\)'
    cmd_start_pattern = r'^\s*router\.RegisterCmd\(&dgc\.Command\{\s*$'
    cmd_end_pattern = r'^\s*\}\)\s*$'
    old_struct_pattern = r'^type\s+ActivityCommands\s+struct\s*\{\s*\}'

    top_section = []
    bottom_section = []

    found_register_func = False
    finished_register_func = False

    cat_name = None
    cat_desc = None
    commands = []
    cur_cmd_block = []
    collecting_cmd_block = False

    for line in lines:
        stripped = line.rstrip('\n')

        # Replace the package line
        if stripped.startswith('package '):
            top_section.append(f'package {package_name}\n')
            continue

        # Skip empty struct
        if re.match(old_struct_pattern, stripped):
            continue

        # Detect the old RegisterCommands function start
        if re.match(func_pattern, stripped):
            found_register_func = True
            continue

        # Detect the end of old RegisterCommands function
        if found_register_func and stripped.startswith('}'):
            finished_register_func = True
            found_register_func = False
            continue

        # Inside old RegisterCommands block: detect category and commands
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

        # Lines outside old RegisterCommands
        if not finished_register_func:
            top_section.append(line)
        else:
            bottom_section.append(line)

    # Build new struct + constructor + Register method
    new_code = []
    new_code.append('type ActivityCommands struct { services *services.Services }\n')
    new_code.append('func NewActivityCommands(s *services.Services) *ActivityCommands {\n')
    new_code.append('\treturn &ActivityCommands{services: s}\n')
    new_code.append('}\n')
    new_code.append('func (c *ActivityCommands) Register() {\n')
    new_code.append('\tcmdRouter := c.services.Discord.GetRouter()\n\n')

    if cat_name and cat_desc and commands:
        new_code.append('\tcmdRouter.RegisterCategory(&router.Category{\n')
        new_code.append(f'\t\tName:        "{cat_name}",\n')
        new_code.append(f'\t\tDescription: "{cat_desc}",\n')
        new_code.append('\t}, []*router.Command{\n')
        for cmd_block in commands:
            new_code.append('\t\t{\n')
            for line_cmd in cmd_block:
                # Replace library references from dgc => router
                line_cmd = line_cmd.replace('dgc', 'router')
                # Convert Usage => Arguments and remove [p]
                if 'Usage:' in line_cmd:
                    line_cmd = line_cmd.replace('Usage:', 'Arguments:')
                    line_cmd = line_cmd.replace('[p]', '')
                # Convert []string => []router.CommandFlag
                if 'Flags:' in line_cmd and '[]string{' in line_cmd:
                    line_cmd = line_cmd.replace('[]string{', '[]router.CommandFlag{')
                # Convert CmdFlagDangerous => router.CmdFlagDangerous
                line_cmd = line_cmd.replace('CmdFlagDangerous', 'router.CmdFlagDangerous')
                # Convert a. handler => c. handler
                if 'Handler:' in line_cmd and 'a.' in line_cmd:
                    line_cmd = line_cmd.replace('Handler:     a.', 'Handler:     c.')

                new_code.append(f'\t\t\t{line_cmd}\n')
            new_code.append('\t\t},\n')
        new_code.append('\t})\n')

    new_code.append('}\n\n')

    # Merge everything
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
        # We're only interested in directories that contain "xxx.go" and "xxx-old.go"
        # For each file ending in "-old.go", do the transformation & renaming
        for filename in files:
            if filename.endswith("-old.go"):
                old_path = os.path.join(root, filename)
                # The new file name is the same minus "-old"
                new_filename = filename.replace("-old.go", ".go")
                new_path = os.path.join(root, new_filename)

                # If there's an existing .go file, delete it
                if os.path.exists(new_path):
                    os.remove(new_path)

                # Read the old file
                with open(old_path, 'r') as f:
                    old_lines = f.readlines()

                # The package name is the directory's name
                folder_name = os.path.basename(root)

                # Transform
                new_contents = transform_file_contents(old_lines, folder_name)

                # Write back to the old file
                with open(old_path, 'w') as f:
                    f.writelines(new_contents)

                # Rename old -> new
                os.rename(old_path, new_path)

                print(f"Processed: {old_path} -> {new_path}")

if __name__ == "__main__":
    main()