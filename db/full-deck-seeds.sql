-- Full 52-card deck plus jokers seed data

-- Insert all 52 playing cards
INSERT INTO cards (uuid, name, description, card_type_id, front_image_url, back_image_url, alt_text) VALUES
-- Spades
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Spades', 'Ace of Spades', 2, 'https://example.com/cards/ace-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000002', '-', '')), '2 of Spades', 'Two of Spades', 2, 'https://example.com/cards/2-spades.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000003', '-', '')), '3 of Spades', 'Three of Spades', 2, 'https://example.com/cards/3-spades.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000004', '-', '')), '4 of Spades', 'Four of Spades', 2, 'https://example.com/cards/4-spades.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000005', '-', '')), '5 of Spades', 'Five of Spades', 2, 'https://example.com/cards/5-spades.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000006', '-', '')), '6 of Spades', 'Six of Spades', 2, 'https://example.com/cards/6-spades.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000007', '-', '')), '7 of Spades', 'Seven of Spades', 2, 'https://example.com/cards/7-spades.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000008', '-', '')), '8 of Spades', 'Eight of Spades', 2, 'https://example.com/cards/8-spades.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000009', '-', '')), '9 of Spades', 'Nine of Spades', 2, 'https://example.com/cards/9-spades.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000010', '-', '')), '10 of Spades', 'Ten of Spades', 2, 'https://example.com/cards/10-spades.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Spades', 'Jack of Spades', 2, 'https://example.com/cards/jack-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Spades', 'Queen of Spades', 2, 'https://example.com/cards/queen-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Spades playing card'),
(UNHEX(REPLACE('01000000-0000-0000-0000-000000000013', '-', '')), 'King of Spades', 'King of Spades', 2, 'https://example.com/cards/king-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Spades playing card'),

-- Hearts
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Hearts', 'Ace of hearts', 2, 'https://example.com/cards/ace-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000002', '-', '')), '2 of Hearts', 'Two of Hearts', 2, 'https://example.com/cards/2-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000003', '-', '')), '3 of Hearts', 'Three of Hearts', 2, 'https://example.com/cards/3-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000004', '-', '')), '4 of Hearts', 'Four of Hearts', 2, 'https://example.com/cards/4-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000005', '-', '')), '5 of Hearts', 'Five of Hearts', 2, 'https://example.com/cards/5-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000006', '-', '')), '6 of Hearts', 'Six of Hearts', 2, 'https://example.com/cards/6-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000007', '-', '')), '7 of Hearts', 'Seven of Hearts', 2, 'https://example.com/cards/7-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000008', '-', '')), '8 of Hearts', 'Eight of Hearts', 2, 'https://example.com/cards/8-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000009', '-', '')), '9 of Hearts', 'Nine of Hearts', 2, 'https://example.com/cards/9-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000010', '-', '')), '10 of Hearts', 'Ten of Hearts', 2, 'https://example.com/cards/10-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Hearts', 'Jack of hearts', 2, 'https://example.com/cards/jack-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Hearts', 'Queen of hearts', 2, 'https://example.com/cards/queen-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Hearts playing card'),
(UNHEX(REPLACE('02000000-0000-0000-0000-000000000013', '-', '')), 'King of Hearts', 'King of Hearts', 2, 'https://example.com/cards/king-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Hearts playing card'),

-- Diamonds
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Diamonds', 'Ace of diamonds', 2, 'https://example.com/cards/ace-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000002', '-', '')), '2 of Diamonds', 'Two of Diamonds', 2, 'https://example.com/cards/2-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000003', '-', '')), '3 of Diamonds', 'Three of Diamonds', 2, 'https://example.com/cards/3-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000004', '-', '')), '4 of Diamonds', 'Four of Diamonds', 2, 'https://example.com/cards/4-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000005', '-', '')), '5 of Diamonds', 'Five of Diamonds', 2, 'https://example.com/cards/5-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000006', '-', '')), '6 of Diamonds', 'Six of Diamonds', 2, 'https://example.com/cards/6-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000007', '-', '')), '7 of Diamonds', 'Seven of Diamonds', 2, 'https://example.com/cards/7-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000008', '-', '')), '8 of Diamonds', 'Eight of Diamonds', 2, 'https://example.com/cards/8-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000009', '-', '')), '9 of Diamonds', 'Nine of Diamonds', 2, 'https://example.com/cards/9-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000010', '-', '')), '10 of Diamonds', 'Ten of Diamonds', 2, 'https://example.com/cards/10-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Diamonds', 'Jack of diamonds', 2, 'https://example.com/cards/jack-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Diamonds', 'Queen of Diamonds', 2, 'https://example.com/cards/queen-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Diamonds playing card'),
(UNHEX(REPLACE('03000000-0000-0000-0000-000000000013', '-', '')), 'King of Diamonds', 'King of Diamonds', 2, 'https://example.com/cards/king-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Diamonds playing card'),

-- Clubs
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Clubs', 'Ace of clubs', 2, 'https://example.com/cards/ace-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000002', '-', '')), '2 of Clubs', 'Two of Clubs', 2, 'https://example.com/cards/2-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000003', '-', '')), '3 of Clubs', 'Three of Clubs', 2, 'https://example.com/cards/3-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000004', '-', '')), '4 of Clubs', 'Four of Clubs', 2, 'https://example.com/cards/4-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000005', '-', '')), '5 of Clubs', 'Five of Clubs', 2, 'https://example.com/cards/5-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000006', '-', '')), '6 of Clubs', 'Six of Clubs', 2, 'https://example.com/cards/6-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000007', '-', '')), '7 of Clubs', 'Seven of Clubs', 2, 'https://example.com/cards/7-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000008', '-', '')), '8 of Clubs', 'Eight of Clubs', 2, 'https://example.com/cards/8-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000009', '-', '')), '9 of Clubs', 'Nine of Clubs', 2, 'https://example.com/cards/9-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000010', '-', '')), '10 of Clubs', 'Ten of Clubs', 2, 'https://example.com/cards/10-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Clubs', 'Jack of clubs', 2, 'https://example.com/cards/jack-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Clubs', 'Queen of clubs', 2, 'https://example.com/cards/queen-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Clubs playing card'),
(UNHEX(REPLACE('04000000-0000-0000-0000-000000000013', '-', '')), 'King of Clubs', 'King of clubs', 2, 'https://example.com/cards/king-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Clubs playing card'),

-- Jokers (as image cards)
(UNHEX(REPLACE('99000000-0000-0000-0000-000000000001', '-', '')), 'Red Joker', 'Red Wild Card', 1, 'https://example.com/cards/joker-red.jpg', 'https://example.com/cards/card-back.jpg', 'Red Joker wild card'),
(UNHEX(REPLACE('99000000-0000-0000-0000-000000000002', '-', '')), 'Black Joker', 'Black Wild Card', 1, 'https://example.com/cards/joker-black.jpg', 'https://example.com/cards/card-back.jpg', 'Black Joker wild card');

-- Insert playing card specific data for the 52 cards
INSERT INTO playing_cards (id, suit, ranking) VALUES
-- Spades
('01000000-0000-0000-0000-000000000001', 'spades', 1),
('01000000-0000-0000-0000-000000000002', 'spades', 2),
('01000000-0000-0000-0000-000000000003', 'spades', 3),
('01000000-0000-0000-0000-000000000004', 'spades', 4),
('01000000-0000-0000-0000-000000000005', 'spades', 5),
('01000000-0000-0000-0000-000000000006', 'spades', 6),
('01000000-0000-0000-0000-000000000007', 'spades', 7),
('01000000-0000-0000-0000-000000000008', 'spades', 8),
('01000000-0000-0000-0000-000000000009', 'spades', 9),
('01000000-0000-0000-0000-000000000010', 'spades', 10),
('01000000-0000-0000-0000-000000000011', 'spades', 11),
('01000000-0000-0000-0000-000000000012', 'spades', 12),
('01000000-0000-0000-0000-000000000013', 'spades', 13),

-- Hearts
('02000000-0000-0000-0000-000000000001', 'hearts', 1),
('02000000-0000-0000-0000-000000000002', 'hearts', 2),
('02000000-0000-0000-0000-000000000003', 'hearts', 3),
('02000000-0000-0000-0000-000000000004', 'hearts', 4),
('02000000-0000-0000-0000-000000000005', 'hearts', 5),
('02000000-0000-0000-0000-000000000006', 'hearts', 6),
('02000000-0000-0000-0000-000000000007', 'hearts', 7),
('02000000-0000-0000-0000-000000000008', 'hearts', 8),
('02000000-0000-0000-0000-000000000009', 'hearts', 9),
('02000000-0000-0000-0000-000000000010', 'hearts', 10),
('02000000-0000-0000-0000-000000000011', 'hearts', 11),
('02000000-0000-0000-0000-000000000012', 'hearts', 12),
('02000000-0000-0000-0000-000000000013', 'hearts', 13),

-- Diamonds
('03000000-0000-0000-0000-000000000001', 'diamonds', 1),
('03000000-0000-0000-0000-000000000002', 'diamonds', 2),
('03000000-0000-0000-0000-000000000003', 'diamonds', 3),
('03000000-0000-0000-0000-000000000004', 'diamonds', 4),
('03000000-0000-0000-0000-000000000005', 'diamonds', 5),
('03000000-0000-0000-0000-000000000006', 'diamonds', 6),
('03000000-0000-0000-0000-000000000007', 'diamonds', 7),
('03000000-0000-0000-0000-000000000008', 'diamonds', 8),
('03000000-0000-0000-0000-000000000009', 'diamonds', 9),
('03000000-0000-0000-0000-000000000010', 'diamonds', 10),
('03000000-0000-0000-0000-000000000011', 'diamonds', 11),
('03000000-0000-0000-0000-000000000012', 'diamonds', 12),
('03000000-0000-0000-0000-000000000013', 'diamonds', 13),

-- Clubs
('04000000-0000-0000-0000-000000000001', 'clubs', 1),
('04000000-0000-0000-0000-000000000002', 'clubs', 2),
('04000000-0000-0000-0000-000000000003', 'clubs', 3),
('04000000-0000-0000-0000-000000000004', 'clubs', 4),
('04000000-0000-0000-0000-000000000005', 'clubs', 5),
('04000000-0000-0000-0000-000000000006', 'clubs', 6),
('04000000-0000-0000-0000-000000000007', 'clubs', 7),
('04000000-0000-0000-0000-000000000008', 'clubs', 8),
('04000000-0000-0000-0000-000000000009', 'clubs', 9),
('04000000-0000-0000-0000-000000000010', 'clubs', 10),
('04000000-0000-0000-0000-000000000011', 'clubs', 11),
('04000000-0000-0000-0000-000000000012', 'clubs', 12),
('04000000-0000-0000-0000-000000000013', 'clubs', 13);
