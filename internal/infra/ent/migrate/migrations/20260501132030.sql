-- Create "todos" table
CREATE TABLE `todos` (
  `id` char(36) NOT NULL,
  `title` varchar(255) NOT NULL,
  `status` enum('pending','in_progress','done') NOT NULL DEFAULT "pending",
  `created_at` timestamp NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
