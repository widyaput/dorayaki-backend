CREATE TABLE IF NOT EXISTS `toko` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `nama` longtext NOT NULL,
  `jalan` longtext NOT NULL,
  `kecamatan` longtext NOT NULL,
  `provinsi` longtext NOT NULL,
  `created_at` bigint(20) DEFAULT NULL,
  `updated_at` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
