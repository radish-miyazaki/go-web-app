CREATE TABLE `user`
(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの拡張子',
    `name` VARCHAR(20) NOT NULL COMMENT 'ユーザー名',
    `password` VARCHAR(80) NOT NULL COMMENT 'パスワードハッシュ',
    `role` VARCHAR(80) NOT NULL COMMENT 'ロール',
    `created_at` DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
    `updated_at` DATETIME(6) NOT NULL COMMENT 'レコード更新日時',
    PRIMARY KEY(`id`),
    UNIQUE KEY `uix_name` (`name`) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

CREATE TABLE `task`
(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'タスクの拡張子',
    `title` VARCHAR(128) NOT NULL COMMENT 'タスクのタイトル',
    `status` VARCHAR(80) NOT NULL COMMENT 'タスクの状態',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'タスクを作成したユーザの識別子',
    `created_at` DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
    `updated_at` DATETIME(6) NOT NULL COMMENT 'レコード更新日時',
    PRIMARY KEY(`id`),
    CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
        ON DELETE RESTRICT ON UPDATE RESTRICT
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='タスク';
