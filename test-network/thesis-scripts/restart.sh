cd ..
./network.sh down
docker system prune --volumes -f
./network.sh up createChannel -ca -s couchdb