
import os
import simplejson
import subprocess

# (Use wget --user-agent "Firefox")
#  http://www.gutenberg.org/dirs/5/55/55.txt

def next_sentence():
    fin = open('wizard_of_oz.txt')
    buf = fin.read(4096)

    while True:

        ind1 = buf.find('\n\n')
        ind2 = buf.find('  ')
        ind3 = buf.find('"  "')

        if ind1 == -1 and ind2 == -1:
            nextbuf = fin.read(4096)
            if nextbuf == '':
                break
            buf += nextbuf
            continue

        ind = ind1 if (ind1 != -1 and (ind1 < ind2 or ind2 == -1)) else ind2

        sentence = buf[:ind]

        sentence = sentence.replace(',', ' , ')
        sentence = sentence.replace('"', ' " ')
        sentence = sentence.replace("'", " '")
        sentence = sentence.replace(';', ' ; ')
        sentence = sentence.replace(':', ' : ')
        sentence = sentence.replace('.', ' . ')
        sentence = sentence.replace('!', ' ! ')
        sentence = sentence.replace('--', ' -- ')

        sentence = sentence.split()
        buf = buf[ind+2:]

        if sentence:
            yield ' '.join(sentence)


def run_ccl_parser(sentences):

    os.system("mv Unnamed.lexicon last_ccl_parser.lexicon")
    os.system("rm Unnamed.*")

    open('/tmp/ccl_input.txt', 'w').write('\n'.join(sentences))
    open('/tmp/ccl_options.txt', 'w').write(
        "TraceBits 65535\nPrintingMode extra_parse\nLexMinPrint 0\nStatisticsTopListMaxLen 1000\nMaxLabels 1000\n")
    open('/tmp/ccl_exec.txt', 'w').write(
        "/tmp/ccl_input.txt line learn+parse -p -G /tmp/ccl_options.txt")

    os.system("/home/johng/cclparser/main/UnknownOS/cclparser /tmp/ccl_exec.txt")

    p = subprocess.Popen("python lexicon_to_json.py Unnamed.lexicon",
            shell = True, stdout = subprocess.PIPE)

    for line in sorted(p.stdout.readlines()):
        yield line.strip()

def run_sibyl_parser(sentences):

    open('/tmp/sibyl_input.txt', 'w').write('\n'.join(sentences))
    p = subprocess.Popen("/home/johng/sibyl/testing/main /tmp/sibyl_input.txt",
        shell = True, stdout = subprocess.PIPE)

    for line in sorted(p.stdout.readlines()):
        yield line.strip()

sentences = []
for ind, sent in enumerate(next_sentence()):
    print sent

"""
    ccl_lines = list(run_ccl_parser(sentences))
    sibyl_lines = list(run_sibyl_parser(sentences))

    if ccl_lines != sibyl_lines:

        for line1, line2 in zip(ccl_lines, sibyl_lines):
            print "%s\t\t%s\t\t%s" % ('' if line1 == line2 else '!!!', line1, line2)
        break
"""

