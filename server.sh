# Build server binary and run it
cd ./bot-battle
npm run build
cd ..
go build -o server ./server
./server/server