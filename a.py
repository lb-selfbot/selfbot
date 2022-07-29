import re



with open("image.txt", "rb") as f:
    raw_commands = f.read().decode("utf-8")
# r'@commands.command\(help="(.*?)", usage="(.*?)"\)' + '\n' + r'async def (.*?)\('
commands = re.findall(r'@commands.command\(help="(.*?)", usage="(.*?)"\)' + '\r\n' + r'async def (.*?)\(', raw_commands)
print(commands)

def get_code(name):
    index = raw_commands.index(name + "(")
    new_cmds = raw_commands[index:]
    new_cmds = new_cmds[new_cmds.index("):\r\n"):].split("\r\n", 1)[-1]
    code = []
    split = new_cmds.split("\r\n")
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
    #code = get_code(name)
    goCode += \
f'''router.RegisterCmd(&dgc.Command{{
    Name:        "{name}",
    Description: "{description}",
    Usage:       "[p]{name}{' '+usage if usage else ''}",
    Handler:     u.{name},
}})\n'''

with open("new_image.txt", "wb") as f:
    f.write(goCode.encode("utf-8"))