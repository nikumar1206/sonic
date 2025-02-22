# sonic

jsonic but actually called sonic, cuz its supposed to be fast

- SonicConfig: should theoretically be able to be as strict as possible or as lazy as possible

- Parse should be an iter.Seq
- Inverse implementation that uses channels ? can we do goroutines here?
- Marshal/Unmarshal from/to structs

- should we allow like " and ' to close one another?

### Benchmarks

`bash test -bench=. -benchtime=1000000x`
BenchmarkParsing-14 1000000 2790 ns/op 95.32 MB/s 7088 B/op 105 allocs/op
