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

insert into sites (url, category, has_content, country, location)
values ('https://www.rosario3.com/rss.html', 'actualidad', false, 'Argentina', 'Rosario'),
        ('https://www.pagina12.com.ar/rss/secciones/ciencia/notas', 'ciencia', true, 'Argentina', 'Buenos Aires'),
        ('https://www.eltribuno.com/salta/rss-new/portada.rss', 'actualidad', true, 'Argentina', 'Salta'),
        ('https://www.losandes.com.ar/arcio/rss/', 'actualidad', true, 'Argentina', 'Mendoza'),
        ('https://www.pagina12.com.ar/rss/secciones/deportes/notas', 'deportes', true, 'Argentina', 'Argentina'),
        ('https://www.pagina12.com.ar/rss/secciones/el-mundo/notas,mundo', 'mundo', true, 'Argentina', 'Argentina'),
        ('https://www.pagina12.com.ar/rss/suplementos/rosario12/notas', 'actualidad', true, 'Argentina', 'Rosario'),
        ('https://www.pagina12.com.ar/rss/secciones/economia/notas', 'economia', true, 'Argentina', 'Argentina'),
        ('https://www.paparazzi.com.ar/feed/', 'espectaculos', true, 'Argentina', 'Argentina'),
        ('https://www.eldiarioar.com/rss', 'actualidad', false, 'Argentina', 'Argentina'),
        ('https://www.infobae.com/argentina-rss.xml', 'actualidad', true, 'Argentina', 'Argentina'),
        ('https://www.lanacion.com.ar/arc/outboundfeeds/rss/?outputType=xml', 'actualidad', true, 'Argentina',
         'Argentina'),
        ('https://www.impulsonegocios.com/feed/', 'actualidad', true, 'Argentina', 'Rosario'),
        ('https://www.eldia.com/.rss', 'actualidad', true, 'Argentina', 'La Plata'),
        ('https://elsolnoticias.com.ar/feed/', 'actualidad', true, 'Argentina', 'Quilmes'),
        ('https://www.lacapitalmdp.com/feed/', 'actualidad', true, 'Argentina', 'MDQ'),
        ('https://primerobahia.com.ar/feed/', 'actualidad', true, 'Argentina', 'Bahia Blanca'),
        ('https://www.actualidaddemercedes.com/feed/', 'actualidad', true, 'Argentina', 'Mercedes'),
        ('https://www.bigbangnews.com/rss/actualidad.xml', 'actualidad', true, 'Argentina', 'CABA'),
        ('https://inforama.com.ar/feed/', 'actualidad', true, 'Argentina', 'Catamarca'),
        ('https://www.diarioepoca.com/rss', 'actualidad', true, 'Argentina', 'Corrientes'),
        ('http://vivocomodoro.com.ar/feed/', 'actualidad', true, 'Argentina', 'Comodoro'),
        ('https://www.eldiario.com.ar/feed/', 'actualidad', true, 'Argentina,''Entre' 'Rios'),
        ('https://agenciasanluis.com/feed/', 'actualidad', true, 'Argentina', 'San Luis'),
        ('https://www.diariodecuyo.com.ar/rss/rss.xml', 'actualidad', true, 'Argentina', 'San Juan'),
        ('https://www.lavoz.com.ar/arc/outboundfeeds/feeds/rss/?outputType=xml,actualidad,true', 'Argentina',
         'Cordoba'),
        ('https://www.rionegro.com.ar/feed/,actualidad', true, 'Argentina', 'Rio Negro'),
        ('https://www.eltribuno.com/salta/rss-new/portada.rss', 'actualidad', true, 'Argentina', 'Salta'),
        ('https://www.eltribuno.com/jujuy/rss-new/portada.rss', 'actualidad', true, 'Argentina', 'Jujuy');
