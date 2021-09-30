import os
import argparse
import unittest
import re

from string import digits

pipe_grep_error = "| grep -o -P \"(?<=The exit code is )([0-9]*)\""
stderr_to_null = "2>/dev/null"
ignore_addresses = r"0x[0-9a-fA-F]{5}"
def del_ws(s): return "".join(s.split())


stdin = "12a34567890YN\n\t"


class TestChunks(unittest.TestCase):

    def error_ref_output(self, flags, filepath):
        return os.popen(f"ruby ./tests/refCompile {flags} {filepath}" + " | sed -n '/^-- Compiling...$/ { :a; n; p; ba; }'  | sed -e 's/-- //g'").read().rstrip()

    def code_ref_output(self, flags, filepath: str):
        ref_output = "0"
        if "syntaxErr" in filepath:
            ref_output = "100"
        elif "semanticErr" in filepath:
            ref_output = "200"
        return ref_output

    def code_our_output(self, flags, filepath):
        return os.popen(f"./compile {flags} {filepath} >/dev/null 2>&1 && echo $? || echo $?").read().rstrip()

    def tree_ref_output(self, flags, filepath):
        output = os.popen(
            f"echo \"{stdin}\" | ruby ./tests/refCompile {flags} {filepath}").read().strip()
        tree = output.split(
            "===========================================================")[1]
        return "".join(line.lstrip(digits) for line in tree.splitlines(keepends=True)).strip().replace("\t", "")

    def exe_ref_output(self, flags, filepath):
        output = os.popen(
            f"echo \"{stdin}\" | ruby ./tests/refCompile {flags} {filepath}").read().strip()
        tree = output.split(
            "===========================================================")[1]
        output = "".join(line.strip() for line in tree.splitlines(
            keepends=True)).strip().replace("\t", "")
        return re.sub(ignore_addresses, "", output)

    def tree_our_output(self, flags, filepath):
        output = os.popen(f"./compile {flags} {filepath}").read().strip()
        tree = output.split(
            "===========================================================")[1]
        t = "".join(i for i in tree.splitlines(keepends=True)).strip()
        return t

    def gen_ref_code(self, flags, filepath):
        raw = os.popen(
            f"echo \"\" | ruby tests/refCompile {flags} {filepath} 2>&1 {pipe_grep_error}").read().rstrip()
        return raw

    def gen_our_code(self, flags, filepath):
        cmd = f"./compile {filepath}"
        input_path = f"{'/'.join(filepath.split('/')[:-1])}/input.s"
        cmd += f"&& arm-linux-gnueabi-gcc -o input -mcpu=arm1176jzf-s -pthread -mtune=arm1176jzf-s {input_path}"
        cmd += f"&& echo \"{stdin}\" | qemu-arm -L /usr/arm-linux-gnueabi/ input && echo $? || echo $?"
        return os.popen(cmd).read().rstrip()

    def default_our_output(self, flags, filepath):
        return os.popen(f"./compile {flags} {filepath}").read().rstrip()

    # Assumes the required executable has been compiled
    def our_generated_output(self, flags, filepath):
        cmd = f"./compile {filepath}"
        input_path = f"{'/'.join(filepath.split('/')[:-1])}/input.s"
        cmd += f"&& arm-linux-gnueabi-gcc -o input -mcpu=arm1176jzf-s -pthread -mtune=arm1176jzf-s {input_path}"
        cmd += f"&& echo \"{stdin}\" | qemu-arm -L /usr/arm-linux-gnueabi/ input 2>&1"
        output = os.popen(cmd).read().rstrip()
        output = "".join(line.strip() for line in output.splitlines(
            keepends=True)).strip().replace("\t", "")
        return re.sub(ignore_addresses, "", output)

    # Adding echo $? for receiving the exit status
    def show_exit_code(self, flags, compiler_path, filepath):
        return os.popen(f"echo $?").read().rstrip()

    def check(self, flags, filepath: str):
        our_output = self.our(flags, filepath)
        ref_output = self.ref(flags, filepath)
        if self.del_ws_flag:
            ref_output = del_ws(ref_output)
            our_output = del_ws(our_output)
        if ref_output != our_output:
            if args.log_fail:
                with open("logs/fail.txt", "a") as file:
                    file.write(filepath + "\n")
            else:
                msg = f"Expected output was\n{list(ref_output)}\n and Actual output was\n{list(our_output)}"

                print(self.our.__name__)
                self.fail(msg)

    def extension_output(self, flags, filepath):
        comments = [line[1:].strip() for line in open(filepath)
                    if re.match(r'^#.*', line)]
        try:
            target_index = comments.index("Output:")
            our_output = comments[target_index + 1:]
        except ValueError:
            return ""
        return "".join(our_output)

    def extension_code(self, flags, filepath):
        comments = [line[1:].strip()
                    for line in open(filepath) if re.match(r'#.*', line)]
        try:
            target_index = comments.index("Exit Code:")
            our_code = comments[target_index+1]
        except ValueError:
            if "semanticErr" in filepath:
                return "200"
            if "syntaxErr" in filepath:
                return "100"
            our_code = "0"
        return our_code

    def test_chunks(self):
        self.del_ws_flag = False
        TEST_DIR = "./tests"
        EXT_TEST_DIR = "./tests/extensions"
        chunks = sorted(
            list(map(lambda c: f"{TEST_DIR}/{c}", next(os.walk(TEST_DIR))[1])))
        
        logTags = []
        
        if args.extension:
            logTags.append("Extension")
            chunks = sorted(
                list(map(lambda c: f"{EXT_TEST_DIR}/{c}", next(os.walk(EXT_TEST_DIR))[1])))
            self.exe_ref_output = self.extension_output
            self.code_ref_output = self.extension_code

        chunks = chunks[:args.chunk_number+1]
        if args.only:
            chunks = [chunks[args.only]]

        folders = set(["valid"])
        tests = set()

        flags = ""
        if args.parse:
            flags += "-p"
            logTags.append("Syntactic")
            folders.add("syntaxErr")
            tests.add((self.code_ref_output, self.code_our_output))

        if args.tree_ast:
            flags += "-t"
            folders = set(["valid"])
            self.del_ws_flag = True
            tests.add((self.tree_ref_output, self.tree_our_output))

        if args.semantic_analysis:
            flags += "-s"
            logTags.append("Semantic")
            folders.add("semanticErr")
            folders.add("syntaxErr")
            tests.add((self.code_ref_output, self.code_our_output))

        if args.execute:
            folders = set(["valid"])
            logTags.append("Valid")
            flags += "-x"
            tests = [(self.exe_ref_output, self.our_generated_output)]
            tests += [(self.code_ref_output, self.code_our_output)]
        
        # LOGGING INFO
        passed = 0
        total = 0

        for chunk in chunks:
            chunk_no = chunk.split("/")[-1]

            for path, _, files in os.walk(chunk):
                if not path.split("/")[-1] in folders:
                    continue
                for name in files:
                    total+= 1
                    if name.endswith(".wacc"):
                        filepath = os.path.join(path, name)
                        print(f"Tested ./compile {flags} {filepath}")
                        for (ref, our) in tests:
                            self.our = our
                            self.ref = ref
                            self.check(flags, filepath)
                    passed += 1
            print(f"=================================== PASSED {chunk_no.upper()} =======================================")
        with open("logs/log.txt", "a") as file:

            file.write(f"\nPassing {passed}/ {total} tests for {', '.join(logTags)}")

            print(f"Passing {passed}/ {total} tests!")
            print(f"These tests are tagged as {', '.join(logTags)}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--chunk-number", help="everything is tested uptil this chunk folder", action="store", type=int, default=15, choices=range(16))
    parser.add_argument("-o", "--only", help="only this chunk folder is tested", action="store", type=int, choices=range(16))
    parser.add_argument("-p", "--parse", help="Compares the generated parse trees", action="store_true")
    parser.add_argument("-t", "--tree-ast", help="Compares the generated ASTs", action="store_true")
    parser.add_argument("-s", "--semantic-analysis", help="Compares the semantic analysis checks", action="store_true")
    parser.add_argument("-x", "--execute", help="Compares the output of the generated assembly files.", action="store_true")
    parser.add_argument("-e", "--extension", help="Tests the extension", action="store_true")
    parser.add_argument("-lf", "--log-fail", help="run all tests and log the failures", action="store_true")
    
    global args
    args = parser.parse_args()

    assert args.chunk_number >= 0
    suite = unittest.TestSuite()
    suite.addTest(TestChunks("test_chunks"))
    runner = unittest.TextTestRunner()
    result = runner.run(suite)
    if not result.wasSuccessful():
        exit(1)
    exit(0)
