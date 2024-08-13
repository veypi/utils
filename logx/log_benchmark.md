# 日志性能分析

## zap

```bash

# 输出控制台 带路径
BenchmarkNoContext-16            94934             65178 ns/op             538 B/op         10 allocs/op

# 写入文件
带路径
BenchmarkZapNoContext-16         1043810              5965 ns/op             329 B/op          4 allocs/op
不带路径
BenchmarkZapNoContext-16         1000000              6057 ns/op              32 B/op          1 allocs/op
```

## zerolog

```bash
输出console
    带路径
        232824             26089 ns/op             296 B/op          3 allocs/op
    不带路径
        400622             15201 ns/op               0 B/op          0 allocs/op

写入文件
不带路径
BenchmarkZeroNoContext-16        1305238              4573 ns/op               0 B/op          0 allocs/op
带路径
BenchmarkZeroNoContext-16         735138              7082 ns/op             296 B/op          3 allocs/op

```

```bash
# 原有库测试
goos: darwin
goarch: amd64
pkg: github.com/rs/zerolog
BenchmarkLogEmpty-16            	910646728	        6.33 ns/op	      0 B/op	      0 allocs/op
BenchmarkDisabled-16            	1000000000	        0.485 ns/op	      0 B/op	      0 allocs/op
BenchmarkInfo-16                	320548521	       19.5 ns/op	      0 B/op	      0 allocs/op
BenchmarkContextFields-16       	307581892	       19.8 ns/op	      0 B/op	      0 allocs/op
BenchmarkContextAppend-16       	865624342	        6.77 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFields-16           	72584544	       84.4 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogArrayObject-16      	12815616	      403 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Strs-16   	262370852	       24.1 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Interfaces-16         	24392785	      273 ns/op	    256 B/op	      2 allocs/op
BenchmarkLogFieldType/Interface-16          	60081538	       87.6 ns/op	     80 B/op	      2 allocs/op
BenchmarkLogFieldType/Interface(Objects)-16 	23136542	      291 ns/op	    256 B/op	      2 allocs/op
BenchmarkLogFieldType/Bool-16               	427726760	       13.7 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Bools-16              	342280291	       17.0 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Floats-16             	80542653	       85.9 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Str-16                	414047431	       15.0 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Errs-16               	99231771	       62.2 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Time-16               	98538722	       65.6 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Times-16              	12264505	      470 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Durs-16               	43145167	      141 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Int-16                	368180292	       16.6 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Ints-16               	201922300	       29.3 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Float-16              	287201074	       21.5 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Interface(Object)-16  	126588674	       49.0 ns/op	     48 B/op	      1 allocs/op
BenchmarkLogFieldType/Err-16                	277919796	       21.9 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Dur-16                	272078011	       21.8 ns/op	      0 B/op	      0 allocs/op
BenchmarkLogFieldType/Object-16             	135381682	       45.1 ns/op	     48 B/op	      1 allocs/op
BenchmarkContextFieldType/Interface-16      	25623534	      240 ns/op	    594 B/op	      3 allocs/op
BenchmarkContextFieldType/Interfaces-16     	14194680	      420 ns/op	    770 B/op	      3 allocs/op
BenchmarkContextFieldType/Bool-16           	33644347	      182 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Ints-16           	35774460	      167 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Floats-16         	37536844	      185 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Err-16            	32963114	      172 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Timestamp-16      	36575163	      171 ns/op	    529 B/op	      2 allocs/op
BenchmarkContextFieldType/Float-16          	36080226	      172 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Str-16            	34642272	      169 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Strs-16           	35303484	      166 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Object-16         	32752846	      183 ns/op	    577 B/op	      3 allocs/op
BenchmarkContextFieldType/Bools-16          	36114555	      172 ns/op	    512 B/op	      1 allocs/op
BenchmarkContextFieldType/Errs-16           	36229861	      172 ns/op	    513 B/op	      1 allocs/op
BenchmarkContextFieldType/Durs-16           	30578263	      235 ns/op	    512 B/op	      1 allocs/op
```