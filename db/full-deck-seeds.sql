-- Full 52-card deck plus jokers seed data

-- Insert all 52 playing cards + 2 jokers
INSERT INTO cards (id, uuid, name, description, card_type_id, front_image_url, back_image_url, alt_text) VALUES
-- Spades (IDs 1-13)
(1, UNHEX(REPLACE('01000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Spades', 'Ace of Spades', 2, 'https://example.com/cards/ace-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Spades playing card'),
(2, UNHEX(REPLACE('01000000-0000-0000-0000-000000000002', '-', '')), '2 of Spades', 'Two of Spades', 2, 'https://example.com/cards/2-spades.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Spades playing card'),
(3, UNHEX(REPLACE('01000000-0000-0000-0000-000000000003', '-', '')), '3 of Spades', 'Three of Spades', 2, 'https://example.com/cards/3-spades.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Spades playing card'),
(4, UNHEX(REPLACE('01000000-0000-0000-0000-000000000004', '-', '')), '4 of Spades', 'Four of Spades', 2, 'https://example.com/cards/4-spades.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Spades playing card'),
(5, UNHEX(REPLACE('01000000-0000-0000-0000-000000000005', '-', '')), '5 of Spades', 'Five of Spades', 2, 'https://example.com/cards/5-spades.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Spades playing card'),
(6, UNHEX(REPLACE('01000000-0000-0000-0000-000000000006', '-', '')), '6 of Spades', 'Six of Spades', 2, 'https://example.com/cards/6-spades.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Spades playing card'),
(7, UNHEX(REPLACE('01000000-0000-0000-0000-000000000007', '-', '')), '7 of Spades', 'Seven of Spades', 2, 'https://example.com/cards/7-spades.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Spades playing card'),
(8, UNHEX(REPLACE('01000000-0000-0000-0000-000000000008', '-', '')), '8 of Spades', 'Eight of Spades', 2, 'https://example.com/cards/8-spades.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Spades playing card'),
(9, UNHEX(REPLACE('01000000-0000-0000-0000-000000000009', '-', '')), '9 of Spades', 'Nine of Spades', 2, 'https://example.com/cards/9-spades.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Spades playing card'),
(10, UNHEX(REPLACE('01000000-0000-0000-0000-000000000010', '-', '')), '10 of Spades', 'Ten of Spades', 2, 'https://example.com/cards/10-spades.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Spades playing card'),
(11, UNHEX(REPLACE('01000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Spades', 'Jack of Spades', 2, 'https://example.com/cards/jack-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Spades playing card'),
(12, UNHEX(REPLACE('01000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Spades', 'Queen of Spades', 2, 'https://example.com/cards/queen-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Spades playing card'),
(13, UNHEX(REPLACE('01000000-0000-0000-0000-000000000013', '-', '')), 'King of Spades', 'King of Spades', 2, 'https://example.com/cards/king-spades.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Spades playing card'),

-- Hearts (IDs 14-26)
(14, UNHEX(REPLACE('02000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Hearts', 'Ace of Hearts', 2, 'https://example.com/cards/ace-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Hearts playing card'),
(15, UNHEX(REPLACE('02000000-0000-0000-0000-000000000002', '-', '')), '2 of Hearts', 'Two of Hearts', 2, 'https://example.com/cards/2-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Hearts playing card'),
(16, UNHEX(REPLACE('02000000-0000-0000-0000-000000000003', '-', '')), '3 of Hearts', 'Three of Hearts', 2, 'https://example.com/cards/3-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Hearts playing card'),
(17, UNHEX(REPLACE('02000000-0000-0000-0000-000000000004', '-', '')), '4 of Hearts', 'Four of Hearts', 2, 'https://example.com/cards/4-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Hearts playing card'),
(18, UNHEX(REPLACE('02000000-0000-0000-0000-000000000005', '-', '')), '5 of Hearts', 'Five of Hearts', 2, 'https://example.com/cards/5-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Hearts playing card'),
(19, UNHEX(REPLACE('02000000-0000-0000-0000-000000000006', '-', '')), '6 of Hearts', 'Six of Hearts', 2, 'https://example.com/cards/6-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Hearts playing card'),
(20, UNHEX(REPLACE('02000000-0000-0000-0000-000000000007', '-', '')), '7 of Hearts', 'Seven of Hearts', 2, 'https://example.com/cards/7-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Hearts playing card'),
(21, UNHEX(REPLACE('02000000-0000-0000-0000-000000000008', '-', '')), '8 of Hearts', 'Eight of Hearts', 2, 'https://example.com/cards/8-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Hearts playing card'),
(22, UNHEX(REPLACE('02000000-0000-0000-0000-000000000009', '-', '')), '9 of Hearts', 'Nine of Hearts', 2, 'https://example.com/cards/9-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Hearts playing card'),
(23, UNHEX(REPLACE('02000000-0000-0000-0000-000000000010', '-', '')), '10 of Hearts', 'Ten of Hearts', 2, 'https://example.com/cards/10-hearts.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Hearts playing card'),
(24, UNHEX(REPLACE('02000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Hearts', 'Jack of Hearts', 2, 'https://example.com/cards/jack-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Hearts playing card'),
(25, UNHEX(REPLACE('02000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Hearts', 'Queen of Hearts', 2, 'https://example.com/cards/queen-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Hearts playing card'),
(26, UNHEX(REPLACE('02000000-0000-0000-0000-000000000013', '-', '')), 'King of Hearts', 'King of Hearts', 2, 'https://example.com/cards/king-hearts.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Hearts playing card'),

-- Diamonds (IDs 27-39)
(27, UNHEX(REPLACE('03000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Diamonds', 'Ace of Diamonds', 2, 'https://example.com/cards/ace-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Diamonds playing card'),
(28, UNHEX(REPLACE('03000000-0000-0000-0000-000000000002', '-', '')), '2 of Diamonds', 'Two of Diamonds', 2, 'https://example.com/cards/2-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Diamonds playing card'),
(29, UNHEX(REPLACE('03000000-0000-0000-0000-000000000003', '-', '')), '3 of Diamonds', 'Three of Diamonds', 2, 'https://example.com/cards/3-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Diamonds playing card'),
(30, UNHEX(REPLACE('03000000-0000-0000-0000-000000000004', '-', '')), '4 of Diamonds', 'Four of Diamonds', 2, 'https://example.com/cards/4-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Diamonds playing card'),
(31, UNHEX(REPLACE('03000000-0000-0000-0000-000000000005', '-', '')), '5 of Diamonds', 'Five of Diamonds', 2, 'https://example.com/cards/5-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Diamonds playing card'),
(32, UNHEX(REPLACE('03000000-0000-0000-0000-000000000006', '-', '')), '6 of Diamonds', 'Six of Diamonds', 2, 'https://example.com/cards/6-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Diamonds playing card'),
(33, UNHEX(REPLACE('03000000-0000-0000-0000-000000000007', '-', '')), '7 of Diamonds', 'Seven of Diamonds', 2, 'https://example.com/cards/7-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Diamonds playing card'),
(34, UNHEX(REPLACE('03000000-0000-0000-0000-000000000008', '-', '')), '8 of Diamonds', 'Eight of Diamonds', 2, 'https://example.com/cards/8-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Diamonds playing card'),
(35, UNHEX(REPLACE('03000000-0000-0000-0000-000000000009', '-', '')), '9 of Diamonds', 'Nine of Diamonds', 2, 'https://example.com/cards/9-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Diamonds playing card'),
(36, UNHEX(REPLACE('03000000-0000-0000-0000-000000000010', '-', '')), '10 of Diamonds', 'Ten of Diamonds', 2, 'https://example.com/cards/10-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Diamonds playing card'),
(37, UNHEX(REPLACE('03000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Diamonds', 'Jack of Diamonds', 2, 'https://example.com/cards/jack-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Diamonds playing card'),
(38, UNHEX(REPLACE('03000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Diamonds', 'Queen of Diamonds', 2, 'https://example.com/cards/queen-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Diamonds playing card'),
(39, UNHEX(REPLACE('03000000-0000-0000-0000-000000000013', '-', '')), 'King of Diamonds', 'King of Diamonds', 2, 'https://example.com/cards/king-diamonds.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Diamonds playing card'),

-- Clubs (IDs 40-52)
(40, UNHEX(REPLACE('04000000-0000-0000-0000-000000000001', '-', '')), 'Ace of Clubs', 'Ace of Clubs', 2, 'https://example.com/cards/ace-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'Ace of Clubs playing card'),
(41, UNHEX(REPLACE('04000000-0000-0000-0000-000000000002', '-', '')), '2 of Clubs', 'Two of Clubs', 2, 'https://example.com/cards/2-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '2 of Clubs playing card'),
(42, UNHEX(REPLACE('04000000-0000-0000-0000-000000000003', '-', '')), '3 of Clubs', 'Three of Clubs', 2, 'https://example.com/cards/3-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '3 of Clubs playing card'),
(43, UNHEX(REPLACE('04000000-0000-0000-0000-000000000004', '-', '')), '4 of Clubs', 'Four of Clubs', 2, 'https://example.com/cards/4-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '4 of Clubs playing card'),
(44, UNHEX(REPLACE('04000000-0000-0000-0000-000000000005', '-', '')), '5 of Clubs', 'Five of Clubs', 2, 'https://example.com/cards/5-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '5 of Clubs playing card'),
(45, UNHEX(REPLACE('04000000-0000-0000-0000-000000000006', '-', '')), '6 of Clubs', 'Six of Clubs', 2, 'https://example.com/cards/6-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '6 of Clubs playing card'),
(46, UNHEX(REPLACE('04000000-0000-0000-0000-000000000007', '-', '')), '7 of Clubs', 'Seven of Clubs', 2, 'https://example.com/cards/7-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '7 of Clubs playing card'),
(47, UNHEX(REPLACE('04000000-0000-0000-0000-000000000008', '-', '')), '8 of Clubs', 'Eight of Clubs', 2, 'https://example.com/cards/8-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '8 of Clubs playing card'),
(48, UNHEX(REPLACE('04000000-0000-0000-0000-000000000009', '-', '')), '9 of Clubs', 'Nine of Clubs', 2, 'https://example.com/cards/9-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '9 of Clubs playing card'),
(49, UNHEX(REPLACE('04000000-0000-0000-0000-000000000010', '-', '')), '10 of Clubs', 'Ten of Clubs', 2, 'https://example.com/cards/10-clubs.jpg', 'https://example.com/cards/playing-back.jpg', '10 of Clubs playing card'),
(50, UNHEX(REPLACE('04000000-0000-0000-0000-000000000011', '-', '')), 'Jack of Clubs', 'Jack of Clubs', 2, 'https://example.com/cards/jack-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'Jack of Clubs playing card'),
(51, UNHEX(REPLACE('04000000-0000-0000-0000-000000000012', '-', '')), 'Queen of Clubs', 'Queen of Clubs', 2, 'https://example.com/cards/queen-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'Queen of Clubs playing card'),
(52, UNHEX(REPLACE('04000000-0000-0000-0000-000000000013', '-', '')), 'King of Clubs', 'King of Clubs', 2, 'https://example.com/cards/king-clubs.jpg', 'https://example.com/cards/playing-back.jpg', 'King of Clubs playing card'),

-- Jokers (IDs 53-54, as image cards)
(53, UNHEX(REPLACE('99000000-0000-0000-0000-000000000001', '-', '')), 'Red Joker', 'Red Wild Card', 1, 'https://example.com/cards/joker-red.jpg', 'https://example.com/cards/card-back.jpg', 'Red Joker wild card'),
(54, UNHEX(REPLACE('99000000-0000-0000-0000-000000000002', '-', '')), 'Black Joker', 'Black Wild Card', 1, 'https://example.com/cards/joker-black.jpg', 'https://example.com/cards/card-back.jpg', 'Black Joker wild card');

-- Insert playing card specific data for the 52 cards
-- Using explicit card IDs for clear relationships
INSERT INTO playing_cards (card_id, suit, ranking) VALUES
-- Spades (card_id 1-13)
(1, 'spades', 1), 
(2, 'spades', 2), 
(3, 'spades', 3), 
(4, 'spades', 4),
(5, 'spades', 5), 
(6, 'spades', 6), 
(7, 'spades', 7), 
(8, 'spades', 8),
(9, 'spades', 9), 
(10, 'spades', 10), 
(11, 'spades', 11), 
(12, 'spades', 12), 
(13, 'spades', 13),

-- Hearts (card_id 14-26)
(14, 'hearts', 1), 
(15, 'hearts', 2), 
(16, 'hearts', 3), 
(17, 'hearts', 4),
(18, 'hearts', 5), 
(19, 'hearts', 6), 
(20, 'hearts', 7), 
(21, 'hearts', 8),
(22, 'hearts', 9), 
(23, 'hearts', 10), 
(24, 'hearts', 11), 
(25, 'hearts', 12), 
(26, 'hearts', 13),

-- Diamonds (card_id 27-39)
(27, 'diamonds', 1), 
(28, 'diamonds', 2), 
(29, 'diamonds', 3), 
(30, 'diamonds', 4),
(31, 'diamonds', 5),
(32, 'diamonds', 6), 
(33, 'diamonds', 7), 
(34, 'diamonds', 8),
(35, 'diamonds', 9), 
(36, 'diamonds', 10), 
(37, 'diamonds', 11), 
(38, 'diamonds', 12), 

(39, 'diamonds', 13),

-- Clubs (card_id 40-52)
(40, 'clubs', 1), 
(41, 'clubs', 2), 
(42, 'clubs', 3), 
(43, 'clubs', 4),
(44, 'clubs', 5), 
(45, 'clubs', 6), 
(46, 'clubs', 7), 
(47, 'clubs', 8),
(48, 'clubs', 9), 
(49, 'clubs', 10), 
(50, 'clubs', 11), 
(51, 'clubs', 12), 
(52, 'clubs', 13);
