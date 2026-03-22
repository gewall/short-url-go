CREATE TABLE IF NOT EXISTS clicks (
  id         BIGSERIAL PRIMARY KEY,
  link_id    UUID REFERENCES links(id) ON DELETE CASCADE,
  clicked_at TIMESTAMPTZ DEFAULT NOW(),
  ip_hash    TEXT,
  country    VARCHAR(2),
  city       TEXT,
  device     TEXT,
  browser    TEXT,
  os         TEXT,
  referrer   TEXT
);
CREATE TABLE IF NOT EXISTS click_stats_hourly (
  id          BIGSERIAL PRIMARY KEY,
  link_id     UUID REFERENCES links(id) ON DELETE CASCADE,
  hour        TIMESTAMPTZ NOT NULL,
  click_count INTEGER DEFAULT 0,

  UNIQUE(link_id, hour)
);
CREATE TABLE IF NOT EXISTS click_stats_country (
  id          BIGSERIAL PRIMARY KEY,
  link_id     UUID REFERENCES links(id) ON DELETE CASCADE,
  country     VARCHAR(2) NOT NULL,
  click_count INTEGER DEFAULT 0,
  date        DATE NOT NULL,

  UNIQUE(link_id, country, date)
);
CREATE TABLE IF NOT EXISTS click_stats_device (
  id          BIGSERIAL PRIMARY KEY,
  link_id     UUID REFERENCES links(id) ON DELETE CASCADE,
  device      TEXT NOT NULL,   -- 'mobile' | 'desktop' | 'tablet'
  browser     TEXT NOT NULL,   -- 'chrome' | 'safari' | 'firefox'
  os          TEXT NOT NULL,   -- 'android' | 'ios' | 'windows'
  click_count INTEGER DEFAULT 0,
  date        DATE NOT NULL,

  UNIQUE(link_id, device, browser, os, date)
);

ALTER TABLE clicks
ADD CONSTRAINT fk_click_link
FOREIGN KEY (link_id)
REFERENCES links(id)
ON DELETE CASCADE;

ALTER TABLE click_stats_country
ADD CONSTRAINT fk_stats_country_link
FOREIGN KEY (link_id)
REFERENCES links(id)
ON DELETE CASCADE;

ALTER TABLE click_stats_device
ADD CONSTRAINT fk_stats_device_link
FOREIGN KEY (link_id)
REFERENCES links(id)
ON DELETE CASCADE;

ALTER TABLE click_stats_hourly
ADD CONSTRAINT fk_stats_hourly_link
FOREIGN KEY (link_id)
REFERENCES links(id)
ON DELETE CASCADE;
