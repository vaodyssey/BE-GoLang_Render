CREATE TABLE `users` (
    `id` varchar(36) PRIMARY KEY,
    `username` varchar(255),
    `password` varchar(255),
    `email` varchar(360),
    `phone` varchar(10),
    `address` varchar(255)
);

CREATE TABLE `products` (
    `id` varchar(36) PRIMARY KEY,
    `name` varchar(255) NOT NULL ,
    `image` varchar(255) NOT NULL ,
    `description` varchar(255),
    `price` float NOT NULL ,
    `created_at` timestamp
);

CREATE TABLE `order_details` (
    `order_id` varchar(36),
    `product_id` varchar(36),
    `quantity` integer
);

CREATE TABLE `orders` (
    `id` varchar(36) PRIMARY KEY,
    `amount` float,
    `status` int,
    `user_id` varchar(36),
    `created_at` timestamp
);