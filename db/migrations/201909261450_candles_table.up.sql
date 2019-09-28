CREATE TABLE candles
(
  id serial,
  asset_id int references assets(id),
  close_time timestamp NOT NULL,
  open_price numeric NOT NULL,
  high_price numeric NOT NULL,
  low_price numeric NOT NULL,
  close_price numeric NOT NULL,
  volume  numeric NOT NULL,
  CONSTRAINT candles_pk PRIMARY KEY (id),
  CONSTRAINT candles_unique_key UNIQUE (asset_id, close_time)
);
