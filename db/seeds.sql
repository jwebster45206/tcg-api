-- TCG API Sample Data

-- Card types lookup data
INSERT INTO card_types (id, name, description) VALUES
(1, 'image', 'Image cards for artwork and illustrations'),
(2, 'playing', 'Traditional playing cards with suit and rank');

-- Sample image cards (stored entirely in cards table)
INSERT INTO cards (id, name, description, card_type_id, image_url, alt_text) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Dragon Artwork', 'Beautiful dragon illustration card', 1, 'https://example.com/images/dragon.jpg', 'Red dragon breathing fire'),
('550e8400-e29b-41d4-a716-446655440002', 'Castle Scene', 'Medieval castle landscape artwork', 1, 'https://example.com/images/castle.jpg', 'Stone castle on hill'),
('550e8400-e29b-41d4-a716-446655440003', 'Phoenix Rising', 'Majestic phoenix in flames', 1, 'https://example.com/images/phoenix.jpg', 'Phoenix with spread wings');

-- Sample playing cards (with image properties for consistency)
INSERT INTO cards (id, name, description, card_type_id, image_url, alt_text) VALUES
('550e8400-e29b-41d4-a716-446655440010', 'Ace of Spades', 'The death card', 2, 'https://example.com/cards/ace-spades.jpg', 'Ace of Spades playing card'),
('550e8400-e29b-41d4-a716-446655440011', 'King of Hearts', 'The suicide king', 2, 'https://example.com/cards/king-hearts.jpg', 'King of Hearts playing card'),
('550e8400-e29b-41d4-a716-446655440012', 'Queen of Diamonds', 'The lady of wealth', 2, 'https://example.com/cards/queen-diamonds.jpg', 'Queen of Diamonds playing card'),
('550e8400-e29b-41d4-a716-446655440013', 'Jack of Clubs', 'The knave of clubs', 2, 'https://example.com/cards/jack-clubs.jpg', 'Jack of Clubs playing card');

INSERT INTO playing_cards (id, suit, rank) VALUES
('550e8400-e29b-41d4-a716-446655440010', 'spades', 1),     -- Ace of Spades
('550e8400-e29b-41d4-a716-446655440011', 'hearts', 13),    -- King of Hearts
('550e8400-e29b-41d4-a716-446655440012', 'diamonds', 12),  -- Queen of Diamonds
('550e8400-e29b-41d4-a716-446655440013', 'clubs', 11);     -- Jack of Clubs
