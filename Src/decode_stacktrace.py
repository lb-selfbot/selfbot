import os

print("Input stacktrace (CTRL+Z to end): ", end="", flush=True)

stacktrace = []
while True:
    try:
        line = input()
    except EOFError:
        break
    stacktrace.append(line)


with open("stack.txt", "w+") as f:
    f.write("\n".join(stacktrace))

try:
    os.system(
        "garble -seed=U4Nxq2YSVm6ivvOejQbElA== reverse -tags desktop github.com/LightningDev1/LB-Selfbot stack.txt"
    )
finally:
    os.remove("stack.txt")
