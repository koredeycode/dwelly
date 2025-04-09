-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()

);

CREATE TABLE listings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  intent TEXT NOT NULL, -- 'buy', 'rent', 'sell', 'lease'
  title TEXT NOT NULL,
  description TEXT,
  price_range TEXT, -- e.g., "100k - 150k", flexible for search
  location TEXT NOT NULL,
  category TEXT NOT NULL, -- e.g., apartment, house, land
  status TEXT DEFAULT 'active', -- 'active', 'negotiation', 'completed'
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE listing_images (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
  url TEXT NOT NULL
);

CREATE TABLE inquiries (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  listing_id UUID REFERENCES listings(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
  status TEXT DEFAULT 'open', -- 'open', 'negotiating', 'resolved', 'closed'
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  inquiry_id UUID REFERENCES inquiries(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS inquiries;
DROP TABLE IF EXISTS listing_images;
DROP TABLE IF EXISTS listings;
DROP TABLE IF EXISTS users;
