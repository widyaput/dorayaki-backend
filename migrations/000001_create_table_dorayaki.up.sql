CREATE TABLE `dorayaki` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rasa` longtext NOT NULL,
  `deskripsi` longtext NOT NULL,
  `image_url` longtext NOT NULL,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
