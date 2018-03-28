CREATE TABLE `bookmarks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `url` varchar(255) NOT NULL,
  `title` varchar(100) NOT NULL,
  `author_name` varchar(100) NOT NULL,
  `added_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `width` int(11) NOT NULL DEFAULT 0,
  `height` int(11) NOT NULL DEFAULT 0,
  `duration` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `keywords` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `bookmark_keywords` (
    `bookmark_id` int(10) unsigned NOT NULL,
    `keyword_id` int(10) unsigned NOT NULL,
    KEY `bookmark_keywords_bookmark_id` (`bookmark_id`),
    KEY `bookmark_keywords_keyword_id` (`keyword_id`),
    CONSTRAINT `fk_bookmark_keywords_bookmark_id` FOREIGN KEY (`bookmark_id`) REFERENCES `bookmarks` (`id`),
    CONSTRAINT `fk_bookmark_keywords_keywords_id` FOREIGN KEY (`keyword_id`) REFERENCES `keywords` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
