# Fortuna

 [![GoDoc](https://godoc.org/github.com/jacoelho/fortuna?status.svg)](http://godoc.org/github.com/jacoelho/fortuna)

Learning implementation of fortuna algorithm described in cryptography engineering.

Not suitable for real world application.

## fortunactl

simple http oracle service.

Numbers
```
curl "http://localhost:8080/numbers?min=1&max=55&count=7"
48
48
51
30
19
51
29
```

Sequence (without duplicates)
```
curl "http://localhost:8080/sequence?min=1&max=55&count=7"
27
10
51
26
49
12
5
```


## fortunatst

generate file filled with random data to be used with [dieharder](https://webhome.phy.duke.edu/~rgb/General/dieharder.php)

```
dieharder -a -t 10000000 -k 1 -f test.dat| tee test.report
```