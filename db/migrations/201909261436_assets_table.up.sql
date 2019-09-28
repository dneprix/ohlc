CREATE TABLE assets
(
  id serial,
  coin_from character varying(50) NOT NULL,
  coin_to character varying(50) NOT NULL,
  exchange character varying(50) NOT NULL,
  downloader character varying(50) NOT NULL,
  url character varying(100) NOT NULL,
  CONSTRAINT assets_pk PRIMARY KEY (id),
  CONSTRAINT assets_unique_key UNIQUE (coin_from, coin_to, exchange)
);

CREATE INDEX assets_downloader_index ON assets(downloader);
