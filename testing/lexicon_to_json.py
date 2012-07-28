import re
import sys

toplevel_re = re.compile(r' "([^"]+)" \(\d+\)')
position_re = re.compile(r'   (Left|Right) (\d+):')
stat1_re = re.compile(r'    Learned: ([\d\.]+)\s+Blocked: ([\d\.]+)\s+(?:In: (-?[\d\.]+)\s+)?Out: (-?[\d\.]+)\s+')
stat2_re = re.compile(r'    Derived\s+In: (-?[\d\.]+)\s+')
seen_re = re.compile(r'     Seen:')
label_re = re.compile(r'\s+"(\.?)([^"]+)" \((-?[\d\.]+)\)')

file_input = open(sys.argv[1])

def parse_position(token, line):

    # parse adjacency's position
    side, pos = position_re.match(line).groups()
    position = -int(pos) if side == 'Left' else int(pos)

    # parse bootstrapped stats
    line = file_input.readline()
    count, stop, in_raw, out = stat1_re.match(line).groups()
    in_, = stat2_re.match(file_input.readline()).groups()

    print "%s\t%d\tcount\t%s" % (token, position, count)
    print "%s\t%d\tstop\t%s" % (token, position, stop)
    print "%s\t%d\tin_raw\t%s" % (token, position, in_raw if in_raw else 0.0)
    print "%s\t%d\tout\t%s" % (token, position, out)
    print "%s\t%d\tin\t%s" % (token, position, in_ if in_ else 0.0)

    # parse Seen:
    line = file_input.readline()
    if not seen_re.match(line):
        return line

    while True:
        line = file_input.readline()
        found = False

        if toplevel_re.match(line):
            return line

        for m in label_re.finditer(line):
            found = True
            is_adj, label, weight = m.groups()

            if is_adj:
                print "%s\t%d\t\".%s\"\t%s" % (token, position, label, weight)
            else:
                print "%s\t%d\t\"%s\"\t%s" % (token, position, label, weight)

        if not found:
            return line

def parse_toplevel(line):
    token = toplevel_re.match(line).group(1)
    line = file_input.readline()

    while True:
        line = parse_position(token, line)

        if toplevel_re.match(line) or not line:
            return line

line = file_input.readline()
while line:
    line = parse_toplevel(line)

