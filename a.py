import re



with open("image.txt", "r") as f:
    raw_commands = f.read()

commands = re.findall(r'@commands.command\(help="(.*?)", usage="(.*?)"\)' + '\n' + r'async def (.*?)\(', raw_commands)


def get_code(name):
    index = raw_commands.index(name + "(")
    new_cmds = raw_commands[index:]
    new_cmds = new_cmds[new_cmds.index("):\n"):].split("\n", 1)[-1]
    code = []
    split = new_cmds.splitlines()
    for line in split:
        if line == "":
            if len(split) < split.index(line) + 2:
                break
            if split[split.index(line) + 1].startswith("@commands.command"):
                break
        code.append(line)
    return ''.join('	\n// ' + a for a in code)

goCode = ""

for command in commands:
    description, usage, name = command
    code = get_code(name)
    goCode += \
f'''func (ImageCommands) {name}(ctx *dgc.Ctx) {{{code}
}}
\n'''

with open("new_image.txt", "w") as f:
    f.write(goCode)