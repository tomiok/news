CREATE TABLE articles
(
    id          INT          NOT NULL AUTO_INCREMENT,
    uid         VARCHAR(100) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    link        VARCHAR(255) NOT NULL,
    lang        VARCHAR(10)  NOT NULL,
    description TEXT         NOT NULL,
    pub_date    TIMESTAMP    NOT NULL,
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