all:
	time go run . 21 > out1.txt
# real    0m9.732s
# user    0m56.404s
# sys     0m0.535s

# [online version](https://benchmarksgame-team.pages.debian.net/benchmarksgame/program/binarytrees-go-8.html):
	time go run ./19s/main.go 21 > out2.txt
# real    0m21.050s
# user    2m10.221s
# sys     0m2.067s

	diff out1.txt out2.txt

twofold speedup:
# real 	21.050/9.732 = 2.16x
# user  130.221/56.404 = 2.3x
# sys   2.067/0.535 = 3.86x
