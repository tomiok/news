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
    saved_at    BIGINT       NOT NULL,
    categories  TEXT[] NOT NULL DEFAULT '{}'::text[]
);

CREATE TABLE sites
(
    id          SERIAL PRIMARY KEY,
    url         VARCHAR(250) NOT NULL,
    category    VARCHAR(150),
    has_content BOOL,
    country     VARCHAR(150),
    location    VARCHAR(150),
    enabled     BOOL DEFAULT true
);

-- https://feeds.elpais.com/mrss-s/pages/ep/site/elpais.com/section/america/portada
-- https://www.terra.com/rss

insert into sites (url, category, has_content, country, location)
values ('https://www.rosario3.com/rss.html', 'actualidad', false, 'argentina', 'rosario'),
       ('https://www.pagina12.com.ar/rss/secciones/ciencia/notas', 'ciencia', true, 'argentina', 'buenos aires'),
       ('https://www.eltribuno.com/salta/rss-new/portada.rss', 'actualidad', true, 'argentina', 'salta'),
       ('https://www.pagina12.com.ar/rss/secciones/deportes/notas', 'deportes', true, 'argentina', 'argentina'),
       ('https://www.pagina12.com.ar/rss/secciones/el-mundo/notas', 'mundo', true, 'argentina', 'argentina'),
       ('https://www.pagina12.com.ar/rss/suplementos/rosario12/notas', 'actualidad', true, 'argentina', 'rosario'),
       ('https://www.pagina12.com.ar/rss/secciones/economia/notas', 'economia', true, 'argentina', 'argentina'),
       ('https://www.paparazzi.com.ar/feed/', 'espectaculos', true, 'argentina', 'argentina'),
       ('https://www.eldiarioar.com/rss', 'actualidad', false, 'argentina', 'argentina'),
       ('https://www.infobae.com/arc/outboundfeeds/rss/', 'actualidad', true, 'argentina', 'argentina'),
       ('https://www.lanacion.com.ar/arc/outboundfeeds/rss/?outputType=xml', 'actualidad', true, 'argentina',
        'argentina'),
       ('https://www.impulsonegocios.com/feed/', 'actualidad', true, 'argentina', 'rosario'),
       ('https://www.eldia.com/.rss', 'actualidad', true, 'argentina', 'la plata'),
       ('https://elsolnoticias.com.ar/feed/', 'actualidad', true, 'argentina', 'quilmes'),
       ('https://www.lacapitalmdp.com/feed/', 'actualidad', true, 'argentina', 'mar del plata'),
       ('https://primerobahia.com.ar/feed/', 'actualidad', true, 'argentina', 'bahia blanca'),
       ('https://www.actualidaddemercedes.com/feed/', 'actualidad', true, 'argentina', 'mercedes'),
       ('https://www.bigbangnews.com/feed/', 'actualidad', true, 'argentina', 'caba'),
       ('https://inforama.com.ar/feed/', 'actualidad', true, 'argentina', 'catamarca'),
       ('https://www.diarioepoca.com/rss', 'actualidad', true, 'argentina', 'corrientes'),
       ('http://vivocomodoro.com.ar/feed/', 'actualidad', true, 'argentina', 'comodoro'),
       ('https://www.eldiario.com.ar/feed/', 'actualidad', true, 'argentina', 'entre rios'),
       ('https://agenciasanluis.com/feed/', 'actualidad', true, 'argentina', 'san luis'),
       ('https://www.diariodecuyo.com.ar/rss/rss.xml', 'actualidad', true, 'argentina', 'san juan'),
       ('https://www.lavoz.com.ar/arc/outboundfeeds/feeds/rss/?outputType=xml', 'actualidad', true, 'argentina',
        'cordoba'),
       ('https://www.rionegro.com.ar/feed/', 'actualidad', true, 'argentina', 'rio negro'),
       ('https://www.eltribuno.com/salta/rss-new/portada.rss', 'actualidad', true, 'argentina', 'salta'),
       ('https://www.eltribuno.com/jujuy/rss-new/portada.rss', 'actualidad', true, 'argentina', 'jujuy'),
       ('https://ojodeprensa.com.ar/feed/', 'actualidad', true, 'argentina', 'rosario'),
       ('https://diarioconurbano.com.ar/feed', 'actualidad', true, 'argentina', 'conurbano'),
       ('https://conurbanodiario.com.ar/?feed=rss2', 'actualidad', true, 'argentina', 'conurbano'),
       ('https://www.infoban.com.ar/feed/', 'actualidad', true, 'argentina', 'conurbano'),
       ('https://www.inforegion.com.ar/feed/', 'actualidad', false, 'argentina', 'conurbano');

