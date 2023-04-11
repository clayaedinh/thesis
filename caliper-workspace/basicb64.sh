maxTps() {
    npx caliper launch manager \
        --caliper-workspace . \
        --caliper-benchconfig benchmarks/basicb64/config-max-tps.yaml \
        --caliper-networkconfig networks/thesis-network.yaml
}

fixedTps() {
    npx caliper launch manager \
        --caliper-workspace . \
        --caliper-benchconfig benchmarks/basicb64/config.yaml \
        --caliper-networkconfig networks/thesis-network.yaml
}

if [ "$1" = "max" ]; then
    maxTps
elif [ "$1" = "fixed" ]; then
    fixedTps
else
    echo "arguments: max | fixed"
fi