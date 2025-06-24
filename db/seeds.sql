-- TCG API Sample Data

-- Card types lookup data
INSERT INTO card_types (id, name, description) VALUES
(1, 'image', 'Image cards for artwork and illustrations'),
(2, 'playing', 'Traditional playing cards with suit and ranking');

-- Sample image cards (stored entirely in cards table)
INSERT INTO cards (uuid, name, description, card_type_id, front_image_url, back_image_url, alt_text) VALUES
(UNHEX(REPLACE('550e8400-e29b-41d4-a716-446655440001', '-', '')), 'Dragon Artwork', 'Beautiful dragon illustration card', 1, 'https://example.com/images/dragon-front.jpg', 'https://example.com/images/card-back.jpg', 'Red dragon breathing fire'),
(UNHEX(REPLACE('550e8400-e29b-41d4-a716-446655440002', '-', '')), 'Castle Scene', 'Medieval castle landscape artwork', 1, 'https://example.com/images/castle-front.jpg', 'https://example.com/images/card-back.jpg', 'Stone castle on hill'),
(UNHEX(REPLACE('550e8400-e29b-41d4-a716-446655440003', '-', '')), 'Phoenix Rising', 'Majestic phoenix in flames', 1, 'https://example.com/images/phoenix-front.jpg', 'https://example.com/images/card-back.jpg', 'Phoenix with spread wings');
