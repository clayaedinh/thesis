if [ $# -eq 0 ]; then
    echo "Please enter pid"
else
./rsa -user=CaliperAdmin updatep $1 DRUG_BRAND DRUG_DOSAGE DRUG_NAME DRUG_ADDR DRUG_DOC 1234567 7
fi