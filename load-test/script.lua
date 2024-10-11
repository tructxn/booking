-- dynamic_load_test.lua

-- Function to generate a random hotel ID
local function randomHotelID()
    return "ID_hotel_" .. math.random(0, 5)  -- Random hotel ID from ID_hotel_0 to ID_hotel_10
end

-- Function to generate a random room type ID
local function randomRoomTypeID()
    return math.random(1, 10)  -- Random RoomTypeId between 1 and 20
end

-- Function to generate a random check-in and check-out date
-- Function to generate a random check-in and check-out date
local function randomDates()
    local checkInDate = os.date("!%Y-%m-%dT%H:%M:%S.%f+07:00", os.time() + math.random(0, 86400))  -- Correctly format the check-in date
    local checkOutDate = os.date("!%Y-%m-%dT%H:%M:%S.%f+07:00", os.time() + math.random(86401, 172800))  -- Correctly format the check-out date
    return checkInDate, checkOutDate
end


wrk.method = "POST"  -- Specify the request method as POST

-- Function to prepare the request body
local function prepareRequestBody()
    local hotelID = randomHotelID()
    local checkInDate, checkOutDate = randomDates()
    local roomTypeID = randomRoomTypeID()
    -- Create the JSON body for the request
    return string.format(
        '{"HotelID": "%s","CheckInDate": "2024-10-12T16:43:24.527031+07:00","CheckOutDate": "2024-10-15T16:43:24.527031+07:00","RoomTypeId": %d,"description": ""}',
        hotelID, roomTypeID
    )
end

wrk.body = prepareRequestBody()  -- Set the request body
wrk.headers["Content-Type"] = "application/json"  -- Set the content type to JSON
