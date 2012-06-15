
import os

# (Use wget --user-agent "Firefox")
#  http://www.gutenberg.org/dirs/5/55/55.txt

def next_sentence():
    fin = open('wizard_of_oz.txt')
    buf = fin.read(4096)

    while True:
    
        ind1 = buf.find('\r\n\r\n')
        ind2 = buf.find('  ')
    
        if ind1 == -1 and ind2 == -1:
            buf += fin.read(4096)
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
        "TraceBits 65535\nPrintingMode extra_parse\nLexMinPrint 0")
    open('/tmp/ccl_exec.txt', 'w').write(
        "/tmp/ccl_input.txt line learn+parse -p -G /tmp/ccl_options.txt")

    os.system("/home/johng/cclparser/main/UnknownOS/cclparser /tmp/ccl_exec.txt")

def run_sibyl_parser(sentences):

    open('/tmp/sibyl_input.txt', 'w').write('\n'.join(sentences))
    os.system("/home/johng/sibyl/testing/main /tmp/sibyl_input.txt")


sentences = []
for ind, sent in enumerate(next_sentence()):
    sentences.append(sent)
    if ind == 10:
        break

run_ccl_parser(sentences)
run_sibyl_parser(sentences)

