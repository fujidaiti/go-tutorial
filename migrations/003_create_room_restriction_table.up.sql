CREATE TABLE room_restrictions (
  id SERIAL PRIMARY KEY,
  arrival_date DATE NOT NULL,
  departure_date DATE NOT NULL
    CHECK(departure_date > arrival_date),
  room_id INTEGER NOT NULL REFERENCES rooms(id),
  booking_id INTEGER NOT NULL REFERENCES bookings(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);