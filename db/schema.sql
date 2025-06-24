-- TCG API Database Schema

-- Card types lookup table
CREATE TABLE card_types (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- Cards table: base for all card types
CREATE TABLE cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid BINARY(16) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    card_type_id INT NOT NULL,
    front_image_url VARCHAR(500),
    back_image_url VARCHAR(500),
    alt_text VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Playing cards specific attributes  
CREATE TABLE playing_cards (
    id CHAR(36) PRIMARY KEY,
    suit ENUM('hearts', 'diamonds', 'clubs', 'spades') NOT NULL,
    ranking INT NOT NULL
);

-- TODO: Add game_cards table
-- TODO: Add decks table
-- TODO: Add deck_cards relationship table

-- Performance indexes
CREATE INDEX idx_cards_uuid ON cards(uuid);
CREATE INDEX idx_cards_type_id ON cards(card_type_id);
CREATE INDEX idx_cards_created_at ON cards(created_at);
