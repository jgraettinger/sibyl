import re
import sys
import simplejson

toplevel_re = re.compile(r' "([^"]+)" \(\d+\)')
position_re = re.compile(r'   (Left|Right) (\d+):')
stat1_re = re.compile(r'    Learned: ([\d\.]+)\s+Blocked: ([\d\.]+)\s+In: (-?[\d\.]+)\s+Out: (-?[\d\.]+)\s+')
stat2_re = re.compile(r'    Derived\s+In: (-?[\d\.]+)\s+')
seen_re = re.compile(r'     Seen:')
label_re = re.compile(r'\s+"(\.?)([^"]+)" \((-?[\d\.]+)\)')

def parse_position(token, line):

    adj_point = {'token': token}

    # parse adjacency's position
    side, pos = position_re.match(line).groups()
    adj_point['position'] = -int(pos) if side == 'Left' else int(pos)

    # parse bootstrapped stats
    line = sys.stdin.readline()
    count, stop, in_raw, out = stat1_re.match(line).groups()
    in_, = stat2_re.match(sys.stdin.readline()).groups()

    adj_point['update_count'] = int(float(count))
    adj_point['stop'] = int(float(stop))
    adj_point['in_raw'] = float(in_raw)
    adj_point['out'] = float(out)
    adj_point['in'] = float(in_)

    labels = adj_point['labels'] = {}

    # parse Seen:
    line = sys.stdin.readline()
    if not seen_re.match(line):
        print simplejson.dumps(adj_point)
        return line

    while True:
        line = sys.stdin.readline()
        found = False

        if toplevel_re.match(line):
            print simplejson.dumps(adj_point)
            return line

        for m in label_re.finditer(line):
            found = True
            is_adj, token, weight = m.groups()

            label = labels.setdefault(token,
                {'class_weight': 0.0, 'adjacency_weight': 0.0})

            if is_adj:
                label['adjacency_weight'] = float(weight)
            else:
                label['class_weight'] = float(weight)

        if not found:
            print simplejson.dumps(adj_point)
            return line

def parse_toplevel(line):
    token = toplevel_re.match(line).group(1)
    line = sys.stdin.readline()

    while True:
        line = parse_position(token, line)

        if toplevel_re.match(line) or not line:
            return line

line = sys.stdin.readline()
while line:
    line = parse_toplevel(line)

