if [ $# -eq 0 ]; then
    echo "Please enter pid"
else
./rsa updatep $1 sample_brand sample_dosage sample_name sample_addr sample_doc 1234567 7
fi