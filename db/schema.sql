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
    deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Playing cards specific attributes  
CREATE TABLE playing_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_id INT NOT NULL,
    suit ENUM('hearts', 'diamonds', 'clubs', 'spades') NOT NULL,
    ranking INT NOT NULL,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
);

-- Deck types lookup table
CREATE TABLE deck_types (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- Decks table
CREATE TABLE decks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid BINARY(16) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type_id INT NOT NULL DEFAULT 1,
    sleeve_image_url VARCHAR(500),
    deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (type_id) REFERENCES deck_types(id)
);

-- Deck cards relationship table
CREATE TABLE deck_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    deck_id INT NOT NULL,
    card_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
);

CREATE INDEX idx_cards_uuid ON cards(uuid);
CREATE INDEX idx_cards_type_id ON cards(card_type_id);
CREATE INDEX idx_cards_created_at ON cards(created_at);
CREATE INDEX idx_cards_updated_at ON cards(updated_at);
CREATE INDEX idx_playing_cards_card_id ON playing_cards(card_id);
CREATE INDEX idx_decks_uuid ON decks(uuid);
CREATE INDEX idx_decks_type_id ON decks(type_id);
CREATE INDEX idx_deck_cards_deck_id ON deck_cards(deck_id);
CREATE INDEX idx_deck_cards_card_id ON deck_cards(card_id);
