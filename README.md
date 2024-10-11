# go-rest-api

# Book room
``curl -X POST http://localhost:3000/booking/ \
-H "Content-Type: application/json" \
-d '{
"HotelID": "ID_hotel_0",
"CheckInDate": "2024-10-12T16:43:24.527031+07:00",
"CheckOutDate": "2024-10-15T16:43:24.527031+07:00",
"RoomTypeId": 10,
"description": ""
}'
``
## Load Test book room
``wrk -t2 -c4 -d5s -s script.lua http://localhost:3000/booking/``

# Run 
``docker-compose up --build -d `