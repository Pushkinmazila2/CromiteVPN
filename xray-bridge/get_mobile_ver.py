import sys, json

data = sys.stdin.read()
for line in data.split('\n}\n'):
    try:
        m = json.loads(line + '}')
        if m.get('Path') == 'golang.org/x/mobile':
            print(m['Version'])
            break
    except:
        pass