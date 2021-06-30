CREATE TABLE IF EXISTS `toko_dorayaki` (
  `toko_id` bigint(20) NOT NULL,
  `dorayaki_id` bigint(20) NOT NULL,
  `stok` bigint(20) DEFAULT 0,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`toko_id`,`dorayaki_id`),
  KEY `fk_toko_dorayaki_dorayaki` (`dorayaki_id`),
  CONSTRAINT `fk_toko_dorayaki_dorayaki` FOREIGN KEY (`dorayaki_id`) REFERENCES `dorayaki` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_toko_dorayaki_toko` FOREIGN KEY (`toko_id`) REFERENCES `toko` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_toko_stok` FOREIGN KEY (`toko_id`) REFERENCES `toko` (`id`) ON DELETE CASCADE,
  CONSTRAINT `chk_toko_dorayaki_stok` CHECK (`stok` >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
