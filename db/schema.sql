CREATE TABLE `lgtm_images`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `filename`   varchar(255) NOT NULL COMMENT '画像ファイル名になっているuniqueな文字列が格納される .e.g. aaaaaaaa-aaaa-aaaa-aaaa-123456788qqq',
    `path`       varchar(255) NOT NULL COMMENT 'S3に保存されているパス .e.g. YYYY/MM/DD/HH',
    `created_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_lgtm_images_01` (`filename`),
    KEY          `idx_lgtm_images_01` (`path`)
);
