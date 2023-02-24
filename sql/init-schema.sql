CREATE TABLE articles
(
    id          INT          NOT NULL AUTO_INCREMENT,
    uid         VARCHAR(100) NOT NULL UNIQUE,
    title       VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL,
    content     TEXT         NOT NULL,
    link        VARCHAR(255) NOT NULL,
    country     VARCHAR(10)  NOT NULL,
    location    VARCHAR(100) NOT NULL,
    lang        VARCHAR(10)  NOT NULL,
    pub_date    BIGINT       NOT NULL,
    saved_at    BIGINT       NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE feed_lock
(
    id        INT     NOT NULL AUTO_INCREMENT,
    is_locked TINYINT NOT NULL,
    timestamp BIGINT  NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE article_categories
(
    article_id  INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (article_id, category_id),
    FOREIGN KEY (article_id) REFERENCES articles (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE categories
(
    id   INT          NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);