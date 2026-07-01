BEGIN;

INSERT INTO restaurant (
    id,
    name,
    email,
    phone,
    street_address,
    opening_time,
    closing_time,
    delivery_fee,
    delivery_radius,
    is_open,
    version
)
OVERRIDING SYSTEM VALUE
VALUES
(
    1,
    'Restaurant',
    'restaurant@example.com',
    '9876543210',
    '221B Baker Street',
    '2000-01-01 14:30:00+05:30',
    '2000-01-02 03:30:00+05:30',
    50,
    5,
    TRUE,
    1
)
ON CONFLICT (id) DO NOTHING;


INSERT INTO users (
    id,
    name,
    email,
    password_hash,
    phone,
    role,
    activated,
    version
)
OVERRIDING SYSTEM VALUE
VALUES
(
    1,
    'John Doe',
    'customer@example.com',
    '\x24326124313224324a7567306f75476d384342556142334737743450656332424c63595a557359697a675559634a77486f7a5232454b497472776969',
    '9876543210',
    'customer',
    TRUE,
    1
),
(
    2,
    'Restaurant Owner',
    'owner@example.com',
    '\x243261243132244e494e444266373662684d6f685a566a304845556a655559677673725a7a4935424737347762547478516a4a37626e653552513447',
    '9123456789',
    'owner',
    TRUE,
    1
)
ON CONFLICT (id) DO NOTHING;


INSERT INTO categories (
    id,
    name,
    display_order,
    version
)
OVERRIDING SYSTEM VALUE
VALUES
(
    2,
    'Burgers',
    2,
    2
)
ON CONFLICT (id) DO NOTHING;


INSERT INTO dishes (
    id,
    category_id,
    name,
    description,
    price,
    image_url,
    is_available,
    is_vegetarian,
    is_featured,
    preparation_time,
    version
)
OVERRIDING SYSTEM VALUE
VALUES
(
    2,
    2,
    'Margherita Pizza',
    'Classic cheese pizza',
    299,
    'https://example.com/pizza.jpg',
    TRUE,
    TRUE,
    TRUE,
    20,
    1
)
ON CONFLICT (id) DO NOTHING;


SELECT setval('restaurant_id_seq', (SELECT MAX(id) FROM restaurant));
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));
SELECT setval('dishes_id_seq', (SELECT MAX(id) FROM dishes));

COMMIT;
