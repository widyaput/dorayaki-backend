CREATE TABLE IF NOT EXISTS `dorayaki` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rasa` longtext NOT NULL,
  `deskripsi` longtext NOT NULL,
  `base64` longtext NOT NULL,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
