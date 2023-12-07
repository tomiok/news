CREATE TABLE articles
(
    id          SERIAL PRIMARY KEY,
    uid         VARCHAR(100) NOT NULL UNIQUE,
    title       VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL,
    content     TEXT         NOT NULL,
    raw_content TEXT         NOT NULL,
    link        VARCHAR(255) NOT NULL,
    country     VARCHAR(10)  NOT NULL,
    location    VARCHAR(100) NOT NULL,
    lang        VARCHAR(10)  NOT NULL,
    pub_date    BIGINT       NOT NULL,
    source      varchar(255) NOT NULL,
    saved_at    BIGINT       NOT NULL
);

CREATE TABLE feed_lock
(
    id        SERIAL PRIMARY KEY,
    is_locked BOOL   NOT NULL,
    timestamp BIGINT NOT NULL
);

CREATE TABLE categories
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE article_categories
(
    article_id  INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (article_id, category_id),
    FOREIGN KEY (article_id) REFERENCES articles (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE sites
(
    id          SERIAL PRIMARY KEY,
    url         VARCHAR(250) NOT NULL,
    category    VARCHAR(150),
    has_content BOOL,
    country_id  INT NOT NULL,
    location_id INT NOT NULL
);

-- https://feeds.elpais.com/mrss-s/pages/ep/site/elpais.com/section/america/portada
-- https://www.terra.com/rss

insert into sites (url, category, has_content, country_id, location_id)
values ('https://www.rosario3.com/rss.html', 'actualidad', false, 'argentina', 'rosario'),
       ('https://www.pagina12.com.ar/rss/secciones/ciencia/notas', 'ciencia', true, 'argentina', 'buenosaires'),
       ('https://www.eltribuno.com/salta/rss-new/portada.rss', 'actualidad', true, 'argentina', 'Salta'),
       ('https://www.losandes.com.ar/arcio/rss/', 'actualidad', true, 'argentina', 'Mendoza'),
       ('https://www.pagina12.com.ar/rss/secciones/deportes/notas', 'deportes', true, 'argentina', 'argentina'),
       ('https://www.pagina12.com.ar/rss/secciones/el-mundo/notas', 'mundo', true, 'argentina', 'argentina'),
       ('https://www.pagina12.com.ar/rss/suplementos/rosario12/notas', 'actualidad', true, 'argentina', 'rosario'),
       ('https://www.pagina12.com.ar/rss/secciones/economia/notas', 'economia', true, 'argentina', 'argentina'),
       ('https://www.paparazzi.com.ar/feed/', 'espectaculos', true, 'argentina', 'argentina'),
       ('https://www.eldiarioar.com/rss', 'actualidad', false, 'argentina', 'argentina'),
       ('https://www.infobae.com/argentina-rss.xml', 'actualidad', true, 'argentina', 'argentina'),
       ('https://www.lanacion.com.ar/arc/outboundfeeds/rss/?outputType=xml', 'actualidad', true, 'argentina',
        'argentina'),
       ('https://www.impulsonegocios.com/feed/', 'actualidad', true, 'argentina', 'rosario'),
       ('https://www.eldia.com/.rss', 'actualidad', true, 'argentina', 'La Plata'),
       ('https://elsolnoticias.com.ar/feed/', 'actualidad', true, 'argentina', 'Quilmes'),
       ('https://www.lacapitalmdp.com/feed/', 'actualidad', true, 'argentina', 'MDQ'),
       ('https://primerobahia.com.ar/feed/', 'actualidad', true, 'argentina', 'Bahia Blanca'),
       ('https://www.actualidaddemercedes.com/feed/', 'actualidad', true, 'argentina', 'Mercedes'),
       ('https://www.bigbangnews.com/rss/actualidad', 'actualidad', true, 'argentina', 'CABA'),
       ('https://inforama.com.ar/feed/', 'actualidad', true, 'argentina', 'Catamarca'),
       ('https://www.diarioepoca.com/rss', 'actualidad', true, 'argentina', 'Corrientes'),
       ('http://vivocomodoro.com.ar/feed/', 'actualidad', true, 'argentina', 'Comodoro'),
       ('https://www.eldiario.com.ar/feed/', 'actualidad', true, 'argentina', 'Entre Rios'),
       ('https://agenciasanluis.com/feed/', 'actualidad', true, 'argentina', 'San Luis'),
       ('https://www.diariodecuyo.com.ar/rss/rss.xml', 'actualidad', true, 'argentina', 'San Juan'),
       ('https://www.lavoz.com.ar/arc/outboundfeeds/feeds/rss/?outputType=xml', 'actualidad', true, 'argentina',
        'Cordoba'),
       ('https://www.rionegro.com.ar/feed/', 'actualidad', true, 'argentina', 'Rio Negro'),
       ('https://www.eltribuno.com/salta/rss-new/portada.rss', 'actualidad', true, 'argentina', 'Salta'),
       ('https://www.eltribuno.com/jujuy/rss-new/portada.rss', 'actualidad', true, 'argentina', 'Jujuy');
