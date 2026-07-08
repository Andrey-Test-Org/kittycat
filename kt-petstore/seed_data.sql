USE DATABASE TEST;
USE SCHEMA PETSTORE;

-- 1. customers (15 rows)
INSERT INTO customers (id, full_name, email, phone, address, created_at, updated_at) VALUES
('cust001', 'Alice Johnson',      'alice.johnson@example.com',    '555-0001', '101 Felidae St, Meowville',     '2024-01-15 10:00:00', '2024-01-15 10:00:00'),
('cust002', 'Bob Smith',          'bob.smith@example.com',        '555-0002', '202 Purr Lane, Catherton',      '2024-01-20 11:30:00', '2024-01-20 11:30:00'),
('cust003', 'Carol Williams',     'carol.w@example.com',          '555-0003', '303 Whisker Ave, Kittyburg',    '2024-02-01 09:15:00', '2024-02-01 09:15:00'),
('cust004', 'David Brown',        'david.brown@example.com',      '555-0004', '404 Tabby Rd, Pawsport',         '2024-02-05 14:20:00', '2024-02-05 14:20:00'),
('cust005', 'Emma Davis',         'emma.davis@example.com',       '555-0005', '505 Calico Ct, Feline Falls',   '2024-02-10 08:45:00', '2024-02-10 08:45:00'),
('cust006', 'Frank Miller',       'frank.miller@example.com',     '555-0006', '606 Scratch St, Clawton',       '2024-02-15 16:00:00', '2024-02-15 16:00:00'),
('cust007', 'Grace Wilson',       'grace.wilson@example.com',     '555-0007', '707 Meow Blvd, Catnap City',    '2024-03-01 12:10:00', '2024-03-01 12:10:00'),
('cust008', 'Henry Taylor',       'henry.taylor@example.com',     '555-0008', '808 Hiss Way, Purrtown',        '2024-03-05 10:30:00', '2024-03-05 10:30:00'),
('cust009', 'Iris Anderson',      'iris.anderson@example.com',    '555-0009', '909 Kitten Dr, Moggiewood',     '2024-03-10 15:45:00', '2024-03-10 15:45:00'),
('cust010', 'Jack Thomas',        'jack.thomas@example.com',      '555-0010', '110 Pawsome Pl, Catville',      '2024-03-15 11:00:00', '2024-03-15 11:00:00'),
('cust011', 'Karen Jackson',      'karen.j@example.com',          '555-0011', '121 Litter Ln, Tailsworth',     '2024-03-20 09:30:00', '2024-03-20 09:30:00'),
('cust012', 'Leo Martinez',       'leo.martinez@example.com',     '555-0012', '132 Furry Rd, Catsby Cove',     '2024-04-01 13:15:00', '2024-04-01 13:15:00'),
('cust013', 'Mia Robinson',       'mia.robinson@example.com',     '555-0013', '143 Purrfect St, Meowmeadow',   '2024-04-05 10:00:00', '2024-04-05 10:00:00'),
('cust014', 'Nathan Clark',       'nathan.clark@example.com',     '555-0014', '154 Clawsome Ave, Kittyfield',  '2024-04-10 14:30:00', '2024-04-10 14:30:00'),
('cust015', 'Olivia Lewis',       'olivia.lewis@example.com',     '555-0015', '165 Catitude Ln, Feline Park',  '2024-04-15 08:00:00', '2024-04-15 08:00:00');

-- 2. breeds (12 rows)
INSERT INTO breeds (id, name, description, origin_country, avg_lifespan_years, created_at, updated_at) VALUES
('brd001', 'Maine Coon',     'Large, gentle, fluffy, dog-like personality',        'USA',      13, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd002', 'Siamese',        'Vocal, sleek, social, striking blue eyes',           'Thailand', 15, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd003', 'Persian',        'Long-haired, calm, flat-faced, dignified',           'Iran',     14, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd004', 'Bengal',         'Wild-looking, energetic, spotted or marbled coat',  'USA',      12, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd005', 'British Shorthair','Stocky, round face, calm and independent',        'UK',       14, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd006', 'Sphynx',         'Hairless, extroverted, warm to touch, wrinkly',     'Canada',   13, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd007', 'Ragdoll',        'Large, floppy when held, blue-eyed, affectionate',  'USA',      15, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd008', 'Scottish Fold',  'Folded ears, round face, sweet-tempered',           'UK',       13, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd009', 'Abyssinian',     'Lean, active, ticked coat, curious and playful',    'Ethiopia', 13, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd010', 'Norwegian Forest','Large, thick double coat, loves climbing, hardy',  'Norway',   14, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd011', 'Russian Blue',   'Silvery-blue coat, green eyes, quiet and loyal',    'Russia',   15, '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('brd012', 'Domestic Shorthair','Mixed breed, varied coat, healthy and sturdy',   'Mixed',    15, '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 3. shelters (10 rows)
INSERT INTO shelters (id, name, address, phone, email, created_at, updated_at) VALUES
('shl001', 'Whiskers Rescue',         '123 Purr Lane, Catherton',      '555-1001', 'adopt@whiskers.org',         '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl002', 'Happy Tails Shelter',     '456 Wagging Way, Meowville',    '555-1002', 'info@happytails.org',        '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl003', 'Purrfect Home Rescue',    '789 Kitten Ave, Kittyburg',     '555-1003', 'contact@purrfecthome.org',   '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl004', 'Feline Friends Sanctuary','321 Tabby St, Pawsport',        '555-1004', 'hello@felinefriends.org',    '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl005', 'Cat Care Center',         '654 Calico Ct, Feline Falls',   '555-1005', 'admin@catcare.org',          '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl006', 'Second Chance Cats',      '987 Scratch Rd, Clawton',       '555-1006', 'adopt@secondchancecats.org', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl007', 'Meow Manor',              '147 Meow Blvd, Catnap City',    '555-1007', 'info@meowmanor.org',         '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl008', 'Safe Paws Rescue',        '258 Hiss Way, Purrtown',        '555-1008', 'help@safepaws.org',          '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl009', 'Kitty Haven',             '369 Kitten Dr, Moggiewood',     '555-1009', 'adopt@kittyhaven.org',       '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
('shl010', 'Paws and Whiskers',       '741 Pawsome Pl, Catville',      '555-1010', 'contact@pawswhiskers.org',   '2024-01-01 00:00:00', '2024-01-01 00:00:00');

-- 4. cats (20 rows)
INSERT INTO cats (id, name, breed_id, shelter_id, birth_date, price_cents, currency, status, description, created_at, updated_at) VALUES
('cat001', 'Luna',          'brd001', 'shl001', '2023-05-14', 25000,  'USD', 'ADOPTED',   'Friendly, loves cuddles and laser pointers',             '2024-01-15 10:00:00', '2024-03-20 10:00:00'),
('cat002', 'Simba',         'brd002', 'shl001', '2022-08-01', 30000,  'USD', 'AVAILABLE', 'Very talkative, loves warmth and sunbeams',              '2024-01-15 10:05:00', '2024-01-15 10:05:00'),
('cat003', 'Luna Bell',     'brd001', 'shl001', '2024-01-20', 20000,  'USD', 'ADOPTED',   'Playful kitten with a big personality',                   '2024-01-15 10:10:00', '2024-03-20 10:05:00'),
('cat004', 'Mochi',         'brd003', 'shl002', '2023-03-10', 35000,  'USD', 'AVAILABLE', 'Calm and regal, loves being brushed',                    '2024-01-16 09:00:00', '2024-01-16 09:00:00'),
('cat005', 'Nala',          'brd004', 'shl002', '2023-07-22', 40000,  'USD', 'AVAILABLE', 'Energetic, loves climbing and playing fetch',            '2024-01-16 09:15:00', '2024-01-16 09:15:00'),
('cat006', 'Oliver',        'brd005', 'shl003', '2022-11-05', 28000,  'USD', 'RESERVED',  'Stocky and dignified, prefers quiet evenings',           '2024-01-17 14:00:00', '2024-02-10 14:00:00'),
('cat007', 'Cleo',          'brd006', 'shl003', '2023-09-15', 32000,  'USD', 'AVAILABLE', 'Hairless but warm, very social and cuddly',              '2024-01-17 14:30:00', '2024-01-17 14:30:00'),
('cat008', 'Milo',          'brd007', 'shl004', '2023-01-30', 38000,  'USD', 'AVAILABLE', 'Goes limp when picked up, super affectionate',           '2024-01-18 11:00:00', '2024-01-18 11:00:00'),
('cat009', 'Zoe',           'brd008', 'shl004', '2023-06-12', 26000,  'USD', 'AVAILABLE', 'Adorable folded ears, loves chin scratches',             '2024-01-18 11:20:00', '2024-01-18 11:20:00'),
('cat010', 'Leo',           'brd009', 'shl005', '2023-04-18', 27000,  'USD', 'ADOPTED',   'Active and curious, always exploring',                   '2024-01-19 10:00:00', '2024-04-01 10:00:00'),
('cat011', 'Willow',        'brd010', 'shl005', '2022-12-03', 33000,  'USD', 'AVAILABLE', 'Big and fluffy, loves climbing cat trees',               '2024-01-19 10:30:00', '2024-01-19 10:30:00'),
('cat012', 'Shadow',        'brd011', 'shl006', '2023-02-08', 36000,  'USD', 'AVAILABLE', 'Quiet and loyal, stunning silver-blue coat',             '2024-01-20 09:00:00', '2024-01-20 09:00:00'),
('cat013', 'Pumpkin',       'brd012', 'shl006', '2023-10-31', 15000,  'USD', 'AVAILABLE', 'Orange tabby, loves food and naps in equal measure',     '2024-01-20 09:30:00', '2024-01-20 09:30:00'),
('cat014', 'Pearl',         'brd003', 'shl007', '2023-08-25', 37000,  'USD', 'RESERVED',  'White Persian, very calm, loves gentle petting',         '2024-01-21 15:00:00', '2024-03-15 15:00:00'),
('cat015', 'Max',           'brd004', 'shl007', '2023-03-17', 42000,  'USD', 'AVAILABLE', 'Bold and athletic, striking marble pattern',             '2024-01-21 15:30:00', '2024-01-21 15:30:00'),
('cat016', 'Bella',         'brd007', 'shl008', '2023-11-11', 39000,  'USD', 'AVAILABLE', 'Gentle giant, follows humans room to room',              '2024-01-22 08:00:00', '2024-01-22 08:00:00'),
('cat017', 'Ginger',        'brd012', 'shl008', '2022-06-20', 12000,  'USD', 'AVAILABLE', 'Tortoiseshell with a sassy personality',                 '2024-01-22 08:30:00', '2024-01-22 08:30:00'),
('cat018', 'Tiger',         'brd009', 'shl009', '2023-05-05', 29000,  'USD', 'ADOPTED',   'Sleek and athletic, loves interactive toys',             '2024-01-23 13:00:00', '2024-04-10 13:00:00'),
('cat019', 'Misty',         'brd011', 'shl009', '2023-07-07', 34000,  'USD', 'AVAILABLE', 'Shy at first but very loyal once comfortable',           '2024-01-23 13:30:00', '2024-01-23 13:30:00'),
('cat020', 'Coco',          'brd005', 'shl010', '2023-12-01', 31000,  'USD', 'RETIRED',   'Round and plush, now living the good life in shelter',   '2024-01-24 16:00:00', '2024-05-01 16:00:00');

-- 5. carts (15 rows)
INSERT INTO carts (id, customer_id, created_at, updated_at) VALUES
('cart001', 'cust001', '2024-02-01 10:00:00', '2024-02-01 10:30:00'),
('cart002', 'cust002', '2024-02-05 11:00:00', '2024-02-05 11:15:00'),
('cart003', 'cust003', '2024-02-10 09:30:00', '2024-02-10 10:00:00'),
('cart004', 'cust004', '2024-02-15 14:00:00', '2024-02-15 14:20:00'),
('cart005', 'cust005', '2024-02-20 08:45:00', '2024-02-20 09:00:00'),
('cart006', 'cust006', '2024-03-01 16:10:00', '2024-03-01 16:30:00'),
('cart007', 'cust007', '2024-03-05 12:00:00', '2024-03-05 12:10:00'),
('cart008', 'cust008', '2024-03-10 10:45:00', '2024-03-10 11:00:00'),
('cart009', 'cust009', '2024-03-15 15:30:00', '2024-03-15 15:45:00'),
('cart010', 'cust010', '2024-03-20 11:00:00', '2024-03-20 11:20:00'),
('cart011', 'cust011', '2024-03-25 09:15:00', '2024-03-25 09:30:00'),
('cart012', 'cust012', '2024-04-01 13:30:00', '2024-04-01 14:00:00'),
('cart013', 'cust013', '2024-04-05 10:15:00', '2024-04-05 10:30:00'),
('cart014', 'cust014', '2024-04-10 14:45:00', '2024-04-10 15:00:00'),
('cart015', 'cust015', '2024-04-15 08:30:00', '2024-04-15 08:45:00');

-- 6. cart_items (20 rows)
INSERT INTO cart_items (id, cart_id, cat_id, quantity, created_at) VALUES
('citem001', 'cart001', 'cat002', 1, '2024-02-01 10:05:00'),
('citem002', 'cart001', 'cat003', 1, '2024-02-01 10:10:00'),
('citem003', 'cart002', 'cat005', 1, '2024-02-05 11:05:00'),
('citem004', 'cart003', 'cat004', 1, '2024-02-10 09:35:00'),
('citem005', 'cart003', 'cat009', 1, '2024-02-10 09:40:00'),
('citem006', 'cart004', 'cat006', 1, '2024-02-15 14:05:00'),
('citem007', 'cart004', 'cat007', 1, '2024-02-15 14:10:00'),
('citem008', 'cart005', 'cat008', 1, '2024-02-20 08:50:00'),
('citem009', 'cart005', 'cat011', 1, '2024-02-20 08:55:00'),
('citem010', 'cart006', 'cat012', 1, '2024-03-01 16:15:00'),
('citem011', 'cart007', 'cat013', 1, '2024-03-05 12:05:00'),
('citem012', 'cart007', 'cat015', 1, '2024-03-05 12:08:00'),
('citem013', 'cart008', 'cat014', 1, '2024-03-10 10:50:00'),
('citem014', 'cart008', 'cat016', 1, '2024-03-10 10:55:00'),
('citem015', 'cart009', 'cat017', 1, '2024-03-15 15:35:00'),
('citem016', 'cart010', 'cat019', 1, '2024-03-20 11:05:00'),
('citem017', 'cart011', 'cat005', 1, '2024-03-25 09:20:00'),
('citem018', 'cart012', 'cat002', 1, '2024-04-01 13:35:00'),
('citem019', 'cart013', 'cat011', 1, '2024-04-05 10:20:00'),
('citem020', 'cart014', 'cat015', 1, '2024-04-10 14:50:00');

-- 7. orders (15 rows)
INSERT INTO orders (id, customer_id, status, total_cents, currency, ship_address, bill_address, notes, created_at, updated_at, placed_at, shipped_at, cancelled_at) VALUES
('ord001', 'cust001', 'SHIPPED',   45000, 'USD', '101 Felidae St, Meowville',   '101 Felidae St, Meowville', 'Please deliver in the morning',       '2024-02-01 10:30:00', '2024-02-05 12:00:00', '2024-02-01 10:30:00', '2024-02-05 12:00:00', NULL),
('ord002', 'cust002', 'CANCELLED', 40000, 'USD', '202 Purr Lane, Catherton',    '202 Purr Lane, Catherton',  NULL,                                  '2024-02-05 11:20:00', '2024-02-06 09:00:00', '2024-02-05 11:20:00', NULL,                 '2024-02-06 09:00:00'),
('ord003', 'cust003', 'SHIPPED',   61000, 'USD', '303 Whisker Ave, Kittyburg',  '303 Whisker Ave, Kittyburg','Leave at front door',                 '2024-02-10 10:05:00', '2024-02-15 14:00:00', '2024-02-10 10:05:00', '2024-02-15 14:00:00', NULL),
('ord004', 'cust004', 'PAID',      60000, 'USD', '404 Tabby Rd, Pawsport',       '404 Tabby Rd, Pawsport',    'Gift wrapping please',                '2024-02-15 14:25:00', '2024-02-16 10:00:00', '2024-02-15 14:25:00', NULL,                 NULL),
('ord005', 'cust005', 'PENDING',   71000, 'USD', '505 Calico Ct, Feline Falls',  '505 Calico Ct, Feline Falls', NULL,                                '2024-02-20 09:05:00', '2024-02-20 09:05:00', '2024-02-20 09:05:00', NULL,                 NULL),
('ord006', 'cust006', 'SHIPPED',   36000, 'USD', '606 Scratch St, Clawton',      '606 Scratch St, Clawton',   NULL,                                  '2024-03-01 16:35:00', '2024-03-05 11:00:00', '2024-03-01 16:35:00', '2024-03-05 11:00:00', NULL),
('ord007', 'cust007', 'PAID',      57000, 'USD', '707 Meow Blvd, Catnap City',   '707 Meow Blvd, Catnap City','Call before delivery',                '2024-03-05 12:15:00', '2024-03-06 09:00:00', '2024-03-05 12:15:00', NULL,                 NULL),
('ord008', 'cust008', 'CANCELLED', 70000, 'USD', '808 Hiss Way, Purrtown',       '808 Hiss Way, Purrtown',    'Changed my mind',                     '2024-03-10 11:05:00', '2024-03-11 10:00:00', '2024-03-10 11:05:00', NULL,                 '2024-03-11 10:00:00'),
('ord009', 'cust009', 'SHIPPED',   12000, 'USD', '909 Kitten Dr, Moggiewood',    '909 Kitten Dr, Moggiewood', NULL,                                  '2024-03-15 15:50:00', '2024-03-20 10:00:00', '2024-03-15 15:50:00', '2024-03-20 10:00:00', NULL),
('ord010', 'cust010', 'PAID',      34000, 'USD', '110 Pawsome Pl, Catville',     '110 Pawsome Pl, Catville',  NULL,                                  '2024-03-20 11:25:00', '2024-03-21 09:00:00', '2024-03-20 11:25:00', NULL,                 NULL),
('ord011', 'cust011', 'PENDING',   40000, 'USD', '121 Litter Ln, Tailsworth',    '121 Litter Ln, Tailsworth', NULL,                                  '2024-03-25 09:35:00', '2024-03-25 09:35:00', '2024-03-25 09:35:00', NULL,                 NULL),
('ord012', 'cust012', 'CANCELLED', 30000, 'USD', '132 Furry Rd, Catsby Cove',    '132 Furry Rd, Catsby Cove', 'Found a cat elsewhere',               '2024-04-01 14:05:00', '2024-04-02 10:00:00', '2024-04-01 14:05:00', NULL,                 '2024-04-02 10:00:00'),
('ord013', 'cust013', 'SHIPPED',   33000, 'USD', '143 Purrfect St, Meowmeadow',  '143 Purrfect St, Meowmeadow', NULL,                                '2024-04-05 10:35:00', '2024-04-10 12:00:00', '2024-04-05 10:35:00', '2024-04-10 12:00:00', NULL),
('ord014', 'cust014', 'PAID',      42000, 'USD', '154 Clawsome Ave, Kittyfield', '154 Clawsome Ave, Kittyfield', NULL,                               '2024-04-10 15:05:00', '2024-04-11 09:00:00', '2024-04-10 15:05:00', NULL,                 NULL),
('ord015', 'cust015', 'PENDING',   31000, 'USD', '165 Catitude Ln, Feline Park', '165 Catitude Ln, Feline Park', 'Excited to adopt!',                  '2024-04-15 08:50:00', '2024-04-15 08:50:00', '2024-04-15 08:50:00', NULL,                 NULL);

-- 8. order_items (20 rows)
INSERT INTO order_items (id, order_id, cat_id, quantity, price_cents, currency, created_at) VALUES
('oitem001', 'ord001', 'cat001', 1, 25000, 'USD', '2024-02-01 10:30:00'),
('oitem002', 'ord001', 'cat003', 1, 20000, 'USD', '2024-02-01 10:30:00'),
('oitem003', 'ord002', 'cat005', 1, 40000, 'USD', '2024-02-05 11:20:00'),
('oitem004', 'ord003', 'cat004', 1, 35000, 'USD', '2024-02-10 10:05:00'),
('oitem005', 'ord003', 'cat009', 1, 26000, 'USD', '2024-02-10 10:05:00'),
('oitem006', 'ord004', 'cat006', 1, 28000, 'USD', '2024-02-15 14:25:00'),
('oitem007', 'ord004', 'cat007', 1, 32000, 'USD', '2024-02-15 14:25:00'),
('oitem008', 'ord005', 'cat008', 1, 38000, 'USD', '2024-02-20 09:05:00'),
('oitem009', 'ord005', 'cat011', 1, 33000, 'USD', '2024-02-20 09:05:00'),
('oitem010', 'ord006', 'cat012', 1, 36000, 'USD', '2024-03-01 16:35:00'),
('oitem011', 'ord007', 'cat013', 1, 15000, 'USD', '2024-03-05 12:15:00'),
('oitem012', 'ord007', 'cat015', 1, 42000, 'USD', '2024-03-05 12:15:00'),
('oitem013', 'ord008', 'cat014', 1, 37000, 'USD', '2024-03-10 11:05:00'),
('oitem014', 'ord008', 'cat016', 1, 39000, 'USD', '2024-03-10 11:05:00'),
('oitem015', 'ord009', 'cat017', 1, 12000, 'USD', '2024-03-15 15:50:00'),
('oitem016', 'ord010', 'cat019', 1, 34000, 'USD', '2024-03-20 11:25:00'),
('oitem017', 'ord011', 'cat005', 1, 40000, 'USD', '2024-03-25 09:35:00'),
('oitem018', 'ord012', 'cat002', 1, 30000, 'USD', '2024-04-01 14:05:00'),
('oitem019', 'ord013', 'cat011', 1, 33000, 'USD', '2024-04-05 10:35:00'),
('oitem020', 'ord014', 'cat015', 1, 42000, 'USD', '2024-04-10 15:05:00');

-- 9. inventory (20 rows)
INSERT INTO inventory (cat_id, available, reserved, updated_at) VALUES
('cat001', 0, 0, '2024-03-20 10:00:00'),
('cat002', 1, 0, '2024-01-15 10:05:00'),
('cat003', 0, 0, '2024-03-20 10:05:00'),
('cat004', 1, 0, '2024-01-16 09:00:00'),
('cat005', 1, 0, '2024-01-16 09:15:00'),
('cat006', 0, 1, '2024-02-10 14:00:00'),
('cat007', 1, 0, '2024-01-17 14:30:00'),
('cat008', 1, 0, '2024-01-18 11:00:00'),
('cat009', 1, 0, '2024-01-18 11:20:00'),
('cat010', 0, 0, '2024-04-01 10:00:00'),
('cat011', 1, 0, '2024-01-19 10:30:00'),
('cat012', 1, 0, '2024-01-20 09:00:00'),
('cat013', 1, 0, '2024-01-20 09:30:00'),
('cat014', 0, 1, '2024-03-15 15:00:00'),
('cat015', 1, 0, '2024-01-21 15:30:00'),
('cat016', 1, 0, '2024-01-22 08:00:00'),
('cat017', 1, 0, '2024-01-22 08:30:00'),
('cat018', 0, 0, '2024-04-10 13:00:00'),
('cat019', 1, 0, '2024-01-23 13:30:00'),
('cat020', 0, 0, '2024-05-01 16:00:00');

-- 10. audit_log (20 rows)
INSERT INTO audit_log (id, actor, action, target, created_at) VALUES
('aud001', 'cust001', 'order.place',    'ord001', '2024-02-01 10:30:00'),
('aud002', 'system',  'order.pay',      'ord001', '2024-02-02 09:00:00'),
('aud003', 'system',  'order.ship',     'ord001', '2024-02-05 12:00:00'),
('aud004', 'cust002', 'order.place',    'ord002', '2024-02-05 11:20:00'),
('aud005', 'cust002', 'order.cancel',   'ord002', '2024-02-06 09:00:00'),
('aud006', 'cust003', 'order.place',    'ord003', '2024-02-10 10:05:00'),
('aud007', 'system',  'order.pay',      'ord003', '2024-02-11 09:00:00'),
('aud008', 'system',  'order.ship',     'ord003', '2024-02-15 14:00:00'),
('aud009', 'cust004', 'order.place',    'ord004', '2024-02-15 14:25:00'),
('aud010', 'system',  'order.pay',      'ord004', '2024-02-16 10:00:00'),
('aud011', 'cust005', 'order.place',    'ord005', '2024-02-20 09:05:00'),
('aud012', 'cust006', 'order.place',    'ord006', '2024-03-01 16:35:00'),
('aud013', 'system',  'order.ship',     'ord006', '2024-03-05 11:00:00'),
('aud014', 'cust007', 'order.place',    'ord007', '2024-03-05 12:15:00'),
('aud015', 'system',  'order.pay',      'ord007', '2024-03-06 09:00:00'),
('aud016', 'cust008', 'order.place',    'ord008', '2024-03-10 11:05:00'),
('aud017', 'cust008', 'order.cancel',   'ord008', '2024-03-11 10:00:00'),
('aud018', 'cust009', 'order.place',    'ord009', '2024-03-15 15:50:00'),
('aud019', 'system',  'order.ship',     'ord009', '2024-03-20 10:00:00'),
('aud020', 'cust010', 'order.place',    'ord010', '2024-03-20 11:25:00');
