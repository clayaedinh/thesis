if [ $# -eq 0 ]; then
    echo "Please enter pid"
else
./basicb64 -user=user0003 -org=org2 -port=localhost:9051 setfillp $1 7
fi