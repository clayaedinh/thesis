if [ "$1" = "standard" ] || [ $# -eq 0 ]; then
    go test -bench=BenchmarkStandard -benchtime=3x
elif [ "$1" = "split" ]; then
    go test -bench=BenchmarkSplit -benchtime=10x
elif [ "$1" = "all" ]; then
    go test -bench=Bench -benchtime=10x
elif [ "$1" = "scalereport" ]; then
    go test -bench=BenchmarkPrescriptionAmountAndReportRead -benchtime=100x -timeout=10h
elif [ $# -gt 0 ]; then
    echo "unrecognized argument"
fi

