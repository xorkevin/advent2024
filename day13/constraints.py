import sys
import z3


lines = [[int(i) for i in l.rstrip().split(",")] for l in sys.stdin]

total = 0

for line in lines:
    a = z3.Int("a")
    b = z3.Int("b")

    s = z3.Optimize()
    s.add(a >= 0)
    s.add(b >= 0)
    s.add(line[0] == a * line[2] + b * line[4])
    s.add(line[1] == a * line[3] + b * line[5])

    s.minimize(3 * a + b)

    if s.check() == z3.sat:
        model = s.model()
        total += model.eval(3 * a + b).as_long()

print("Part 2:", total)
