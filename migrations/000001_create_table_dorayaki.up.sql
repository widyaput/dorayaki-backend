CREATE TABLE IF NOT EXISTS `dorayaki` (
  `rasa` varchar(191) NOT NULL,
  `toko_id` bigint(20) NOT NULL,
  `deskripsi` longtext NOT NULL,
  `image_url` longtext NOT NULL,
  `stok` bigint(20) DEFAULT 0,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`rasa`),
  CONSTRAINT `chk_dorayaki_stok` CHECK (`stok` >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
