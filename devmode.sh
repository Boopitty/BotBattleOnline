# Launch server with dev frontend
go build -o server ./server
./server/server
cd ./bot-battle || exit 1
npm run dev