CREATE TABLE bookings (
  id SERIAL PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT,
  arrival_date DATE NOT NULL,
  departure_date DATE NOT NULL
    CHECK (departure_date > arrival_date),
  room_id INTEGER NOT NULL REFERENCES rooms(id),
  status TEXT NOT NULL DEFAULT 'ongoing'
    CHECK (status IN ('processed', 'ongoing', 'cancelled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
