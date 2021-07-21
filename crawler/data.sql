CREATE TABLE `sp_boss_jobs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `job_name` varchar(32) NOT NULL DEFAULT '' COMMENT '工作名称',
  `salary` varchar(30) NOT NULL COMMENT '薪资',
  `job_type` varchar(4) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '类型',
  `city` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '城市',
  `href` varchar(255) NOT NULL COMMENT '详情连接',
  `company_name` varchar(32) DEFAULT '' COMMENT '公司名称',
  `company_address` varchar(64) NOT NULL DEFAULT '' COMMENT '公司地址',
  `work_years` varchar(16) DEFAULT '' COMMENT '工作年限',
  `education` varchar(16) DEFAULT '' COMMENT '学历要求',
  `company_label` varchar(16) DEFAULT '' COMMENT '公司所属行业',
  `financing_stage` varchar(16) DEFAULT '' COMMENT '融资阶段',
  `staff_number` varchar(16) DEFAULT '' COMMENT '公司规模-员工人数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='boss招聘信息表';