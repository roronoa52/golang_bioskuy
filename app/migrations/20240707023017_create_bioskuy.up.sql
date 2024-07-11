CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('user', 'admin', 'super admin');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE movie_status AS ENUM ('notShowing', 'showing', 'alreadyShowing');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE seat_booking_status AS ENUM ('active', 'pending', 'success');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE payment_status AS ENUM ('paid', 'unpaid');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    token VARCHAR,
    role user_role DEFAULT 'user'
);

CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    name VARCHAR NOT NULL
);

CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    title VARCHAR NOT NULL,
    description VARCHAR,
    price INT NOT NULL,
    duration INT NOT NULL,
    status movie_status DEFAULT 'notShowing'
);

CREATE TABLE studios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    name VARCHAR NOT NULL,
    capacity INT
);

CREATE TABLE seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    seat_name VARCHAR NOT NULL,
    isAvailable BOOL NOT NULL DEFAULT true,
    studio_id UUID NOT NULL,
    FOREIGN KEY (studio_id) REFERENCES studios(id)
);

CREATE TABLE showtimes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    studio_id UUID NOT NULL,
    movie_id UUID NOT NULL,
    show_start TIMESTAMP NOT NULL,
    show_end TIMESTAMP NOT NULL,
    FOREIGN KEY (studio_id) REFERENCES studios(id),
    FOREIGN KEY (movie_id) REFERENCES movies(id)
);

CREATE TABLE seat_bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    status seat_booking_status DEFAULT 'pending' NOT NULL,
    user_id UUID NOT NULL,
    showtime_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (showtime_id) REFERENCES showtimes(id)
);

CREATE TABLE seat_detail_for_bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    seatBooking_id UUID,
    seat_id UUID,
    FOREIGN KEY (seatBooking_id) REFERENCES seat_bookings(id),
    FOREIGN KEY (seat_id) REFERENCES seats(id)
);

CREATE TABLE genre_to_movies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    genre_id UUID,
    movie_id UUID,
    FOREIGN KEY (genre_id) REFERENCES genres(id),
    FOREIGN KEY (movie_id) REFERENCES movies(id)
);

-- Modified payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NOT NULL,
    seatdetailforbooking_id UUID NOT NULL,
    total_seat INT NOT NULL,
    total_price INT NOT NULL,
    status payment_status DEFAULT 'unpaid' NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (seatdetailforbooking_id) REFERENCES seat_detail_for_bookings(id)
);

-- Modified payments_history table
CREATE TABLE payments_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NOT NULL,
    seatdetailforbooking_id UUID NOT NULL,
    title VARCHAR NOT NULL,
    status payment_status DEFAULT 'unpaid' NOT NULL,
    price INT NOT NULL,
    duration INT NOT NULL,
    studio_name VARCHAR NOT NULL,
    show_start TIMESTAMP,
    show_end TIMESTAMP,
    total_seat INT NOT NULL,
    total_price INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (seatdetailforbooking_id) REFERENCES seat_detail_for_bookings(id)
);
